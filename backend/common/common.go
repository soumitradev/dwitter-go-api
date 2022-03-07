// Package common stores useful common global variables and functions used across packages.
package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/soumitradev/Dwitter/backend/prisma/db"

	"cloud.google.com/go/storage"
	"github.com/functionalfoundry/graphqlws"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/sendgrid/sendgrid-go"
)

const DefaultPFPURL = "https://storage.googleapis.com/download/storage/v1/b/dwitter-72e9d.appspot.com/o/pfp%2Fdefault.jpg?alt=media"

var Client *db.PrismaClient
var BaseCtx context.Context
var Bucket *storage.BucketHandle
var MediaCreatedButNotUsed map[string]bool
var AccountCreatedButNotVerified map[string]string
var SubscriptionManager graphqlws.SubscriptionManager
var GraphqlwsHandler http.Handler
var Validate *validator.Validate
var SendgridClient *sendgrid.Client

type HTTPError struct {
	Error string `json:"error"`
}

func init() {
	BaseCtx = context.Background()
	MediaCreatedButNotUsed = make(map[string]bool)
	AccountCreatedButNotVerified = make(map[string]string)
}

func InitSendgrid() {
	SendgridClient = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
}

// Check given credentials and return true if valid
func CheckCreds(username string, password string) (bool, error) {
	user, err := Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(BaseCtx)
	if err != nil {
		return false, errors.New("username/password error")
	}

	if !user.Verified {
		return false, errors.New("account not verified: please check your email for a verification link")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("username/password error")
	}
	return true, nil
}

// Delete a Dweet
func InternalDeleteDweet(postID string) (*db.DweetModel, error) {
	// Get all the replies to the post (these need to be deleted first since they depend on the root Dweet)
	post, err := Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	// If the Dweet itself is a reply, remove the reply from the original post
	if post.IsReply {
		// Find the dweet that was replied to
		id, exist := post.OriginalReplyID()
		if !exist {
			return nil, errors.New("original Dweet not found")
		}

		// Remove the Reply from the post
		_, err := Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(id),
		).Update(
			db.Dweet.ReplyCount.Decrement(1),
			db.Dweet.ReplyDweets.Unlink(
				db.Dweet.ID.Equals(postID),
			),
		).Exec(BaseCtx)
		if err != nil {
			return nil, err
		}
	}

	dweet, err := Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.RedweetDweets.Fetch().With(
			db.Redweet.Author.Fetch(),
		).OrderBy(
			db.Redweet.RedweetTime.Order(db.DESC),
		),
	).Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	for _, redweet := range dweet.RedweetDweets() {
		InternalDeleteRedweet(redweet.OriginalRedweetID, redweet.Author().Username)
	}

	dweet, err = Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.ReplyDweets.Fetch().OrderBy(
			db.Dweet.LikeCount.Order(db.DESC),
		),
	).Exec(BaseCtx)
	if err != nil {
		return nil, err
	}
	for _, daughterDweet := range dweet.ReplyDweets() {
		InternalDeleteDweet(daughterDweet.ID)
	}

	_, err = Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Delete().Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	// The following comment block is kept as a homage to the great recursive SQL function that once resided here.
	// May the soul of this legendary query rest in peace. It was a honor to use you.

	// Delete all the dependent posts (this includes redweets and replies to the post) recursively using RAW SQL
	// We use RAW SQL here because prisma-go-client doesn't support cascade deletes yet:
	// Link: https://github.com/prisma/prisma-client-go/issues/201
	// Recursive SQL function with modifications from: https://stackoverflow.com/q/10381243
	// delQuery := `WITH RECURSIVE all_posts (id, parentid1, root_id) AS (SELECT t1.db_id, t1.original_reply_id AS parentid1, t1.db_id AS root_id FROM public."Dweet" t1 UNION ALL SELECT c1.db_id, c1.original_reply_id AS parentid1, p.root_id FROM public."Dweet" c1 JOIN all_posts p ON (p.id = c1.original_reply_id) ) DELETE FROM public."Dweet" WHERE db_id IN ( SELECT id FROM all_posts WHERE root_id = $1);`
	// _, err = Client.Prisma.ExecuteRaw(delQuery, post.DbID).Exec(BaseCtx)

	return post, err
}

