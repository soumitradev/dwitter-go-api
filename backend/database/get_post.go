package database

import (
	"fmt"

	"github.com/soumitradev/Dwitter/backend/common"
	"github.com/soumitradev/Dwitter/backend/prisma/db"
	"github.com/soumitradev/Dwitter/backend/schema"
	"github.com/soumitradev/Dwitter/backend/util"
)

// Get dweet when not authenticated
func GetPostUnauth(postID string, repliesToFetch int, replyOffset int) (schema.DweetType, error) {
	// Validate params
	err := common.Validate.Var(postID, "required,alphanum,eq=10")
	if err != nil {
		return schema.DweetType{}, err
	}

	err = common.Validate.Var(replyOffset, "gte=0")
	if err != nil {
		return schema.DweetType{}, err
	}

	// Check params and return data accordingly
	var post *db.DweetModel
	if repliesToFetch < 0 {
		post, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(common.BaseCtx)
	} else {
		post, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(repliesToFetch).Skip(replyOffset),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(common.BaseCtx)
	}
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	npost := schema.NoAuthFormatAsDweetType(post)
	return npost, err
}

// Get dweet when authenticated
func GetPost(postID string, repliesToFetch int, replyOffset int, viewerUsername string) (schema.DweetType, error) {
	// Validate params
	err := common.Validate.Var(postID, "required,alphanum,eq=10")
	if err != nil {
		return schema.DweetType{}, err
	}

	err = common.Validate.Var(replyOffset, "gte=0")
	if err != nil {
		return schema.DweetType{}, err
	}

	err = common.Validate.Var(viewerUsername, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return schema.DweetType{}, err
	}

	// Get your own following-list
	viewUser, err := common.Client.User.FindUnique(
		db.User.Username.Equals(viewerUsername),
	).With(
		db.User.Following.Fetch(),
	).Exec(common.BaseCtx)
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	following := viewUser.Following()

	var post *db.DweetModel

	// Check params and return data accordingly
	if repliesToFetch < 0 {
		post, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
			db.Dweet.RedweetUsers.Fetch(),
		).Exec(common.BaseCtx)
	} else {
		post, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(repliesToFetch).Skip(replyOffset),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
			db.Dweet.RedweetUsers.Fetch(),
		).Exec(common.BaseCtx)
	}
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// If the dweet is liked by requesting user, include the requesting user in the like_users list
	likes := post.LikeUsers()
	selfLike := false
	for _, user := range likes {
		if user.Username == viewerUsername {
			selfLike = true
		}
	}
	// Find known people that liked the dweet
	mutualLikes := util.HashIntersectUsers(likes, following)

	// Add requesting user to like_users list
	if selfLike {
		mutualLikes = append(mutualLikes, *viewUser)
	}

	// If the dweet is redweeted by requesting user, include the requesting user in the redweet_users list
	redweetUsers := post.RedweetUsers()
	selfRedweet := false
	for _, user := range redweetUsers {
		if user.Username == viewerUsername {
			selfRedweet = true
		}
	}
	// Find known people that redweeted the dweet
	mutualRedweets := util.HashIntersectUsers(redweetUsers, following)

	// Add requesting user to redweet_users list
	if selfRedweet {
		mutualRedweets = append(mutualRedweets, *viewUser)
	}

	// Send back the dweet requested, along with like_users
	npost := schema.AuthFormatAsDweetType(post, mutualLikes, mutualRedweets)
	return npost, err
}
