// Package database provides some functions to interface with the posstgresql database
package database

import (
	"dwitter_go_graphql/cdn"
	"dwitter_go_graphql/common"
	"dwitter_go_graphql/prisma/db"
	"dwitter_go_graphql/schema"
	"dwitter_go_graphql/util"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Get dweet when not authenticated
func GetPostUnauth(postID string, repliesToFetch int, replyOffset int) (schema.DweetType, error) {
	// Check params and return data accordingly
	var post *db.DweetModel
	var err error
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
	mutuals := util.HashIntersectUsers(likes, following)

	// Add requesting user to like_users list
	if selfLike {
		mutuals = append(mutuals, *viewUser)
	}

	// Send back the dweet requested, along with like_users
	npost := schema.AuthFormatAsDweetType(post, mutuals)
	return npost, err
}

// Get user when not authenticated
func GetUserUnauth(username string, repliesToFetch int, dweetOffset int) (schema.UserType, error) {

	var user *db.UserModel
	var err error

	// Check params and return data accordingly
	if repliesToFetch < 0 {
		user, err = common.Client.User.FindUnique(
			db.User.Username.Equals(username),
		).With(
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(common.BaseCtx)
	} else {
		user, err = common.Client.User.FindUnique(
			db.User.Username.Equals(username),
		).With(
			db.User.Dweets.Fetch().Take(repliesToFetch).Take(dweetOffset).With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(common.BaseCtx)
	}
	if err == db.ErrNotFound {
		return schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	nuser := schema.NoAuthFormatAsUserType(user)
	return nuser, err
}

// Search users when not authenticated
func SearchUsersUnauth(text string, numberToFetch int, numOffset int, numDweets int, dweetOffset int) ([]schema.UserType, error) {
	var users []db.UserModel
	var err error

	// Check params and return data accordingly
	if numberToFetch < 0 {
		if numDweets < 0 {
			users, err = common.Client.User.FindMany(
				db.User.Username.Contains(text),
			).With(
				db.User.Dweets.Fetch(),
			).Exec(common.BaseCtx)
		} else {
			users, err = common.Client.User.FindMany(
				db.User.Username.Contains(text),
			).With(
				db.User.Dweets.Fetch().Take(numDweets).Skip(dweetOffset),
			).Exec(common.BaseCtx)
		}
	} else {
		if numDweets < 0 {
			users, err = common.Client.User.FindMany(
				db.User.Username.Contains(text),
			).With(
				db.User.Dweets.Fetch(),
			).Take(numberToFetch).Skip(numOffset).Exec(common.BaseCtx)
		} else {
			users, err = common.Client.User.FindMany(
				db.User.Username.Contains(text),
			).With(
				db.User.Dweets.Fetch().Take(numDweets).Skip(dweetOffset),
			).Take(numberToFetch).Skip(numOffset).Exec(common.BaseCtx)
		}
	}

	if err == db.ErrNotFound {
		return []schema.UserType{}, fmt.Errorf("users not found: %v", err)
	}
	if err != nil {
		return []schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Format
	var formatted []schema.UserType
	for _, user := range users {
		nuser := schema.NoAuthFormatAsUserType(&user)
		formatted = append(formatted, nuser)
	}
	return formatted, err
}

// Search users when authenticated
func SearchUsers(query string, numberToFetch int, numOffset int, numDweets int, dweetOffset int, viewerUsername string) ([]schema.UserType, error) {
	var users []db.UserModel
	var err error

	// Get your own following-list
	viewUser, err := common.Client.User.FindUnique(
		db.User.Username.Equals(viewerUsername),
	).With(
		db.User.Following.Fetch(),
	).Exec(common.BaseCtx)
	if err != nil {
		return []schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	following := viewUser.Following()

	// Check params and return data accordingly
	if numberToFetch < 0 {
		if numDweets < 0 {
			users, err = common.Client.User.FindMany(
				db.User.Username.Contains(query),
			).With(
				db.User.Dweets.Fetch(),
				db.User.Followers.Fetch(),
			).Exec(common.BaseCtx)
		} else {
			users, err = common.Client.User.FindMany(
				db.User.Username.Contains(query),
			).With(
				db.User.Dweets.Fetch().Take(numDweets).Skip(numOffset),
				db.User.Followers.Fetch(),
			).Exec(common.BaseCtx)
		}
	} else {
		if numDweets < 0 {
			users, err = common.Client.User.FindMany(
				db.User.Username.Contains(query),
			).With(
				db.User.Dweets.Fetch(),
				db.User.Followers.Fetch(),
			).Take(numberToFetch).Skip(dweetOffset).Exec(common.BaseCtx)
		} else {
			users, err = common.Client.User.FindMany(
				db.User.Username.Contains(query),
			).With(
				db.User.Dweets.Fetch().Take(numDweets).Skip(numOffset),
				db.User.Followers.Fetch(),
			).Take(numberToFetch).Skip(dweetOffset).Exec(common.BaseCtx)
		}
	}

	if err == db.ErrNotFound {
		return []schema.UserType{}, fmt.Errorf("users not found: %v", err)
	}
	if err != nil {
		return []schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Get common followers and format and return
	var formatted []schema.UserType
	for _, user := range users {
		mutuals := util.HashIntersectUsers(user.Followers(), following)
		nuser := schema.AuthFormatAsUserType(&user, mutuals)
		formatted = append(formatted, nuser)
	}
	return formatted, err
}

// Search dweets when authenticated
func SearchPosts(query string, numberToFetch int, numOffset int, repliesToFetch int, replyOffset int, viewerUsername string) ([]schema.DweetType, error) {
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
		mutuals := util.HashIntersectUsers(likes, following)

		// Add requesting user to like_users list
		if selfLike {
			mutuals = append(mutuals, *viewUser)
		}

		// Send back the dweet requested, along with like_users
		npost := schema.AuthFormatAsDweetType(&post, mutuals)
		formatted = append(formatted, npost)
	}

	return formatted, err
}

// Search dweets when not authenticated
func SearchPostsUnauth(query string, numberToFetch int, numOffset int, repliesToFetch int, replyOffset int) ([]schema.DweetType, error) {
	var posts []db.DweetModel
	var err error

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

// Get user when authenticated
func GetUser(username string, dweetsToFetch int, dweetOffset int, viewerUsername string) (schema.UserType, error) {
	var user *db.UserModel
	var mutuals []db.UserModel
	var err error

	if viewerUsername == username {
		// Fetch the user requested
		if dweetsToFetch < 0 {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(username),
			).With(
				db.User.Dweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.User.Followers.Fetch(),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		} else {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(username),
			).With(
				db.User.Dweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				).Take(dweetsToFetch).Skip(dweetOffset),
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.User.Followers.Fetch(),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		}
		if err == db.ErrNotFound {
			return schema.UserType{}, fmt.Errorf("user not found: %v", err)
		}

		if err != nil {
			return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
		}

		nuser := schema.FormatAsUserType(user)
		return nuser, err
	} else {
		// Get your own following-list
		viewUser, err := common.Client.User.FindUnique(
			db.User.Username.Equals(viewerUsername),
		).With(
			db.User.Following.Fetch(),
		).Exec(common.BaseCtx)
		if err != nil {
			return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
		}

		following := viewUser.Following()

		// Fetch the user requested with followers so we get the mutuals
		if dweetsToFetch < 0 {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(username),
			).With(
				db.User.Dweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.User.Followers.Fetch(),
			).Exec(common.BaseCtx)
		} else {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(username),
			).With(
				db.User.Dweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				).Take(dweetsToFetch).Skip(dweetOffset),
				db.User.Followers.Fetch(),
			).Exec(common.BaseCtx)
		}

		if err == db.ErrNotFound {
			return schema.UserType{}, fmt.Errorf("user not found: %v", err)
		}

		if err != nil {
			return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
		}

		// Get mutuals
		followers := user.Followers()
		mutuals = util.HashIntersectUsers(followers, following)
		// Send back the user requested, along with mutuals in the followers field
		nuser := schema.AuthFormatAsUserType(user, mutuals)
		return nuser, err
	}
}

// Create a User
func SignUpUser(username string, password string, firstName string, lastName string, bio string, email string) (schema.UserType, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Check if user with username or email already exists
	_, err1 := common.Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(common.BaseCtx)
	_, err2 := common.Client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(common.BaseCtx)
	if (err1 == db.ErrNotFound) || (err2 == db.ErrNotFound) {
		// Create user if no such user exists
		createdUser, err := common.Client.User.CreateOne(
			db.User.Username.Set(username),
			db.User.PasswordHash.Set(string(passwordHash)),
			db.User.FirstName.Set(firstName),
			db.User.Email.Set(email),
			db.User.Bio.Set(bio),
			db.User.ProfilePicURL.Set(common.DefaultPFPURL),
			db.User.TokenVersion.Set(rand.Intn(10000)),
			db.User.CreatedAt.Set(time.Now()),
			db.User.LastName.Set(lastName),
		).With(
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(common.BaseCtx)

		if err != nil {
			return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
		}

		nuser := schema.AuthFormatAsUserType(createdUser, []db.UserModel{})
		return nuser, err
	} else {
		return schema.UserType{}, errors.New("username/email already taken")
	}
}

// Update a dweet
func UpdateDweet(postID string, username string, body string, mediaLinks []string, repliesToFetch int, replyOffset int) (schema.DweetType, error) {
	post, err := common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch(),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Check if user owns dweet
	if post.Author().Username != username {
		return schema.DweetType{}, fmt.Errorf("authorization error: %v", errors.New("not authorized to edit dweet"))
	}

	// Delete the media that isn't used anymore
	oldMedia := post.Media
	toDelete := util.HashDifference(oldMedia, mediaLinks)
	for _, mediaLink := range toDelete {
		loc, err := cdn.LinkToLocation(mediaLink)
		if err != nil {
			return schema.DweetType{}, err
		}
		err = cdn.DeleteLocation(loc, true)
		if err != nil {
			return schema.DweetType{}, err
		}
	}

	// Check params and return data accordingly
	if repliesToFetch < 0 {
		post, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
		).Update(
			db.Dweet.DweetBody.Set(body),
			db.Dweet.Media.Set(mediaLinks),
			db.Dweet.LastUpdatedAt.Set(time.Now()),
		).Exec(common.BaseCtx)
	} else {
		post, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(repliesToFetch).Skip(replyOffset),
			db.Dweet.LikeUsers.Fetch(),
		).Update(
			db.Dweet.DweetBody.Set(body),
			db.Dweet.Media.Set(mediaLinks),
			db.Dweet.LastUpdatedAt.Set(time.Now()),
		).Exec(common.BaseCtx)
	}
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Mark media as used to prevent auto-deletion on expiry
	for _, link := range mediaLinks {
		delete(common.MediaCreatedButNotUsed, link)
	}

	// Add common likes and format
	user, err := common.Client.User.FindUnique(
		db.User.Username.Equals(username),
	).With(
		db.User.Following.Fetch(),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	mutuals := util.HashIntersectUsers(user.Following(), post.LikeUsers())
	npost := schema.AuthFormatAsDweetType(post, mutuals)
	return npost, err
}

// Update a user
func UpdateUser(username string, firstName string, lastName string, email string, bio string, PfpUrl string, dweetsToFetch int,
	dweetsOffset int, followersToFetch int, followersOffset int, followingToFetch int, followingOffset int) (schema.UserType, error) {
	var user *db.UserModel
	var err error

	// Check params and return data accordingly
	if followingToFetch < 0 {
		if followersToFetch < 0 {
			if dweetsToFetch < 0 {
				user, err = common.Client.User.FindUnique(
					db.User.Username.Equals(username),
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
					db.User.ProfilePicURL.Set(PfpUrl),
				).Exec(common.BaseCtx)
			} else {
				user, err = common.Client.User.FindUnique(
					db.User.Username.Equals(username),
				).With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(dweetsToFetch).Skip(dweetsOffset),
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
					db.User.ProfilePicURL.Set(PfpUrl),
				).Exec(common.BaseCtx)
			}
		} else {
			if dweetsToFetch < 0 {
				user, err = common.Client.User.FindUnique(
					db.User.Username.Equals(username),
				).With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.LikedDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch().Take(followersToFetch).Skip(followersOffset),
					db.User.Following.Fetch(),
				).Update(
					db.User.FirstName.Set(firstName),
					db.User.LastName.Set(lastName),
					db.User.Email.Set(email),
					db.User.Bio.Set(bio),
					db.User.ProfilePicURL.Set(PfpUrl),
				).Exec(common.BaseCtx)
			} else {
				user, err = common.Client.User.FindUnique(
					db.User.Username.Equals(username),
				).With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(dweetsToFetch).Skip(dweetsOffset),
					db.User.LikedDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch().Take(followersToFetch).Skip(followersOffset),
					db.User.Following.Fetch(),
				).Update(
					db.User.FirstName.Set(firstName),
					db.User.LastName.Set(lastName),
					db.User.Email.Set(email),
					db.User.Bio.Set(bio),
					db.User.ProfilePicURL.Set(PfpUrl),
				).Exec(common.BaseCtx)
			}
		}
	} else {
		if followersToFetch < 0 {
			if dweetsToFetch < 0 {
				user, err = common.Client.User.FindUnique(
					db.User.Username.Equals(username),
				).With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.LikedDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch(),
					db.User.Following.Fetch().Take(followingToFetch).Skip(followingOffset),
				).Update(
					db.User.FirstName.Set(firstName),
					db.User.LastName.Set(lastName),
					db.User.Email.Set(email),
					db.User.Bio.Set(bio),
					db.User.ProfilePicURL.Set(PfpUrl),
				).Exec(common.BaseCtx)
			} else {
				user, err = common.Client.User.FindUnique(
					db.User.Username.Equals(username),
				).With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(dweetsToFetch).Skip(dweetsOffset),
					db.User.LikedDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch(),
					db.User.Following.Fetch().Take(followingToFetch).Skip(followingOffset),
				).Update(
					db.User.FirstName.Set(firstName),
					db.User.LastName.Set(lastName),
					db.User.Email.Set(email),
					db.User.Bio.Set(bio),
					db.User.ProfilePicURL.Set(PfpUrl),
				).Exec(common.BaseCtx)
			}
		} else {
			if dweetsToFetch < 0 {
				user, err = common.Client.User.FindUnique(
					db.User.Username.Equals(username),
				).With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.LikedDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch().Take(followersToFetch).Skip(followersOffset),
					db.User.Following.Fetch().Take(followingToFetch).Skip(followingOffset),
				).Update(
					db.User.FirstName.Set(firstName),
					db.User.LastName.Set(lastName),
					db.User.Email.Set(email),
					db.User.Bio.Set(bio),
					db.User.ProfilePicURL.Set(PfpUrl),
				).Exec(common.BaseCtx)
			} else {
				user, err = common.Client.User.FindUnique(
					db.User.Username.Equals(username),
				).With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(dweetsToFetch).Skip(dweetsOffset),
					db.User.LikedDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch().Take(followersToFetch).Skip(followersOffset),
					db.User.Following.Fetch().Take(followingToFetch).Skip(followingOffset),
				).Update(
					db.User.FirstName.Set(firstName),
					db.User.LastName.Set(lastName),
					db.User.Email.Set(email),
					db.User.Bio.Set(bio),
					db.User.ProfilePicURL.Set(PfpUrl),
				).Exec(common.BaseCtx)
			}
		}
	}
	if err == db.ErrNotFound {
		return schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	nuser := schema.FormatAsUserType(user)
	return nuser, err
}

