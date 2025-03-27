package compressdir

import (
	"archive/zip"
	"bufio"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func check(e error, message string, panickedAF bool) {
	if e != nil {
		if panickedAF {
			log.Fatalf("%s: %s\n", message, e)
			panic(e)
		} else {
			log.Printf("%s: %s\n", message, e)
		}
	}
}

func compressFile(path string, d fs.DirEntry, compressWriter *zip.Writer, dirToCompress string) error {
	info, err := d.Info()

	check(err, "couldn't get file info", true)

	header, err := zip.FileInfoHeader(info)
	check(err, "couldn't create compression header", true)

	header.Method = zip.Deflate
	header.Name = filepath.Join(dirToCompress, d.Name())

	log.Printf("\t\theader path: %s\n", header.Name)

	headerWriter, err := compressWriter.CreateHeader(header)
	check(err, "couldn't create header", false)

	file, err := os.Open(path)
	check(err, "couldn't open file", false)
	defer file.Close()

	_, err = io.Copy(headerWriter, file)
	return err
}

func walkDir(
	source string,
	ignoreSet mapset.Set[string],
	compressWriter *zip.Writer,
	dirToCompress string) error {

	log.Printf("walking directory: %s\n", source)

	dirContents, err := os.ReadDir(source)
	check(err, "couldn't read directory", true)

	ignoreFileExists := slices.ContainsFunc(dirContents,
		func(d fs.DirEntry) bool { return d.Name() == ".gitignore" })
	ignorePath := filepath.Join(source, ".gitignore")

	log.Printf("\tcreating ignore set from: %s\n", ignorePath)
	log.Printf("\t\tmerging: %t\n", ignoreSet != nil)
	if ignoreFileExists {
		ignoreSet = getIgnoreSetFromFile(ignorePath, ignoreSet)
	}

	for _, entry := range dirContents {
		path := filepath.Join(source, entry.Name())
		log.Printf("\tentry: %s\n", path)

		if ignoreMe(ignoreSet, entry.Name()) {
			log.Printf("\t\tignoring entry!\n")
			continue
		}

		if entry.IsDir() {
			err = walkDir(path, ignoreSet, compressWriter, filepath.Join(dirToCompress, entry.Name()))
			check(err, "couldn't walk directory", false)
		} else {
			err = compressFile(path, entry, compressWriter, dirToCompress)
			check(err, "couldn't compress file", false)
			log.Printf("\t\tcompressed file!\n")
		}
	}

	return nil
}

func CompressDir(source string, target string, dirToCompress string) {
	zipFile, err := os.Create(target)
	check(err, "couldn't create zip file", true)
	defer zipFile.Close()

	compressWriter := zip.NewWriter(zipFile)
	defer compressWriter.Close()

	walkDir(source, nil, compressWriter, dirToCompress)

	check(err, "couldn't compress file", false)
}
