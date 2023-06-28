package helpers

type IndexFile struct {
	File       string `toml:"file"`
	Hash       string `toml:"hash,omitempty"`
	HashFormat string `toml:"hash-format,omitempty"`
	Alias      string `toml:"alias,omitempty"`
	MetaFile   bool   `toml:"metafile,omitempty"`
	Preserve   bool   `toml:"preserve,omitempty"`
	fileFound  bool
}

type IndexPathHolder interface {
	UpdateHash(hash string, format string)
	MarkFound()
	MarkMetaFile()
	MarkedFound() bool
	IsMetaFile() bool
}

func (i *IndexFile) UpdateHash(hash string, format string) {
	i.Hash = hash
	i.HashFormat = format
}

func (i *IndexFile) MarkFound() {
	i.fileFound = true
}

func (i *IndexFile) MarkMetaFile() {
	i.MetaFile = true
}

func (i *IndexFile) MarkedFound() bool {
	return i.fileFound
}

func (i *IndexFile) IsMetaFile() bool {
	return i.MetaFile
}
