package models

type FileMetadata struct {
	BwtIndex    uint32
	Padding     uint8
	TreeSize    uint16
	TreePadding uint8
	Tree        string
}
