package service

import (
	"math/rand/v2"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

// generateShortCode returns a pseudo-random string of codeLength characters,
// composed of characters from alphabet. It does not guarantee uniqueness;
// callers must handle collisions when persisting the code.
func generateShortCode() string {
	r := make([]byte, codeLength)
	l := len(alphabet)

	for i := range r {
		alphabetIndex := rand.IntN(l)
		r[i] = alphabet[alphabetIndex]
	}

	return string(r)
}