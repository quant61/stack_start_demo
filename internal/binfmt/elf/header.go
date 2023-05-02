package elf

import (
	"debug/elf"
	"encoding/binary"
)

type ElfHeader64 struct {
	elf.Header64
}

func (h ElfHeader64) ByteOrder() binary.ByteOrder { return ByteOrderFromElfByte(h.Ident[elf.EI_DATA]) }

type ElfHeader32 struct {
	elf.Header32
}

func (h ElfHeader32) ByteOrder() binary.ByteOrder { return ByteOrderFromElfByte(h.Ident[elf.EI_DATA]) }
