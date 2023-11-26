// Package server implements the web server
package server

import (
	"net"

	"github.com/leonardoce/backup-adapter/pkg/adapter"
	"github.com/leonardoce/backup-filepath/internal/logging"
	"github.com/leonardoce/backup-filepath/internal/walmanager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// Cmd is the "server" subcommand
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "This command starts the WAL management server",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.FromContext(cmd.Context())
			grpcServer := grpc.NewServer()

			basePath := viper.GetString("base-path")
			walManagerImpl := walmanager.NewWalManagerImplementation(basePath)
			adapter.RegisterWalManagerServer(
				grpcServer,
				walManagerImpl)

			listener, err := net.Listen(
				viper.GetString("listening-network"),
				viper.GetString("listening-address"),
			)
			if err != nil {
				logger.Error(err, "While starting server")
				return err
			}

			logger.Info("Starting WAL Management server", "basePath", basePath)
			err = grpcServer.Serve(listener)
			if err != nil {
				logger.Error(err, "While terminatind server")
			}

			return err
		},
	}

	cmd.Flags().String(
		"base-path",
		"/wal-path",
		"The path from where to store WAL files",
	)
	_ = viper.BindPFlag("base-path", cmd.Flags().Lookup("base-path"))

	cmd.Flags().String(
		"listening-network",
		"unix",
		`network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket"`,
	)
	_ = viper.BindPFlag("listening-network", cmd.Flags().Lookup("listening-network"))

	cmd.Flags().String(
		"listening-address",
		"/controller/walmanager.socket",
		`listening address, whose meaning depends on "listening-network"`,
	)
	_ = viper.BindPFlag("listening-address", cmd.Flags().Lookup("listening-address"))

	return cmd
}
