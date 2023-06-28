package helpers

import "github.com/packwiz/packwiz/core"

func Refresh() error {
	Setup()

	pack, err := core.LoadPack()

	if err != nil {
		return err
	}

	index, err := LoadIndex(pack)

	if err != nil {
		return err
	}

	err = index.Refresh()

	if err != nil {
		return err
	}

	err = index.Write()

	if err != nil {
		return err
	}

	err = pack.UpdateIndexHash()

	if err != nil {
		return err
	}

	err = pack.Write()

	if err != nil {
		return err
	}

	return nil
}
