package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ansel1/merry"
	mUser "github.com/joshsziegler/zauth/models/user"
	"github.com/joshsziegler/zauth/models/user2group"
)

type userGroupsData struct {
	Message string
	Error   string
	// RequestingUser is the one who asked for this page.
	RequestingUser mUser.User
	// RequestedUser is the User they want to view on this page.
	RequestedUser mUser.User
	// GroupMembership holds all Groups, and whether this User is a member.
	GroupMembership []user2group.Group
}

// userGroupsGet is a sub-handler that shows the groups the user belongs to
func userGroupsGet(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Get the requested username from the URL
	requestedUsername := c.GetRouteVarTrim("username")
	// Handle the request (RequestedUser is not necessarily RequestingUser!)
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

	data := userGroupsData{
		RequestingUser:  *c.User,
		RequestedUser:   requestedUser,
		Message:         c.NormalFlashMessage,
		Error:           c.ErrorFlashMessage,
		GroupMembership: groupMembership,
	}

	// User is viewing this user (or viewing the edit results)
	Render(w, "user_groups.html", data)
	return nil
}

//  userGroupsPost

// userAddRemoveGroups is a sub-handler that handles adding or removing a single
// user from a single group.
func userAddRemoveGroups(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
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
		// TODO: Add message specific to THIS argument
		return ErrRequestArgument.Here()
	}
	// Set flash message indicating result
	if err != nil {
		addErrorFlashMessage(w, r, flash+"failed.")
		return err
	}
	err = addNormalFlashMessage(w, r, flash+"succeeded.")
	if err != nil {
		return err
	}
	// Redirect them to the requested user's details page
	http.Redirect(w, r, "/users/"+requestedUsername+"/groups", http.StatusFound)
	return nil
}
