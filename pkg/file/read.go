package file

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

// toUtf8 reads the byte array and converts each rune to UTF-8 to avoid
// character encoding issues.
//
// For example, a Latin-1 encoded file with the character '–', would normally
// result in an error when converted to a string via `string(bytes)`.
func toUtf8(bytes []byte) string {
	buf := make([]rune, len(bytes))
	for i, b := range bytes {
		buf[i] = rune(b)
	}
	return string(buf)
}

// ReadAsBytes joins the path, and returns the file's contents as a byte array.
func ReadAsBytes(path ...string) ([]byte, error) {
	joinedPath := filepath.Join(path...) // relative path
	bytes, err := ioutil.ReadFile(joinedPath)
	return bytes, err
}

// ReadAsBytes joins the path, and returns the file's contents as a string.
//
// If the file uses a character encoding other than UTF-8, this should handle it
// correctly. For example, a latin-1 file with the character '–', would normally
// result in an error when converted to a string via `string(bytes)`. This
// function will convert it correclty, however.
func ReadAsString(path ...string) (string, error) {
	b, err := ReadAsBytes(path...)
	if err != nil {
		return "", err
	}
	s := toUtf8(b)
	return s
}
