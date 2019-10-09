package httpserver

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/csrf"
	mail "github.com/joshsziegler/zauth/models/email"
	mUser "github.com/joshsziegler/zauth/models/user"
)

type formNewUser struct {
	FirstName string
	LastName  string
	Email     string
}

type newUserPageData struct {
	User         *mUser.User
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
func NewUserGet(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
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
func NewUserPost(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	data := newUserPageData{User: c.User, CSRFField: csrf.TemplateField(r)}
	form := newFormNewUser(r)
	newUser, err := mUser.NewUser(c.Tx, form.FirstName, form.LastName, form.Email)
	if err != nil {
		data.Form = form // Show current form values along with error
		data.ErrorMessage = merry.UserMessage(err)
		Render(w, "user_new.html", data)
		return nil
	}

	// Send new user an email asking them to login and set their password
	// TODO: Allow this email to be configured? - JZ
	resetLink := newUser.GetPasswordResetToken(2) // TODO: Set expiration via config?
	err = mail.Send("MindModeling", "no-reply@mindmodeling.org",
		newUser.CommonName(), newUser.Email,
		"Your New MindModeling Account",
		"A new account has been created",
		`<p>Hello `+newUser.CommonName()+`,</p>
		<p>A new account has been created for you on MindModeling. To complete
		the setup, you need to 
		<a href="http://localhost:8888/reset-password/`+resetLink+`">
		set your password here</a>. This link is valid for the next `+strconv.Itoa(2)+`
		hours.</p>`)
	if err != nil {
		// User was created successfully, but the password email failed
		msg := fmt.Sprintf("User %s successfully created, but their password reset email failed to send. %s",
			newUser.Username, err)
		// Show an error flash message, but ignore any error this might return
		addErrorFlashMessage(w, r, msg)
		return err
	}
	// New User created successfully, redirect them to its page
	msg := fmt.Sprintf("User %s successfully created. They were sent an email to set their password.",
		newUser.Username)
	// Show an error flash message, but ignore any error this might return
	addNormalFlashMessage(w, r, msg)
	http.Redirect(w, r, "/users/"+newUser.Username, 302)
	return nil
}
