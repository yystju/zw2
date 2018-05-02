# Zip File System for HTTP

## This is an exercise of learning Go language.

NOTE: All files in the zip file will be loaded into memory. It is assumed that the content wouldn't be too large.

Example code:

```go
package main

import (
	"archive/zip"
	"bytes"
	"github.com/yystju/zw2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var filePathName string = "./test.zip"

	zipFile, err := os.OpenFile(filePathName, os.O_RDONLY, 0666)

	if err != nil {
		panic(err)
	}

	zipBytes, err := ioutil.ReadAll(zipFile)

	if err != nil {
		panic(err)
	}

	zipFile.Close()

	zb := zw2.NewZipBytes(zipBytes)

	reader, err := zip.NewReader(zb, int64(zb.Length()))

	if err != nil {
		panic(err)
	}

	root := zw2.NewZipFileNode(nil, "/", nil, nil)

	for _, f := range reader.File {
		buf := new(bytes.Buffer)

		rc, err := f.Open()

		if err != nil {
			panic(err)
		}

		io.Copy(buf, rc)

		rc.Close()

		root.AddDescendants(f.Name, buf.Bytes(), f.FileInfo())
	}

	root.Walk(0, func(indent int, node *zw2.ZipFileNode) bool {
		log.Println(">> ", indent, ", ", node.Name)
		return true
	})

	fs, err := zw2.NewHttpFileSystem(root)

	if err != nil {
		panic(err)
	}

	http.Handle("/api/", http.StripPrefix("/api/", http.FileServer(fs)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

```
