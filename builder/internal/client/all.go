package client

import "github.com/packwiz/packwiz/core"

func Bundle(outDir string, pack *core.Pack) error {
	err := CurseForge(outDir, pack)

	if err != nil {
		return err
	}

	err = Modrinth(outDir, pack)

	if err != nil {
		return err
	}

	return nil
}
