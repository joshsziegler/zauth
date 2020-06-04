package httpserver

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ansel1/merry"
	"github.com/gorilla/csrf"
	"github.com/joshsziegler/zauth/pkg/filesharing"
	"github.com/joshsziegler/zauth/pkg/user"
	// "github.com/joshsziegler/zauth/pkg/user"
)

// FileListGet shows the user a list of all files in the specific group.
func FileListGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Get the requested group name from the URL
	name := c.GetRouteVarTrim("name")
	fmt.Printf("group name param %v\n", name)
	group, err := user.GetGroupWithUsers(c.Tx, name)
	if err != nil {
		return merry.Wrap(err)
	}
	fmt.Printf("group %+v\n", group)

	// TODO: Check permissions. Must either be in this group, or an admin to view this page.
	// if XXXXXX {
	//		return ErrPermissionDenied.Here()
	// else if !c.User.IsAdmin() {
	//		return ErrPermissionDenied.Here()
	// }

	files, err := filesharing.GetFiles(c.Tx, name)
	if err != nil {
		return ErrInternal.Here()
	}

	data := struct {
		User      user.User
		Group     user.Group
		Files     []filesharing.File
		CSRFField template.HTML
	}{
		*c.User,
		group,
		files,
		csrf.TemplateField(r),
	}
	Render(w, "file_list.html", data)
	return nil
}
