package version

import (
	"fmt"
	"slices"
	"strings"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/spf13/cobra"
)

var ExpandCommand = CreateExpandCommand()

var MappedLoaders = map[string]string{
	"fabric":     "Fabric",
	"forge":      "Forge",
	"neoforge":   "NeoForge",
	"quilt":      "Quilt",
	"liteloader": "LiteLoader",
}

var LowerLoaders = []string{"fabric", "forge", "neoforge", "quilt", "liteloader"}

func CreateExpandCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "expand [minecraft version] [loader]",
		Short:   "Expand a pack stub into its full form.",
		Long:    `Expand a pack stub into its full form.`,
		Aliases: []string{"x"},
		Args:    cobra.MinimumNArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetConfig()

			if err != nil {
				return err
			}

			if !slices.Contains(LowerLoaders, strings.ToLower(args[1])) {
				return fmt.Errorf("invalid loader: %s", args[1])
			}

			stub, err := multiver.LoadPackStub(args[0], MappedLoaders[strings.ToLower(args[1])])

			if err != nil {
				return err
			}

			_, err = stub.InitRealPack(*cfg)

			return err
		},
	}

	return &cmd
}
