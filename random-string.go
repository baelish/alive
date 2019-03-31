package main

import (
	"math/rand"
)

const randomBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randomBytes[rand.Int63()%int64(len(randomBytes))]
	}
	return string(b)
}
