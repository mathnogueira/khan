// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/topfreegames/khan/api"
	"github.com/uber-go/zap"
)

var host string
var port int
var debug bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the khan API server",
	Long: `Starts khan server with the specified arguments. You can use
environment variables to override configuration keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		ll := zap.InfoLevel
		if debug {
			ll = zap.DebugLevel
		}
		l := zap.NewJSON(ll)

		cmdL := l.With(
			zap.String("source", "startCmd"),
			zap.String("operation", "Run"),
			zap.String("host", host),
			zap.Int("port", port),
			zap.Bool("debug", debug),
		)

		cmdL.Debug("Creating application...")
		app := api.GetApp(
			host,
			port,
			ConfigFile,
			debug,
			l,
		)
		cmdL.Debug("Application created successfully.")

		cmdL.Debug("Starting application...")
		app.Start()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&host, "bind", "b", "0.0.0.0", "Host to bind khan to")
	startCmd.Flags().IntVarP(&port, "port", "p", 8888, "Port to bind khan to")
	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
}
