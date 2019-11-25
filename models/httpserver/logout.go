package httpserver

import (
	"net/http"
)

// LogoutGet handles a user's request to logout of zauth.
func LogoutGet(c *Context, w http.ResponseWriter, r *http.Request) error {
	// Always returns a session, even if it's empty
	session, err := store.Get(r, sessionName)
	if err != nil {
		return ErrGetSecureSession.Here()
	}
	// Delete the Username key-value pair
	delete(session.Values, "Username")
	// Save the updated session BEFORE writing the response so it's sent
	err = session.Save(r, w)
	if err != nil {
		return ErrInternal.Here()
	}
	addNormalFlashMessage(w, r, "Successfully logged out.")
	http.Redirect(w, r, urlLogin, http.StatusFound) // StatusFound ~ 302
	return nil
}
