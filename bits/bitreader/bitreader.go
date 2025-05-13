package bitreader

type BitReader struct {
	Buff        []byte
	CurrentByte int
	BitPos      uint8
}

func New(b []byte) *BitReader {
	return &BitReader{
		Buff:        b,
		CurrentByte: 0,
		BitPos:      0,
	}
}

func (br *BitReader) Next() bool {

	b := br.Buff[br.CurrentByte]
	bit := (b >> (7 - br.BitPos) & 1) == 1

	br.BitPos += 1
	if br.BitPos == 8 {
		br.BitPos = 0
		br.CurrentByte += 1
	}

	return bit
}

func (br *BitReader) ReadNBits(n int) []bool {
	bits := make([]bool, 0, n)

	for range n {
		bits = append(bits, br.Next())
	}

	return bits
}

func (br *BitReader) ReadByte() (byte, error) {
	var b byte

	for i := range 8 {
		if br.Next() {
			b |= 1 << (7 - i)
		}
	}
	return b, nil
}
