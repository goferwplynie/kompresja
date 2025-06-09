package compression

import (
	"encoding/binary"
	"math"
	"os"

	"github.com/goferwplynie/kompresja/archive"
	"github.com/goferwplynie/kompresja/bits/bitbuffer"
	"github.com/goferwplynie/kompresja/internal/algorithm/bwt"
	"github.com/goferwplynie/kompresja/internal/algorithm/huffman"
	"github.com/goferwplynie/kompresja/internal/algorithm/mtf"
)

const partSize = 10 * 1024 * 1024

func Compress(file archive.File) []byte {
	f, err := os.Open(file.Path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		panic(err)
	}
	maxParts := math.Ceil(float64(fInfo.Size()) / float64(partSize))
	size := int64(partSize)
	if file.Part == int(maxParts-1) {
		size = fInfo.Size() % partSize
	}

	data := make([]byte, size)
	_, err = f.ReadAt(data, int64(file.Part)*partSize)

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
