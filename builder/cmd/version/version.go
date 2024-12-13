package version

import (
	"github.com/spf13/cobra"
)

var VersionCommand = CreateVersionCommand()

func CreateVersionCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "version",
		Short:   "Version-related targets.",
		Long:    `Version-related targets.`,
		Aliases: []string{"v", "ver"},
	}

	cmd.AddCommand(InitCommand)
	cmd.AddCommand(MigrateCommand)
	cmd.AddCommand(ExpandCommand)

	return &cmd
}
