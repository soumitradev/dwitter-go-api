package main

import (
	"dwitter_go_graphql/prisma/db"
	"errors"

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

*/

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

	nuser := NoAuthFormatAsUserType(user)
	return nuser, err
}

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

func LoginUser(username string, password string) (LoginResponse, error) {
	user, err := client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(ctx)
	if err != nil {
		return LoginResponse{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return LoginResponse{}, errors.New("invalid password")
	}

	JWT, err := CreateToken(username)
	if err != nil {
		return LoginResponse{}, errors.New("internal server error while authenticating")
	}

	return LoginResponse{
		AccessToken: JWT,
	}, err
}
