package alias

import (
	"math/rand"
	"time"
)

const (
	Chars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Length = 6
)

func Generate() string {
	rand.NewSource(time.Now().UnixNano())

	alias := make([]byte, Length)

	for i := range alias {
		alias[i] = Chars[rand.Intn(len(Chars))]
	}

	return string(alias)
}
