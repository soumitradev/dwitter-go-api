package main

import (
	"dwitter_go_graphql/prisma/db"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

/*

This code is getting a bit messy, and I think I know how to solve it.

Use cases.

So far, I have been looking at whatever I'm doing as an API. I was looking at it from a wrong perspective.

This is an app. This is Dwitter. I don't care which data someone wants in what format.

They can do that using multiple queries, and GraphQL handles a lot of it for them.

My job here is to build Dwitter.

I need to look at it like an app.

When I'm looking at a User's profile, does it matter what posts they've liked?

I'll have to start building my API based on how it'll be used, not on some weird hypothetical 3rd party.

What do I need?

When on the homepage (when logged in), I need:
- A list of latest dweets from the people you follow
- A create dweet button

When viewing a User (when not logged in):
- I need their basic info: Bio, Name, username
- Followers and Following counts
- Some of their Dweets (more can be loaded later on scrolling)

When viewing a User when logged in, I need the same info, except I also need who follows them so I can show mutuals.
Also, a button to follow/unfollow them.

When viewing a Dweet (when not logged in):
- I need the basic dweet info: Body, Author
- Likes, Redweets and reply counts
- Some replies (more can be loaded on scrolling)

When viewing a Dweet when logged in, I need the same info except I also need the users that liked the Dweet
(so I can show which people that the User follows liked the dweet)
Also, a button to like/unlike it.
Also, a button to redweet/unredweet it.
Also, a button to create a reply to it

If the dweet is your own, a button to edit it.


When viewing your own profile when logged in:
- I need their basic info: Bio, Name, username
- Followers and Following counts
- Some of their Dweets (more can be loaded later on scrolling)

Here, we have 4 buttons:
- Liked Dweets (to view dweets you liked)
- Followers (to view people that follow you)
- Following (to view people that you follow)
- Edit Profile button to update the user

Additionally, you can:
- Delete a user
- Delete a dweet

*/

// Get dweet when not authenticated
func NoAuthGetPost(postID string, replies_to_fetch int) (DweetType, error) {
	// When viewing a Dweet (when not logged in):
	// - I need the basic dweet info: Body, Author
	// - Likes, Redweets and reply counts
	// - Some replies (more can be loaded on scrolling)

	var post *db.DweetModel
	var err error
	if replies_to_fetch < 0 {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch(),
			db.Dweet.ReplyTo.Fetch(),
			db.Dweet.RedweetOf.Fetch(),
		).Exec(ctx)
	} else {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch().Take(replies_to_fetch),
			db.Dweet.ReplyTo.Fetch(),
			db.Dweet.RedweetOf.Fetch(),
		).Exec(ctx)
	}

	npost := NoAuthFormatAsDweetType(post)
	return npost, err
}

// Get dweet when authenticated
func AuthGetPost(postID string, replies_to_fetch int, viewUserID string) (DweetType, error) {
	// When viewing a Dweet (when not logged in):
	// - I need the basic dweet info: Body, Author
	// - Likes, Redweets and reply counts
	// - Some replies (more can be loaded on scrolling)

	// Get your own following-list
	viewUser, err := client.User.FindUnique(
		db.User.Username.Equals(viewUserID),
	).With(
		db.User.Following.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	following := viewUser.Following()

	var post *db.DweetModel

	// Fetch the user requested with like_users so we see who liked the dweet
	if replies_to_fetch < 0 {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.RedweetOf.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
		).Exec(ctx)
	} else {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch().Take(replies_to_fetch).With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.RedweetOf.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
		).Exec(ctx)
	}
	if err != nil {
		return DweetType{}, err
	}

	// If the dweet is liked by requesting user, include the requesting user in the like_users list
	likes := post.LikeUsers()
	selfLike := false
	for _, user := range likes {
		if user.Username == viewUserID {
			selfLike = true
		}
	}
	// Find known people that liked thw dweet
	mutuals := HashIntersectUsers(likes, following)

	// Add requesting user to like_users list
	if selfLike {
		mutuals = append(mutuals, *viewUser)
	}

	// Send back the dweet requested, along with like_users
	npost := AuthFormatAsDweetType(post, mutuals)
	return npost, err
}

// Get user when not authenticated
func NoAuthGetUser(userID string, dweets_to_fetch int) (UserType, error) {
	// When viewing a User (when not logged in):
	// - I need their basic info: Bio, Name, username
	// - Followers and Following counts
	// - Some of their Dweets (more can be loaded later on scrolling)

	var user *db.UserModel
	var err error

	if dweets_to_fetch < 0 {
		user, err = client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(ctx)
	} else {
		user, err = client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Dweets.Fetch().Take(dweets_to_fetch).With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(ctx)
	}
	if err != nil {
		return UserType{}, err
	}

	nuser := NoAuthFormatAsUserType(user)
	return nuser, err
}

