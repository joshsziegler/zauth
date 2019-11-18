package httpserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/joshsziegler/zauth/models/user"
)

// zauthContext is a struct passed to zauthHandlers with additional information
// not in the standard HTTP handlers, such as User object.
type zauthContext struct {
	// Router is the Gorilla-based router. It's included here so we can route
	// names can be reversed (e.g. what is the URI for a user's details page?).
	Router *mux.Router
	// Tx is the database transaction that is started for you.
	Tx *sqlx.Tx
	// User is the person making this HTTP request.
	User               *user.User
	NormalFlashMessage string
	ErrorFlashMessage  string
	RouteVariables     map[string]string
}

// get the username from the secure session. If it doesn't exist, redirect to
// the login page.
//
// Does NOT use the database!
func getUsername(w http.ResponseWriter, r *http.Request) *string {
	// Always returns a session, even if it's empty
	session, err := store.Get(r, sessionName)
	if err != nil {
		log.Debug("secure session exists, but could not be decoded")
	}
	val := session.Values["Username"]
	if val == nil { // User not logged in
		return nil
	}
	username := fmt.Sprintf("%v", val)
	return &username
}

// GetUser returns the specified User. This potentially avoids a second DB call
// if the HTTP request is being made by this user, and is thus already loaded
// into memory.
func (c *zauthContext) GetUser(username string) (user.User, error) {
	if c.User != nil && c.User.Username == username {
		return *c.User, nil
	}
	return user.GetUserWithGroups(c.Tx, username)
}

// GetRouteVarTrim returns the whitespace trimmed Gorilla Mux Route Variable.
func (c *zauthContext) GetRouteVarTrim(varName string) string {
	return strings.Trim(c.RouteVariables[varName], " ")
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
//
// TODO: Defer logging the request until we have the result, so we can log them
//       on the same line like so: josh POST /group/new -> error: duplicate name
func wrapHandler(router *mux.Router, subHandler zauthHandler, requireLogin bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create the context that we pass to the subHandler
		c := zauthContext{
			Router:         router,
			RouteVariables: mux.Vars(r),
		}
		// Create a DB transaction for our context struct
		tx, err := DB.Beginx()
		if err != nil {
			Error(w, 500, "Error",
				"Sorry, but the server encountered an error.", nil)
		}
		c.Tx = tx
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
			err = addNormalFlashMessage(w, r, "Sorry, but that page requires "+
				"you to login first. If you were previously logged in, your "+
				"session has expired.")
			if err != nil {
				log.Error(err)
			}
			http.Redirect(w, r, urlLogin, http.StatusFound)
			c.Tx.Commit()
			return
		}
		// Get and save user struct if they are logged in
		if username != nil {
			tempUser, err := user.GetUserWithGroups(tx, *username)
			if err != nil {
				Error(w, 500, "Error",
					"Sorry, but the server encountered an error.", nil)
				c.Tx.Commit()
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
			// Rollback if the sub-handler failed
			c.Tx.Rollback()
			return
		}
		c.Tx.Commit()
	})
}
