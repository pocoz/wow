package tools

import (
	"encoding/base64"
	"math/rand"
)

const defaultLength = 16

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandomString() string {
	b := make([]byte, defaultLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return base64.StdEncoding.EncodeToString(b)
}