// Get User's liked dweets
func GetLikedDweets(userID string, numberToFetch int, numOffset int, repliesToFetch int, replyOffset int) ([]schema.DweetType, error) {
	var user *db.UserModel
	var err error

	// Check params and return data accordingly
	if numberToFetch < 0 {
		if repliesToFetch < 0 {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.ReplyDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.LikeUsers.Fetch(),
				),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		} else {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.ReplyDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(repliesToFetch).Skip(replyOffset),
					db.Dweet.LikeUsers.Fetch(),
				),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		}
	} else {
		if repliesToFetch < 0 {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.ReplyDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.LikeUsers.Fetch(),
				).Take(numberToFetch).Skip(numOffset),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		} else {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.LikedDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.ReplyDweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(repliesToFetch).Skip(replyOffset),
					db.Dweet.LikeUsers.Fetch(),
				).Take(numberToFetch).Skip(numOffset),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		}
	}
	if err == db.ErrNotFound {
		return []schema.DweetType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return []schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Add common likes and return formatted
	var liked []schema.DweetType
	for _, dweet := range user.LikedDweets() {
		likes := dweet.LikeUsers()

		// Find known people that liked thw dweet
		mutuals := util.HashIntersectUsers(likes, user.Following())

		// Add requesting user to likeUsers list
		mutuals = append(mutuals, *user)

		liked = append(liked, schema.AuthFormatAsDweetType(&dweet, mutuals))
	}
	return liked, err
}

