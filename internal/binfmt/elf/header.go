package elf

import (
	"debug/elf"
	"encoding/binary"
	"unsafe"
)

type ElfHeaderCommon interface {
	ByteOrder() binary.ByteOrder

	SetBaseData()
}

var _ ElfHeaderCommon = (*ElfHeader32)(nil)
var _ ElfHeaderCommon = (*ElfHeader64)(nil)

type ElfHeader64 struct {
	elf.Header64
}

func (h ElfHeader64) ByteOrder() binary.ByteOrder { return ByteOrderFromElfByte(h.Ident[elf.EI_DATA]) }

func (h *ElfHeader64) SetByteOrder(ord binary.ByteOrder) *ElfHeader64 {
	h.Ident[elf.EI_DATA] = byte(ByteOrderToElfByte(ord))
	return h
}

func (h *ElfHeader64) SetBaseData() {
	h.Ident = [elf.EI_NIDENT]byte{0x7f, 'E', 'L', 'F'}
	h.Ident[elf.EI_CLASS] = byte(elf.ELFCLASS32)

	h.Ehsize = uint16(unsafe.Sizeof(elf.Header64{}))
	h.Phentsize = uint16(unsafe.Sizeof(elf.Prog64{}))
	h.Shentsize = uint16(unsafe.Sizeof(elf.Section64{}))
	h.Phoff = uint64(h.Ehsize)
}


type ElfHeader32 struct {
	elf.Header32
}

func (h ElfHeader32) ByteOrder() binary.ByteOrder { return ByteOrderFromElfByte(h.Ident[elf.EI_DATA]) }

func (h *ElfHeader32) SetByteOrder(ord binary.ByteOrder) *ElfHeader32 {
	h.Ident[elf.EI_DATA] = byte(ByteOrderToElfByte(ord))
	return h
}




func (h *ElfHeader32) SetBaseData() {
	h.Ident = [elf.EI_NIDENT]byte{0x7f, 'E', 'L', 'F'}
	h.Ident[elf.EI_CLASS] = byte(elf.ELFCLASS32)

	h.Ehsize = uint16(unsafe.Sizeof(elf.Header32{}))
	h.Phentsize = uint16(unsafe.Sizeof(elf.Prog32{}))
	h.Shentsize = uint16(unsafe.Sizeof(elf.Section32{}))
	h.Phoff = uint32(h.Ehsize)
}


