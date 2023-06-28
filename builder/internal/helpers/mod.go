package helpers

import "github.com/packwiz/packwiz/core"

type Mod struct {
	MetaFile   string
	Name       string                            `toml:"name"`
	FileName   string                            `toml:"filename"`
	Side       string                            `toml:"side,omitempty"`
	Download   core.ModDownload                  `toml:"download"`
	Update     map[string]map[string]interface{} `toml:"update"`
	UpdateData map[string]interface{}
	Option     *core.ModOption `toml:"option,omitempty"`
}

func (m Mod) GetParsedUpdateData(updaterName string) (interface{}, bool) {
	upd, ok := m.UpdateData[updaterName]

	return upd, ok
}
