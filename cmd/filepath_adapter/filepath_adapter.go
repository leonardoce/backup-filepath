// Package filepath_adapter contain the implementation of the main command
package filepath_adapter

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leonardoce/backup-filepath/cmd/filepath_adapter/server"
	"github.com/leonardoce/backup-filepath/internal/logging"
)

// Cmd is the "filepath_adapter" command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "filepath_adapter",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := logging.IntoContext(
				cmd.Context(),
				viper.GetBool("debug"))
			cmd.SetContext(ctx)
		},
	}

	cmd.AddCommand(server.Cmd())

	cmd.PersistentFlags().Bool(
		"debug",
		true,
		"Enable debugging mode",
	)
	_ = viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))

	cmd.PersistentFlags().String(
		"dsn",
		"dbname=filepath_adapter",
		"The DSN to be used to connect to the PostgreSQL database",
	)
	_ = viper.BindPFlag("dsn", cmd.PersistentFlags().Lookup("dsn"))

	viper.SetEnvPrefix("filepath_adapter")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return cmd
}
