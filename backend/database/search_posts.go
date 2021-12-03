package database

import (
	"fmt"

	"github.com/soumitradev/Dwitter/backend/common"
	"github.com/soumitradev/Dwitter/backend/prisma/db"
	"github.com/soumitradev/Dwitter/backend/schema"
	"github.com/soumitradev/Dwitter/backend/util"
)

// Search dweets when not authenticated
func SearchPostsUnauth(query string, numberToFetch int, numOffset int, repliesToFetch int, replyOffset int) ([]schema.DweetType, error) {
	// Validate params
	err := common.Validate.Var(query, "required,gt=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	err = common.Validate.Var(numOffset, "gte=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	err = common.Validate.Var(replyOffset, "gte=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	var posts []db.DweetModel

	// Check params and return data accordingly
	if numberToFetch < 0 {
		if repliesToFetch < 0 {
			posts, err = common.Client.Dweet.FindMany(
				db.Dweet.DweetBody.Contains(query),
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
			posts, err = common.Client.Dweet.FindMany(
				db.Dweet.DweetBody.Contains(query),
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
	} else {
		if repliesToFetch < 0 {
			posts, err = common.Client.Dweet.FindMany(
				db.Dweet.DweetBody.Contains(query),
			).Take(numberToFetch).Skip(numOffset).With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
			).Exec(common.BaseCtx)
		} else {
			posts, err = common.Client.Dweet.FindMany(
				db.Dweet.DweetBody.Contains(query),
			).Take(numberToFetch).Skip(numOffset).With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				).Take(repliesToFetch).Skip(replyOffset),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
			).Exec(common.BaseCtx)
		}
	}

	if err == db.ErrNotFound {
		return []schema.DweetType{}, fmt.Errorf("dweets not found: %v", err)
	}
	if err != nil {
		return []schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Format
	var formatted []schema.DweetType
	for _, post := range posts {
		npost := schema.NoAuthFormatAsDweetType(&post)
		formatted = append(formatted, npost)
	}
	return formatted, err
}

// Search dweets when authenticated
func SearchPosts(query string, numberToFetch int, numOffset int, repliesToFetch int, replyOffset int, viewerUsername string) ([]schema.DweetType, error) {
	// Validate params
	err := common.Validate.Var(query, "required,gt=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	err = common.Validate.Var(numOffset, "gte=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	err = common.Validate.Var(replyOffset, "gte=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	err = common.Validate.Var(viewerUsername, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	// Get your own following-list
	viewUser, err := common.Client.User.FindUnique(
		db.User.Username.Equals(viewerUsername),
	).With(
		db.User.Following.Fetch(),
	).Exec(common.BaseCtx)
	if err != nil {
		return []schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	following := viewUser.Following()

	var posts []db.DweetModel

	// Check params and return data accordingly
	if numberToFetch < 0 {
		if repliesToFetch < 0 {
			posts, err = common.Client.Dweet.FindMany(
				db.Dweet.DweetBody.Contains(query),
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
			posts, err = common.Client.Dweet.FindMany(
				db.Dweet.DweetBody.Contains(query),
			).With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				).Take(repliesToFetch).Skip(repliesToFetch),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.LikeUsers.Fetch(),
				db.Dweet.RedweetUsers.Fetch(),
			).Exec(common.BaseCtx)
		}
	} else {
		if repliesToFetch < 0 {
			posts, err = common.Client.Dweet.FindMany(
				db.Dweet.DweetBody.Contains(query),
			).Take(numberToFetch).Skip(numOffset).With(
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
			posts, err = common.Client.Dweet.FindMany(
				db.Dweet.DweetBody.Contains(query),
			).Take(numberToFetch).Skip(numOffset).With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				).Take(repliesToFetch).Skip(repliesToFetch),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.LikeUsers.Fetch(),
				db.Dweet.RedweetUsers.Fetch(),
			).Exec(common.BaseCtx)
		}
	}
	if err == db.ErrNotFound {
		return []schema.DweetType{}, fmt.Errorf("dweets not found: %v", err)
	}
	if err != nil {
		return []schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	var formatted []schema.DweetType

	for _, post := range posts {
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

		// If the dweet is liked by requesting user, include the requesting user in the like_users list
		redweets := post.RedweetUsers()
		selfRedweet := false
		for _, user := range redweets {
			if user.Username == viewerUsername {
				selfRedweet = true
			}
		}
		// Find known people that liked the dweet
		mutualRedweets := util.HashIntersectUsers(redweets, following)

		// Add requesting user to like_users list
		if selfRedweet {
			mutualRedweets = append(mutualRedweets, *viewUser)
		}

		// Send back the dweet requested, along with like_users
		npost := schema.AuthFormatAsDweetType(&post, mutualLikes, mutualRedweets)
		formatted = append(formatted, npost)
	}

	return formatted, err
}
