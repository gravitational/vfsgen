// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/go-bindata"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func main() {
	var fs http.FileSystem = httpfs.New(mapfs.New(map[string]string{
		"sample-file.txt":                "This file compresses well. Blaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaah.",
		"not-worth-compressing-file.txt": "Its normal contents are here.",
		"folderA/file1.txt":              "Stuff.",
		"folderA/file2.txt":              "Stuff.",
		"folderB/folderC/file3.txt":      "Stuff.",
	}))

	cfg := bindata.NewConfig()
	cfg.Input = fs
	cfg.Package = "bindata_test"
	cfg.Output = "./main_bindata_test.go"

	err := bindata.Translate(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}