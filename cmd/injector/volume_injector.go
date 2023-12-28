package injector

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leonardoce/backup-filepath/internal/injector"
	"github.com/leonardoce/backup-filepath/internal/logging"
)

// Cmd is the "volume_injector" command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "volume_injector",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := logging.IntoContext(
				cmd.Context(),
				viper.GetBool("debug"))
			cmd.SetContext(ctx)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			webhook := injector.New()
			return webhook.Run(cmd.Context())
		},
	}

	cmd.PersistentFlags().Bool(
		"debug",
		true,
		"Enable debugging mode",
	)
	_ = viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))

	viper.SetEnvPrefix("volume_injector")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return cmd
}
