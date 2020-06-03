package httpserver

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/joshsziegler/zauth/pkg/filesharing"
	"github.com/joshsziegler/zauth/pkg/user"
	// "github.com/joshsziegler/zauth/pkg/user"
)

// FileListGet shows the user a list of all files in the specific group.
func FileListGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Get the requested group name from the URL
	groupName := c.GetRouteVarTrim("groupname")

	// TODO: Check that this Group EXISTS first

	// TODO: Check permissions. Must either be in this group, or an admin to view this page.
	// if XXXXXX {
	//		return ErrPermissionDenied.Here()
	// else if !c.User.IsAdmin() {
	//		return ErrPermissionDenied.Here()
	// }

	files, err := filesharing.GetFiles(c.Tx, groupName)
	if err != nil {
		return ErrInternal.Here()
	}

	data := struct {
		User      user.User
		Files     []filesharing.File
		CSRFField template.HTML
	}{
		*c.User,
		files,
		csrf.TemplateField(r),
	}
	Render(w, "file_list.html", data)
	return nil
}
