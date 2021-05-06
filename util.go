package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	math_rand "math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func initRandom() {
	// More secure random seeding than usual: https://stackoverflow.com/a/54491783
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

func genID(n int) string {
	// Make a random string of length n
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[math_rand.Intn(len(letterBytes))]
	}
	return string(b)
}
