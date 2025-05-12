package bitbuffer

type BitBuffer struct {
	Bytes      []byte
	CurrentBit uint8
	Buffer     byte
}

func New(bytes ...[]byte) *BitBuffer {
	var b []byte
	if len(bytes) > 0 {
		b = bytes[0]
	}
	return &BitBuffer{
		Bytes:      b,
		CurrentBit: 0,
		Buffer:     0,
	}
}

func (bb *BitBuffer) AddBit(b bool) {
	if b {
		bb.Buffer <<= 1
		bb.Buffer |= 1
	} else {
		bb.Buffer <<= 1
	}
	bb.CurrentBit += 1

	if bb.CurrentBit == 8 {
		bb.Flush()
	}

}

func (bb *BitBuffer) Flush() {
	bb.Bytes = append(bb.Bytes, bb.Buffer)
	bb.Buffer = 0
	bb.CurrentBit = 0
}

func (bb *BitBuffer) Finalize() {
	if bb.CurrentBit > 0 {
		bb.Buffer <<= (8 - bb.CurrentBit)
		bb.Flush()
	}
}

func (bb *BitBuffer) AddBits(bitstring string) {
	for _, r := range bitstring {
		bb.AddBit(r == '1')
	}
}

func (bb *BitBuffer) MergeLeft(bb2 *BitBuffer) {
	if bb2.CurrentBit != 0 {
		bb2.Finalize()
	}
	bb.Bytes = append(bb2.Bytes, bb.Bytes...)
}

func (bb *BitBuffer) MergeRight(bb2 *BitBuffer) {
	if bb2.CurrentBit != 0 {
		bb2.Finalize()
	}
	bb.Bytes = append(bb.Bytes, bb2.Bytes...)
}

func (bb *BitBuffer) AddByte(b byte) {
	if bb.CurrentBit == 0 {
		bb.Buffer = b
		bb.Flush()
	} else {
		for i := 7; i > 0; i-- {
			bb.AddBit(b>>i&1 == 1)
		}
	}
}

func (bb *BitBuffer) AddBytes(bytes []byte) {
	for _, b := range bytes {
		bb.AddByte(b)
	}
}
