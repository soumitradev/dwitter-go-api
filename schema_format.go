package main

import (
	"dwitter_go_graphql/prisma/db"
)

func FormatAsBasicDweetType(dweet *db.DweetModel) BasicDweetType {
	reply_id, present := dweet.OriginalReplyID()
	if !present {
		reply_id = ""
	}
	redweet_id, present := dweet.OriginalRedweetID()
	if !present {
		redweet_id = ""
	}
	return BasicDweetType{
		DweetBody:         dweet.DweetBody,
		ID:                dweet.ID,
		Author:            FormatAsBasicUserType(dweet.Author()),
		AuthorID:          dweet.AuthorID,
		PostedAt:          dweet.PostedAt,
		LastUpdatedAt:     dweet.LastUpdatedAt,
		LikeCount:         dweet.LikeCount,
		IsReply:           dweet.IsReply,
		OriginalReplyID:   reply_id,
		ReplyCount:        dweet.ReplyCount,
		IsRedweet:         dweet.IsRedweet,
		OriginalRedweetID: redweet_id,
		RedweetCount:      dweet.RedweetCount,
		Media:             dweet.Media,
	}
}

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

	redweet_id, present := dweet.OriginalRedweetID()
	if !present {
		redweet_id = ""
	}
	original_redweet_dweet, present := dweet.RedweetOf()
	var redweet_of BasicDweetType
	if present {
		reply_to = FormatAsBasicDweetType(original_redweet_dweet)
	} else {
		redweet_of = BasicDweetType{}
	}

	var reply_dweets []BasicDweetType
	reply_dweets_db_schema := dweet.ReplyDweets()
	for i := 0; i < len(like_users_db_schema); i++ {
		reply_dweets = append(reply_dweets, FormatAsBasicDweetType(&reply_dweets_db_schema[i]))
	}

	var redweet_dweets []BasicDweetType
	redweet_dweets_db_schema := dweet.ReplyDweets()
	for i := 0; i < len(redweet_dweets_db_schema); i++ {
		reply_dweets = append(reply_dweets, FormatAsBasicDweetType(&redweet_dweets_db_schema[i]))
	}

	return DweetType{
		DweetBody:         dweet.DweetBody,
		ID:                dweet.ID,
		Author:            author,
		AuthorID:          dweet.AuthorID,
		PostedAt:          dweet.PostedAt,
		LastUpdatedAt:     dweet.LastUpdatedAt,
		LikeCount:         dweet.LikeCount,
		LikeUsers:         like_users,
		IsReply:           dweet.IsReply,
		OriginalReplyID:   reply_id,
		ReplyTo:           reply_to,
		ReplyCount:        dweet.ReplyCount,
		ReplyDweets:       reply_dweets,
		IsRedweet:         dweet.IsRedweet,
		OriginalRedweetID: redweet_id,
		RedweetOf:         redweet_of,
		RedweetCount:      dweet.RedweetCount,
		RedweetDweets:     redweet_dweets,
		Media:             dweet.Media,
	}
}

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
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
		LastName:       lastName,
	}
}

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
		Dweets:         dweets,
		LikedDweets:    liked_dweets,
		FollowerCount:  user.FollowerCount,
		Followers:      followers,
		FollowingCount: user.FollowingCount,
		Following:      following,
		CreatedAt:      user.CreatedAt,
	}
}

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
		Dweets:         dweets,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
	}
}

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

	redweet_id, present := dweet.OriginalRedweetID()
	if !present {
		redweet_id = ""
	}
	original_redweet_dweet, present := dweet.RedweetOf()
	var redweet_of BasicDweetType
	if present {
		redweet_of = FormatAsBasicDweetType(original_redweet_dweet)
	} else {
		redweet_of = BasicDweetType{}
	}

	var reply_dweets []BasicDweetType
	reply_dweets_db_schema := dweet.ReplyDweets()
	for i := 0; i < len(reply_dweets_db_schema); i++ {
		reply_dweets = append(reply_dweets, FormatAsBasicDweetType(&reply_dweets_db_schema[i]))
	}

	return DweetType{
		DweetBody:         dweet.DweetBody,
		ID:                dweet.ID,
		Author:            author,
		AuthorID:          dweet.AuthorID,
		PostedAt:          dweet.PostedAt,
		LastUpdatedAt:     dweet.LastUpdatedAt,
		LikeCount:         dweet.LikeCount,
		LikeUsers:         []BasicUserType{},
		IsReply:           dweet.IsReply,
		OriginalReplyID:   reply_id,
		ReplyTo:           reply_to,
		ReplyCount:        dweet.ReplyCount,
		ReplyDweets:       reply_dweets,
		IsRedweet:         dweet.IsRedweet,
		OriginalRedweetID: redweet_id,
		RedweetOf:         redweet_of,
		RedweetCount:      dweet.RedweetCount,
		RedweetDweets:     []BasicDweetType{},
		Media:             dweet.Media,
	}
}
