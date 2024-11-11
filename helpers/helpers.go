package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const SessionTokenBytes = 32

// returns a random string using crypto/rand.
// n is the number of bytes being used to generate the random string.
func GenerateSessionToken() (string, error) {
	// generate a 32 byte for session token
	b := make([]byte, SessionTokenBytes)
	nRead, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("bytes: %w", err)
	}
	if nRead < 32 {
		return "", fmt.Errorf("bytes: didn't read enough random bytes")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
