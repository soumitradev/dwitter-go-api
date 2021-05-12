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
	if err != nil {
		return DweetType{}, err
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
func SignUpUser(username string, password string, firstName string, lastName string, bio string, email string) (UserType, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserType{}, err
	}

	createdUser, err := client.User.CreateOne(
		db.User.Username.Set(username),
		db.User.PasswordHash.Set(string(passwordHash)),
		db.User.FirstName.Set(firstName),
		db.User.Email.Set(email),
		db.User.Bio.Set(bio),
		db.User.CreatedAt.Set(time.Now()),
		db.User.LastName.Set(lastName),
	).With(
		db.User.Dweets.Fetch(),
	).Exec(ctx)

	if err != nil {
		return UserType{}, err
	}

	nuser := AuthFormatAsUserType(createdUser, []db.UserModel{})
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
		db.Dweet.ReplyTo.Fetch(),
		db.Dweet.ReplyDweets.Fetch(),
		db.Dweet.RedweetOf.Fetch(),
		db.Dweet.LikeUsers.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	user, err := client.User.FindUnique(
		db.User.Username.Equals(postID),
	).With(
		db.User.Following.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	mutuals := HashIntersectUsers(user.Following(), post.LikeUsers())
	npost := AuthFormatAsDweetType(post, mutuals)
	return npost, err
}

// Update a user
func AuthUpdateUser(userID, firstName, lastName, email, bio string) (UserType, error) {
	user, err := client.User.FindUnique(
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
	).Update(
		db.User.FirstName.Set(firstName),
		db.User.LastName.Set(lastName),
		db.User.Email.Set(email),
		db.User.Bio.Set(bio),
	).Exec(ctx)
	if err != nil {
		return UserType{}, err
	}

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
				db.User.Following.Fetch(),
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
				db.User.Following.Fetch(),
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
				db.User.Following.Fetch(),
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
				db.User.Following.Fetch(),
			).Exec(ctx)
		}
	}
	if err != nil {
		return []DweetType{}, err
	}

	var liked []DweetType
	for _, dweet := range user.LikedDweets() {
		likes := dweet.LikeUsers()

		// Find known people that liked thw dweet
		mutuals := HashIntersectUsers(likes, user.Following())

		// Add requesting user to like_users list
		mutuals = append(mutuals, *user)

		liked = append(liked, AuthFormatAsDweetType(&dweet, mutuals))
	}
	return liked, err
}

// Get users that follow user
func FetchFollowers(userID string, numberToFetch int, dweetsToFetch int) ([]UserType, error) {
	var user *db.UserModel
	var err error
	if numberToFetch < 0 {
		if dweetsToFetch < 0 {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Followers.Fetch().With(
					db.User.Dweets.Fetch(),
					db.User.Followers.Fetch(),
				),
				db.User.Following.Fetch(),
			).Exec(ctx)
		} else {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Followers.Fetch().With(
					db.User.Dweets.Fetch().Take(dweetsToFetch),
					db.User.Followers.Fetch(),
				),
				db.User.Following.Fetch(),
			).Exec(ctx)
		}
	} else {
		if dweetsToFetch < 0 {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Followers.Fetch().Take(numberToFetch).With(
					db.User.Dweets.Fetch(),
					db.User.Followers.Fetch(),
				),
				db.User.Following.Fetch(),
			).Exec(ctx)
		} else {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Followers.Fetch().Take(numberToFetch).With(
					db.User.Dweets.Fetch().Take(dweetsToFetch),
					db.User.Followers.Fetch(),
				),
				db.User.Following.Fetch(),
			).Exec(ctx)
		}
	}
	if err != nil {
		return []UserType{}, err
	}

	var followers []UserType
	for _, follower := range user.Followers() {
		followerFollowers := follower.Followers()

		mutuals := HashIntersectUsers(followerFollowers, user.Following())

		mutuals = append(mutuals, *user)
		followers = append(followers, AuthFormatAsUserType(&follower, mutuals))
	}
	return followers, err
}

