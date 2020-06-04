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
	"github.com/joshsziegler/zauth/pkg/user"
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
	// TODO: Get the optional folder they are uploading to
	// Get the requested group name from the URL and verify it's valid
	name := c.GetRouteVarTrim("name")
	// Check permissions BEFORE getting group for speed (must belong to group or be and admin)
	if !c.User.IsInGroup(name) {
		if !c.User.IsAdmin() {
			return merry.Here(ErrPermissionDenied)
		}
	}
	group, err := user.GetGroupWithUsers(c.Tx, name)
	if err != nil {
		return merry.Wrap(err) // TODO: Return more descriptive error (SQL error, or no group with that name?)
	}
	fmt.Printf("Group: %+v\n", group) // TODO: Remove this...

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
	c.AddNormalFlash(fmt.Sprintf("Files successfully uploaded to %s", name))
	http.Redirect(w, r, fmt.Sprintf("/groups/%s/files", name), http.StatusFound)
	return nil
}
