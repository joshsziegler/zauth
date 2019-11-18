package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ansel1/merry"
	mUser "github.com/joshsziegler/zauth/models/user"
)

// userSetEnabled is a sub-handler that handles enabling or disabling a User
func userSetEnabled(c *zauthContext, w http.ResponseWriter, r *http.Request) (err error) {
	// Get the requested username from the URL
	requestedUsername := c.GetRouteVarTrim("username")
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	operation := c.GetRouteVarTrim("isEnabled")
	if operation == "enable" {
		err = mUser.UserEnable(requestedUsername)
	} else if operation == "disable" {
		err = mUser.UserDisable(requestedUsername)
	} else {
		return merry.Here(ErrRequestArgument).
			WithMessagef("invalid operation '%s' (must be 'enable' or 'disable')",
				operation)
	}
	// Set flash message indicating result
	if err != nil {
		addErrorFlashMessage(w, r, fmt.Sprintf("Failed to %s user.", operation))
	} else {
		addNormalFlashMessage(w, r, fmt.Sprintf("User successfully %sd.",
			operation))
	}
	// Redirect them to the requested user's details page
	http.Redirect(w, r, "/users/"+requestedUsername, http.StatusFound)
	return nil
}
