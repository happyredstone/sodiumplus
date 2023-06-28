package server

import (
	internalServer "github.com/NoSadBeHappy/SodiumPlus/builder/internal/server"
	"github.com/spf13/cobra"
)

var CleanCommand = &cobra.Command{
	Use:     "clean",
	Short:   "Clean bundled server files.",
	Long:    `Clean bundled server files.`,
	Aliases: []string{"c"},

	RunE: func(cmd *cobra.Command, args []string) error {
		return internalServer.Clean()
	},
}
