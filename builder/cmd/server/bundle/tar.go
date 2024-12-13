package bundle

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/server"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var TarCommand = &cobra.Command{
	Use:          "tar",
	Aliases:      []string{"t"},
	Short:        "Create the server Tar bundle.",
	Long:         `Create the server Tar bundle.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		pack, err := core.LoadPack()

		if err != nil {
			return err
		}

		return server.Tar(&pack)
	},
}
