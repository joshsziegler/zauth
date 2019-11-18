package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ansel1/merry"
)

const (
	flashMessageSessionName = `zauth-flashes`
	errorFlashKey           = `error-flash`
	normalFlashKey          = `message-flash`
)

var (
	// ErrGettingFlashMessages denotes an error when retrieving flash storage.
	ErrGettingFlashMessages = merry.
		New("error getting flash messages").
		WithUserMessage("Failed to retrieve flash messages.")
)

// getFlashMessages returns a single error flash message and a single normal
// flash message. This is to restrict the HTTP handlers to a single, most-
// important message of each type.
//
// Normal flash messages are not errors, and are typically informing the user
// that an operation was successful such as logging out, or creating a new user.
func getFlashMessages(w http.ResponseWriter, r *http.Request) (messageFlash string, errorFlash string) {
	session, err := store.Get(r, flashMessageSessionName)
	if err != nil {
		// This could be an error, but more likely there simply aren't any
		// flash message set. So simply return empty strings for both.
		return
	}
	// Get any error flash message
	fm := session.Flashes(errorFlashKey)
	if fm != nil {
		errorFlash = fmt.Sprintf("%v", fm[0])
	}
	// Get any normal flash message
	fm = session.Flashes(normalFlashKey)
	if fm != nil {
		messageFlash = fmt.Sprintf("%v", fm[0])
	}
	// Always save after retrieving flash messages so they are removed
	err = session.Save(r, w)
	if err != nil {
		log.Error(err)
	}
	return
}

// addNormalFlashMessage adds an message to the flash messages store to be
// viewed by the user upon next page request.
//
// If you need to add a flash message, you should do so before writing to the
// response(Writer)! This is due to gorilla's session.Save().
//
// Normal flash messages are not errors, and are typically informing the user
// that an operation was successful such as logging out, or creating a new user.
func addNormalFlashMessage(w http.ResponseWriter, r *http.Request, message string) error {
	session, err := store.Get(r, flashMessageSessionName)
	if err != nil {
		return ErrGettingFlashMessages.Here()
	}
	session.AddFlash(message, normalFlashKey)
	err = session.Save(r, w)
	if err != nil {
		log.Error(err)
	}
	return nil
}

// addErrorFlashMessage adds an error message to the flash messages store to be
// viewed by the user upon next page request.
//
// If you need to add a flash message, you should do so before writing to the
// response(Writer)! This is due to gorilla's session.Save().
func addErrorFlashMessage(w http.ResponseWriter, r *http.Request, message string) error {
	session, err := store.Get(r, flashMessageSessionName)
	if err != nil {
		return ErrGettingFlashMessages.Here()
	}
	session.AddFlash(message, errorFlashKey)
	err = session.Save(r, w)
	if err != nil {
		log.Error(err)
	}
	return nil
}
