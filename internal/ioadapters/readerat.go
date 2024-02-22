package ioadapters

import (
	"errors"
	"io"
	"slices"
)

type readerAtAdapter struct {
	r         io.Reader
	readBytes []byte
}

func ToReaderAt(r io.Reader) io.ReaderAt {
	ra, ok := r.(io.ReaderAt)
	if ok {
		return ra
	}
	return &readerAtAdapter{
		r: r,
	}
}

func (r *readerAtAdapter) ReadAt(p []byte, off int64) (n int, err error) {
	if int(off)+len(p) > len(r.readBytes) {
		err := r.expandBuffer(int(off) + len(p))
		if err != nil {
			return 0, err
		}
	}
	return bytesReaderAt(r.readBytes).ReadAt(p, off)
}

func (r *readerAtAdapter) expandBuffer(newSize int) error {
	if cap(r.readBytes) < newSize {
		r.readBytes = slices.Grow(r.readBytes, newSize-cap(r.readBytes))
	}

	newPart := r.readBytes[len(r.readBytes):newSize]
	n, err := r.r.Read(newPart)
	switch {
	case err == nil:
		r.readBytes = r.readBytes[:newSize]
	case errors.Is(err, io.EOF):
		r.readBytes = r.readBytes[:len(r.readBytes)+n]
	default:
		return err
	}
	return nil
}

func BytesReadAt(src []byte, dst []byte, off int64) (n int, err error) {
	return bytesReaderAt(src).ReadAt(dst, off)
}

type bytesReaderAt []byte

func (bra bytesReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	idx := 0
	for i := int(off); i < len(bra) && idx < len(p); i, idx = i+1, idx+1 {
		p[idx] = bra[i]
	}
	if idx != len(p) {
		return idx, io.EOF
	}
	return idx, nil
}
