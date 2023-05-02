package elf

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"github.com/quant61/stack_start_demo/internal"
)

type ElfProg64 = elf.Prog64
type ElfProg32 = elf.Prog32

func ByteOrderFromElfByte(b byte) binary.ByteOrder {
	switch elf.Data(b) {
	case elf.ELFDATA2LSB:
		return binary.LittleEndian
	case elf.ELFDATA2MSB:
		return binary.BigEndian
	default:
		panic(fmt.Sprintf("unknown ELF data encoding: %x", b))
	}
}

func ByteOrderToElfByte(ord binary.ByteOrder) elf.Data {
	switch ord {
	case binary.LittleEndian:
		return elf.ELFDATA2LSB
	case binary.BigEndian:
		return elf.ELFDATA2MSB
	case nil:
		return elf.ELFDATANONE
	default:
		panic(fmt.Sprintf("unknown byte order: %T", ord))
	}
}


func Amd64HeaderLinux() ElfHeader64 {
	var h ElfHeader64
	h.SetBaseData()
	h.SetByteOrder(binary.LittleEndian)

	h.Ident[elf.EI_VERSION] = byte(elf.EV_CURRENT)
	h.Ident[elf.EI_OSABI] = byte(elf.ELFOSABI_NONE)
	h.Ident[elf.EI_ABIVERSION] = 0

	h.Type = uint16(elf.ET_EXEC)
	h.Machine = uint16(elf.EM_X86_64)
	h.Version = uint32(elf.EV_CURRENT)

	return h
}

func X86_32HeaderLinux() ElfHeader32 {
	var h ElfHeader32
	h.SetBaseData()
	h.SetByteOrder(binary.LittleEndian)

	h.Ident[elf.EI_VERSION] = byte(elf.EV_CURRENT)
	h.Ident[elf.EI_OSABI] = byte(elf.ELFOSABI_NONE)
	h.Ident[elf.EI_ABIVERSION] = 0

	h.Type = uint16(elf.ET_EXEC)
	h.Machine = uint16(elf.EM_X86_64)
	h.Version = uint32(elf.EV_CURRENT)


	return h
}

func BuildElfBinary64() ([]byte, binary.ByteOrder) {
	h := Amd64HeaderLinux()
	h.Phnum = 1

	const start_addr = 0x10000
	fileAddrStart := start_addr + int(h.Ehsize) + int(h.Phentsize)*1

	prog0 := elf.Prog64{
		Type:   uint32(elf.PT_LOAD),
		Flags:  uint32(elf.PF_R | elf.PF_X),
		Vaddr:  start_addr,
		Paddr:  start_addr,
		Align:  0x10000,
		Filesz: 512,
		Memsz:  512,
	}
	h.Entry = uint64(fileAddrStart)

	b := internal.MustBytes(h, h.ByteOrder())
	b = append(b, internal.MustBytes(prog0, h.ByteOrder())...)
	// int 3
	b = append(b, 0xcc)
	return b, h.ByteOrder()
}

func BuildElfBinary32() ([]byte, binary.ByteOrder) {
	h := X86_32HeaderLinux()
	h.Phnum = 1

	const start_addr = 0x10000
	fileAddrStart := start_addr + int(h.Ehsize) + int(h.Phentsize)*1

	prog0 := elf.Prog32{
		Type:   uint32(elf.PT_LOAD),
		Flags:  uint32(elf.PF_R | elf.PF_X),
		Vaddr:  start_addr,
		Paddr:  start_addr,
		Align:  0x10000,
		Filesz: 512,
		Memsz:  512,
	}
	h.Entry = uint32(fileAddrStart)

	b := internal.MustBytes(h, h.ByteOrder())
	b = append(b, internal.MustBytes(prog0, h.ByteOrder())...)
	// int 3
	b = append(b, 0xcc)
	return b, h.ByteOrder()
}