// Get users that follow user
func GetFollowers(userID string, numberToFetch int, numOffset int, dweetsToFetch int, dweetOffset int) ([]schema.UserType, error) {
	var user *db.UserModel
	var err error

	// Check params and return data accordingly
	if numberToFetch < 0 {
		if dweetsToFetch < 0 {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Followers.Fetch().With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch(),
				),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		} else {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Followers.Fetch().With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(dweetsToFetch).Skip(dweetOffset),
					db.User.Followers.Fetch(),
				),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		}
	} else {
		if dweetsToFetch < 0 {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Followers.Fetch().With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch(),
				).Take(numberToFetch).Skip(numOffset),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		} else {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Followers.Fetch().With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(dweetsToFetch).Skip(dweetOffset),
					db.User.Followers.Fetch(),
				).Take(numberToFetch).Skip(numOffset),
				db.User.Following.Fetch(),
			).Exec(common.BaseCtx)
		}
	}
	if err == db.ErrNotFound {
		return []schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return []schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Add common followers and format
	var followers []schema.UserType
	for _, follower := range user.Followers() {
		followerFollowers := follower.Followers()

		mutuals := util.HashIntersectUsers(followerFollowers, user.Following())

		mutuals = append(mutuals, *user)
		followers = append(followers, schema.AuthFormatAsUserType(&follower, mutuals))
	}
	return followers, err
}

