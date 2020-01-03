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

type formNewGroup struct {
	Name        string
	Description string
}

type newGroupPageData struct {
	User         *user.User
	ErrorMessage string
	Form         formNewGroup
	CSRFField    template.HTML
}

func newFormNewGroup(r *http.Request) formNewGroup {
	f := formNewGroup{}
	f.Name = strings.Trim(r.FormValue("Name"), " ")
	f.Description = strings.Trim(r.FormValue("Description"), " ")
	return f
}

// NewGroupGet is a sub-handler that shows the Group creation page.
func NewGroupGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	data := newGroupPageData{User: c.User, CSRFField: csrf.TemplateField(r)}
	Render(w, "group_new.html", data)
	return nil
}

// NewGroupPost is a sub-handler that processes the Group creation form.
func NewGroupPost(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	data := newGroupPageData{User: c.User, CSRFField: csrf.TemplateField(r)}
	form := newFormNewGroup(r)
	err := user.AddGroup(c.Tx, form.Name, form.Description)
	if err != nil {
		data.Form = form // Show current form values along with error
		//data.ErrorMessage = merry.UserMessage(err)
		data.ErrorMessage = merry.Details(err)
		Render(w, "group_new.html", data)
		return nil
	}

	// New group created successfully, redirect them to the list page?
	msg := fmt.Sprintf("Group %s successfully created.", form.Name)
	c.AddNormalFlash(msg)
	http.Redirect(w, r, "/groups", 302)
	return nil
}
