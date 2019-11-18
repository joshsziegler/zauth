package httpserver

import (
	"html/template"
	"net/http"

	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/csrf"
	"github.com/joshsziegler/zauth/models/user"
)

type LoginPageData struct {
	Message   string
	Error     string
	Username  string
	CSRFField template.HTML
}

// LoginGetPost handles a user's request to view the login page (GET and POST).
func LoginGetPost(c *zauthContext, w http.ResponseWriter, r *http.Request) error {
	if c.User != nil {
		// User is already logged in, so redirect them to their details page
		url, err := c.Router.Get("userDetail").URL("username", c.User.Username)
		if err != nil {
			return err
		}
		http.Redirect(w, r, url.String(), http.StatusFound) // StatusFound ~ 302
		return nil
	}

	// Create page data here so we don't forget to create the CSRF token
	data := LoginPageData{CSRFField: csrf.TemplateField(r),
		Message: c.NormalFlashMessage, Error: c.ErrorFlashMessage}

	switch r.Method {
	case "GET":
		Render(w, "login.html", data)
	case "POST":
		// Make sure they provided a username and password
		username := strings.Trim(r.FormValue("username"), " ")
		password := r.FormValue("password") // don't trim since spaces are allowed
		if username == "" || strings.Trim(password, " ") == "" {
			data.Error = "Please provide a username and password."
			Render(w, "login.html", data)
			return nil
		}

		// Authenticate using the provided username and password
		err := user.Login(DB, username, password)
		if err != nil { // error, or invalid username and/or password
			log.Info(err)
			data.Username = username
			if merry.Is(err, user.ErrorLoginDisabled) {
				data.Error = "This account has been disabled."
			} else {
				data.Error = "Invalid username and/or password."
			}
			Render(w, "login.html", data)
			return nil
		}

		// Always returns a session, even if it's empty
		session, err := store.Get(r, sessionName)
		if err != nil {
			log.Debug("secure session exists, but could not be decoded")
		}
		session.Values["Username"] = username
		// Save the updated session BEFORE writing the response so it's sent
		err = session.Save(r, w)
		if err != nil {
			return ErrInternal.Here()
		}

		log.Infof("logged in as %s", username)
		http.Redirect(w, r, "/users/"+username, 302)
	}
	return nil
}
