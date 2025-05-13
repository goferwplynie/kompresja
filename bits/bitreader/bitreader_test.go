package bitreader

import "testing"

func TestReadByte(t *testing.T) {
	br := New([]byte{8})

	b, _ := br.ReadByte()
	if b != 8 {
		t.Errorf("wanted %v but got %v", 8, b)
	}
}
