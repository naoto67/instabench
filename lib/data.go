package lib

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generate(size int64) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}
