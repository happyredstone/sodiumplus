package bundle

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/client"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var ModrinthCommand = &cobra.Command{
	Use:          "modrinth",
	Aliases:      []string{"mr"},
	Short:        "Create the client Modrinth bundle.",
	Long:         `Create the client Modrinth bundle.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		pack, err := core.LoadPack()

		if err != nil {
			return err
		}

		return client.Modrinth("output", &pack)
	},
}
