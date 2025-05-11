package models

type FileMetadata struct {
	BwtIndex    uint32
	Padding     byte
	TreeSize    uint32
	TreePadding byte
	Tree        string
}
