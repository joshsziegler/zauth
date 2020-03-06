package file

import (
	"io/ioutil"
	"path/filepath"
)

// ReadAsBytes joins the path, and returns the file's contents as a byte array.
func ReadAsBytes(path ...string) ([]byte, error) {
	joinedPath := filepath.Join(path...) // relative path
	bytes, err := ioutil.ReadFile(joinedPath)
	return bytes, err
}

// ReadAsBytes joins the path, and returns the file's contents as a string.
func ReadAsString(path ...string) (string, error) {
	b, err := ReadAsBytes(path...)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
