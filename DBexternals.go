package main

import "dwitter_go_graphql/prisma/db"

func GetFullPost(postID string, replies_to_fetch int) (*db.DweetModel, error) {
	// Get full post info
	var post *db.DweetModel
	var err error
	if replies_to_fetch < 0 {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.LikeUsers.Fetch(),
			db.Dweet.ReplyDweets.Fetch(),
			db.Dweet.ReplyTo.Fetch(),
			db.Dweet.RedweetOf.Fetch(),
		).Exec(ctx)
	} else {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.LikeUsers.Fetch(),
			db.Dweet.ReplyDweets.Fetch().Take(replies_to_fetch),
			db.Dweet.ReplyTo.Fetch(),
			db.Dweet.RedweetOf.Fetch(),
		).Exec(ctx)
	}
	return post, err
}

func GetFullUser(userID string, dweets_to_fetch int, liked_dweets_to_fetch int) (*db.UserModel, error) {
	// Get full user info
	var user *db.UserModel
	var err error

	if (dweets_to_fetch < 0) && (liked_dweets_to_fetch < 0) {
		user, err = client.User.FindUnique(
			db.User.Mention.Equals(userID),
		).With(
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.User.LikedDweets.Fetch(),
			db.User.Followers.Fetch(),
			db.User.Following.Fetch(),
		).Exec(ctx)
	} else if (dweets_to_fetch < 0) && (liked_dweets_to_fetch > 0) {
		user, err = client.User.FindUnique(
			db.User.Mention.Equals(userID),
		).With(
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.User.LikedDweets.Fetch().Take(liked_dweets_to_fetch),
			db.User.Followers.Fetch(),
			db.User.Following.Fetch(),
		).Exec(ctx)
	} else if (dweets_to_fetch > 0) && (liked_dweets_to_fetch < 0) {
		user, err = client.User.FindUnique(
			db.User.Mention.Equals(userID),
		).With(
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(dweets_to_fetch),
			db.User.LikedDweets.Fetch(),
			db.User.Followers.Fetch(),
			db.User.Following.Fetch(),
		).Exec(ctx)
	} else {
		user, err = client.User.FindUnique(
			db.User.Mention.Equals(userID),
		).With(
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(dweets_to_fetch),
			db.User.LikedDweets.Fetch().Take(liked_dweets_to_fetch),
			db.User.Followers.Fetch(),
			db.User.Following.Fetch(),
		).Exec(ctx)
	}
	return user, err
}
