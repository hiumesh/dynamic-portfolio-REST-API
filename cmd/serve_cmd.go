package cmd

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/api"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/observability"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/reloader"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func serve(ctx context.Context) {
	if err := config.LoadFile(configFile); err != nil {
		logrus.WithError(err).Fatal("unable to load config")
	}

	if err := config.LoadDirectory(watchDir); err != nil {
		logrus.WithError(err).Fatal("unable to load config from watch dir")
	}

	conf, err := config.LoadGlobalFromEnv()
	if err != nil {
		logrus.WithError(err).Fatal("unable to load config")
	}

	db, err := gorm.Open(postgres.Open(conf.DB.URL), &gorm.Config{Logger: observability.NewGormLogrusLogger(conf.LOGGING.Level, conf.LOGGING.SQL)})
	if err != nil {
		logrus.Fatalf("error opening database: %+v", err)
	}

	logrus.Info(conf.API.Host)
	logrus.Info(conf.API.Port)

	addr := net.JoinHostPort(conf.API.Host, conf.API.Port)
	logrus.Infof("Dynamic Portfolio API started on: %s", addr)

	a := api.NewAPIWithVersion(conf, db, utilities.Version)
	ah := reloader.NewAtomicHandler(a)

	// req := httptest.NewRequest(http.MethodGet, "/health", nil)
	// w := httptest.NewRecorder()

	// a.ServeHTTP(w, req)
	// ah.ServeHTTP(w, req)

	baseCtx, baseCancel := context.WithCancel(context.Background())
	defer baseCancel()

	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           ah,
		ReadHeaderTimeout: 2 * time.Second, // to mitigate a Slowloris attack
		BaseContext: func(net.Listener) context.Context {
			return baseCtx
		},
	}
	log := logrus.WithField("component", "api")

	var wg sync.WaitGroup
	defer wg.Wait() // Do not return to caller until this goroutine is done.

	if watchDir != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			fn := func(latestCfg *config.GlobalConfiguration) {
				log.Info("reloading api with new configuration")
				latestAPI := api.NewAPIWithVersion(
					latestCfg, db, utilities.Version)
				ah.Store(latestAPI)
			}

			rl := reloader.NewReloader(watchDir)
			if err := rl.Watch(ctx, fn); err != nil {
				log.WithError(err).Error("watcher is exiting")
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		defer baseCancel() // close baseContext

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
		defer shutdownCancel()

		if err := httpSrv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			log.WithError(err).Error("shutdown failed")
		}
	}()

	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("http server listen failed")
	}
}
