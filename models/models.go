package models

import "github.com/goferwplynie/kompresja/bits/bitbuffer"

type FileMetadata struct {
	BwtIndex    uint32
	Padding     uint8
	TreeSize    uint16
	TreePadding uint8
	Tree        *bitbuffer.BitBuffer
}

type ArchiveMetaData struct {
	PathLen  uint16
	Path     string
	FileSize uint64
	Part     uint32
	Start    uint
}
