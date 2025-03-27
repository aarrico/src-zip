package ignore

import (
	gi "github.com/sabhiram/go-gitignore"
	"log"
	"path/filepath"
)

type Ignorer interface {
	IgnoreMe(path string) bool
}

type ignoreMatcher struct {
	ignores *gi.GitIgnore
	baseDir string
}

func createIgnoreSet(path string) (Ignorer, error) {
	ignoreSet, err := gi.CompileIgnoreFile(path)
	if err != nil {
		log.Printf("Error compiling ignore file [%s]: %v", path, err)
		return nil, err
	}

	return &ignoreMatcher{
		ignores: ignoreSet,
		baseDir: filepath.Dir(path),
	}, nil
}

func (matcher *ignoreMatcher) IgnoreMe(path string) bool {

	if matcher.ignores == nil {
		return false
	}

	return matcher.ignores.MatchesPath(path)
}
