package tests

import (
	"fmt"
	"math/rand"
)

// TestTXKV returns a text transaction, allong with expected key, value pair
func TestTxKV() ([]byte, []byte, []byte) {
	k := RandAsciiBytes(8)
	v := RandAsciiBytes(8)
	return k, v, []byte(fmt.Sprintf("%s=%s", k, v))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandAsciiBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
