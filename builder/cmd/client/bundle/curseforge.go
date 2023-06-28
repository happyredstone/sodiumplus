package bundle

import (
	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/client"
	"github.com/spf13/cobra"
)

var CurseforgeCommand = &cobra.Command{
	Use:          "curseforge",
	Aliases:      []string{"cf", "curse"},
	Short:        "Create the client CurseForge bundle.",
	Long:         `Create the client CurseForge bundle.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return client.CurseForge()
	},
}
