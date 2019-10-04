package password

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"strings"

	"github.com/ansel1/merry"
	"golang.org/x/crypto/bcrypt"
)

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func hashMD5(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	hashedPassword := hasher.Sum(nil)
	return toBase64(hashedPassword)
}

func validMD5(password string, hashedPassword string) (bool, error) {
	hashPrefix := hashedPassword[0:5]
	hash := hashedPassword[5:]

	if len(hash) < 1 || hashPrefix != "{MD5}" {
		return false, merry.New("not a MD5 password")
	}

	if hashMD5(password) != hash {
		return false, nil // Passwords DO NOT match
	}
	return true, nil
}

func hashSHA1(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	hashedPassword := hasher.Sum(nil)
	return toBase64(hashedPassword)
}

func validSHA1(password string, hashedPassword string) (bool, error) {
	hashPrefix := hashedPassword[0:5]
	hash := hashedPassword[5:]

	if len(hash) < 1 || hashPrefix != "{SHA}" {
		return false, merry.New("not a SHA1 password")
	}

	if hashSHA1(password) != hash {
		return false, nil // Passwords DO NOT match
	}
	return true, nil
}

func hashSSHA(password string, salt []byte) []byte {
	pass := []byte(password)
	str := append(pass[:], salt[:]...)
	sum := sha1.Sum(str)
	result := append(sum[:], salt[:]...)
	return result
}

func validSSHA(password string, hashedPassword string) (bool, error) {
	if len(hashedPassword) < 7 || string(hashedPassword[0:6]) != "{SSHA}" {
		return false, merry.New("not a SSHA password")
	}

	data, err := base64.StdEncoding.DecodeString(hashedPassword[6:])
	if len(data) < 21 || err != nil {
		return false, merry.New("Base64 Decode Failed")
	}

	newhash := hashSSHA(password, data[20:])
	hashedpw := base64.StdEncoding.EncodeToString(newhash)

	if hashedpw == hashedPassword[6:] {
		return true, nil
	}

	return false, nil
}

func hashBCRYPT(password string) (string, error) {
	// Default Cost is currently 10. We had it set to 15, but it took too long
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), merry.Wrap(err)
}

func validBCRYPT(password string, hashedPassword string) (bool, error) {
	// this method only returns an error if the hash doesn't match, unliked
	// our other methods which have different error conditions
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword),
		[]byte(password))
	return err == nil, nil
}

// Valid returns true if the given password matches the password hash.
//
// If the password is valid, it will also return false if it is not using a
// secure hashing algorithm and needs updated (currently bcrypt).
func Valid(password string, hashedPassword string) (valid bool,
	insecure bool, err error) {

	if strings.HasPrefix(hashedPassword, "$2a$") {
		insecure = false
		valid, err = validBCRYPT(password, hashedPassword)
		return
	}

	insecure = true
	switch {
	case strings.HasPrefix(hashedPassword, "{MD5}"):
		valid, err = validMD5(password, hashedPassword)
		return
	case strings.HasPrefix(hashedPassword, "{SHA}"):
		valid, err = validSHA1(password, hashedPassword)
		return
	case strings.HasPrefix(hashedPassword, "{SSHA}"):
		valid, err = validSSHA(password, hashedPassword)
		return
	default:
		return false, false, merry.New("could not identify password hash type")
	}
}

// Hash take a plaintext password and returns a securely hashed version.
//
// Currently uses bcrypt
//
// Use this instead of a specific hashing algorithm so we can change which
// algorithm is used between versions.
func Hash(password string) (string, error) {
	return hashBCRYPT(password)
}
