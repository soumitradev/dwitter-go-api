package main

import (
	crypto_rand "crypto/rand"
	"dwitter_go_graphql/prisma/db"
	"encoding/binary"
	math_rand "math/rand"
	"reflect"
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

// Hash function with modifications from: https://github.com/juliangruber/go-intersect/blob/2e99d8c0a75f6975a52f7efeb81926a19b221214/intersect.go#L42-L62
// Hash has complexity: O(n * x) where x is a factor of hash function efficiency (between 1 and 2)
func HashIntersectUsers(a []db.UserModel, b []db.UserModel) []db.UserModel {
	set := make([]db.UserModel, 0)
	hash := make(map[string]bool)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		elt := el.(db.UserModel).Username
		hash[elt] = true
	}

	for i := 0; i < bv.Len(); i++ {
		el := bv.Index(i).Interface()
		elt := el.(db.UserModel).Username
		if _, found := hash[elt]; found {
			set = append(set, el.(db.UserModel))
		}
	}

	return set
}
