package zw2

import (
	"io"
)

type ZipBytes struct {
	b []byte
}

func (z *ZipBytes) Bytes() []byte {
	return z.b
}

func (z *ZipBytes) Length() int {
	return len(z.b)
}

func (z *ZipBytes) ReadAt(p []byte, off int64) (int, error) {
	//NOTE: This is because the zip file is assumed to be not too large as an in memory file.
	ioff := int(off)

	pos := -1
	n := 0

	if ioff+len(p) < len(z.Bytes()) {
		n = len(p)
		pos = ioff
	} else {
		n = len(z.Bytes()) - ioff
		pos = ioff
	}

	if n > 0 {
		copy(p, z.Bytes()[pos:pos+n])
		return n, nil
	} else {
		return 0, io.EOF
	}
}

func Hello() string {
	return "Hello"
}

func NewZipBytes(b []byte) *ZipBytes {
	z := new(ZipBytes)

	z.b = b

	return z
}
