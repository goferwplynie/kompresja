package bitreader

import "errors"

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

func (br *BitReader) Next() (bool, error) {
	if br.CurrentByte >= len(br.Buff) {
		return false, errors.New("EOF")
	}

	b := br.Buff[br.CurrentByte]
	bit := (b >> (7 - br.BitPos) & 1) == 1

	br.BitPos += 1
	if br.BitPos == 8 {
		br.BitPos = 0
		br.CurrentByte += 1
	}

	return bit, nil
}

func (br *BitReader) ReadNBits(n int) ([]bool, error) {
	bits := make([]bool, 0, n)

	for range n - 1 {
		bit, err := br.Next()
		if err != nil {
			return bits, err
		}
		bits = append(bits, bit)
	}

	return bits, nil
}

func (br *BitReader) ReadByte() (byte, error) {
	var b byte

	for i := range 8 {
		bit, err := br.Next()
		if err != nil {
			return b, err
		}
		if bit {
			b |= 1 << (7 - i)
		}
	}
	return b, nil
}
