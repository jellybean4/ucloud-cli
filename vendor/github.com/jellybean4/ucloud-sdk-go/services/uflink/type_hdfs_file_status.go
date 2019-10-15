package uflink

type HDFSFileStatus struct {
	Path string

	ModificationTime int

	Length int

	BlockSize int
}
