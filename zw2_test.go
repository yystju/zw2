package zw2_test

import (
	"archive/zip"
	"bytes"
	"io"
	// "io/ioutil"
	"log"
	"net/http"
	// "os"
	"testing"
	"zw2"
)

func TestZipFileNode_create(t *testing.T) {
	log.Println("[TestZipFileNode_create]")

	var filePathName string = "./testdata/test.zip"

	reader, err := zip.OpenReader(filePathName)

	if err != nil {
		panic(err)
	}

	defer reader.Close()

	root := zw2.NewZipFileNode(nil, "/", nil, nil)

	for _, f := range reader.File {
		log.Println(">> ", f.Name)

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
		log.Println("", indent, ", ", node.Name)
		return true
	})
}

// func TestZipBytes_createwithzipbytes(t *testing.T) {
// 	log.Println("[TestZipFS_createwithzipbytes]")

// 	var filePathName string = "./testdata/test.zip"

// 	zipFile, err := os.OpenFile(filePathName, os.O_RDONLY, 0666)

// 	if err != nil {
// 		panic(err)
// 	}

// 	zipBytes, err := ioutil.ReadAll(zipFile)

// 	if err != nil {
// 		panic(err)
// 	}

// 	zipFile.Close()

// 	zb := zw2.NewZipBytes(zipBytes)

// 	reader, err := zip.NewReader(zb, int64(zb.Length()))

// 	if err != nil {
// 		panic(err)
// 	}

// 	root := zw2.NewZipFileNode(nil, "/", nil, nil)

// 	for _, f := range reader.File {
// 		log.Println(">> ", f.Name)

// 		buf := new(bytes.Buffer)

// 		rc, err := f.Open()

// 		if err != nil {
// 			panic(err)
// 		}

// 		io.Copy(buf, rc)

// 		rc.Close()

// 		root.AddDescendants(f.Name, buf.Bytes(), f.FileInfo())
// 	}

// 	root.Walk(0, func(indent int, node *zw2.ZipFileNode) bool {
// 		log.Println("", indent, ", ", node.Name)
// 		return true
// 	})
// }

/*
	Manual test.
	http://localhost:8080/api/test
*/
func TestHttpFileSystem_manualhttpserver(t *testing.T) {
	log.Println("[TestHttpFileSystem_manualhttpserver]")
	log.Println("\tTry \"http://localhost:8080/api/test\" in browser...")

	var filePathName string = "./testdata/test.zip"

	reader, err := zip.OpenReader(filePathName)

	if err != nil {
		panic(err)
	}

	defer reader.Close()

	root := zw2.NewZipFileNode(nil, "/", nil, nil)

	for _, f := range reader.File {
		// log.Println("f.Name : ", f.Name)
		buf := new(bytes.Buffer)

		rc, err := f.Open()

		if err != nil {
			panic(err)
		}

		io.Copy(buf, rc)

		rc.Close()

		root.AddDescendants(f.Name, buf.Bytes(), f.FileInfo())
	}

	fs, _ := zw2.NewHttpFileSystem(root)

	http.Handle("/api/", http.StripPrefix("/api/", http.FileServer(fs)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
