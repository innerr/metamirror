package tools

import (
	"bytes"
	"testing"
)

func TestCodec(t *testing.T) {
	keyword := []byte("12345678123456781234567812345678")
	str := "test data 12345678"
	data := []byte(str)
	coded := Encode(bytes.NewReader(data), uint32(len(data)), true, keyword)
	decoded := Decode(coded, keyword)
	println(string(decoded))
}

