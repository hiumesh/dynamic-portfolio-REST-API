package observability

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
	"github.com/sirupsen/logrus"
)

func AddRequestID(globalConfig *config.GlobalConfiguration) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		id := uuid.Must(uuid.NewV4()).String()
		if globalConfig.API.RequestIDHeader != "" {
			id = ctx.Request.Header.Get(globalConfig.API.RequestIDHeader)
		}

		utilities.WithRequestID(ctx, id)

		ctx.Next()
	})
}

func NewStructuredLogger(logger *logrus.Logger, config *config.GlobalConfiguration) gin.HandlerFunc {
	loggerInstance := structuredLogger{Logger: logger, Config: config}
	return gin.HandlerFunc(func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/health" {
			ctx.Next()
		} else {
			start := time.Now()

			entry := loggerInstance.NewLogEntry(ctx)

			ctx.Set("logger", &entry)

			ctx.Next()

			latency := time.Since(start)
			bodySize := ctx.Writer.Size()
			statusCode := ctx.Writer.Status()
			responseHeader := ctx.Writer.Header()

			entry.Write(statusCode, bodySize, responseHeader, latency, nil)
		}
	})

}

type structuredLogger struct {
	Logger *logrus.Logger
	Config *config.GlobalConfiguration
}

func (l *structuredLogger) NewLogEntry(ctx *gin.Context) *logEntry {
	referrer := utilities.GetReferrer(ctx.Request, l.Config)
	e := &logEntry{Entry: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{
		"component":   "api",
		"method":      ctx.Request.Method,
		"path":        ctx.Request.URL.Path,
		"remote_addr": utilities.GetIPAddress(ctx.Request),
		"referer":     referrer,
	}

	if reqID := utilities.GetRequestID(ctx); reqID != "" {
		logFields["request_id"] = reqID
	}

	e.Entry = e.Entry.WithFields(logFields)
	return e
}

type logEntry struct {
	Entry *logrus.Entry
}

func (e *logEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	fields := logrus.Fields{
		"status":   status,
		"duration": elapsed.Nanoseconds(),
	}

	errorCode := header.Get("x-sb-error-code")
	if errorCode != "" {
		fields["error_code"] = errorCode
	}

	entry := e.Entry.WithFields(fields)
	entry.Info("request completed")
	e.Entry = entry
}

func (e *logEntry) Panic(v interface{}, stack []byte) {
	entry := e.Entry.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
	entry.Error("request panicked")
	e.Entry = entry
}

func GetLogEntry(ctx *gin.Context) *logEntry {
	v, exists := ctx.Get("logger")
	if !exists {
		return &logEntry{Entry: logrus.NewEntry(logrus.StandardLogger())}
	}

	l, ok := v.(*logEntry)
	if !ok {
		return &logEntry{Entry: logrus.NewEntry(logrus.StandardLogger())}
	}

	return l
}

func LogEntrySetField(ctx *gin.Context, key string, value interface{}) {
	v, exists := ctx.Get("logger")
	if !exists {
		return
	}

	l, ok := v.(*logEntry)
	if !ok {
		return
	}

	l.Entry = l.Entry.WithField(key, value)
}

func LogEntrySetFields(ctx *gin.Context, fields logrus.Fields) {
	v, exists := ctx.Get("logger")
	if !exists {
		return
	}

	l, ok := v.(*logEntry)
	if !ok {
		return
	}

	l.Entry = l.Entry.WithFields(fields)
}
