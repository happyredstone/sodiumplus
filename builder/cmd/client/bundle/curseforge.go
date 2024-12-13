package bundle

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/client"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var CurseforgeCommand = &cobra.Command{
	Use:          "curseforge",
	Aliases:      []string{"cf", "curse"},
	Short:        "Create the client CurseForge bundle.",
	Long:         `Create the client CurseForge bundle.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		pack, err := core.LoadPack()

		if err != nil {
			return err
		}

		return client.CurseForge("output", &pack)
	},
}
