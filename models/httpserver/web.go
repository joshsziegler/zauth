package httpserver

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"

	"github.com/gobuffalo/packr"
	logging "github.com/op/go-logging"

	"github.com/joshsziegler/zauth/pkg/httpserver"
	"github.com/joshsziegler/zauth/pkg/secrets"
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
	templates = httpserver.MustLoadBoxedTemplates(boxTemplates)

	// Create the HTTP route handler
	r := mux.NewRouter().StrictSlash(true)
	//r.NotFoundHandler = http.HandlerFunc(HandlerPageNotFound)
	r.NotFoundHandler = Wrap(r, pageNotFound, false)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(boxStatic)))
	r.Handle("/", Wrap(r, LoginOrUserPageGet, false)).Methods("GET")
	r.Handle(urlLogin, Wrap(r, LoginGetPost, false)).Methods("GET", "POST").Name("login")
	r.Handle("/logout", Wrap(r, LogoutGet, true)).Methods("GET")
	r.Handle("/user/new", Wrap(r, NewUserGet, true)).Methods("GET")
	r.Handle("/user/new", Wrap(r, NewUserPost, true)).Methods("POST")
	r.Handle("/users", Wrap(r, UserListGet, true)).Methods("GET")
	r.Handle("/users/{username}", Wrap(r, UserDetailGet, true)).Methods("GET").Name("userDetail")
	r.Handle("/users/{username}/password", Wrap(r, userSetPassword, true)).Methods("GET", "POST")
	r.Handle("/users/{username}/{isEnabled:(?:enable|disable)}", Wrap(r, userSetEnabled, true)).Methods("GET")
	r.Handle("/users/{username}/groups/{groupname}/{addOrRemove:(?:add|remove)}", Wrap(r, userAddRemoveGroups, true)).Methods("GET")
	r.Handle("/groups", Wrap(r, GroupListGet, true)).Methods("GET")
	r.Handle("/group/new", Wrap(r, NewGroupGet, true)).Methods("GET")
	r.Handle("/group/new", Wrap(r, NewGroupPost, true)).Methods("POST")
	// /groups/{groupname} - If none, show all if admin or redirect to self TODO: implement
	r.Handle("/reset-password/{token}", Wrap(r, PasswordResetGetPost, false)).Methods("GET", "POST")

	// Start the HTTP servers
	log.Infof("HTTP server listening on: %s", listenTo)
	err := http.ListenAndServe(listenTo,
		csrf.Protect(secrets.CSRFKey(), csrf.Secure(isProduction))(r))
	if err != nil {
		log.Fatalf("error running http server: %s", err)
	}
}
