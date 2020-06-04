package httpserver

import "github.com/ansel1/merry"

var (
	ErrInternal = merry.
			New("internal server error").
			WithUserMessage("Sorry, but the server encountered an error.")
	ErrBadRequest = merry.
			New("bad request").
			WithUserMessage("Your request was bad and/or invalid.")
	ErrPermissionDenied = merry.
				New("permission denied").
				WithUserMessage("You do not have permission to view this page.")
	ErrGettingUser = merry.
			New("failed to retrieve user").
			WithUserMessage("Failed to retrieve user record.")
	ErrGetSecureSession = merry.
				WithMessage(ErrInternal, "secure session exists, but could not be decoded")
	ErrRequestArgument = merry.
				New("invalid HTTP request argument").
				WithUserMessage("One or more of your request arguments was invalid.")
)
