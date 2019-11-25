package httpserver

import (
	"net/http"

	mGroup "github.com/joshsziegler/zauth/models/group"
	mUser "github.com/joshsziegler/zauth/models/user"
)

type groupListData struct {
	Message string
	Error   string
	User    mUser.User
	Groups  []*mGroup.Group
}

// GroupListGet shows the user a list of all current zauth groups.
func GroupListGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Only admins can view this page
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}

	groups, err := mGroup.GetGroupsSliceWithoutUsers(c.Tx)
	if err != nil {
		return err
		// return ErrInternal.Here()
	}
	data := groupListData{User: *c.User, Groups: groups,
		Message: c.NormalFlashMessage, Error: c.ErrorFlashMessage,
	}
	Render(w, "group_list.html", data)
	return nil
}