// Get users that user follows
func GetFollowing(userID string, numberToFetch int, numOffset int, dweetsToFetch int, dweetOffset int) ([]schema.UserType, error) {
	var user *db.UserModel
	var err error

	// Check params and return data accordingly
	if numberToFetch < 0 {
		if dweetsToFetch < 0 {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Following.Fetch().With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch(),
				),
			).Exec(common.BaseCtx)
		} else {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Following.Fetch().With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(dweetsToFetch).Skip(dweetOffset),
					db.User.Followers.Fetch(),
				),
			).Exec(common.BaseCtx)
		}
	} else {
		if dweetsToFetch < 0 {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Following.Fetch().With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.User.Followers.Fetch(),
				).Take(numberToFetch).Skip(numOffset),
			).Exec(common.BaseCtx)
		} else {
			user, err = common.Client.User.FindUnique(
				db.User.Username.Equals(userID),
			).With(
				db.User.Following.Fetch().With(
					db.User.Dweets.Fetch().With(
						db.Dweet.Author.Fetch(),
					).Take(dweetsToFetch).Skip(dweetOffset),
					db.User.Followers.Fetch(),
				).Take(numberToFetch).Skip(numOffset),
			).Exec(common.BaseCtx)
		}
	}
	if err == db.ErrNotFound {
		return []schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return []schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Add common followers and return
	userFullFollowing, err := common.Client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Following.Fetch(),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return []schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return []schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	var following []schema.UserType
	for _, followed := range user.Following() {
		followerFollowers := followed.Followers()

		mutuals := util.HashIntersectUsers(followerFollowers, userFullFollowing.Following())

		mutuals = append(mutuals, *user)
		following = append(following, schema.AuthFormatAsUserType(&followed, mutuals))
	}

	return following, err
}

