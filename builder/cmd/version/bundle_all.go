package version

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/spf13/cobra"
)

var BundleAllCommand = CreateBundleAllCommand()

func CreateBundleAllCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "bundle-all",
		Short:   "Bundle all version of a multiversion pack.",
		Long:    `Bundle all version of a multiversion pack.`,
		Aliases: []string{"B", "export-all"},

		RunE: func(cmd *cobra.Command, args []string) error {
			baseRoot, err := os.Getwd()

			if err != nil {
				return err
			}

			root, err := filepath.Abs(baseRoot)

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

			items, err := multiver.FindVersions(path.Join(root, cfg.MultiVersion.VersionsDir))

			if err != nil {
				return err
			}

			for ver, loaders := range items {
				for _, loader := range loaders {
					fmt.Printf("Bundling for %s %s...\n", loader, ver)

					err = multiver.BundlePack(outDir, cfg, root, ver, loader)

					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	return &cmd
}
