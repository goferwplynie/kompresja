package bitutils

func NumberToBitString(n, width int) string {
	bitString := ""

	for i := width - 1; i >= 0; i-- {
		if (n>>i)&1 == 1 {
			bitString += "1"
		} else {
			bitString += "0"
		}
	}

	return bitString
}

func ByteToBitString(b byte) string {
	bitString := ""
	for i := 7; i >= 0; i-- {
		if (b>>i)&1 == 1 {

			bitString += "1"
		} else {
			bitString += "0"
		}
	}
	return bitString
}

func BitStringToBytes(bitstring string) []byte {
	bytes := make([]byte, 0, len(bitstring)/8)

	for i := 0; i < len(bitstring); i += 8 {
		chunk := bitstring[i : i+8]
		var b byte
		for _, r := range chunk {
			b <<= 1
			if r == '1' {
				b |= 1
			}
		}
		bytes = append(bytes, b)
	}

	return bytes
}

func BytesToBitString(bytes []byte) string {
	bitString := ""
	for _, b := range bytes {
		bitString += ByteToBitString(b)
	}

	return bitString
}