// Get users that user follows
func FetchFollowing(userID string, numberToFetch int, dweetsToFetch int) ([]UserType, error) {
	var user *db.UserModel
	var err error
	if numberToFetch < 0 {
		if dweetsToFetch < 0 {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Following.Fetch().With(
					db.User.Dweets.Fetch(),
					db.User.Followers.Fetch(),
				),
			).Exec(ctx)
		} else {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Following.Fetch().With(
					db.User.Dweets.Fetch().Take(dweetsToFetch),
					db.User.Followers.Fetch(),
				),
			).Exec(ctx)
		}
	} else {
		if dweetsToFetch < 0 {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Following.Fetch().Take(numberToFetch).With(
					db.User.Dweets.Fetch(),
					db.User.Followers.Fetch(),
				),
			).Exec(ctx)
		} else {
			user, err = client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Following.Fetch().Take(numberToFetch).With(
					db.User.Dweets.Fetch().Take(dweetsToFetch),
					db.User.Followers.Fetch(),
				),
			).Exec(ctx)
		}
	}
	if err != nil {
		return []UserType{}, err
	}

	userFullFollowing, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Following.Fetch(),
	).Exec(ctx)
	if err != nil {
		return []UserType{}, err
	}

	var following []UserType
	for _, followed := range user.Following() {
		followerFollowers := followed.Followers()

		mutuals := HashIntersectUsers(followerFollowers, userFullFollowing.Following())

		mutuals = append(mutuals, *user)
		following = append(following, AuthFormatAsUserType(&followed, mutuals))
	}

	return following, err
}

// Delete a dweet
func AuthDeleteDweet(postID string, userID string) (DweetType, error) {
	dweet, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch().With(
			db.User.Following.Fetch(),
		),
		db.Dweet.RedweetOf.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
		db.Dweet.ReplyTo.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
		db.Dweet.LikeUsers.Fetch(),
		db.Dweet.ReplyDweets.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	if dweet.Author().Username == userID {
		dweet, err := DeleteDweet(postID)
		if err != nil {
			return DweetType{}, err
		}

		mutuals := HashIntersectUsers(dweet.LikeUsers(), dweet.Author().Following())
		formatted := AuthFormatAsDweetType(dweet, mutuals)
		return formatted, err
	}

	return DweetType{}, errors.New("Unauthorized")

}

// Delete a dweet
func AuthDeleteRedweet(postID string, userID string) (DweetType, error) {
	dweet, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch().With(
			db.User.Following.Fetch(),
		),
		db.Dweet.RedweetOf.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
		db.Dweet.ReplyTo.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
		db.Dweet.LikeUsers.Fetch(),
		db.Dweet.ReplyDweets.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	if dweet.Author().Username == userID {
		dweet, err := DeleteRedweet(postID)
		if err != nil {
			return DweetType{}, err
		}

		mutuals := HashIntersectUsers(dweet.LikeUsers(), dweet.Author().Following())
		formatted := AuthFormatAsDweetType(dweet, mutuals)
		return formatted, err
	}

	return DweetType{}, errors.New("Unauthorized")
}

// Create a Post
func AuthCreateDweet(body, authorID string, mediaLinks []string) (DweetType, error) {
	now := time.Now()
	createdPost, err := client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.ID.Set(genID(10)),
		db.Dweet.Author.Link(db.User.Username.Equals(authorID)),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).With(
		db.Dweet.Author.Fetch(),
		db.Dweet.ReplyTo.Fetch(),
		db.Dweet.RedweetOf.Fetch(),
		db.Dweet.ReplyDweets.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	post := AuthFormatAsDweetType(createdPost, []db.UserModel{})
	return post, err
}

// Create a Reply
func AuthCreateReply(originalID, body, authorID string, mediaLinks []string) (DweetType, error) {
	now := time.Now()
	// Create a Reply
	createdReply, err := client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.ID.Set(genID(10)),
		db.Dweet.Author.Link(db.User.Username.Equals(authorID)),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.IsReply.Set(true),
		db.Dweet.ReplyTo.Link(
			db.Dweet.ID.Equals(originalID),
		),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).With(
		db.Dweet.Author.Fetch(),
		db.Dweet.ReplyTo.Fetch(),
		db.Dweet.RedweetOf.Fetch(),
		db.Dweet.ReplyDweets.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	// Update original Dweet to show reply
	_, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(originalID),
	).Update(
		db.Dweet.ReplyDweets.Link(
			db.Dweet.ID.Equals(createdReply.ID),
		),
		db.Dweet.ReplyCount.Increment(1),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	post := AuthFormatAsDweetType(createdReply, []db.UserModel{})
	return post, err
}

