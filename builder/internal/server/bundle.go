package server

import "github.com/packwiz/packwiz/core"

func Bundle(outDir string, pack *core.Pack) error {
	err := Tar(pack)

	if err != nil {
		return err
	}

	err = Zip(pack)

	if err != nil {
		return err
	}

	err = Rename(outDir, pack)

	if err != nil {
		return err
	}

	return nil
}

func CleanBundle(outDir string, pack *core.Pack) error {
	err := Clean(outDir, pack)

	if err != nil {
		return err
	}

	err = Bundle(outDir, pack)

	if err != nil {
		return err
	}

	return nil
}
