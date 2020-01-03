package httpserver

import (
	"net/http"

	"github.com/joshsziegler/zauth/pkg/user"
)

type userListData struct {
	User  user.User
	Users map[int64]*(user.User)
}

// UserListGet shows the user a list of all current zauth users.
func UserListGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Only admins can view this page
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}

	users, err := user.GetUsersMapWithoutGroups(c.Tx)
	if err != nil {
		return ErrInternal.Here()
	}
	data := userListData{User: *c.User, Users: users}
	Render(w, "user_list.html", data)
	return nil
}
