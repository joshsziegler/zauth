package httpserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"

	"github.com/gobuffalo/packr"
	logging "github.com/op/go-logging"

	"github.com/joshsziegler/zauth/models/secrets"
)

var log *logging.Logger

// DB is our shared database connection (handles connection pooling, and is
// goroutine-safe)
var DB *sqlx.DB

const (
	sessionName = `zauth-session`
	urlLogin    = `/login`
)

// secure session store
var store *sessions.CookieStore

// templates holds our loaded Go/HTML templates
var templates *template.Template

// templateHelpers allows us to use these custom functions in our templates
var templateHelpers = template.FuncMap{
	"inc":      func(i int) int { return i + 1 },
	"multiply": func(x int, y int) int { return x * y },
	"marshal": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
	"FormatTimeAsRFC822": func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05 MST")
	},
	"HumanizeTime": humanize.Time,
	"ToLower":      strings.ToLower,
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

// Listen performs setup and runs the Web server (blocking)
func Listen(logger *logging.Logger, database *sqlx.DB, listenTo string,
	isProduction bool) {
	log = logger
	DB = database

	// Setup sessions using secure cookies
	store = sessions.NewCookieStore(secrets.AuthKey(), secrets.EncryptionKey())
	// Set Cookie options to expire sessions and protect against some attacks
	store.Options = &sessions.Options{
		Path:     "/",                  // Send cookies with every page request for this domain
		MaxAge:   60 * 15,              // Expire in 15 minutes to force logout
		SameSite: http.SameSiteLaxMode, // XSS protection; See:  https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies#SameSite_cookies
		HttpOnly: true,                 // Prevent JavaScript access to this cookie
		// Secure:   true,                    // Only sent over HTTPS
	}
	// Load static assets (from disk [dev] or binary [build])
	boxStatic := packr.NewBox("../../static")
	boxTemplates := packr.NewBox("../../templates")
	// Load our templates
	templates = MustLoadBoxedTemplates(boxTemplates)

	// Create the HTTP route handler
	r := mux.NewRouter()
	//r.NotFoundHandler = http.HandlerFunc(HandlerPageNotFound)
	r.NotFoundHandler = wrapHandler(r, pageNotFound, false)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(boxStatic)))
	r.Handle("/", wrapHandler(r, LoginOrUserPageGet, false)).Methods("GET")
	r.Handle(urlLogin, wrapHandler(r, LoginGetPost, false)).Methods("GET", "POST").Name("login")
	r.Handle("/logout", wrapHandler(r, LogoutGet, true)).Methods("GET")
	r.Handle("/user/new", wrapHandler(r, NewUserGet, true)).Methods("GET")
	r.Handle("/user/new", wrapHandler(r, NewUserPost, true)).Methods("POST")
	r.Handle("/users", wrapHandler(r, UserListGet, true)).Methods("GET")
	r.Handle("/users/{username}", wrapHandler(r, UserDetailGet, true)).Methods("GET").Name("userDetail")
	r.Handle("/users/{username}/password", wrapHandler(r, userSetPassword, true)).Methods("GET", "POST")
	r.Handle("/users/{username}/{isEnabled:(?:enable|disable)}", wrapHandler(r, userSetEnabled, true)).Methods("GET")
	r.Handle("/users/{username}/groups/{groupname}/{addOrRemove:(?:add|remove)}", wrapHandler(r, userAddRemoveGroups, true)).Methods("GET")
	r.Handle("/groups", wrapHandler(r, GroupListGet, true)).Methods("GET")
	r.Handle("/group/new", wrapHandler(r, NewGroupGet, true)).Methods("GET")
	r.Handle("/group/new", wrapHandler(r, NewGroupPost, true)).Methods("POST")
	// /groups/{groupname} - If none, show all if admin or redirect to self TODO: implement
	r.Handle("/reset-password/{token}", wrapHandler(r, PasswordResetGetPost, false)).Methods("GET", "POST")

	// Start the HTTP servers
	log.Infof("HTTP server listening on: %s", listenTo)
	err := http.ListenAndServe(listenTo,
		csrf.Protect(secrets.CSRFKey(), csrf.Secure(isProduction))(r))
	if err != nil {
		log.Fatalf("error running http server: %s", err)
	}
}
