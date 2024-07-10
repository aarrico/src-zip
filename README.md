# src-zip

## Description

A CLI compression tool for source code folders. It uses `.gitignore` when building the compressed file to exclude the unnecessary stuff. Inspired from doing take home projects for interviews that aren't in a repo.

It is currently in very early stages. Building out core functionality before turning into a CLI.

Current state:
- compresses a folder with a top level `.gitignore`

In Progress:
- consider nested `.gitignore` files

## Run

Built with Go 1.22.5.

To install dependencies, from the root of the repo run
```sh
go get .
```

To run
```sh
go run .
```