// Delete a dweet
func DeleteDweet(postID string, username string, repliesToFetch int, replyOffset int) (schema.DweetType, error) {
	var deleted *db.DweetModel
	var err error

	// Check params and return data accordingly
	if repliesToFetch < 0 {
		deleted, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch().With(
				db.User.Following.Fetch(),
			),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(common.BaseCtx)
	} else {
		deleted, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch().With(
				db.User.Following.Fetch(),
			),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
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
		mutuals := util.HashIntersectUsers(deleted.LikeUsers(), deleted.Author().Following())
		formatted := schema.AuthFormatAsDweetType(deleted, mutuals)
		return formatted, err
	}

	return schema.DweetType{}, fmt.Errorf("internal server error: %v", errors.New("Unauthorized"))
}

// Delete a redweet
func DeleteRedweet(postID string, username string) (schema.RedweetType, error) {
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

// Create a Post
func NewDweet(body, username string, mediaLinks []string) (schema.DweetType, error) {
	// Generate a unique ID
	randID := util.GenID(10)
	_, err := common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(randID),
	).Exec(common.BaseCtx)

	for err != db.ErrNotFound {
		randID := util.GenID(10)

		_, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(randID),
		).Exec(common.BaseCtx)
	}

	now := time.Now()
	createdPost, err := common.Client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.ID.Set(randID),
		db.Dweet.Author.Link(db.User.Username.Equals(username)),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).With(
		db.Dweet.Author.Fetch(),
		db.Dweet.ReplyTo.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
		db.Dweet.ReplyDweets.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
	).Exec(common.BaseCtx)
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Mark media as used to prevent deletion on expiry
	for _, link := range mediaLinks {
		delete(common.MediaCreatedButNotUsed, link)
	}

	// Format and return
	post := schema.AuthFormatAsDweetType(createdPost, []db.UserModel{})
	return post, err
}

