package httpserver

import (
	"net/http"

	mUser "github.com/joshsziegler/zauth/models/user"
)

type userListData struct {
	User  mUser.User
	Users map[int64]*(mUser.User)
}

// UserListGet shows the user a list of all current zauth users.
func UserListGet(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
	// Only admins can view this page
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}

	// Handle the request
	tx, err := DB.Beginx()
	if err != nil {
		return ErrInternal.Here()
	}
	defer tx.Commit()
	users, err := mUser.GetUsersWithoutGroups(tx)
	if err != nil {
		return ErrInternal.Here()
	}

	data := userListData{User: *c.User, Users: users}
	Render(w, "user_list.html", data)

	return nil
}
