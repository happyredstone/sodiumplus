package bundle

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/client"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var AllCommand = &cobra.Command{
	Use:          "all",
	Aliases:      []string{"a"},
	Short:        "Create all client bundles.",
	Long:         `Create all client bundles.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		pack, err := core.LoadPack()

		if err != nil {
			return err
		}

		return client.Bundle("output", &pack)
	},
}