// Create a new Redweet of a Dweet
func AuthCreateRedweet(originalPostID, userID string) (DweetType, error) {
	now := time.Now()
	// Get post and user
	post, err := GetPostBasic(originalPostID)
	if err != nil {
		return DweetType{}, err
	}

	// Create a Redweet
	createdRedweet, err := client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(post.DweetBody),
		db.Dweet.ID.Set(genID(10)),
		db.Dweet.Author.Link(db.User.Username.Equals(userID)),
		db.Dweet.Media.Set(post.Media),
		db.Dweet.IsRedweet.Set(true),
		db.Dweet.RedweetOf.Link(
			db.Dweet.ID.Equals(post.ID),
		),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).With(
		db.Dweet.Author.Fetch(),
		db.Dweet.ReplyTo.Fetch(),
		db.Dweet.RedweetOf.Fetch(),
		db.Dweet.ReplyDweets.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	// Update original Dweet to show redweet
	redweeted, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(originalPostID),
	).With(
		db.Dweet.LikeUsers.Fetch(),
	).Update(
		db.Dweet.RedweetDweets.Link(
			db.Dweet.ID.Equals(createdRedweet.ID),
		),
		db.Dweet.RedweetCount.Increment(1),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	redweet := AuthFormatAsDweetType(createdRedweet, redweeted.LikeUsers())

	return redweet, err
}

// Create a follower relation
func AuthFollow(followedID string, followerID string) (UserType, error) {
	// Check if user already followed this user
	personBeingFollowed, err := client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch(
			db.User.Username.Equals(followerID),
		),
	).Exec(ctx)
	if err != nil {
		return UserType{}, err
	}

	// If yes, then skip following the user
	if len(personBeingFollowed.Followers()) > 0 {
		personBeingFollowed, err = client.User.FindUnique(
			db.User.Username.Equals(followedID),
		).With(
			db.User.Followers.Fetch(),
		).Exec(ctx)
		if err != nil {
			return UserType{}, err
		}

		authenticatedUser, err := client.User.FindUnique(
			db.User.Username.Equals(followerID),
		).With(
			db.User.Following.Fetch(),
		).Exec(ctx)

		if err != nil {
			return UserType{}, err
		}

		mutuals := HashIntersectUsers(personBeingFollowed.Followers(), authenticatedUser.Following())
		return AuthFormatAsUserType(authenticatedUser, mutuals), err
	}

	// Add follower to followed's follower list
	personBeingFollowed, err = client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch(),
	).Update(
		db.User.FollowerCount.Increment(1),
		db.User.Followers.Link(
			db.User.Username.Equals(followerID),
		),
	).Exec(ctx)
	if err != nil {
		return UserType{}, err
	}

	// Add followed to follower's following list
	authenticatedUser, err := client.User.FindUnique(
		db.User.Username.Equals(followerID),
	).With(
		db.User.Following.Fetch(),
	).Update(
		db.User.FollowingCount.Increment(1),
		db.User.Following.Link(
			db.User.Username.Equals(followedID),
		),
	).Exec(ctx)
	if err != nil {
		return UserType{}, err
	}

	mutuals := HashIntersectUsers(personBeingFollowed.Followers(), authenticatedUser.Following())
	formatted := AuthFormatAsUserType(personBeingFollowed, mutuals)

	return formatted, err
}

// Add a like to a dweet
func AuthLike(likedPostID, userID string) (DweetType, error) {
	// Check if user already liked this dweet
	likedPost, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(likedPostID),
	).With(
		db.Dweet.LikeUsers.Fetch(
			db.User.Username.Equals(userID),
		),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	// If yes, then skip liking the dweet
	if len(likedPost.LikeUsers()) > 0 {
		likedPost, err := client.Dweet.FindUnique(
			db.Dweet.ID.Equals(likedPostID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyTo.Fetch(),
			db.Dweet.RedweetOf.Fetch(),
			db.Dweet.ReplyDweets.Fetch(),
			db.Dweet.LikeUsers.Fetch(),
		).Exec(ctx)
		if err != nil {
			return DweetType{}, err
		}

		// Add post to user's liked dweets
		user, err := client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Following.Fetch(),
		).Exec(ctx)
		if err != nil {
			return DweetType{}, err
		}

		// Find known people that liked thw dweet
		mutuals := HashIntersectUsers(likedPost.LikeUsers(), user.Following())

		formatted := AuthFormatAsDweetType(likedPost, mutuals)
		return formatted, err
	}

	// Else, if not already liked,
	// Create a Like on the post if not created already
	like, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(likedPostID),
	).With(
		db.Dweet.Author.Fetch(),
		db.Dweet.ReplyTo.Fetch(),
		db.Dweet.RedweetOf.Fetch(),
		db.Dweet.ReplyDweets.Fetch(),
		db.Dweet.LikeUsers.Fetch(),
	).Update(
		db.Dweet.LikeCount.Increment(1),
		db.Dweet.LikeUsers.Link(
			db.User.Username.Equals(userID),
		),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	// Add post to user's liked dweets
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Following.Fetch(),
	).Update(
		db.User.LikedDweets.Link(
			db.Dweet.ID.Equals(like.ID),
		),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	// Find known people that liked thw dweet
	mutuals := HashIntersectUsers(like.LikeUsers(), user.Following())

	mutuals = append(mutuals, *user)

	formatted := AuthFormatAsDweetType(like, mutuals)

	return formatted, err
}

