package zw2

import (
	"log"
	"net/http"
	"os"
)

type HttpFileSystem struct {
	root *ZipFileNode
}

func (fs *HttpFileSystem) Open(name string) (http.File, error) {
	log.Println("name : ", name)

	node := fs.root.FindChildByPathName(name)

	// log.Println("node : ", node)

	if node == nil {
		return nil, os.ErrNotExist
	}

	// log.Println("wrapper : ", wrapper)

	return NewHttpFileWrapper(node)
}

func NewHttpFileSystem(root *ZipFileNode) (*HttpFileSystem, error) {
	fs := new(HttpFileSystem)

	fs.root = root

	return fs, nil
}
