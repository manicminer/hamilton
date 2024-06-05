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
	alpha := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numeric := []rune("0123456789")
	chars := append(alpha, numeric...)

	s := make([]rune, 8)
	for i := range s {
		if i == 0 {
			s[i] = alpha[source.Intn(len(alpha))]
		}

		s[i] = chars[source.Intn(len(chars))]
	}
	return string(s)
}
