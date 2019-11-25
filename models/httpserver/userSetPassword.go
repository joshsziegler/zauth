package httpserver

import (
	"html/template"
	"net/http"

	"github.com/ansel1/merry"
	"github.com/gorilla/csrf"
	"github.com/joshsziegler/zauth/models/user"
	mUser "github.com/joshsziegler/zauth/models/user"
	"github.com/joshsziegler/zauth/pkg/password"
)

type userSetPasswordPageData struct {
	Message string
	Error   string
	// RequestingUser is the one who asked for this page.
	RequestingUser mUser.User
	// RequestedUser is the User they want to view on this page.
	RequestedUser     mUser.User
	CSRFField         template.HTML
	PasswordMinLength int
	PasswordMaxLength int
}

// userSetPassword is a sub-handler that handles password changes
func userSetPassword(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Get the requested username from the URL
	requestedUsername := c.GetRouteVarTrim("username")
	// Check permissions
	if !c.User.CanEditUser(requestedUsername) {
		return ErrPermissionDenied.Here()
	}
	// Handle the request (RequestedUser is not necessarily RequestingUser!)
	requestedUser, err := c.GetUser(requestedUsername)
	if err != nil {
		return merry.Wrap(err)
	}
	data := userSetPasswordPageData{
		Message:           c.NormalFlashMessage,
		Error:             c.ErrorFlashMessage,
		RequestingUser:    *c.User,
		RequestedUser:     requestedUser,
		CSRFField:         csrf.TemplateField(r),
		PasswordMinLength: password.MinLength,
		PasswordMaxLength: password.MaxLength,
	}

	switch r.Method {
	case "GET":
		Render(w, "user_set_password.html", data)
		return nil
	case "POST":
		newPassword := r.FormValue("NewPassword")
		// Try to set the user's password (will check password rules)
		err = user.SetUserPassword(c.Tx, requestedUsername, newPassword)
		if err != nil {
			data.Error = merry.UserMessage(err)
			Render(w, "user_set_password.html", data)
			return nil
		}
		addNormalFlashMessage(w, r, "Password changed successfully.")
		http.Redirect(w, r, "/users/"+requestedUsername, http.StatusFound)
		return nil
	}
	return nil
}
