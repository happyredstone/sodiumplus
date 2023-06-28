package bundle

import (
	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/client"
	"github.com/spf13/cobra"
)

var ModrinthCommand = &cobra.Command{
	Use:          "modrinth",
	Aliases:      []string{"mr"},
	Short:        "Create the client Modrinth bundle.",
	Long:         `Create the client Modrinth bundle.`,
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Modrinth()
	},
}
