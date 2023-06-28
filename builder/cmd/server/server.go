package server

import (
	"github.com/NoSadBeHappy/SodiumPlus/builder/cmd/server/bundle"
	"github.com/spf13/cobra"
)

var ServerCommand = CreateServerCommand()

func CreateServerCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "server",
		Short:   "Server build targets.",
		Long:    `Server build targets.`,
		Aliases: []string{"s"},
	}

	cmd.AddCommand(bundle.BundleCommand)
	cmd.AddCommand(CleanCommand)

	return &cmd
}
