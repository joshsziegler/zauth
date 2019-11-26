package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ansel1/merry"
	"github.com/joshsziegler/zauth/models/user2group"
)

// userAddRemoveGroups is a sub-handler that handles adding or removing a single
// user from a single group.
func userAddRemoveGroups(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	var err error
	// Get the requested username from the URL
	requestedUsername := c.GetRouteVarTrim("username")
	group := c.GetRouteVarTrim("groupname")

	// Handle the request
	var flash string
	operation := c.GetRouteVarTrim("addOrRemove")
	if operation == "add" {
		flash = fmt.Sprintf("Adding user %s to group %s ", requestedUsername, group)
		err = user2group.AddUserToGroup(c.Tx, requestedUsername, group)
	} else if operation == "remove" {
		flash = fmt.Sprintf("Removing user %s from group %s ", requestedUsername, group)
		err = user2group.RemoveUserFromGroup(c.Tx, requestedUsername, group)
	} else {
		return merry.Here(ErrRequestArgument).
			WithMessagef("invalid operation '%s' (must be 'add' or 'remove')",
				operation)
	}
	// Set flash message indicating result
	if err != nil {
		c.AddNormalFlash(flash + "failed.")
		return err
	}
	c.AddNormalFlash(flash + "succeeded.")
	// Redirect them to the requested user's details page
	http.Redirect(w, r, "/users/"+requestedUsername, http.StatusFound)
	return nil
}
