package httpserver

import (
	"net/http"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"

	"github.com/joshsziegler/zauth/models/user"
)

// LogRequest except those that start with '/static/'
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.RequestURI, "/static/") {
			log.Info(r.RequestURI)
		}
		next.ServeHTTP(w, r)
	})
}

// RequireLogin redirects the user to the login page IFF not already logged in
func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := getUsername(w, r)
		if username == nil { // Not logged in
			http.Redirect(w, r, urlLogin, http.StatusFound)
			return
		}
		// User IS logged in, allow them to continue
		next.ServeHTTP(w, r)
	})
}

// zauthContext is a struct passed to zauthHandlers with additional information
// not in the standard HTTP handlers, such as User object.
type zauthContext struct {
	// Router is the Gorilla-based router. It's included here so we can route
	// names can be reversed (e.g. what is the URI for a user's details page?).
	Router *mux.Router
	// User is the person making this HTTP request.
	User               *user.User
	NormalFlashMessage string
	ErrorFlashMessage  string
}

// GetUser returns the specified User. This potentially avoids a second DB call
// if the HTTP request is being made by this user, and is thus already loaded
// into memory.
func (c *zauthContext) GetUser(username string) (user.User, error) {
	if c.User.Username == username {
		return *c.User, nil
	}
	return user.GetUserWithGroups(DB, username)
}

// zauthHandler adds a context argument and is used with handleWrap
type zauthHandler = func(c *zauthContext, w http.ResponseWriter, r *http.Request) error

// wrapHandler provides a wrapper to page-specific handlers. It handles logging,
// getting and checking user authentication, flash messages and error rendering.
// The sub handler MUST implement zauthHandler.
//
// Order of operations:
//  - Get the user object if logged in
//  - Log request (before auth req redirection)
//  - Redirect to login IFF auth is required
//  - Get flash messages if any
//  -- run page-specific sub-handler
//  - If there is an error, STOP continuing and render a proper error
//  - If there are new flash messages, SAVE them
//  - Render the page OR render the error
func wrapHandler(router *mux.Router, subHandler zauthHandler, requireLogin bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create the context that we pass to the subHandler
		c := zauthContext{Router: router}
		var err error
		// Get username of user (or nil if they are not logged in)
		username := getUsername(w, r)
		// Log this request, including their username if they are logged in
		if username == nil {
			log.Infof("anonymous %s %s", r.Method, r.RequestURI)
		} else {
			log.Infof("%s %s %s", *username, r.Method, r.RequestURI)
		}
		// Redirect if this page requires authentication
		if requireLogin && username == nil { // Not logged in
			err = addNormalFlashMessage(w, r, "Sorry, but that page requires you to login first.")
			if err != nil {
				log.Error(err)
			}
			http.Redirect(w, r, urlLogin, http.StatusFound)
			return
		}
		// Get and save user struct if they are logged in
		if username != nil {
			tempUser, err := user.GetUserWithGroups(DB, *username)
			if err != nil {
				Error(w, 500, "Error", "Sorry, but the server encountered an error.", nil)
				return
			}
			// Convert to pointer to allow us to check for an empty User using nil
			c.User = &tempUser
		}
		// Get flash messages, if any
		c.NormalFlashMessage, c.ErrorFlashMessage = getFlashMessages(w, r)
		// Run the page-specific handler, which renders the page unless there is an error
		err = subHandler(&c, w, r)
		if err != nil {
			if merry.Is(err, ErrPermissionDenied) {
				ErrorUnauthorized(w, merry.UserMessage(err), c.User)
			} else if merry.Is(err, ErrInternal) {
				Error(w, 500, "Error", merry.UserMessage(err), c.User)
			} else {
				// We can't guarantee this error has a nice UserMessage
				Error(w, 500, "Error", merry.Details(err), nil)
			}
			return
		}
	})
}