// Create a Reply
func NewReply(originalPostID string, body string, authorUsername string, mediaLinks []string) (schema.DweetType, error) {
	// Generate unique ID
	randID := util.GenID(10)
	_, err := common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(randID),
	).Exec(common.BaseCtx)

	for err != db.ErrNotFound {
		randID := util.GenID(10)

		_, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(randID),
		).Exec(common.BaseCtx)
	}

	now := time.Now()
	// Create a Reply
	createdReply, err := common.Client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.ID.Set(randID),
		db.Dweet.Author.Link(db.User.Username.Equals(authorUsername)),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.IsReply.Set(true),
		db.Dweet.ReplyTo.Link(
			db.Dweet.ID.Equals(originalPostID),
		),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).With(
		db.Dweet.Author.Fetch(),
		db.Dweet.ReplyTo.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
		db.Dweet.ReplyDweets.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
	).Exec(common.BaseCtx)
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}
	for _, link := range mediaLinks {
		delete(common.MediaCreatedButNotUsed, link)
	}

	// Update original Dweet to show reply
	_, err = common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(originalPostID),
	).Update(
		db.Dweet.ReplyDweets.Link(
			db.Dweet.ID.Equals(createdReply.ID),
		),
		db.Dweet.ReplyCount.Increment(1),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("original dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	post := schema.AuthFormatAsDweetType(createdReply, []db.UserModel{})
	return post, err
}

// Create a new Redweet of a Dweet
func Redweet(originalPostID, username string) (schema.RedweetType, error) {
	// Create a Redweet
	user, err := common.Client.User.FindUnique(
		db.User.Username.Equals(username),
	).With(
		db.User.Redweets.Fetch(
			db.Redweet.OriginalRedweetID.Equals(originalPostID),
		),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.RedweetType{}, fmt.Errorf("original dweet not found: %v", err)
	}
	if err != nil {
		return schema.RedweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// If already redweeted, return redweet
	if len(user.Redweets()) > 0 {
		redweet, err := common.Client.Redweet.FindUnique(
			db.Redweet.DbID.Equals(user.Redweets()[0].DbID),
		).With(
			db.Redweet.Author.Fetch(),
			db.Redweet.RedweetOf.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(common.BaseCtx)
		return schema.FormatAsRedweetType(redweet), err
	}

	// Create a Redweet
	createdRedweet, err := common.Client.Redweet.CreateOne(
		db.Redweet.Author.Link(
			db.User.Username.Equals(username),
		),
		db.Redweet.RedweetOf.Link(
			db.Dweet.ID.Equals(originalPostID),
		),
		db.Redweet.RedweetTime.Set(time.Now()),
	).With(
		db.Redweet.Author.Fetch(),
		db.Redweet.RedweetOf.Fetch().With(
			db.Dweet.Author.Fetch(),
		),
	).Exec(common.BaseCtx)
	if err != nil {
		return schema.RedweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Update original Dweet to show redweet
	_, err = common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(originalPostID),
	).Update(
		db.Dweet.RedweetDweets.Link(
			db.Redweet.DbID.Equals(createdRedweet.DbID),
		),
		db.Dweet.RedweetCount.Increment(1),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.RedweetType{}, fmt.Errorf("original dweet not found: %v", err)
	}
	if err != nil {
		return schema.RedweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	return schema.FormatAsRedweetType(createdRedweet), err
}

// Create a follower relation
func Follow(followedID string, followerID string, dweetsToFetch int, dweetOffset int) (schema.UserType, error) {
	// Check if user already followed this user
	personBeingFollowed, err := common.Client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch(
			db.User.Username.Equals(followerID),
		),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// If yes, then skip following the user
	if len(personBeingFollowed.Followers()) > 0 {
		if dweetsToFetch < 0 {
			personBeingFollowed, err = common.Client.User.FindUnique(
				db.User.Username.Equals(followedID),
			).With(
				db.User.Followers.Fetch(),
				db.User.Dweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
			).Exec(common.BaseCtx)
		} else {
			personBeingFollowed, err = common.Client.User.FindUnique(
				db.User.Username.Equals(followedID),
			).With(
				db.User.Followers.Fetch(),
				db.User.Dweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				).Take(dweetsToFetch).Skip(dweetOffset),
			).Exec(common.BaseCtx)
		}
		if err == db.ErrNotFound {
			return schema.UserType{}, fmt.Errorf("user not found: %v", err)
		}
		if err != nil {
			return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
		}

		authenticatedUser, err := common.Client.User.FindUnique(
			db.User.Username.Equals(followerID),
		).With(
			db.User.Following.Fetch(),
		).Exec(common.BaseCtx)
		if err == db.ErrNotFound {
			return schema.UserType{}, fmt.Errorf("user not found: %v", err)
		}
		if err != nil {
			return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
		}

		mutuals := util.HashIntersectUsers(personBeingFollowed.Followers(), authenticatedUser.Following())
		return schema.AuthFormatAsUserType(personBeingFollowed, mutuals), err
	}

	// Add follower to followed's follower list
	if dweetsToFetch < 0 {
		personBeingFollowed, err = common.Client.User.FindUnique(
			db.User.Username.Equals(followedID),
		).With(
			db.User.Followers.Fetch(),
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Update(
			db.User.FollowerCount.Increment(1),
			db.User.Followers.Link(
				db.User.Username.Equals(followerID),
			),
		).Exec(common.BaseCtx)
	} else {
		personBeingFollowed, err = common.Client.User.FindUnique(
			db.User.Username.Equals(followedID),
		).With(
			db.User.Followers.Fetch(),
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(dweetsToFetch),
		).Update(
			db.User.FollowerCount.Increment(1),
			db.User.Followers.Link(
				db.User.Username.Equals(followerID),
			),
		).Exec(common.BaseCtx)
	}
	if err == db.ErrNotFound {
		return schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Add followed to follower's following list
	authenticatedUser, err := common.Client.User.FindUnique(
		db.User.Username.Equals(followerID),
	).With(
		db.User.Following.Fetch(),
	).Update(
		db.User.FollowingCount.Increment(1),
		db.User.Following.Link(
			db.User.Username.Equals(followedID),
		),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	mutuals := util.HashIntersectUsers(personBeingFollowed.Followers(), authenticatedUser.Following())
	formatted := schema.AuthFormatAsUserType(personBeingFollowed, mutuals)

	return formatted, err
}

// Add a like to a dweet
func Like(likedPostID, userID string, repliesToFetch int, replyOffset int) (schema.DweetType, error) {
	// Check if user already liked this dweet
	likedPost, err := common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(likedPostID),
	).With(
		db.Dweet.LikeUsers.Fetch(
			db.User.Username.Equals(userID),
		),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// If yes, then skip liking the dweet
	if len(likedPost.LikeUsers()) > 0 {
		if repliesToFetch < 0 {
			likedPost, err = common.Client.Dweet.FindUnique(
				db.Dweet.ID.Equals(likedPostID),
			).With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.LikeUsers.Fetch(),
			).Exec(common.BaseCtx)
			if err == db.ErrNotFound {
				return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
			}
			if err != nil {
				return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
			}
		} else {
			likedPost, err = common.Client.Dweet.FindUnique(
				db.Dweet.ID.Equals(likedPostID),
			).With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				).Take(repliesToFetch).Skip(replyOffset),
				db.Dweet.LikeUsers.Fetch(),
			).Exec(common.BaseCtx)
			if err == db.ErrNotFound {
				return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
			}
			if err != nil {
				return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
			}
		}

		// Add post to user's liked dweets
		user, err := common.Client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Following.Fetch(),
		).Exec(common.BaseCtx)
		if err == db.ErrNotFound {
			return schema.DweetType{}, fmt.Errorf("user not found: %v", err)
		}
		if err != nil {
			return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
		}

		// Find known people that liked thw dweet
		mutuals := util.HashIntersectUsers(likedPost.LikeUsers(), user.Following())
		mutuals = append(mutuals, *user)

		formatted := schema.AuthFormatAsDweetType(likedPost, mutuals)
		return formatted, err
	}

	// Else, if not already liked,
	// Create a Like on the post if not created already
	var like *db.DweetModel
	if repliesToFetch < 0 {
		like, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(likedPostID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
		).Update(
			db.Dweet.LikeCount.Increment(1),
			db.Dweet.LikeUsers.Link(
				db.User.Username.Equals(userID),
			),
		).Exec(common.BaseCtx)
		if err == db.ErrNotFound {
			return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
		}
		if err != nil {
			return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
		}
	} else {
		like, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(likedPostID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(repliesToFetch),
			db.Dweet.LikeUsers.Fetch(),
		).Update(
			db.Dweet.LikeCount.Increment(1),
			db.Dweet.LikeUsers.Link(
				db.User.Username.Equals(userID),
			),
		).Exec(common.BaseCtx)
		if err == db.ErrNotFound {
			return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
		}
		if err != nil {
			return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
		}
	}

	// Add post to user's liked dweets
	user, err := common.Client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Following.Fetch(),
	).Update(
		db.User.LikedDweets.Link(
			db.Dweet.ID.Equals(like.ID),
		),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Find known people that liked thw dweet
	mutuals := util.HashIntersectUsers(like.LikeUsers(), user.Following())

	mutuals = append(mutuals, *user)

	formatted := schema.AuthFormatAsDweetType(like, mutuals)

	return formatted, err
}

// Remove a like from a dweet
func Unlike(postID string, userID string, repliesToFetch int, replyOffset int) (schema.DweetType, error) {

	likedPost, err := common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.LikeUsers.Fetch(
			db.User.Username.Equals(userID),
		),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// If yes, then skip unliking the dweet
	if len(likedPost.LikeUsers()) == 0 {
		var post *db.DweetModel
		if repliesToFetch < 0 {
			post, err = common.Client.Dweet.FindUnique(
				db.Dweet.ID.Equals(postID),
			).With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.LikeUsers.Fetch(),
			).Exec(common.BaseCtx)
			if err == db.ErrNotFound {
				return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
			}
			if err != nil {
				return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
			}
		} else {
			post, err = common.Client.Dweet.FindUnique(
				db.Dweet.ID.Equals(postID),
			).With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				).Take(repliesToFetch).Skip(replyOffset),
				db.Dweet.LikeUsers.Fetch(),
			).Exec(common.BaseCtx)
			if err == db.ErrNotFound {
				return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
			}
			if err != nil {
				return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
			}
		}

		user, err := common.Client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Following.Fetch(),
		).Exec(common.BaseCtx)
		if err == db.ErrNotFound {
			return schema.DweetType{}, fmt.Errorf("user not found: %v", err)
		}
		if err != nil {
			return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
		}

		// Find known people that liked the dweet
		mutuals := util.HashIntersectUsers(post.LikeUsers(), user.Following())

		formatted := schema.AuthFormatAsDweetType(post, mutuals)

		return formatted, err
	}

	// Find the post and decrease its likes by 1
	var post *db.DweetModel
	if repliesToFetch < 0 {
		post, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.LikeUsers.Fetch(),
		).Update(
			db.Dweet.LikeCount.Decrement(1),
			db.Dweet.LikeUsers.Unlink(
				db.User.Username.Equals(userID),
			),
		).Exec(common.BaseCtx)
	} else {
		post, err = common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyTo.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
			db.Dweet.ReplyDweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(repliesToFetch),
			db.Dweet.LikeUsers.Fetch(),
		).Update(
			db.Dweet.LikeCount.Decrement(1),
			db.Dweet.LikeUsers.Unlink(
				db.User.Username.Equals(userID),
			),
		).Exec(common.BaseCtx)
	}
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("dweet not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	user, err := common.Client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Following.Fetch(),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.DweetType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Find known people that liked thw dweet
	mutuals := util.HashIntersectUsers(post.LikeUsers(), user.Following())

	mutuals = append(mutuals, *user)

	formatted := schema.AuthFormatAsDweetType(post, mutuals)

	return formatted, err
}

