package utils

import (
	"math/rand"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StringWithCharset(lrand *rand.Rand, length int, charset string) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func String(rand *rand.Rand, length int) string {
	return StringWithCharset(rand, length, charset)
}
