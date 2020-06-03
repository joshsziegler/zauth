package file

import "testing"

func TestReadingLatin1File(t *testing.T) {
	s, err := ReadAsString("testdata", "latin1-example.input")
	if err != nil {
		t.Error(err)
	}
	b, err := ReadAsBytes("testdata", "latin1-example.output")
	if err != nil {
		t.Error(err)
	}
	e := string(b)
	if e != s {
		t.Error("Expected: ", e)
		t.Error("Result: ", s)
	}
}
