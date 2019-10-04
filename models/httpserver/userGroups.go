package httpserver

import (
	"net/http"

	"github.com/ansel1/merry"
	mUser "github.com/joshsziegler/zauth/models/user"
)

type userGroupsData struct {
	Message string
	Error   string
	// RequestingUser is the one who asked for this page.
	RequestingUser mUser.User
	// RequestedUser is the User they want to view on this page.
	RequestedUser mUser.User
}

// userGroupsGet is a sub-handler that shows the groups the user belongs to
func userGroupsGet(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
	// Get the requested username from the URL
	requestedUsername := c.GetRouteVarTrim("username")
	// Check permissions
	if !c.User.CanEditUser(requestedUsername) {
		return ErrPermissionDenied.Here()
	}
	// Handle the request (RequestedUser is not necessarily RequestingUser!)
	requestedUser, err := c.GetUser(requestedUsername)
	if err != nil {
		return merry.Wrap(err)
	}
	data := userGroupsData{
		RequestingUser: *c.User,
		RequestedUser:  requestedUser,
		Message:        c.NormalFlashMessage,
		Error:          c.ErrorFlashMessage,
	}

	// User is viewing this user (or viewing the edit results)
	Render(w, "user_groups.html", data)
	return nil
}
