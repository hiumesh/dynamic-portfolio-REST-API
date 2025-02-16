package cmd

import (
	"context"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/observability"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configFile = ""
	watchDir   = ""
)

var rootCmd = cobra.Command{
	Use: "dp",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&serveCmd, &seedCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "base configuration file to load")
	rootCmd.PersistentFlags().StringVarP(&watchDir, "config-dir", "d", "", "directory containing a sorted list of config files to watch for changes")
	return &rootCmd
}

func loadGlobalConfig(ctx context.Context) *config.GlobalConfiguration {
	if ctx == nil {
		panic("context must not be nil")
	}

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

	if err := observability.ConfigureLogging(&conf.LOGGING); err != nil {
		logrus.WithError(err).Error("unable to configure logging")
	}

	return conf
}

func execWithConfigAndArgs(cmd *cobra.Command, fn func(config *config.GlobalConfiguration, args []string), args []string) {
	fn(loadGlobalConfig(cmd.Context()), args)
}
