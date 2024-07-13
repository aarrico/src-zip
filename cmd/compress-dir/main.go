package main

import (
	"path/filepath"
	"os"
	"github.com/aarrico/src-zip/internal/compressdir"
)
func main() {

	source := filepath.Clean(os.Args[1])
	compressType := "zip" // os.Args[2]
	dirToCompress := filepath.Base(source)
	target := source + "." + compressType
	compressdir.CompressDir(source, target, dirToCompress)
}
