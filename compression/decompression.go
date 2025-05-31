package compression

import (
	"encoding/binary"

	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/internal/algorithm/bwt"
	"github.com/goferwplynie/kompresja/internal/algorithm/huffman"
	"github.com/goferwplynie/kompresja/internal/algorithm/mtf"
	"github.com/goferwplynie/kompresja/models"
)

func Decompress(by []byte) []byte {
	meta, data := extractData(by)
	b := huffman.Decode(data, meta)
	b = mtf.Decode(b)
	b = bwt.Decode(b, int(meta.BwtIndex))
	return b
}

func extractData(b []byte) (metadata models.FileMetadata, data *bitbuffer.BitBuffer) {
	byteCount := 0

	metadata.Padding = b[byteCount]
	byteCount += 1

	metadata.TreeSize = binary.BigEndian.Uint16(b[byteCount : byteCount+2])
	byteCount += 2

	metadata.TreePadding = b[byteCount]
	byteCount += 1

	treeBytes := b[byteCount : byteCount+int(metadata.TreeSize)]
	metadata.Tree = bitbuffer.New(treeBytes)

	byteCount += len(metadata.Tree.Bytes)

	metadata.BwtIndex = binary.BigEndian.Uint32(b[byteCount : byteCount+4])
	byteCount += 4

	data = bitbuffer.New(b[byteCount:])

	return metadata, data
}
