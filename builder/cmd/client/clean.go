package client

import (
	internalClient "github.com/HappyRedstone/SodiumPlus/builder/internal/client"
	"github.com/spf13/cobra"
)

var CleanCommand = &cobra.Command{
	Use:     "clean",
	Short:   "Clean bundled client files.",
	Long:    `Clean bundled client files.`,
	Aliases: []string{"c"},

	RunE: func(cmd *cobra.Command, args []string) error {
		return internalClient.Clean()
	},
}
