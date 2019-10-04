package httpserver

import (
	"net/http"

	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"

	mUser "github.com/joshsziegler/zauth/models/user"
)

type userDetailData struct {
	Message string
	Error   string
	// RequestingUser is the one who asked for this page.
	RequestingUser mUser.User
	// RequestedUser is the User they want to view on this page.
	RequestedUser mUser.User
}

// UserDetailGet is a sub-handler that shows the details for a specific user.
func UserDetailGet(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
	// Get the requested username from the URL
	vars := mux.Vars(r)
	requestedUsername := strings.Trim(vars["username"], " ")
	// Check permissions
	if !c.User.CanEditUser(requestedUsername) {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	// Note: RequestedUser is not necessarily the same as RequestingUser!
	requestedUser, err := c.GetUser(requestedUsername)
	if err != nil {
		return merry.Wrap(err)
	}
	data := userDetailData{
		RequestingUser: *c.User,
		RequestedUser:  requestedUser,
		Message:        c.NormalFlashMessage,
		Error:          c.ErrorFlashMessage,
	}

	// User is viewing this user (or viewing the edit results)
	Render(w, "user_detail.html", data)
	return nil
}
