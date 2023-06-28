package helpers

import "golang.org/x/exp/slices"

type IndexFiles map[string]IndexPathHolder

func (f *IndexFiles) ToTomlRep() IndexFilesTomlRepresentation {
	rep := make(IndexFilesTomlRepresentation, 0, len(*f))

	for _, v := range *f {
		if file, ok := v.(*IndexFile); ok {
			rep = append(rep, *file)
		} else if file, ok := v.(*IndexFileMultipleAlias); ok {
			for _, alias := range *file {
				rep = append(rep, alias)
			}
		} else {
			panic("Unknown type in IndexFiles")
		}
	}

	slices.SortFunc(rep, func(a IndexFile, b IndexFile) bool {
		if a.File == b.File {
			return a.Alias < b.Alias
		} else {
			return a.File < b.File
		}
	})

	return rep
}

func (f *IndexFiles) UpdateFileEntry(path string, format string, hash string, markAsMetaFile bool) {
	if *f == nil {
		*f = make(IndexFiles)
	}

	file, found := (*f)[path]

	if found {
		file.MarkFound()
		file.UpdateHash(hash, format)

		if markAsMetaFile {
			file.MarkMetaFile()
		}
	} else {
		newFile := IndexFile{
			File:       path,
			Hash:       hash,
			HashFormat: format,
			MetaFile:   markAsMetaFile,
			fileFound:  true,
		}

		(*f)[path] = &newFile
	}
}
