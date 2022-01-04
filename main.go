package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os/exec"
)

type Endian struct {
	binary.ByteOrder
}

func mustBytes(v interface{}, ord binary.ByteOrder) []byte {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, ord, v)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

type startProcessArgs struct {
	// TODO: args, env...
}

// type guards
var _ func(cmd *exec.Cmd) (*ProcMemReader, error) = NewProcMemReader
var _ io.ReaderAt = (*ProcMemReader)(nil)

