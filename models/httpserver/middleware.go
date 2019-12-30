package httpserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/joshsziegler/zauth/models/user"
	"github.com/joshsziegler/zauth/pkg/log"
)

const (
	flashMessageSessionName = `zauth-flashes`
	errorFlashKey           = `error-flash`
	normalFlashKey          = `message-flash`
	errGetFlashMessages     = "error getting flash messages"
	errAddFlashMessage      = "error while adding a normal flash message"
)

// Context is a struct passed to Handlers with additional information not
// in the standard HTTP handlers, such as User object.
type Context struct {
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
	Response           http.ResponseWriter
	Request            *http.Request
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
func (c *Context) GetUser(username string) (user.User, error) {
	if c.User != nil && c.User.Username == username {
		return *c.User, nil
	}
	return user.GetUserWithGroups(c.Tx, username)
}

// GetRouteVarTrim returns the whitespace trimmed Gorilla Mux Route Variable.
func (c *Context) GetRouteVarTrim(varName string) string {
	return strings.Trim(c.RouteVariables[varName], " ")
}

// Handler adds a context argument and is used with handleWrap
type Handler = func(c *Context, w http.ResponseWriter, r *http.Request) error

// Wrap provides a wrapper to page-specific handlers. It handles logging,
// getting and checking user authentication, flash messages and error rendering.
// The sub handler MUST implement Handler.
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
func Wrap(router *mux.Router, subHandler Handler, requireLogin bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create the context that we pass to the subHandler
		c := Context{
			Router:         router,
			RouteVariables: mux.Vars(r),
			Response:       w,
			Request:        r,
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
			c.AddNormalFlash("Sorry, but that page requires you to " +
				"login first. If you were previously logged in, your session " +
				"has expired.")
			http.Redirect(w, r, urlLogin, http.StatusFound)
			err = c.Tx.Commit()
			if err != nil {
				log.Error(err)
			}
			return
		}
		// Get and save user struct if they are logged in
		if username != nil {
			tempUser, err := user.GetUserWithGroups(tx, *username)
			if err != nil {
				Error(w, 500, "Error",
					"Sorry, but the server encountered an error.", nil)
				err = c.Tx.Commit()
				if err != nil {
					log.Error(err)
				}
				return
			}
			// Convert to pointer to allow us to check for an empty User using nil
			c.User = &tempUser
		}
		// Get flash messages, if any
		c.getFlashMessages()
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
			err = c.Tx.Rollback()
			if err != nil {
				log.Error(err)
			}
			return
		}
		err = c.Tx.Commit()
		if err != nil {
			log.Error(err)
		}
	})
}

// getFlashMessages gets a single error flash message and a single normal
// flash message. This is to restrict the HTTP handlers to a single, most-
// important message of each type.
//
// Normal flash messages are not errors, and are typically informing the user
// that an operation was successful such as logging out, or creating a new user.
//
// This method isn't public, because Wrap() should be the only one calling it.
func (c *Context) getFlashMessages() {
	session, err := store.Get(c.Request, flashMessageSessionName)
	if err != nil {
		// This could be an error, but more likely there simply aren't any
		// flash message set. So simply return empty strings for both.
		return
	}
	// Get any error flash message
	fm := session.Flashes(errorFlashKey)
	if fm != nil {
		c.ErrorFlashMessage = fmt.Sprintf("%v", fm[0])
	}
	// Get any normal flash message
	fm = session.Flashes(normalFlashKey)
	if fm != nil {
		c.NormalFlashMessage = fmt.Sprintf("%v", fm[0])
	}
	// Always save after retrieving flash messages so they are removed
	err = session.Save(c.Request, c.Response)
	if err != nil {
		log.Error(err)
	}
}

// AddNormalFlash adds a flash message to the store to be viewed by the user
// upon next page request.
//
// If you need to add a flash message, you should do so before writing to the
// response(Writer)! This is due to gorilla's session.Save().
//
// Normal flash messages are not errors, and are typically informing the user
// that an operation was successful such as logging out, or creating a new user.
func (c *Context) AddNormalFlash(message string) {
	session, err := store.Get(c.Request, flashMessageSessionName)
	if err != nil {
		log.Error(merry.Prepend(err, errGetFlashMessages))
	}
	session.AddFlash(message, normalFlashKey)
	err = session.Save(c.Request, c.Response)
	if err != nil {
		log.Error(merry.Prepend(err, errAddFlashMessage))
	}
}

// AddErrorFlash adds an error flash message to the store to be viewed by the
// user upon next page request.
//
// If you need to add a flash message, you should do so before writing to the
// response(Writer)! This is due to gorilla's session.Save().
func (c *Context) AddErrorFlash(message string) {
	session, err := store.Get(c.Request, flashMessageSessionName)
	if err != nil {
		log.Error(merry.Prepend(err, errGetFlashMessages))
	}
	session.AddFlash(message, errorFlashKey)
	err = session.Save(c.Request, c.Response)
	if err != nil {
		log.Error(merry.Prepend(err, errAddFlashMessage))
	}
}
