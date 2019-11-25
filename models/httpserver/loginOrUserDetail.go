package httpserver

import (
	"fmt"
	"net/http"
)

// LoginOrUserPageGet asks the user to login if they aren't already. If they
// are, it will redirect them to their user details page.
func LoginOrUserPageGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	if c.User == nil { // nil if not logged in
		http.Redirect(w, r, urlLogin, 302)
		return nil
	}
	uri := fmt.Sprintf("/users/%s", c.User.Username)
	http.Redirect(w, r, uri, 302)
	return nil
}
