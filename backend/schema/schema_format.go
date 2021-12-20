// Package schema provides useful custom types and functions to format database objects into these types
package schema

import (
	"errors"

	"github.com/soumitradev/Dwitter/backend/prisma/db"
)

// Format as BasicDweet
func FormatAsBasicDweetType(dweet *db.DweetModel) BasicDweetType {
	reply_id, present := dweet.OriginalReplyID()
	if !present {
		reply_id = ""
	}
	return BasicDweetType{
		DweetBody:       dweet.DweetBody,
		ID:              dweet.ID,
		Author:          FormatAsBasicUserType(dweet.Author()),
		AuthorID:        dweet.AuthorID,
		PostedAt:        dweet.PostedAt,
		LastUpdatedAt:   dweet.LastUpdatedAt,
		LikeCount:       dweet.LikeCount,
		IsReply:         dweet.IsReply,
		OriginalReplyID: reply_id,
		ReplyCount:      dweet.ReplyCount,
		RedweetCount:    dweet.RedweetCount,
		Media:           dweet.Media,
	}
}

// Format as Dweet
func FormatAsDweetType(dweet *db.DweetModel, likeUsers []db.UserModel, redweetUsers []db.UserModel) DweetType {
	author := FormatAsBasicUserType(dweet.Author())

	reply_id, present := dweet.OriginalReplyID()
	if !present {
		reply_id = ""
	}
	original_reply_dweet, present := dweet.ReplyTo()
	var reply_to BasicDweetType
	if present {
		reply_to = FormatAsBasicDweetType(original_reply_dweet)
	} else {
		reply_to = BasicDweetType{}
	}

	var reply_dweets []BasicDweetType
	reply_dweets_db_schema := dweet.ReplyDweets()
	for i := 0; i < len(reply_dweets_db_schema); i++ {
		reply_dweets = append(reply_dweets, FormatAsBasicDweetType(&reply_dweets_db_schema[i]))
	}

	var likes []BasicUserType
	for i := 0; i < len(likeUsers); i++ {
		likes = append(likes, FormatAsBasicUserType((&likeUsers[i])))
	}

	var redweet_users []BasicUserType
	for i := 0; i < len(redweetUsers); i++ {
		redweet_users = append(redweet_users, FormatAsBasicUserType((&redweetUsers[i])))
	}

	return DweetType{
		DweetBody:       dweet.DweetBody,
		ID:              dweet.ID,
		Author:          author,
		AuthorID:        dweet.AuthorID,
		PostedAt:        dweet.PostedAt,
		LastUpdatedAt:   dweet.LastUpdatedAt,
		LikeCount:       dweet.LikeCount,
		LikeUsers:       likes,
		IsReply:         dweet.IsReply,
		OriginalReplyID: reply_id,
		ReplyTo:         reply_to,
		ReplyCount:      dweet.ReplyCount,
		ReplyDweets:     reply_dweets,
		RedweetCount:    dweet.RedweetCount,
		RedweetUsers:    redweet_users,
		Media:           dweet.Media,
	}
}

// Format as BasicUser
func FormatAsBasicUserType(user *db.UserModel) BasicUserType {
	return BasicUserType{
		Username:       user.Username,
		Name:           user.Name,
		Email:          user.Email,
		Bio:            user.Bio,
		PfpURL:         user.ProfilePicURL,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
	}
}

// Format as User
func FormatAsUserType(user *db.UserModel, alsoFollowedBy []db.UserModel, alsoFollowing []db.UserModel, objectsToFetch string, objectList []interface{}, showEmail bool) (UserType, error) {
	var email string
	var feedObjects []interface{}
	var dweets []BasicDweetType
	var redweets []RedweetType
	var liked_dweets []BasicDweetType
	var redweeted_dweets []BasicDweetType

	switch objectsToFetch {
	case "feed":
		feedObjects = make([]interface{}, len(objectList))
		for index, obj := range objectList {
			if dweet, ok := obj.(db.DweetModel); ok {
				feedObjects[index] = FormatAsBasicDweetType(&dweet)
			} else if redweet, ok := obj.(db.RedweetModel); ok {
				feedObjects[index] = FormatAsRedweetType(&redweet)
			} else {
				return UserType{}, errors.New("internal server error")
			}
		}
	case "dweet":
		dweets = make([]BasicDweetType, len(objectList))
		for index, obj := range objectList {
			if dweet, ok := obj.(db.DweetModel); ok {
				dweets[index] = FormatAsBasicDweetType(&dweet)
			} else {
				return UserType{}, errors.New("internal server error")
			}
		}
	case "redweet":
		redweets = make([]RedweetType, len(objectList))
		for index, obj := range objectList {
			if redweet, ok := obj.(db.RedweetModel); ok {
				redweets[index] = FormatAsRedweetType(&redweet)
			} else {
				return UserType{}, errors.New("internal server error")
			}
		}
	case "redweetedDweet":
		redweeted_dweets = make([]BasicDweetType, len(objectList))
		for index, obj := range objectList {
			if dweet, ok := obj.(db.DweetModel); ok {
				redweeted_dweets[index] = FormatAsBasicDweetType(&dweet)
			} else {
				return UserType{}, errors.New("internal server error")
			}
		}
	case "liked":
		liked_dweets = make([]BasicDweetType, len(objectList))
		for index, obj := range objectList {
			if dweet, ok := obj.(db.DweetModel); ok {
				liked_dweets[index] = FormatAsBasicDweetType(&dweet)
			} else {
				return UserType{}, errors.New("internal server error")
			}
		}
	default:
		break
	}

	var followers []BasicUserType
	for i := 0; i < len(alsoFollowedBy); i++ {
		followers = append(followers, FormatAsBasicUserType(&alsoFollowedBy[i]))
	}

	var following []BasicUserType
	for i := 0; i < len(alsoFollowing); i++ {
		following = append(following, FormatAsBasicUserType(&alsoFollowing[i]))
	}

	if showEmail {
		email = user.Email
	} else {
		email = ""
	}

	return UserType{
		Username:        user.Username,
		Name:            user.Name,
		Email:           email,
		Bio:             user.Bio,
		PfpURL:          user.ProfilePicURL,
		Dweets:          dweets,
		Redweets:        redweets,
		RedweetedDweets: redweeted_dweets,
		FeedObjects:     feedObjects,
		LikedDweets:     liked_dweets,
		FollowerCount:   user.FollowerCount,
		Followers:       followers,
		FollowingCount:  user.FollowingCount,
		Following:       following,
		CreatedAt:       user.CreatedAt,
	}, nil
}

// Format as Redweet
func FormatAsRedweetType(redweet *db.RedweetModel) RedweetType {
	return RedweetType{
		Author:            FormatAsBasicUserType(redweet.Author()),
		AuthorID:          redweet.AuthorID,
		RedweetOf:         FormatAsBasicDweetType(redweet.RedweetOf()),
		OriginalRedweetID: redweet.OriginalRedweetID,
		RedweetTime:       redweet.RedweetTime,
	}
}
