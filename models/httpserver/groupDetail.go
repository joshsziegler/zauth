package httpserver

import (
	"net/http"

	"github.com/ansel1/merry"
	"github.com/joshsziegler/zauth/pkg/user"
)

type groupDetailData struct {
	Message        string
	Error          string
	RequestingUser user.User  // User who asked for this page.
	Group          user.Group // Group they want to view
	Members        []string   // All of the Usernames in this Group.
}

// GroupDetailGet is a sub-handler that shows the details for a specific user.
func GroupDetailGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Get the requested group name from the URL
	name := c.GetRouteVarTrim("name")
	// TODO: Check the group exists
	group, err := user.GetGroupWithUsers(c.Tx, name)
	if err != nil {
		return merry.Wrap(err)
	}

	// TODO: Check permissions
	// if !c.User.CanEditUser(name) {
	// 	return ErrPermissionDenied.Here()
	// }

	// Handle the request
	// TODO: Get the group details
	// TODO: Get group members

	data := groupDetailData{
		RequestingUser: *c.User,
		Group:          group,
		Message:        c.NormalFlashMessage,
		Error:          c.ErrorFlashMessage,
		// TODO: Members
	}

	Render(w, "group_detail.html", data)
	return nil
}
