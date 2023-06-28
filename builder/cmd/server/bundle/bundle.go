package bundle

import "github.com/spf13/cobra"

var BundleCommand = CreateBundleCommand()

func CreateBundleCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "bundle",
		Short:   "Bundle server targets.",
		Long:    `Bundle server targets.`,
		Aliases: []string{"b"},
	}

	cmd.AddCommand(AllCommand)
	cmd.AddCommand(TarCommand)
	cmd.AddCommand(ZipCommand)

	return &cmd
}
