package main

import (
	"encoding/binary"
	"io"
	"os"
	"os/exec"
)

// Unit
// WIP
// TODO: wrap all logic into class
type Unit struct {
	Order binary.ByteOrder
	Command *exec.Cmd
	Proc os.Process
	PtrSize int

	Mem io.ReaderAt

}




