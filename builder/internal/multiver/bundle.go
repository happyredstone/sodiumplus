package multiver

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/all"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
)

func BundlePack(outDir string, cfg *config.Config, root string, version string, loader string) error {
	if !slices.Contains(LowerLoaders, strings.ToLower(loader)) {
		return fmt.Errorf("invalid loader: %s", loader)
	}

	stub, err := LoadPackStub(version, MappedLoaders[strings.ToLower(loader)])

	if err != nil {
		return err
	}

	pack, err := stub.InitRealPack(*cfg)

	if err != nil {
		return err
	}

	err = all.BundleFor(outDir, root, cfg, pack)

	if err != nil {
		return err
	}

	return os.Chdir(root)
}
