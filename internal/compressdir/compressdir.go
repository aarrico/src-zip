package compressdir

import (
	"archive/zip"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	header.Name = dirToCompress

	log.Printf("\t\theader path: %s\n", header.Name)

	headerWriter, err := compressWriter.CreateHeader(header)
	check(err, "couldn't create header", false)

	file, err := os.Open(path)
	check(err, "couldn't open file", false)
	defer file.Close()

	_, err = io.Copy(headerWriter, file)
	return err
}

func CompressDir(source string, target string) {
	zipFile, err := os.Create(target)
	check(err, "couldn't create zip file", true)
	defer zipFile.Close()

	compressWriter := zip.NewWriter(zipFile)
	defer compressWriter.Close()

	cmd := exec.Command("git", "ls-files")
	cmd.Dir = source
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	files := strings.Split(string(stdout), "\n")

	for _, file := range files {
		if file == "" {
			continue
		}
		log.Printf("info: %s\n", file)
	
		fi, err := os.Stat(filepath.Join(source, file))
		if err != nil {
			log.Fatal(err)
		}

		de := fs.FileInfoToDirEntry(fi)
		compressFile(filepath.Join(source, file), de, compressWriter, file)
	}
}
