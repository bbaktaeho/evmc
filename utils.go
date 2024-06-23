package evmc

import "unicode/utf8"

func removeInvalidUTF8Bytes(data []byte) []byte {
	if utf8.Valid(data) {
		return data
	}
	valid := make([]byte, 0, len(data))
	for i := 0; i < len(data); {
		r, size := utf8.DecodeRune(data[i:])
		if r != utf8.RuneError || size != 1 {
			valid = append(valid, data[i:i+size]...)
		}
		i += size
	}
	return valid
}
