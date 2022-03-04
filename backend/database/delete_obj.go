package database

import (
	"errors"
	"fmt"

	"github.com/soumitradev/Dwitter/backend/cdn"
	"github.com/soumitradev/Dwitter/backend/common"
	"github.com/soumitradev/Dwitter/backend/prisma/db"
	"github.com/soumitradev/Dwitter/backend/schema"
	"github.com/soumitradev/Dwitter/backend/util"
)

// Delete a dweet
func DeleteDweet(postID string, username string, repliesToFetch int, replyOffset int) (schema.DweetType, error) {
	// Validate params
	err := common.Validate.Var(postID, "required,alphanum,len=10")
	if err != nil {
		return schema.DweetType{}, err
	}

	err = common.Validate.Var(username, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return schema.DweetType{}, err
	}

	err = common.Validate.Var(replyOffset, "gte=0")
	if err != nil {
		return schema.DweetType{}, err
	}

	var deleted *db.DweetModel

	// Check params and return data accordingly
	if repliesToFetch < 0 {
		deleted, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch().With(
				db.User.Following.Fetch().OrderBy(
					db.User.FollowerCount.Order(db.DESC),
				),
			),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch().OrderBy(
				db.User.FollowerCount.Order(db.DESC),
			),
			db.Dweet.RedweetUsers.Fetch().OrderBy(
				db.User.FollowerCount.Order(db.DESC),
			),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).OrderBy(
				db.Dweet.LikeCount.Order(db.DESC),
			),
		).Exec(common.BaseCtx)
	} else {
		deleted, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch().With(
				db.User.Following.Fetch().OrderBy(
					db.User.FollowerCount.Order(db.DESC),
				),
			),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch().OrderBy(
				db.User.FollowerCount.Order(db.DESC),
			),
			db.Dweet.RedweetUsers.Fetch().OrderBy(
				db.User.FollowerCount.Order(db.DESC),
			),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).OrderBy(
				db.Dweet.LikeCount.Order(db.DESC),
			).Take(repliesToFetch).Skip(replyOffset),
		).Exec(common.BaseCtx)
	}
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Check if authorized to delete dweet
	if deleted.Author().Username == username {
		_, err := deleteDweet(postID)

		// Delete the media that isn't used anymore
		oldMedia := deleted.Media
		for _, mediaLink := range oldMedia {
			loc, err := cdn.LinkToLocation(mediaLink)
			if err != nil {
				return schema.DweetType{}, err
			}
			err = cdn.DeleteLocation(loc, true)
			if err != nil {
				return schema.DweetType{}, err
			}
		}

		if err != nil {
			return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
		}

		// Format and return with common likes
		knownUsers := deleted.Author().Following()
		knownUsers = append(knownUsers, *deleted.Author())

		mutualLikes := util.HashIntersectUsers(deleted.LikeUsers(), knownUsers)
		mutualRedweets := util.HashIntersectUsers(deleted.RedweetUsers(), knownUsers)

		formatted := schema.FormatAsDweetType(deleted, mutualLikes, mutualRedweets)
		return formatted, err
	}

	return schema.DweetType{}, fmt.Errorf("internal server error: %v", errors.New("Unauthorized"))
}

// Delete a redweet
func DeleteRedweet(postID string, username string) (schema.RedweetType, error) {
	// Validate params
	err := common.Validate.Var(postID, "required,alphanum,len=10")
	if err != nil {
		return schema.RedweetType{}, err
	}

	err = common.Validate.Var(username, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return schema.RedweetType{}, err
	}

	redweet, err := deleteRedweet(postID, username)
	if err == db.ErrNotFound {
		return schema.RedweetType{}, fmt.Errorf("redweet not found: %v", err)
	}
	if err != nil {
		return schema.RedweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	formatted := schema.FormatAsRedweetType(redweet)
	return formatted, err
}

// TODO: DeleteUser