// Delete a User
func InternalDeleteUser(username string) (*db.UserModel, error) {
	// Get the user
	basicUser, err := Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	user, err := Client.User.FindUnique(
		db.User.Username.Equals(username),
	).With(
		db.User.Dweets.Fetch(),
		db.User.RedweetedDweets.Fetch(),
		db.User.Redweets.Fetch(),
		db.User.LikedDweets.Fetch(),
		db.User.Followers.Fetch(),
		db.User.Following.Fetch(),
	).Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	// Delete all dependent objects, and adjust all relations
	for _, dweet := range user.Dweets() {
		if _, err := InternalDeleteDweet(dweet.ID); err != nil {
			return nil, err
		}
	}

	for _, redweet := range user.Redweets() {
		if _, err := InternalDeleteRedweet(redweet.OriginalRedweetID, user.Username); err != nil {
			return nil, err
		}
	}

	for _, liked := range user.LikedDweets() {
		if _, err := InternalUnlike(liked.ID, user.Username); err != nil {
			return nil, err
		}
	}

	for _, follower := range user.Followers() {
		if _, err := InternalUnfollow(user.Username, follower.Username); err != nil {
			return nil, err
		}
	}

	for _, followed := range user.Following() {
		if _, err := InternalUnfollow(followed.Username, user.Username); err != nil {
			return nil, err
		}
	}

	_, err = Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Delete().Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	return basicUser, err
}

// Remove a Redweet
func InternalDeleteRedweet(postID string, username string) (*db.RedweetModel, error) {
	// Get the redweet
	user, err := Client.User.FindUnique(
		db.User.Username.Equals(username),
	).With(
		db.User.Redweets.Fetch(
			db.Redweet.OriginalRedweetID.Equals(postID),
		).With(
			db.Redweet.Author.Fetch(),
			db.Redweet.RedweetOf.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		),
	).Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	// If no such redweet exists, return
	if len(user.Redweets()) == 0 {
		return nil, db.ErrNotFound
	}

	// Remove the Redweet from the post
	_, err = Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Update(
		db.Dweet.RedweetCount.Decrement(1),
	).Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	_, err = Client.Redweet.FindUnique(
		db.Redweet.DbID.Equals(user.Redweets()[0].DbID),
	).Delete().Exec(BaseCtx)
	if err != nil {
		return nil, err
	}

	return &user.Redweets()[0], err
}

// Remove a like from a dweet
func InternalUnlike(postID string, userID string) (*db.DweetModel, error) {
	// Validate params
	err := Validate.Var(postID, "required,alphanum,len=10")
	if err != nil {
		return nil, err
	}

	err = Validate.Var(userID, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return nil, err
	}

	// Get basic version of post to return
	basicPost, err := Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Exec(BaseCtx)
	if err == db.ErrNotFound {
		return nil, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	// Check if user liked the dweet or not
	likedPost, err := Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.LikeUsers.Fetch(
			db.User.Username.Equals(userID),
		),
	).Exec(BaseCtx)
	if err == db.ErrNotFound {
		return nil, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	// If not, then skip unliking the dweet
	if len(likedPost.LikeUsers()) == 0 {
		return basicPost, nil
	}

	// Find the post and decrease its likes by 1
	_, err = Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.LikeUsers.Fetch().OrderBy(
			db.User.FollowerCount.Order(db.DESC),
		),
	).Update(
		db.Dweet.LikeCount.Decrement(1),
		db.Dweet.LikeUsers.Unlink(
			db.User.Username.Equals(userID),
		),
	).Exec(BaseCtx)
	if err == db.ErrNotFound {
		return nil, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	return basicPost, nil
}

// Delete a follower relation
func InternalUnfollow(followedID string, followerID string) (*db.UserModel, error) {
	// Validate params
	err := Validate.Var(followedID, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return nil, err
	}

	err = Validate.Var(followerID, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return nil, err
	}

	// Get basic user to return
	basicUser, err := Client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).Exec(BaseCtx)
	if err == db.ErrNotFound {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	// Check if user doesn't follow this user in the first place
	personBeingUnfollowed, err := Client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch(
			db.User.Username.Equals(followerID),
		),
	).Exec(BaseCtx)
	if err == db.ErrNotFound {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	// If yes, then skip unfollowing the user
	if len(personBeingUnfollowed.Followers()) == 0 {
		return basicUser, nil
	}

	_, err = Client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch().OrderBy(
			db.User.FollowerCount.Order(db.DESC),
		),
	).Update(
		db.User.FollowerCount.Decrement(1),
		db.User.Followers.Unlink(
			db.User.Username.Equals(followerID),
		),
	).Exec(BaseCtx)
	if err == db.ErrNotFound {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return nil, err
	}

	return basicUser, nil
}
