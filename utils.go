package main

import (
	"os"
	
	"mvdan.cc/sh/v3/shell"
)

func realPath(path string) string {
	out, err := shell.Fields(path, nil)
	panicError(err)
	return out[0]
}

// Panics if the error is not nil.
func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	panicError(err)
	return true
}
