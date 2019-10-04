package password

import (
	"testing"
)

const (
	correctPassword = `Testing123!`
)

func TestSHA1(t *testing.T) {
	// correctSHA1b := "e1NIQX1IaElQL3lPdTJ2STZzVS9SVG84dEF4V1R3aUk9"
	correctSHA1p := "{SHA}HhIP/yOu2vI6sU/RTo8tAxWTwiI="
	hashedPassword := hashSHA1(correctPassword)
	if correctSHA1p != "{SHA}"+hashedPassword {
		t.Errorf("SHA1 Hashing Failed: %v != %v", correctSHA1p, hashedPassword)
	}
	ok, err := validSHA1("Testing123!", correctSHA1p)
	if !ok || err != nil {
		t.Errorf("SHA1 Failed: %v, %s", ok, err)
	}
}

func TestMD5(t *testing.T) {
	// correctMD5b := "e01ENX11UFdNTUdlUmE3KzFCMmFxaTkzVUxBPT0="
	correctMD5p := "{MD5}uPWMMGeRa7+1B2aqi93ULA=="
	hashedPassword := hashMD5(correctPassword)
	if correctMD5p != "{MD5}"+hashedPassword {
		t.Errorf("MD5 Hashing Failed: %v != %v", correctMD5p, hashedPassword)
	}

	ok, err := validMD5("Testing123!", correctMD5p)
	if !ok || err != nil {
		t.Errorf("MD5 Comparison Failed: %v, %s", ok, err)
	}
}

func TestSSHA(t *testing.T) {
	// correctSSHAb := "e1NTSEF9ZlBUODN6ZFcxRUpQZ20vZFVpQ2U3VTQ2L0t4UjY5VXY="
	correctSSHAp := "{SSHA}fPT83zdW1EJPgm/dUiCe7U46/KxR69Uv"
	ok, err := validSSHA("Testing123!", correctSSHAp)
	if !ok || err != nil {
		t.Errorf("SSHA Failed: %v, %s", ok, err)
	}
}

func TestBCRYPT(t *testing.T) {
	hashedPassword, err := hashBCRYPT(correctPassword)
	ok, err := validBCRYPT("Testing123!", hashedPassword)
	if !ok || err != nil {
		t.Errorf("BCRYPT Failed: %v, %s", ok, err)
	}
}
