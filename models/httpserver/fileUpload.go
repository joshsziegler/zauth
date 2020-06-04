package httpserver

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/ansel1/merry"
)

const (
	UploadFolder = "uploads"
	KiB          = 1024
	MiB          = KiB * 1024
	GiB          = MiB * 1024
	TiB          = GiB * 1024
)

var (
	ErrUploadBadRequest = merry.WithMessage(ErrBadRequest, "file upload failed due to bad request")
	ErrUploadBadWrite   = merry.WithMessage(ErrInternal, "file upload failed: could not save file to disk")
)

// FileUploadPost handles uploading one or more files from a user, to a Group.
func FileUploadPost(c *Context, w http.ResponseWriter, r *http.Request) error {
	// TODO: Check permissions and arguments

	// Check for uploads directory, and create it if necessary
	// TODO: Is there a potential race-condition here, with the dir being
	//       deleted before the file(s) are created? - JZ
	if _, err := os.Stat(UploadFolder); os.IsNotExist(err) {
		os.Mkdir(UploadFolder, 0700)
	}

	// Parse multipart file upload
	var Buf bytes.Buffer
	r.ParseMultipartForm(GiB * 10) // TODO: Allow this be configurable
	files := r.MultipartForm.File["file"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return merry.Here(ErrUploadBadRequest)
		}
		defer file.Close()
		// Copy the file data to the buffer
		io.Copy(&Buf, file)
		fileBytes := Buf.Bytes()
		// Convert the byte array Hash to a String
		hashString := fmt.Sprintf("%x", sha256.Sum256(fileBytes))
		// Write the upload to disk
		err = ioutil.WriteFile(path.Join(UploadFolder, hashString), fileBytes, 0640)
		if err != nil {
			return merry.Here(ErrUploadBadWrite)
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
