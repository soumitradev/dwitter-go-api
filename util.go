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

// Modified Merge sort for merging dweets and redweets adapted from: https://www.golangprograms.com/golang-program-for-implementation-of-mergesort.html
func MergeDweetRedweetList(dweets []db.DweetModel, redweets []db.RedweetModel) []interface{} {
	result := make([]interface{}, len(dweets)+len(redweets))

	i := 0
	for len(dweets) > 0 && len(redweets) > 0 {
		if dweets[0].PostedAt.Unix() > redweets[0].RedweetTime.Unix() {
			result[i] = dweets[0]
			dweets = dweets[1:]
		} else {
			result[i] = redweets[0]
			redweets = redweets[1:]
		}
		i++
	}

	for j := 0; j < len(dweets); j++ {
		result[i] = dweets[j]
		i++
	}
	for j := 0; j < len(redweets); j++ {
		result[i] = redweets[j]
		i++
	}

	return result
}

// Modified Merge sort for merging dweets and redweets adapted from: https://www.golangprograms.com/golang-program-for-implementation-of-mergesort.html
func MergeDweetLists(dweetsA []db.DweetModel, dweetsB []db.DweetModel) []db.DweetModel {
	result := make([]db.DweetModel, len(dweetsA)+len(dweetsB))

	i := 0
	for len(dweetsA) > 0 && len(dweetsB) > 0 {
		if dweetsA[0].PostedAt.Unix() > dweetsB[0].PostedAt.Unix() {
			result[i] = dweetsA[0]
			dweetsA = dweetsA[1:]
		} else {
			result[i] = dweetsB[0]
			dweetsB = dweetsB[1:]
		}
		i++
	}

	for j := 0; j < len(dweetsA); j++ {
		result[i] = dweetsA[j]
		i++
	}
	for j := 0; j < len(dweetsB); j++ {
		result[i] = dweetsB[j]
		i++
	}

	return result
}

// Modified Merge sort for merging dweets and redweets adapted from: https://www.golangprograms.com/golang-program-for-implementation-of-mergesort.html
func MergeRedweetLists(redweetsA []db.RedweetModel, redweetsB []db.RedweetModel) []db.RedweetModel {
	result := make([]db.RedweetModel, len(redweetsA)+len(redweetsB))

	i := 0
	for len(redweetsA) > 0 && len(redweetsB) > 0 {
		if redweetsA[0].RedweetTime.Unix() > redweetsB[0].RedweetTime.Unix() {
			result[i] = redweetsA[0]
			redweetsA = redweetsA[1:]
		} else {
			result[i] = redweetsB[0]
			redweetsB = redweetsB[1:]
		}
		i++
	}

	for j := 0; j < len(redweetsA); j++ {
		result[i] = redweetsA[j]
		i++
	}
	for j := 0; j < len(redweetsB); j++ {
		result[i] = redweetsB[j]
		i++
	}

	return result
}
