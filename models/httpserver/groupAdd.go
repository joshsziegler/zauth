package httpserver

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/csrf"
	mGroup "github.com/joshsziegler/zauth/models/group"
	mUser "github.com/joshsziegler/zauth/models/user"
)

type formNewGroup struct {
	Name        string
	Description string
}

type newGroupPageData struct {
	User         *mUser.User
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
func NewGroupGet(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
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
func NewGroupPost(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
	// Check permissions
	if !c.User.IsAdmin() {
		return ErrPermissionDenied.Here()
	}
	// Handle the request
	data := newGroupPageData{User: c.User, CSRFField: csrf.TemplateField(r)}
	form := newFormNewGroup(r)
	err := mGroup.Add(c.Tx, form.Name, form.Description)
	if err != nil {
		data.Form = form // Show current form values along with error
		//data.ErrorMessage = merry.UserMessage(err)
		data.ErrorMessage = merry.Details(err)
		Render(w, "group_new.html", data)
		return nil
	}

	// New group created successfully, redirect them to the list page?
	msg := fmt.Sprintf("Group %s successfully created.", form.Name)
	addNormalFlashMessage(w, r, msg)
	http.Redirect(w, r, "/groups", 302)
	return nil
}
