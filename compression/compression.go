package compression

import (
	"encoding/binary"
	"os"

	"github.com/goferwplynie/kompresja/archive"
	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/internal/algorithm/bwt"
	"github.com/goferwplynie/kompresja/internal/algorithm/huffman"
	"github.com/goferwplynie/kompresja/internal/algorithm/mtf"
)

func Compress(file archive.File) []byte {
	data, _ := os.ReadFile(file.Path)

	b, bwtIndex := bwt.Encode(data)
	b = mtf.Encode(b)
	bb, treeBuffer := huffman.Encode(b)
	addMetaData(bb, treeBuffer, bwtIndex)

	return bb.Bytes
}

func addMetaData(bitBuffer *bitbuffer.BitBuffer, treeBuffer *bitbuffer.BitBuffer, bwtIndex int) {
	meta := bitbuffer.New()

	padding := (8 - bitBuffer.CurrentBit) % 8
	bitBuffer.Finalize()

	bwtIBytes := binary.BigEndian.AppendUint32(make([]byte, 0, 4), uint32(bwtIndex))

	treePadding := (8 - treeBuffer.CurrentBit) % 8
	treeBuffer.Finalize()

	treeSize := len(treeBuffer.Bytes)
	treeSizeB := binary.BigEndian.AppendUint16(make([]byte, 0, 2), uint16(treeSize))

	meta.AddByte(byte(padding))
	meta.AddBytes(treeSizeB)
	meta.AddByte(byte(treePadding))

	meta.MergeRight(treeBuffer)

	meta.AddBytes(bwtIBytes)

	bitBuffer.MergeLeft(meta)

	//[padding][tree size][tree padding][tree][bwtIndex][data]

}
