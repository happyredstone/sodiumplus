package version

import (
	"os"
	"path/filepath"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/spf13/cobra"
)

var BundleCommand = CreateBundleCommand()

func CreateBundleCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "bundle [minecraft version] [loader]",
		Short:   "Bundle a version of a multiversion pack.",
		Long:    `Bundle a version of a multiversion pack.`,
		Aliases: []string{"b", "export"},
		Args:    cobra.MinimumNArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := os.Getwd()

			if err != nil {
				return err
			}

			cfg, err := config.GetConfig()

			if err != nil {
				return err
			}

			outDir, err := filepath.Abs(cfg.MultiVersion.OutDir)

			if err != nil {
				return err
			}

			return multiver.BundlePack(outDir, cfg, root, args[0], args[1])
		},
	}

	return &cmd
}
