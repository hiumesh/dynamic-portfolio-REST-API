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
	Use: "gotrue",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&serveCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "base configuration file to load")
	rootCmd.PersistentFlags().StringVarP(&watchDir, "config-dir", "d", "", "directory containing a sorted list of config files to watch for changes")
	return &rootCmd
}

func loadGlobalConfig(ctx context.Context) *config.GlobalConfiguration {
	if ctx == nil {
		panic("context must not be nil")
	}

	config, err := config.LoadGlobal(configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %+v", err)
	}

	if err := observability.ConfigureLogging(&config.LOGGING); err != nil {
		logrus.WithError(err).Error("unable to configure logging")
	}

	return config
}

func execWithConfigAndArgs(cmd *cobra.Command, fn func(config *config.GlobalConfiguration, args []string), args []string) {
	fn(loadGlobalConfig(cmd.Context()), args)
}
