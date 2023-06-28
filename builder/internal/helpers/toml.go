package helpers

import "path"

type IndexTomlRepresentation struct {
	HashFormat string                       `toml:"hash-format"`
	Files      IndexFilesTomlRepresentation `toml:"files"`
}

type IndexFilesTomlRepresentation []IndexFile

func (rep IndexFilesTomlRepresentation) ToMemoryRep() IndexFiles {
	out := make(IndexFiles)

	for _, v := range rep {
		v := v

		v.File = path.Clean(v.File)
		v.Alias = path.Clean(v.Alias)

		if v.Alias == "." {
			v.Alias = ""
		}

		if existing, ok := out[v.File]; ok {
			if existingFile, ok := existing.(*IndexFile); ok {
				if v.Alias == existingFile.Alias {
					out[v.File] = &v
				} else {
					m := make(IndexFileMultipleAlias)

					m[existingFile.Alias] = *existingFile
					m[v.Alias] = v

					out[v.File] = &m
				}
			} else if existingMap, ok := existing.(*IndexFileMultipleAlias); ok {
				(*existingMap)[v.Alias] = v
			} else {
				panic("Unknown type in IndexFiles")
			}
		} else {
			out[v.File] = &v
		}
	}

	return out
}
