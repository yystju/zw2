package zw2

import (
	"os"
)

type HttpFileWrapper struct {
	data []byte
	pos  int

	node *ZipFileNode
}

// ---- http.File methods ----

func (wrapper *HttpFileWrapper) Close() error {
	wrapper.pos = 0
	return nil
}

func (wrapper *HttpFileWrapper) Read(p []byte) (int, error) {
	if wrapper.node.Info.IsDir() {
		return 0, nil
	}

	n := copy(p, wrapper.data[wrapper.pos:])

	wrapper.pos += n

	return n, nil
}

func (wrapper *HttpFileWrapper) Seek(offset int64, whence int) (int64, error) {
	if wrapper.node.Info.IsDir() {
		return 0, nil
	}

	o := int(offset)
	l := len(wrapper.data)

	var i int = wrapper.pos

	switch whence {
	case 0:
		i = o
	case 1:
		i = (o + wrapper.pos)
	case 2:
		i = (len(wrapper.data) + o)
	}

	if 0 <= i && i < l {
		wrapper.pos = i
	} else if 0 > i {
		wrapper.pos = 0
	} else {
		wrapper.pos = len(wrapper.data)
	}

	return int64(wrapper.pos), nil
}

func (wrapper *HttpFileWrapper) Readdir(count int) ([]os.FileInfo, error) {
	if !wrapper.node.Info.IsDir() {
		return nil, nil
	}

	r := make([]os.FileInfo, 0)

	for _, child := range wrapper.node.Children {
		r = append(r, child.Info)
	}

	return r, nil
}

func (wrapper *HttpFileWrapper) Stat() (os.FileInfo, error) {
	return wrapper.node.Info, nil
}

// ---- Factory method ----

func NewHttpFileWrapper(node *ZipFileNode) (*HttpFileWrapper, error) {
	r := new(HttpFileWrapper)

	r.data = node.Payload
	r.pos = 0

	r.node = node

	return r, nil
}
