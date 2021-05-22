// Package schema provides useful custom types and functions to format database objects into these types
package schema

import (
	"dwitter_go_graphql/prisma/db"
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
func FormatAsDweetType(dweet *db.DweetModel) DweetType {
	author := FormatAsBasicUserType(dweet.Author())

	var like_users []BasicUserType
	like_users_db_schema := dweet.LikeUsers()
	for i := 0; i < len(like_users_db_schema); i++ {
		like_users = append(like_users, FormatAsBasicUserType(&like_users_db_schema[i]))
	}

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
	for i := 0; i < len(like_users_db_schema); i++ {
		reply_dweets = append(reply_dweets, FormatAsBasicDweetType(&reply_dweets_db_schema[i]))
	}

	return DweetType{
		DweetBody:       dweet.DweetBody,
		ID:              dweet.ID,
		Author:          author,
		AuthorID:        dweet.AuthorID,
		PostedAt:        dweet.PostedAt,
		LastUpdatedAt:   dweet.LastUpdatedAt,
		LikeCount:       dweet.LikeCount,
		LikeUsers:       like_users,
		IsReply:         dweet.IsReply,
		OriginalReplyID: reply_id,
		ReplyTo:         reply_to,
		ReplyCount:      dweet.ReplyCount,
		ReplyDweets:     reply_dweets,
		RedweetCount:    dweet.RedweetCount,
		Media:           dweet.Media,
	}
}

// Format as BasicUser
func FormatAsBasicUserType(user *db.UserModel) BasicUserType {
	lastName, exists := user.LastName()
	if !exists {
		lastName = ""
	}
	return BasicUserType{
		Username:       user.Username,
		FirstName:      user.FirstName,
		Email:          user.Email,
		Bio:            user.Bio,
		PfpURL:         user.ProfilePicURL,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
		LastName:       lastName,
	}
}

// Format as User
func FormatAsUserType(user *db.UserModel) UserType {
	var dweets []BasicDweetType
	dweets_db_schema := user.Dweets()
	for i := 0; i < len(dweets_db_schema); i++ {
		dweets = append(dweets, FormatAsBasicDweetType(&dweets_db_schema[i]))
	}

	var liked_dweets []BasicDweetType
	liked_dweets_db_schema := user.LikedDweets()
	for i := 0; i < len(liked_dweets_db_schema); i++ {
		liked_dweets = append(liked_dweets, FormatAsBasicDweetType(&liked_dweets_db_schema[i]))
	}

	var followers []BasicUserType
	followers_db_schema := user.Followers()
	for i := 0; i < len(followers_db_schema); i++ {
		followers = append(followers, FormatAsBasicUserType(&followers_db_schema[i]))
	}

	var following []BasicUserType
	following_db_schema := user.Following()
	for i := 0; i < len(following_db_schema); i++ {
		following = append(following, FormatAsBasicUserType(&following_db_schema[i]))
	}
	lastName, exists := user.LastName()
	if !exists {
		lastName = ""
	}

	return UserType{
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       lastName,
		Email:          user.Email,
		Bio:            user.Bio,
		PfpURL:         user.ProfilePicURL,
		Dweets:         dweets,
		LikedDweets:    liked_dweets,
		FollowerCount:  user.FollowerCount,
		Followers:      followers,
		FollowingCount: user.FollowingCount,
		Following:      following,
		CreatedAt:      user.CreatedAt,
	}
}

// Format as User
func NoAuthFormatAsUserType(user *db.UserModel) UserType {
	var dweets []BasicDweetType
	dweets_db_schema := user.Dweets()
	for i := 0; i < len(dweets_db_schema); i++ {
		dweets = append(dweets, FormatAsBasicDweetType(&dweets_db_schema[i]))
	}

	lastName, exists := user.LastName()
	if !exists {
		lastName = ""
	}

	return UserType{
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       lastName,
		Email:          user.Email,
		Bio:            user.Bio,
		PfpURL:         user.ProfilePicURL,
		Dweets:         dweets,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
	}
}

// Format as User with followers
func AuthFormatAsUserType(user *db.UserModel, mutualUsers []db.UserModel) UserType {
	var dweets []BasicDweetType
	dweets_db_schema := user.Dweets()
	for i := 0; i < len(dweets_db_schema); i++ {
		dweets = append(dweets, FormatAsBasicDweetType(&dweets_db_schema[i]))
	}

	var mutuals []BasicUserType
	for i := 0; i < len(mutualUsers); i++ {
		mutuals = append(mutuals, FormatAsBasicUserType((&mutualUsers[i])))
	}

	lastName, exists := user.LastName()
	if !exists {
		lastName = ""
	}

	return UserType{
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       lastName,
		Email:          user.Email,
		Bio:            user.Bio,
		PfpURL:         user.ProfilePicURL,
		Dweets:         dweets,
		FollowerCount:  user.FollowerCount,
		Followers:      mutuals,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
	}
}

// Format as Dweet
func NoAuthFormatAsDweetType(dweet *db.DweetModel) DweetType {
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

	return DweetType{
		DweetBody:       dweet.DweetBody,
		ID:              dweet.ID,
		Author:          author,
		AuthorID:        dweet.AuthorID,
		PostedAt:        dweet.PostedAt,
		LastUpdatedAt:   dweet.LastUpdatedAt,
		LikeCount:       dweet.LikeCount,
		LikeUsers:       []BasicUserType{},
		IsReply:         dweet.IsReply,
		OriginalReplyID: reply_id,
		ReplyTo:         reply_to,
		ReplyCount:      dweet.ReplyCount,
		ReplyDweets:     reply_dweets,
		RedweetCount:    dweet.RedweetCount,
		Media:           dweet.Media,
	}
}

// Format as Dweet with users that liked it
func AuthFormatAsDweetType(dweet *db.DweetModel, likeUsers []db.UserModel) DweetType {
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
		Media:           dweet.Media,
	}
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
