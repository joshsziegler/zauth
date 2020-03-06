package httpserver

import (
	"fmt"
	"net/http"

	"github.com/joshsziegler/zauth/pkg/user"
	"github.com/joshsziegler/zgo/pkg/log"
)

// Render loads the HTML template 'name' using the provided 'data' struct and
// buffers the output.
func Render(w http.ResponseWriter, name string, data interface{}) {
	// Get the template (may not exist)
	t := templates.Lookup(name)
	if t == nil {
		errStr := fmt.Sprintf("error: could not find template '%s'", name)
		log.Error(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}
	// Execute the template with the provided data
	err := t.Execute(w, data)
	if err != nil {
		errStr := fmt.Sprintf("error rendering template '%s': %s", name, err)
		log.Error(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}
	// Set the content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

type errorData struct {
	User    *user.User
	Header  string
	Message string
}

// Error replies to the request with the specified error message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
func Error(w http.ResponseWriter, code int, header string, message string,
	user *user.User) {
	w.WriteHeader(code)
	Render(w, "error.html", errorData{user, header, message})
}

// ErrorInternal is a helper that response with a generic error 500 message.
func ErrorInternal(w http.ResponseWriter) {
	Error(w, 500, "Error", "Sorry, but the server encountered an error.", nil)
}

// ErrorUnauthorized is a helper that responds with the a 403 error code and a
// custom message (e.g. Sorry, but you don't have permission to view that.)"
//
// See this Stack Overflow question for why 403 is used instead of 401:
//     https://stackoverflow.com/a/6937030
func ErrorUnauthorized(w http.ResponseWriter, message string, user *user.User) {
	Error(w, 403, "Unauthorized", message, user)
}
