package test

import (
	"math/rand"
	"time"
)

var source *rand.Rand

func init() {
	source = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomString returns a random alphanumeric string useful for testing purposes.
func RandomString() string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, 8)
	for i := range s {
		s[i] = chars[source.Intn(len(chars))]
	}
	return string(s)
}
