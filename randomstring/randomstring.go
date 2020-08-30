package randomstring

import (
	"math/rand"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// StringWithCharset returns a random string of given length and charset.
func StringWithCharset(lrand *rand.Rand, length int, charset string) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[lrand.Intn(len(charset))]
	}

	return string(b)
}

// String returns a random string of given length and default charset.
func String(rand *rand.Rand, length int) string {
	return StringWithCharset(rand, length, charset)
}
