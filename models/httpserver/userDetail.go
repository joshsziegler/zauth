package httpserver

import (
	"net/http"

	"github.com/ansel1/merry"

	mUser "github.com/joshsziegler/zauth/models/user"
	"github.com/joshsziegler/zauth/models/user2group"
)

type userDetailData struct {
	Message string
	Error   string
	// RequestingUser is the one who asked for this page.
	RequestingUser mUser.User
	// RequestedUser is the User they want to view on this page.
	RequestedUser mUser.User
	// GroupMembership holds all Groups, and whether RequestedUser is a member.
	GroupMembership []user2group.Group
}

// UserDetailGet is a sub-handler that shows the details for a specific user.
func UserDetailGet(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
	// Get the requested username from the URL
	requestedUsername := c.GetRouteVarTrim("username")
	// Check permissions
	if !c.User.CanEditUser(requestedUsername) {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	// Note: RequestedUser is not necessarily the same as RequestingUser!
	// TODO: This returns the User WITH their Groups, but we get them all below.
	//       Should we get the user without their groups for better performance?
	requestedUser, err := c.GetUser(requestedUsername)
	if err != nil {
		return merry.Wrap(err)
	}
	// Get RequestedUser's group membership matrix
	groupMembership, err := user2group.GetUsersMembership(c.Tx, requestedUser.ID)
	if err != nil {
		return merry.Wrap(err)
	}
	data := userDetailData{
		RequestingUser:  *c.User,
		RequestedUser:   requestedUser,
		Message:         c.NormalFlashMessage,
		Error:           c.ErrorFlashMessage,
		GroupMembership: groupMembership,
	}

	// User is viewing this user (or viewing the edit results)
	Render(w, "user_detail.html", data)
	return nil
}
