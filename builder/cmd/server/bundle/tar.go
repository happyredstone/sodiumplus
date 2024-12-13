package bundle

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/server"
	"github.com/spf13/cobra"
)

var TarCommand = &cobra.Command{
	Use:          "tar",
	Aliases:      []string{"t"},
	Short:        "Create the server Tar bundle.",
	Long:         `Create the server Tar bundle.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return server.Tar()
	},
}
