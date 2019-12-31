package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ansel1/merry"
	"github.com/gorilla/csrf"

	"github.com/joshsziegler/zauth/models/user"
	"github.com/joshsziegler/zauth/pkg/log"
	"github.com/joshsziegler/zauth/pkg/password"
)

// TODO: Should I combine this form&logic with the userDetail change password?
//       That would simplify some things and I could unify the language a bit to
//       "Set Password"
//

// PasswordResetGetPost is a sub-handler that processes password resets via
// secure tokens. These tokens expire after a certain time, and are only valid
// for the username they are created for. This handler does both GET AND POST.
func PasswordResetGetPost(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Only allow anonymous users to use password reset links
	if c.User != nil {
		c.AddNormalFlash("You cannot use a password reset link because you are already logged in.")
		http.Redirect(w, r, fmt.Sprintf("/users/%s", c.User.Username), 302)
		return nil
	}
	// Get the token from the URL
	token := c.GetRouteVarTrim("token")
	// Validate the password reset token
	requestedUsername, err := user.ValidatePasswordResetToken(c.Tx, token)
	if err != nil {
		// Invalid token - Set flash message and redirect to our login page
		log.Errorf("invalid password reset token: %s", err)
		c.AddNormalFlash("Invalid or expired password reset token.")
		http.Redirect(w, r, urlLogin, 302)
		return nil
	}

	// Token is valid: continue
	// Create page data here so we don't forget to create the CSRF token
	requestedUser, err := c.GetUser(requestedUsername)
	if err != nil {
		return err
	}
	data := userSetPasswordPageData{
		//RequestingUser:    nil,
		Message:           c.NormalFlashMessage,
		Error:             c.ErrorFlashMessage,
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

		c.AddNormalFlash("Password successfully changed. Please login.")
		log.Infof("changed password for %s", requestedUsername)
		http.Redirect(w, r, urlLogin, 302)
		return nil
	}
	return nil
}