// Create a follower relation
func Unfollow(followedID string, followerID string, dweetsToFetch int, dweetOffset int) (schema.UserType, error) {
	// Check if user already unfollowed this user
	personBeingFollowed, err := common.Client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch(
			db.User.Username.Equals(followerID),
		),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// If yes, then skip unfollowing the user
	if len(personBeingFollowed.Followers()) == 0 {
		personBeingFollowed, err = common.Client.User.FindUnique(
			db.User.Username.Equals(followedID),
		).With(
			db.User.Followers.Fetch(),
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Exec(common.BaseCtx)
		if err == db.ErrNotFound {
			return schema.UserType{}, fmt.Errorf("user not found: %v", err)
		}
		if err != nil {
			return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
		}

		authenticatedUser, err := common.Client.User.FindUnique(
			db.User.Username.Equals(followerID),
		).With(
			db.User.Following.Fetch(),
		).Exec(common.BaseCtx)

		if err == db.ErrNotFound {
			return schema.UserType{}, fmt.Errorf("user not found: %v", err)
		}
		if err != nil {
			return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
		}

		mutuals := util.HashIntersectUsers(personBeingFollowed.Followers(), authenticatedUser.Following())
		return schema.AuthFormatAsUserType(personBeingFollowed, mutuals), err
	}

	// Add follower to followed's follower list
	if dweetsToFetch < 0 {
		personBeingFollowed, err = common.Client.User.FindUnique(
			db.User.Username.Equals(followedID),
		).With(
			db.User.Followers.Fetch(),
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		).Update(
			db.User.FollowerCount.Decrement(1),
			db.User.Followers.Unlink(
				db.User.Username.Equals(followerID),
			),
		).Exec(common.BaseCtx)
	} else {
		personBeingFollowed, err = common.Client.User.FindUnique(
			db.User.Username.Equals(followedID),
		).With(
			db.User.Followers.Fetch(),
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
			).Take(dweetsToFetch).Skip(dweetOffset),
		).Update(
			db.User.FollowerCount.Decrement(1),
			db.User.Followers.Unlink(
				db.User.Username.Equals(followerID),
			),
		).Exec(common.BaseCtx)
	}
	if err == db.ErrNotFound {
		return schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	// Add followed to follower's following list
	authenticatedUser, err := common.Client.User.FindUnique(
		db.User.Username.Equals(followerID),
	).With(
		db.User.Following.Fetch(),
	).Update(
		db.User.FollowingCount.Decrement(1),
		db.User.Following.Unlink(
			db.User.Username.Equals(followedID),
		),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return schema.UserType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return schema.UserType{}, fmt.Errorf("internal server error: %v", err)
	}

	mutuals := util.HashIntersectUsers(personBeingFollowed.Followers(), authenticatedUser.Following())
	formatted := schema.AuthFormatAsUserType(personBeingFollowed, mutuals)

	return formatted, err
}

// Get feed for authenticated user
func GetFeed(username string) ([]interface{}, error) {
	// grab followed users by username
	// Grab their dweets and redweets
	user, err := common.Client.User.FindUnique(
		db.User.Username.Equals(username),
	).With(
		db.User.Following.Fetch().With(
			db.User.Dweets.Fetch().With(
				db.Dweet.Author.Fetch(),
				db.Dweet.ReplyDweets.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.LikeUsers.Fetch(),
			).OrderBy(
				db.Dweet.PostedAt.Order(db.DESC),
			),

			db.User.Redweets.Fetch().With(
				db.Redweet.Author.Fetch(),
				db.Redweet.RedweetOf.Fetch().With(
					db.Dweet.Author.Fetch(),
					db.Dweet.ReplyTo.Fetch().With(
						db.Dweet.Author.Fetch(),
					),
					db.Dweet.LikeUsers.Fetch(),
				),
			).OrderBy(
				db.Redweet.RedweetTime.Order(db.DESC),
			),
		),
	).Exec(common.BaseCtx)

	if err == db.ErrNotFound {
		return []interface{}{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return []interface{}{}, fmt.Errorf("internal server error: %v", err)
	}

	following := user.Following()

	// Merge the lists, format and return
	var posts []db.DweetModel
	var redweets []db.RedweetModel

	for _, feedUser := range following {
		posts = util.MergeDweetLists(posts, feedUser.Dweets())
		redweets = util.MergeRedweetLists(redweets, feedUser.Redweets())
	}

	merged := util.MergeDweetRedweetList(posts, redweets)

	var formatted []interface{}
	for _, post := range merged {
		var npost interface{}
		if dweet, ok := post.(db.DweetModel); ok {
			npost = schema.AuthFormatAsDweetType(&dweet, util.HashIntersectUsers(dweet.LikeUsers(), user.Following()))
		}
		if redweet, ok := post.(db.RedweetModel); ok {
			npost = schema.FormatAsRedweetType(&redweet)
		}
		formatted = append(formatted, npost)
	}
	return formatted, err
}
