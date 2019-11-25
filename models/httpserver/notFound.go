package httpserver

import (
	"net/http"
)

// pageNotFound is a sub-handler that shows a 404 page with proper top navivation
func pageNotFound(c *Context, w http.ResponseWriter, r *http.Request) error {
	Error(w, http.StatusNotFound, "Error", "Sorry, but that page doesn't exist.", c.User)
	return nil
}
