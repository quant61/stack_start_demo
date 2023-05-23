package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ReaderHelper struct {
	binary.ByteOrder
	PtrSize int
	io.ReaderAt
	Pos int64
}

func (r ReaderHelper) UintAt(pos int64, sz int) (uint64, error) {
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

func (r ReaderHelper) PtrAt(pos int64) (int64, error) {
	b := make([]byte, r.PtrSize)
	_, err := r.ReadAt(b, pos)
	if err != nil {
		return 0, err
	}
	return r.ParsePtr(b)
}

func (r ReaderHelper) ParsePtr(b []byte) (int64, error) {
	switch r.PtrSize {
	case 8:
		return int64(r.ByteOrder.Uint64(b)), nil
	case 4:
		return int64(r.ByteOrder.Uint32(b)), nil
	default:
		return 0, fmt.Errorf("unsupported size %d", r.PtrSize)
	}
}

func (r *ReaderHelper) ReadPtr() (int64, error) {
	v, err := r.PtrAt(r.Pos)
	// TODO: should I advance if err != nil?
	r.Pos += int64(r.PtrSize)
	return v, err
}

func (r ReaderHelper) CStringAt(pos int64, maxSize int) (string, error) {
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
