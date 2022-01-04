package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type readerHelper struct {
	binary.ByteOrder
	PtrParser
	io.ReaderAt
	pos int64
}

func (r readerHelper) UintAt(pos int64, sz int) (uint64, error) {
	var b = make([]byte, sz)
	_, err := r.ReadAt(b, pos)
	if err != nil {
		return 0, err
	}
	switch sz {
	case 8:
		return r.Uint64(b), err
	case 4:
		return uint64(r.Uint32(b)), err
	default:
		return 0, errors.New("unknown size")
	}
}

func (r readerHelper) PtrAt(pos int64) (int64, error) {
	b := make([]byte, r.PtrParser.Len())
	_, err := r.ReadAt(b, pos)
	if err != nil {
		return 0, err
	}
	return r.ParsePtr(b), nil
}

func (r *readerHelper) ReadPtr() (int64, error) {
	v, err := r.PtrAt(r.pos)
	// TODO: should I advance if err != nil?
	r.pos += int64(r.PtrParser.Len())
	return v, err
}

func (r readerHelper) CStringAt(pos int64, maxSize int) (string, error) {
	b := make([]byte, maxSize)
	// it reports error but reads data correctly
	n, err := r.ReadAt(b, pos)
	if n == 0 {
		return "", err
	}
	end := bytes.IndexByte(b, 0)
	if end >= n {
		return string(b[:n]), fmt.Errorf("can't find end of cstring")
	}
	// TODO: turn non-nil err into warning
	return string(b[:end]), nil
}

type PtrParser interface {
	ParsePtr([]byte) int64
	Len() int
}

var parserFactoryByPtrSize = map[int]func(order binary.ByteOrder) PtrParser{
	4: func(order binary.ByteOrder) PtrParser {
		return PtrParser32{order: order}
	},
	8: func(order binary.ByteOrder) PtrParser {
		return PtrParser64{order: order}
	},
}

type PtrParser32 struct {
	order binary.ByteOrder
}

func (p PtrParser32) ParsePtr(b []byte) int64 {
	return int64(p.order.Uint32(b))
}

func (p PtrParser32) Len() int {
	return 4
}

type PtrParser64 struct {
	order binary.ByteOrder
}

func (p PtrParser64) ParsePtr(b []byte) int64 {
	return int64(p.order.Uint64(b))
}

func (p PtrParser64) Len() int {
	return 8
}