// Remove a like from a post
func AuthUnlike(postID string, userID string) (DweetType, error) {

	likedPost, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.LikeUsers.Fetch(
			db.User.Username.Equals(userID),
		),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	// If yes, then skip liking the dweet
	if len(likedPost.LikeUsers()) > 0 {
		// Find the post and decrease its likes by 1
		post, err := client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyTo.Fetch(),
			db.Dweet.RedweetOf.Fetch(),
			db.Dweet.ReplyDweets.Fetch(),
			db.Dweet.LikeUsers.Fetch(),
		).Exec(ctx)
		if err != nil {
			return DweetType{}, err
		}

		user, err := client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Following.Fetch(),
		).Exec(ctx)
		if err != nil {
			return DweetType{}, err
		}

		// Find known people that liked thw dweet
		mutuals := HashIntersectUsers(post.LikeUsers(), user.Following())

		mutuals = append(mutuals, *user)

		formatted := AuthFormatAsDweetType(post, mutuals)

		return formatted, err
	}

	// Find the post and decrease its likes by 1
	post, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch(),
		db.Dweet.ReplyTo.Fetch(),
		db.Dweet.RedweetOf.Fetch(),
		db.Dweet.ReplyDweets.Fetch(),
		db.Dweet.LikeUsers.Fetch(),
	).Update(
		db.Dweet.LikeCount.Decrement(1),
		db.Dweet.LikeUsers.Unlink(
			db.User.Username.Equals(userID),
		),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Following.Fetch(),
	).Exec(ctx)
	if err != nil {
		return DweetType{}, err
	}

	// Find known people that liked thw dweet
	mutuals := HashIntersectUsers(post.LikeUsers(), user.Following())

	mutuals = append(mutuals, *user)

	formatted := AuthFormatAsDweetType(post, mutuals)

	return formatted, err
}

// Create a follower relation
func AuthUnfollow(followedID string, followerID string) (UserType, error) {
	// Check if user already unfollowed this user
	personBeingFollowed, err := client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch(
			db.User.Username.Equals(followerID),
		),
	).Exec(ctx)
	if err != nil {
		return UserType{}, err
	}

	// If yes, then skip unfollowing the user
	if len(personBeingFollowed.Followers()) == 0 {
		personBeingFollowed, err = client.User.FindUnique(
			db.User.Username.Equals(followedID),
		).With(
			db.User.Followers.Fetch(),
		).Exec(ctx)
		if err != nil {
			return UserType{}, err
		}

		authenticatedUser, err := client.User.FindUnique(
			db.User.Username.Equals(followerID),
		).With(
			db.User.Following.Fetch(),
		).Exec(ctx)

		if err != nil {
			return UserType{}, err
		}

		mutuals := HashIntersectUsers(personBeingFollowed.Followers(), authenticatedUser.Following())
		return AuthFormatAsUserType(authenticatedUser, mutuals), err
	}

	// Add follower to followed's follower list
	personBeingFollowed, err = client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch(),
	).Update(
		db.User.FollowerCount.Decrement(1),
		db.User.Followers.Unlink(
			db.User.Username.Equals(followerID),
		),
	).Exec(ctx)
	if err != nil {
		return UserType{}, err
	}

	// Add followed to follower's following list
	authenticatedUser, err := client.User.FindUnique(
		db.User.Username.Equals(followerID),
	).With(
		db.User.Following.Fetch(),
	).Update(
		db.User.FollowingCount.Decrement(1),
		db.User.Following.Unlink(
			db.User.Username.Equals(followedID),
		),
	).Exec(ctx)
	if err != nil {
		return UserType{}, err
	}

	mutuals := HashIntersectUsers(personBeingFollowed.Followers(), authenticatedUser.Following())
	formatted := AuthFormatAsUserType(personBeingFollowed, mutuals)

	return formatted, err
}
