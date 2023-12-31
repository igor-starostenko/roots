package main

import (
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func warn(message string) {
	fmt.Fprintf(os.Stderr, message+"\n")
}

func stop(message string, code int) {
	warn(message)
	os.Exit(code)
}

func fileExists(name string) bool {
	path, err := os.Stat(name)
	if err == nil {
		return !path.IsDir()
	} else if os.IsNotExist(err) {
		return false
	}
	check(err)
	return false
}

func checkFileExists(path string) {
	if !fileExists(path) {
		stop(fmt.Sprintf("Cannot find file \"%s\"\n", path), 1)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
