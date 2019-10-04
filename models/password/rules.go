package password

import (
	"fmt"
	"strings"

	"github.com/ansel1/merry"
)

const (
	// MinLength is the minimum number of characters a password MUST contain.
	MinLength = 10
	// MaxLength is the maximum number of characters a password MUST contain.
	MaxLength = 64
)

var (
	// ErrPasswordWeak indicates the password does not meet our rules.
	ErrPasswordWeak = merry.
			New("password does not meet the requirements and is considered weak").
			WithUserMessage("Password does not meet the requirements and is considered too weak.")
	// ErrPasswordLength indicates the password is not within the min and max
	// length requirement.
	ErrPasswordLength = merry.
				WithMessage(ErrPasswordWeak, fmt.Sprintf("password must be %d-%d characters long", MinLength, MaxLength)).
				WithUserMessage(fmt.Sprintf("Password must be %d-%d characters long.", MinLength, MaxLength))
	// ErrPasswordContainsName indicates the password is not allowed because it
	// contains part or all of their name or username.
	ErrPasswordContainsName = merry.
				WithMessage(ErrPasswordWeak, "password cannot contain the first, last, and/or username").
				WithUserMessage("Password cannot contain your any part of your name or username.")
)

// CheckPasswordRules returns nil if the password meets all of the requirements.
// Otherwise, it returns an error describing which rule it currently violates.
//
// Rules:
//   - Must be at between MinLength and MaxLength
//
// TODO: Check against a list of common passwords - JZ
// TODO: Check this isn't equal to their current password? - JZ
func CheckPasswordRules(username string, firstName string, lastName string,
	password string) error {
	// 1. Check password length
	length := len(password)
	if length < MinLength || length > MaxLength {
		return ErrPasswordLength.Here()
	}
	// 2. Check password against the username, first name, and last name
	lowerCasePassword := strings.ToLower(password)
	if strings.Contains(lowerCasePassword, strings.ToLower(username)) ||
		strings.Contains(lowerCasePassword, strings.ToLower(firstName)) ||
		strings.Contains(lowerCasePassword, strings.ToLower(lastName)) {
		return ErrPasswordContainsName.Here()
	}
	// 3. TODO: Check password against common passwords (reject if found)
	return nil
}
