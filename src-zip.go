package main

import (
	"archive/zip"
	"bufio"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

func createIgnoreSetFromFile(filepath string) mapset.Set[string] {
	ignoreFile, err := os.Open(filepath)
	check(err, "couldn't open file", true)
	defer ignoreFile.Close()

	scanner := bufio.NewScanner(ignoreFile)
	ignoreSet := mapset.NewSet[string]()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line != "" && !strings.HasPrefix(line, "#") {
			ignoreSet.Add(line)
		}
	}

	return ignoreSet
}

func compressFolder(source string, target string, ignoreSet mapset.Set[string]) {
	zipFile, err := os.Create(target)
	check(err, "couldn't create zip file", true)
	defer zipFile.Close()

	compressWriter := zip.NewWriter(zipFile)
	defer compressWriter.Close()

	err = filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		check(err, "couldn't walk path", true)

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		check(err, "couldn't get file info", true)

		header, err := zip.FileInfoHeader(info)
		check(err, "couldn't create compression header", true)
		header.Method = zip.Deflate

		relPath, _ := filepath.Rel(filepath.Dir(source), path)

		header.Name = filepath.ToSlash(relPath)

		headerWriter, err := compressWriter.CreateHeader(header)
		check(err, "couldn't create header", false)

		file, err := os.Open(path)
		check(err, "couldn't open file", false)
		defer file.Close()

		_, err = io.Copy(headerWriter, file)
		return err
	})

	check(err, "couldn't compress file", false)
}

func main() {

	ignoreSet := createIgnoreSetFromFile(".gitignore")
	log.Printf("ignoreSet: %s\n", ignoreSet)

	source := "zipmeup"
	target := source + ".zip"

	compressFolder(source, target, ignoreSet)
}
