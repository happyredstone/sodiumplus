package multiver

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type PackStub struct {
	Minecraft     string `toml:"minecraft"`
	Loader        string `toml:"loader"`
	LoaderVersion string `toml:"loader_version"`
}

func (stub PackStub) Write(path string) error {
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	enc := toml.NewEncoder(f)
	enc.Indent = ""
	err = enc.Encode(stub)

	if err != nil {
		f.Close()
		return err
	}

	return f.Close()
}

func LoadPackStub(ver string, loader string) (*PackStub, error) {
	path := path.Join("versions", ver, loader, "pack.stub.toml")
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var stub PackStub

	err = toml.Unmarshal(content, &stub)

	if err != nil {
		return nil, err
	}

	return &stub, nil
}
