// Package injector implements the sidecar injector
package injector

import (
	"github.com/leonardoce/backup-filepath/internal/injector"

	"github.com/spf13/cobra"
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

	return cmd
}
