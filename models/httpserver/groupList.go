package httpserver

import (
	"net/http"

	"github.com/joshsziegler/zauth/pkg/user"
)

type groupListData struct {
	Message string
	Error   string
	User    user.User
	Groups  []*user.Group
}

// GroupListGet shows the user a list of all current zauth groups.
//
// All users are allows to see the list of groups for discovery purposes.
func GroupListGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	groups, err := user.GetGroupsSliceWithoutUsers(c.Tx)
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
