package cmd

import (
	"strconv"

	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/web"
	"github.com/spf13/cobra"
)

var Port int
var ServeCommand = CreateServeCommand()

func CreateServeCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:          "serve",
		Aliases:      []string{"srv"},
		Short:        "Serve the pack files.",
		Long:         `Serve the pack files.`,
		SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			server := web.CreateServer()

			return web.RunServer(server, "0.0.0.0:"+strconv.Itoa(Port))
		},
	}

	cmd.Flags().IntVarP(&Port, "port", "p", 4000, "The port to serve the files on.")

	return &cmd
}
