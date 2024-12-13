package bundle

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/server"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var ZipCommand = &cobra.Command{
	Use:          "zip",
	Aliases:      []string{"z"},
	Short:        "Create the server Zip bundle.",
	Long:         `Create the server Zip bundle.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		pack, err := core.LoadPack()

		if err != nil {
			return err
		}

		return server.Zip(&pack)
	},
}
