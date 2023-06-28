package bundle

import "github.com/spf13/cobra"

var BundleCommand = CreateBundleCommand()

func CreateBundleCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "bundle",
		Short:   "Bundle client targets.",
		Long:    `Bundle client targets.`,
		Aliases: []string{"b"},
	}

	cmd.AddCommand(AllCommand)
	cmd.AddCommand(CurseforgeCommand)
	cmd.AddCommand(ModrinthCommand)

	return &cmd
}
