package consts

import (
	"context"
	"dwitter_go_graphql/prisma/db"

	"cloud.google.com/go/storage"
	"github.com/functionalfoundry/graphqlws"
	"golang.org/x/crypto/bcrypt"
)

var SubscriptionManager graphqlws.SubscriptionManager

const DefaultPFPURL = "https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media"

var Bucket *storage.BucketHandle

var MediaCreatedButNotUsed map[string]bool

const LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var Client *db.PrismaClient

var BaseCtx context.Context

// Check given credentials and return true if valid
func CheckCreds(username string, password string) (bool, error) {
	user, err := Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(BaseCtx)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
