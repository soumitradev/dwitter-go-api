// Package common stores useful common global variables and functions used across packages.
package common

import (
	"context"
	"dwitter_go_graphql/prisma/db"
	"errors"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/functionalfoundry/graphqlws"
	"golang.org/x/crypto/bcrypt"
)

const LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const DefaultPFPURL = "https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media"

var Client *db.PrismaClient
var BaseCtx context.Context
var Bucket *storage.BucketHandle
var MediaCreatedButNotUsed map[string]bool
var SubscriptionManager graphqlws.SubscriptionManager
var GraphqlwsHandler http.Handler

type HTTPError struct {
	Error string `json:"error"`
}

func init() {
	BaseCtx = context.Background()
	MediaCreatedButNotUsed = make(map[string]bool)
}

// Check given credentials and return true if valid
func CheckCreds(username string, password string) (bool, error) {
	user, err := Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(BaseCtx)
	if err != nil {
		return false, errors.New("username/password error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("username/password error")
	}
	return true, nil
}
