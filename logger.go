package main

import (
	"os"
	"path"
	"fmt"
)

const (
	logFilename = "$HOME/.gomuche/gomuche.log"
)

func getLogFile(isVerbose bool) *os.File {
	filename := os.ExpandEnv(logFilename)
	err := os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		if isVerbose {
			fmt.Printf("Error creating log directory: %v\n", err)
		}
		os.Exit(1)
	}

	file, err := os.OpenFile(filename, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0755)
	if err != nil {
		if isVerbose {
			fmt.Printf("Error opening log file: %v\n", err)
		}
		os.Exit(1)
	}

	return file
}
