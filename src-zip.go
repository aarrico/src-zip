package main

import (
	"bufio"
	mapset "github.com/deckarep/golang-set/v2"
	"log"
	"os"
	"strings"
)

func check(e error, message string, panickedAF bool) {
	if e != nil {
		log.Fatalf("%s: %s\n", message, e)
		if panickedAF {
			panic(e)
		}
	}
}

func createIgnoreSetFromFile(filepath string) mapset.Set[string] {
	ignoreFile, err := os.Open(filepath)
	check(err, "couldn't open file", true)

	scanner := bufio.NewScanner(ignoreFile)
	ignoreSet := mapset.NewSet[string]()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line != "" && !strings.HasPrefix(line, "#") {
			ignoreSet.Add(line)
		}
	}

	ignoreFile.Close()

	return ignoreSet
}

func main() {

	ignoreSet := createIgnoreSetFromFile(".gitignore")
	log.Printf("ignoreSet: %s\n", ignoreSet)
}
