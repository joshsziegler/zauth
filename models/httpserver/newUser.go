package httpserver

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/csrf"
	"github.com/joshsziegler/zauth/pkg/user"
)

type formNewUser struct {
	FirstName string
	LastName  string
	Email     string
}

type newUserPageData struct {
	User         *user.User
	ErrorMessage string
	Form         formNewUser
	CSRFField    template.HTML
}

func newFormNewUser(r *http.Request) formNewUser {
	f := formNewUser{}
	f.FirstName = strings.Trim(r.FormValue("FirstName"), " ")
	f.LastName = strings.Trim(r.FormValue("LastName"), " ")
	f.Email = strings.Trim(r.FormValue("Email"), " ")
	return f
}

// NewUserGet is a sub-handler that shows the User creation page.
func NewUserGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	data := newUserPageData{User: c.User, CSRFField: csrf.TemplateField(r)}
	Render(w, "user_new.html", data)
	return nil
}

// NewUserPost is a sub-handler that processes the User creation form.
func NewUserPost(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	data := newUserPageData{User: c.User, CSRFField: csrf.TemplateField(r)}
	form := newFormNewUser(r)
	newUser, err := user.NewUser(c.Tx, form.FirstName, form.LastName, form.Email)
	if err != nil {
		data.Form = form // Show current form values along with error
		data.ErrorMessage = merry.UserMessage(err)
		Render(w, "user_new.html", data)
		// Don't return the error, since we rendered a custom error page
		return nil
	}
	// Commit here, so that errors further along should not undo this new user
	// operation (e.g. from the Email process)
	// TODO: We don't use the next Tx, but we need it for the wrapper. So what
	//       should we do? -JZ
	c.Tx.Commit()
	c.Tx, err = DB.Beginx()
	if err != nil {
		return err
	}

	// Send new user an email asking them to login and set their password
	newUser.SendPasswordResetEmail()
	if err != nil {
		// User was created successfully, but the password email failed
		msg := fmt.Sprintf("User %s successfully created, but their password reset email failed to send. '%s'",
			newUser.Username, err)
		c.AddNormalFlash(msg)
		http.Redirect(w, r, "/users/"+newUser.Username, 302)
		return nil
	}
	// New User created successfully, redirect them to its page
	msg := fmt.Sprintf("User %s successfully created. They were sent an email to set their password.",
		newUser.Username)
	c.AddNormalFlash(msg)
	http.Redirect(w, r, "/users/"+newUser.Username, 302)
	return nil
}
