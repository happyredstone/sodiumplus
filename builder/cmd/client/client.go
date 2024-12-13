package client

import (
	"github.com/HappyRedstone/SodiumPlus/builder/cmd/client/bundle"
	"github.com/spf13/cobra"
)

var ClientCommand = CreateClientCommand()

func CreateClientCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "client",
		Short:   "Client build targets.",
		Long:    `Client build targets.`,
		Aliases: []string{"cl"},
	}

	cmd.AddCommand(bundle.BundleCommand)
	cmd.AddCommand(CleanCommand)

	return &cmd
}
