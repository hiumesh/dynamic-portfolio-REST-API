package observability

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

const (
	LOG_SQL_ALL       = "all"
	LOG_SQL_NONE      = "none"
	LOG_SQL_STATEMENT = "statement"
)

var (
	loggingOnce sync.Once
)

type CustomFormatter struct {
	logrus.JSONFormatter
}

func NewCustomFormatter() *CustomFormatter {
	return &CustomFormatter{
		JSONFormatter: logrus.JSONFormatter{
			DisableTimestamp: false,
			TimestampFormat:  time.RFC3339,
		},
	}
}

func NewCustomTextFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339,
	}
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// logrus doesn't support formatting the time in UTC so we need to use a custom formatter
	entry.Time = entry.Time.UTC()
	return f.JSONFormatter.Format(entry)
}

func ConfigureLogging(config *config.LoggingConfig) error {
	var err error

	loggingOnce.Do(func() {
		logrus.SetOutput(os.Stdout)
		formatter := NewCustomFormatter()
		logrus.SetFormatter(formatter)
		// logrus.SetFormatter(NewCustomTextFormatter())

		// use a file if you want
		if config.File != "" {
			f, errOpen := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660) //#nosec G302 -- Log files should be rw-rw-r--
			if errOpen != nil {
				err = errOpen
				return
			}
			logrus.SetOutput(f)
			logrus.Infof("Set output file to %s", config.File)
		}

		if config.Level != "" {
			level, errParse := logrus.ParseLevel(config.Level)
			if err != nil {
				err = errParse
				return
			}
			logrus.SetLevel(level)
			logrus.Debug("Set log level to: " + logrus.GetLevel().String())
		}

		f := logrus.Fields{}
		for k, v := range config.Fields {
			f[k] = v
		}
		logrus.WithFields(f)

	})

	return err
}

type GormLogrusLogger struct {
	gormLog          *logrus.Logger
	sqlLog           *logrus.Logger
	shouldLogSQL     bool
	shouldLogSQLArgs bool
	config           logger.Config
}

func NewGormLogrusLogger(lvl string, sql string) *GormLogrusLogger {
	gormLog := logrus.WithField("component", "gorm")
	sqlLog := logrus.WithField("component", "sql")

	logLvl, _ := logrus.ParseLevel(lvl)

	shouldLogSQL := sql == LOG_SQL_STATEMENT || sql == LOG_SQL_ALL
	shouldLogSQLArgs := sql == LOG_SQL_ALL

	return &GormLogrusLogger{
		gormLog:      gormLog.Logger,
		sqlLog:       sqlLog.Logger,
		shouldLogSQL: shouldLogSQL,
		config: logger.Config{
			LogLevel:             logger.LogLevel(logLvl), // Set GORM logging level here
			Colorful:             true,
			ParameterizedQueries: shouldLogSQLArgs,
		},
	}
}

func (l *GormLogrusLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.config.LogLevel = level
	return &newLogger
}

func (l *GormLogrusLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.gormLog.WithContext(ctx).Infof(msg, data...)
}

func (l *GormLogrusLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.gormLog.WithContext(ctx).Warnf(msg, data...)
}

func (l *GormLogrusLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.gormLog.WithContext(ctx).Errorf(msg, data...)
}

func (l *GormLogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := logrus.Fields{
		"duration": elapsed,
		"rows":     rows,
		"sql":      sql,
	}

	if !l.shouldLogSQL {
		return
	}

	switch {
	case err != nil && l.config.LogLevel >= logger.Error:
		l.sqlLog.WithContext(ctx).WithFields(fields).Error(err)
	case elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0 && l.config.LogLevel >= logger.Warn:
		l.sqlLog.WithContext(ctx).WithFields(fields).Warn("slow query")
	case l.config.LogLevel == logger.Info:
		l.sqlLog.WithContext(ctx).WithFields(fields).Info("query executed")
	}
}
