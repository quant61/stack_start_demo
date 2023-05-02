package internal

import (
	"bytes"
	"encoding/binary"
)

func MustBytes(v interface{}, ord binary.ByteOrder) []byte {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, ord, v)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
