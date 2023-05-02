package main

import (
	"encoding/binary"
	"io"
	"os/exec"
)

type Endian struct {
	binary.ByteOrder
}

type startProcessArgs struct {
	// TODO: args, env...
}

// type guards
var _ func(cmd *exec.Cmd) (*ProcMemReader, error) = NewProcMemReader
var _ io.ReaderAt = (*ProcMemReader)(nil)

