package httpserver

import (
	"html/template"
	"net/http"

	"github.com/ansel1/merry"
	"github.com/gorilla/csrf"
	"github.com/joshsziegler/zauth/pkg/filesharing"
	"github.com/joshsziegler/zauth/pkg/user"
)

// FileListGet shows the user a list of all files in the specific group.
func FileListGet(c *Context, w http.ResponseWriter, r *http.Request) error {
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
		return merry.Wrap(err)
	}

	files, err := filesharing.GetFiles(c.Tx, name)
	if err != nil {
		return ErrInternal.Here() // TODO: Return a more descriptive error
	}

	data := struct {
		User      user.User
		Group     user.Group
		Files     []filesharing.File
		CSRFField template.HTML
		Message   string // Flash Message
		Error     string // Flash Message
	}{
		*c.User,
		group,
		files,
		csrf.TemplateField(r),
		c.NormalFlashMessage,
		c.ErrorFlashMessage,
	}
	Render(w, "file_list.html", data)
	return nil
}
