package main

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

var dirToCompress string
var pathToDir string

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

func createIgnoreSet(filepath string) mapset.Set[string] {
	ignoreFile, err := os.Open(filepath)
	check(err, "couldn't open file", true)
	defer ignoreFile.Close()

	scanner := bufio.NewScanner(ignoreFile)
	ignoreSet := mapset.NewSet[string]()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		ignoreSet.Add(line)
	}

	return ignoreSet
}

func ignoreMe(ignoreSet mapset.Set[string], path string) bool {

	if ignoreSet == nil {
		return false
	}

	for _, ignorePattern := range ignoreSet.ToSlice() {

		matched, err := filepath.Match(ignorePattern, path)
		check(err, "couldn't match path", true)
		if matched {
			return true
		}
	}

	return false
}

func compressDir(source string, target string, dirToCompress string) {
	zipFile, err := os.Create(target)
	check(err, "couldn't create zip file", true)
	defer zipFile.Close()

	compressWriter := zip.NewWriter(zipFile)
	defer compressWriter.Close()

	walkDir(source, nil, compressWriter, dirToCompress)

	check(err, "couldn't compress file", false)
}

func walkDir(source string, ignoreSet mapset.Set[string], compressWriter *zip.Writer, dirToCompress string) error {

	log.Printf("walking directory: %s\n", source)

	dirContents, err := os.ReadDir(source)
	check(err, "couldn't read directory", true)

	ignoreFileExists := slices.ContainsFunc(dirContents, func(d fs.DirEntry) bool { return d.Name() == ".gitignore" })

	if ignoreFileExists {
		log.Printf("\tcreating ignore set\n")
		ignoreSet = createIgnoreSet(filepath.Join(source, ".gitignore"))
	}

	for _, entry := range dirContents {
		path := filepath.Join(source, entry.Name())
		log.Printf("\tentry: %s\n", path)

		if entry.IsDir() {
			if !ignoreMe(ignoreSet, entry.Name()) {
				err = walkDir(path, ignoreSet, compressWriter, filepath.Join(dirToCompress, entry.Name()))
				check(err, "couldn't walk directory", false)
			}
		} else {
			if ignoreMe(ignoreSet, entry.Name()) {
				log.Printf("\t\tignoring file: %s\n", entry.Name())
				continue
			}

			err = compressFile(path, entry, compressWriter, dirToCompress)
			check(err, "couldn't compress file", false)
			log.Printf("\t\tcompressed file: %s\n", entry.Name())

		}
	}

	return nil
}

func compressFile(path string, d fs.DirEntry, compressWriter *zip.Writer, dirToCompress string) error {
	info, err := d.Info()

	check(err, "couldn't get file info", true)

	header, err := zip.FileInfoHeader(info)
	check(err, "couldn't create compression header", true)
	
	header.Method = zip.Deflate
	header.Name = filepath.Join(dirToCompress, d.Name())

	log.Printf("\t\tcompressing file: %s\n", header.Name)

	headerWriter, err := compressWriter.CreateHeader(header)
	check(err, "couldn't create header", false)

	file, err := os.Open(path)
	check(err, "couldn't open file", false)
	defer file.Close()

	_, err = io.Copy(headerWriter, file)
	return err
}

func main() {

	source := filepath.Clean(os.Args[1])
	dirToCompress = filepath.Base(source)
	pathToDir = filepath.Dir(source)
	target := source + ".zip"
	compressDir(source, target, dirToCompress)
}
