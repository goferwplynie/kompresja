package mtf

import (
	"bytes"
	"testing"
)

func TestMTFEncodeDecode(t *testing.T) {
	original := []byte("banana_bandana")

	encoded := Encode(original)
	decoded := Decode(encoded)

	if !bytes.Equal(original, decoded) {
		t.Errorf("decoded data does not match original\noriginal: %v\ndecoded: %v", original, decoded)
	}
}
