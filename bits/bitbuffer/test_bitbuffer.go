package bitbuffer

import (
	"testing"
)

func TestAddBit(t *testing.T) {
	bb := New()

	bb.AddBit(true)
	if bb.Buffer != 1 {
		t.Errorf("wanted %v but got %v instead", 1, bb.Buffer)
	}
}
