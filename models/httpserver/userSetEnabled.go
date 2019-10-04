package httpserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	mUser "github.com/joshsziegler/zauth/models/user"
)

// userSetEnabled is a sub-handler that handles enabling or disabling a User
func userSetEnabled(c *zauthContext, w http.ResponseWriter, r *http.Request) (err error) {
	// Get the requested username from the URL
	vars := mux.Vars(r)
	requestedUsername := strings.Trim(vars["username"], " ")
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	operation := strings.Trim(vars["isEnabled"], " ")
	if operation == "enable" {
		err = mUser.UserEnable(requestedUsername)
	} else if operation == "disable" {
		err = mUser.UserDisable(requestedUsername)
	} else {
		// TODO: Add message specific to THIS argument
		return ErrRequestArgument.Here()
	}
	// Set flash message indicating result
	if err != nil {
		addErrorFlashMessage(w, r, fmt.Sprintf("Failed to %s user.", operation))
	}
	err = addNormalFlashMessage(w, r, fmt.Sprintf("User successfully %sd.", operation))
	if err != nil {
		return err
	}
	// Redirect them to the requested user's details page
	http.Redirect(w, r, "/users/"+requestedUsername, http.StatusFound)
	return nil
}
