package io

import (
	"fmt"
	"os"
)

func AssertOverwrite(filename string, overwriteAllowed bool) {
	if !overwriteAllowed && fileExists(filename) {
		panic(fmt.Sprintf("File %s already exists, use -overwrite to allow overwriting.", filename))
	}
}

// Helper function to check if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
