package file

import "os"

// Exists reports whether the path is a regular file and exists as a boolean
func Exists(path string) bool {
	file, err := os.Stat(path)
	if err == nil {
		if file.Mode().IsRegular() {
			return true
		}
	}
	return false
}
