# src-zip

## Description

A CLI compression tool for source code folders. It uses `.gitignore` when building the compressed file to exclude the unnecessary stuff. Inspired from doing take home projects for interviews that aren't in a repo.

It is currently in very early stages. Building out core functionality before turning into a CLI.

Current state:
- compresses a folder with nested `.gitignore` files, but does not recursively consider higher level ignores.

In Progress:
- recursively add to the ignore set when digging deeper into the tree.

Todo:
- support tar.gz compression
- add proper CLI
- support ignore files by command line, `$GIT_DIR/info/exclude`, and env var `core.excludesFile`

## Run

Built with Go 1.22.5.

To install dependencies, from the root of the repo run
```sh
go get . 
```

To run, specificy the path to the directory, absolute or relative.
```sh
go run . {path_to_dir}
```

The zip will be created in the parent of specified directory alongside it.