// Get user when authenticated
func AuthGetUser(userID string, dweets_to_fetch int, viewUserID string) (UserType, error) {
	// When viewing a User when logged in, I need the same info, except I also need who follows them so I can show mutuals.

	var user *db.UserModel
	var mutuals []db.UserModel
	var err error

	if viewUserID == userID {
		// Fetch the user requested
		if dweets_to_fetch < 0 {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Dweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.User.Followers.Fetch(),
				db.User.Following.Fetch(),
			).Exec(ctx)
		} else {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Dweets.Fetch().Take(dweets_to_fetch).With(
					db.Dweet.Author.Fetch(),
				),
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.User.Followers.Fetch(),
				db.User.Following.Fetch(),
			).Exec(ctx)
		}

		if err != nil {
			return UserType{}, err
		}

		nuser := FormatAsUserType(user)
		return nuser, err
	} else {
		// Get your own following-list
		viewUser, err := client.User.FindUnique(
			db.User.Username.Equals(viewUserID),
		).With(
			db.User.Following.Fetch(),
		).Exec(ctx)
		if err != nil {
			return UserType{}, err
		}

		following := viewUser.Following()

		// Fetch the user requested with followers so we get the mutuals
		if dweets_to_fetch < 0 {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Dweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.User.Followers.Fetch(),
			).Exec(ctx)
		} else {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Dweets.Fetch().Take(dweets_to_fetch).With(
					db.Dweet.Author.Fetch(),
				),
				db.User.Followers.Fetch(),
			).Exec(ctx)
		}

		if err != nil {
			return UserType{}, err
		}

		// Get mutuals
		followers := user.Followers()
		mutuals = HashIntersectUsers(followers, following)
		// Send back the user requested, along with mutuals in the followers field
		nuser := AuthFormatAsUserType(user, mutuals)
		return nuser, err
	}
}

// Create a User
func SignUpUser(username string, password string, firstName string, email string) (BasicUserType, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user, err := NewUser(username, string(passwordHash), firstName, "", email, "")
	if err != nil {
		return BasicUserType{}, err
	}

	nuser := FormatAsBasicUserType(user)
	return nuser, err
}

// Check given credentials and return true if valid
func CheckCreds(username string, password string) (bool, error) {
	user, err := client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(ctx)
	if err != nil {
		return false, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false, errors.New("invalid password")
	}
	return true, nil
}

// Update a dweet
func AuthUpdateDweet(postID, userID, body string, mediaLinks []string) (DweetType, error) {
	post, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	if post.Author().Username != userID {
		return DweetType{}, errors.New("not authorized to edit dweet")
	}

	post, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.RedweetDweets.Fetch(),
	).Update(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.LastUpdatedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	redweets := post.RedweetDweets()
	for i := 0; i < len(redweets); i++ {
		_, err := client.Dweet.FindUnique(
			db.Dweet.ID.Equals(redweets[i].ID),
		).Update(
			db.Dweet.DweetBody.Set(body),
		).Exec(ctx)
		if err != nil {
			return DweetType{}, err
		}
	}

	// Return updated post
	post, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	npost := FormatAsDweetType(post)
	return npost, err
}

// Update a user
func AuthUpdateUser(userID, firstName, lastName, email, bio string) (UserType, error) {
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Update(
		db.User.FirstName.Set(firstName),
		db.User.LastName.Set(lastName),
		db.User.Email.Set(email),
		db.User.Bio.Set(bio),
	).Exec(ctx)

	nuser := FormatAsUserType(user)
	return nuser, err
}

// Get User data with dweets that user liked
func FetchLikedDweets(userID string, numberToFetch int, numberOfReplies int) ([]DweetType, error) {
	var user *db.UserModel
	var err error
	if numberToFetch < 0 {
		if numberOfReplies < 0 {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.RedweetOf.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.ReplyDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
				),
			).Exec(ctx)
		} else {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.RedweetOf.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.ReplyDweets.Fetch().Take(numberOfReplies).With(
						db.Dweet.Author.Fetch(),
					),
				),
			).Exec(ctx)
		}
	} else {
		if numberOfReplies < 0 {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.LikedDweets.Fetch().Take(numberToFetch).With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.RedweetOf.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.ReplyDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
				),
			).Exec(ctx)
		} else {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.LikedDweets.Fetch().Take(numberToFetch).With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.RedweetOf.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.ReplyDweets.Fetch().Take(numberOfReplies).With(
						db.Dweet.Author.Fetch(),
					),
				),
			).Exec(ctx)
		}
	}

	var liked []DweetType
	for _, dweet := range user.LikedDweets() {
		liked = append(liked, FormatAsDweetType(&dweet))
	}
	return liked, err
}

// Get users that follow user
func FetchFollowers(userID string, numberToFetch int) ([]UserType, error) {
	var user *db.UserModel
	var err error
	if numberToFetch < 0 {
		user, err = client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Followers.Fetch(),
		).Exec(ctx)
	} else {
		user, err = client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Followers.Fetch().Take(numberToFetch),
		).Exec(ctx)
	}

	var followers []UserType
	for _, follower := range user.Followers() {
		followers = append(followers, FormatAsUserType(&follower))
	}
	return followers, err
}

// Get users that user follows
func FetchFollowing(userID string, numberToFetch int) ([]UserType, error) {
	var user *db.UserModel
	var err error
	if numberToFetch < 0 {
		user, err = client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Following.Fetch(),
		).Exec(ctx)
	} else {
		user, err = client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Following.Fetch().Take(numberToFetch),
		).Exec(ctx)
	}

	var following []UserType
	for _, followed := range user.Following() {
		following = append(following, FormatAsUserType(&followed))
	}
	return following, err
}

// Delete a dweet
func AuthDeleteDweet(postID string, userID string) (DweetType, error) {
	dweet, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	if dweet.Author().Username == userID {
		dweet, err := DeleteDweet(postID)
		if err != nil {
			return DweetType{}, err
		}

		formatted := FormatAsDweetType(dweet)
		return formatted, err
	}

	return DweetType{}, errors.New("Unauthorized")

}

// Delete a dweet
func AuthDeleteRedweet(postID string, userID string) (DweetType, error) {
	dweet, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	if dweet.Author().Username == userID {
		dweet, err := DeleteRedweet(postID)
		if err != nil {
			return DweetType{}, err
		}

		formatted := FormatAsDweetType(dweet)
		return formatted, err
	}

	return DweetType{}, errors.New("Unauthorized")
}
