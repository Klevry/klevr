package common

import (
	"math/rand"
	"time"
)

const (
	charset string = "abcdefghijklmnopqrstuvwxyz" +
		"ABCEDFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789" +
		"!@#$%^&*()-_+=|/?~`:;<>,.'{}[]"

	lenCharset int = len(charset)
)

var seedRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GetKey get key for encryption
func GetKey(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seedRand.Intn(lenCharset)]
	}

	return string(b)
}
