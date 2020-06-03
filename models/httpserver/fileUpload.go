package httpserver

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
)

// FileUploadPost handles uploading one or more files from a user, to a Group.
func FileUploadPost(c *Context, w http.ResponseWriter, r *http.Request) error {
	// TODO: Check permissions and arguments

	var Buf bytes.Buffer

	// Parse multipart file upload,  32 << 20 specifies a maximum size of 32 MB
	r.ParseMultipartForm(GiB * 2)
	files := r.MultipartForm.File["file"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // ~ Code 400
			io.WriteString(w, "Error: Failed to get get file from upload request\n")
			log.Printf("Failed to get file upload from client request:\n %s", err)
			return err
		}
		defer file.Close()

		// Copy the file data to the buffer
		io.Copy(&Buf, file)
		fileBytes := Buf.Bytes()
		// Convert the byte array Hash to a String
		hashString := fmt.Sprintf("%x", sha256.Sum256(fileBytes))
		log.Printf("File: %v    SHA 256: %s\n", fileHeader.Filename, hashString)

		// Write the upload to disk
		err = ioutil.WriteFile(hashString, fileBytes, 0640)
		if err != nil {
			log.Printf("Failed to write file to disk:\n%s", err)
			w.WriteHeader(http.StatusInternalServerError) // ~ Code 500
			io.WriteString(w, "Error: Failed to save uploaded file to disk\n")
			return err
		}

		// Reset buffer so we can use it again (reduces memory allocations)
		Buf.Reset()
	}

	// Success, so send them to the home page
	// TODO: Show a Flash Message indicating the success and take them to the
	// dir of the files using http://www.gorillatoolkit.org/pkg/sessions
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
