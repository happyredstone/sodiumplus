package helpers

import (
	"github.com/packwiz/packwiz/core"
)

func GetPack() (*core.Pack, *Index, error) {
	Setup()

	pack, err := core.LoadPack()

	if err != nil {
		return nil, nil, err
	}

	index, err := LoadIndex(pack)

	if err != nil {
		return nil, nil, err
	}

	err = index.Refresh()

	if err != nil {
		return nil, nil, err
	}

	err = index.Write()

	if err != nil {
		return nil, nil, err
	}

	err = pack.UpdateIndexHash()

	if err != nil {
		return nil, nil, err
	}

	err = pack.Write()

	if err != nil {
		return nil, nil, err
	}

	return &pack, &index, nil
}
