// Package injector implements the sidecar injector
package injector

import (
	"github.com/leonardoce/backup-filepath/internal/injector"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd is the "injector" subcommand
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "injector",
		Short: "This command starts the backup adapter injector",
		RunE: func(cmd *cobra.Command, args []string) error {

			webhook := injector.New()
			return webhook.Run(cmd.Context())
		},
	}

	cmd.Flags().String(
		"listen-addresses",
		":8000",
		"The default port where the web server is listening",
	)
	_ = viper.BindPFlag("listen-addresses", cmd.Flags().Lookup("listen-addresses"))

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
