package main

import "os"

// Exists reports whether the path is a directory and exists as a boolean
func Exists(path string) bool {
	file, err := os.Stat(path)
	if err == nil {
		if file.Mode().IsDir() {
			return true
		}
	}
	return false
}
