package secrets

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ansel1/merry"
	"github.com/gorilla/securecookie"

	"github.com/joshsziegler/zauth/pkg/file"
)

const (
	secretsPath = `secrets.json`
)

var (
	// store is our "global" secrets variable, which is accessible only through
	// the getters defined below
	store secrets
)

// secrets holds four key secrets for running the web server that should be
// saved between runs (to prevent user sessions, cookies, and CSRF tokens from
// being invalidated, and to provide password reset tokens).
//
// These cannot be changed after init(), and are only provided via getters!
type secrets struct {
	AuthKey             []byte
	EncryptionKey       []byte
	CSRFKey             []byte
	PasswordResetSecret []byte
}

// AuthKey is used to authenticate the cookie value using HMAC.
// Gorilla docs suggest 32 or 64 bytes long.
func AuthKey() []byte {
	return store.AuthKey
}

// EncriptionKey is used to encrypt the cookie value.
// Gorilla docs suggest 32 bytes long for AES-256.
func EncryptionKey() []byte {
	return store.EncryptionKey
}

// CSRFKey is used for Cross Site Request Forgery (CSRF) protection
func CSRFKey() []byte {
	return store.CSRFKey
}

// PasswordResetSecret is used for securing password reset tokens.
// For an example, see: github.com/dchest/passwordreset
func PasswordResetSecret() []byte {
	return store.PasswordResetSecret
}

// init loads the secrets JSON file from disk (if it exists), and if any of the
// secrets are missing, it will create them and save the resulting secrets back
// to disk as JSON.
func init() {
	// Read the secrets JSON file from disk (if the file exists)
	if file.Exists(secretsPath) {
		data, err := ioutil.ReadFile(secretsPath)
		if err != nil {
			panic(merry.Prepend(err, "error reading "+secretsPath))
		}
		err = json.Unmarshal(data, &store)
		if err != nil {
			panic(merry.Prepend(err, "error parsing JSON from "+secretsPath))
		}
	}

	// Only write this file to disk if one of the secrets was empty/missing
	writeFile := false
	if len(store.AuthKey) < 1 {
		store.AuthKey = securecookie.GenerateRandomKey(64)
		writeFile = true
	}
	if len(store.EncryptionKey) < 1 {
		store.EncryptionKey = securecookie.GenerateRandomKey(32)
		writeFile = true
	}
	if len(store.CSRFKey) < 1 {
		store.CSRFKey = securecookie.GenerateRandomKey(32)
		writeFile = true
	}
	if len(store.PasswordResetSecret) < 1 {
		store.PasswordResetSecret = securecookie.GenerateRandomKey(32)
		writeFile = true
	}

	if writeFile {
		// Save to disk
		data, err := json.Marshal(&store)
		if err != nil {
			panic(merry.Prepend(err, "error marshaling secrets to JSON"))
		}
		err = ioutil.WriteFile(secretsPath, data, 0644)
		if err != nil {
			panic(merry.Prepend(err, "error writing secrets to "+secretsPath))
		}
	}
}
