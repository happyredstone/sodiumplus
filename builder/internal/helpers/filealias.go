package helpers

type IndexFileMultipleAlias map[string]IndexFile

func (i *IndexFileMultipleAlias) UpdateHash(hash string, format string) {
	for k, v := range *i {
		v.UpdateHash(hash, format)

		(*i)[k] = v
	}
}

func (i *IndexFileMultipleAlias) MarkFound() {
	for k, v := range *i {
		v.MarkFound()

		(*i)[k] = v
	}
}

func (i *IndexFileMultipleAlias) MarkMetaFile() {
	for k, v := range *i {
		v.MarkMetaFile()

		(*i)[k] = v
	}
}

func (i *IndexFileMultipleAlias) MarkedFound() bool {
	for _, v := range *i {
		return v.MarkedFound()
	}

	panic("No entries in indexFileMultipleAlias")
}

func (i *IndexFileMultipleAlias) IsMetaFile() bool {
	for _, v := range *i {
		return v.MetaFile
	}

	panic("No entries in indexFileMultipleAlias")
}
