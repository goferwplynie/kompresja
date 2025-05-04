package bitutils

import "testing"

func TestNumberToBitString(t *testing.T) {
	bits := NumberToBitString(5, 8)
	wanted := "00000101"

	if bits != wanted {
		t.Errorf("wanted: %v\ngot %v instead", wanted, bits)
	}
}

func TestByteToBitString(t *testing.T) {
	charByte := byte('A')
	bits := ByteToBitString(charByte)
	wanted := "01000001"
	if bits != wanted {
		t.Errorf("wanted: %v\ngot %v instead", wanted, bits)
	}
}

func TestBitStringToBytes(t *testing.T) {
	bitString := "0100010001000001"
	wanted := []byte{byte('D'), byte('A')}
	bits := BitStringToBytes(bitString)

	if bits[0] != wanted[0] || bits[1] != wanted[1] {
		t.Errorf("wanted: %v\ngot %v instead", wanted, bits)
	}
}

func TestBytesToBitString(t *testing.T) {
	bytes := []byte{byte('D'), byte('A')}
	wanted := "0100010001000001"
	bitString := BytesToBitString(bytes)

	if bitString != wanted {
		t.Errorf("wanted: %v\ngot %v instead", wanted, bitString)
	}
}
