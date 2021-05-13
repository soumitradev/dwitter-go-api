package main

import (
	"context"
	"errors"
	"time"

	"dwitter_go_graphql/prisma/db"
)

var client *db.PrismaClient

var ctx context.Context

// Connect to the database using prisma
func ConnectDB() {
	client = db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	ctx = context.Background()
}

// Disconnect from DB
func DisconnectDB() {
	if err := client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}

// Get basic User data
func GetUser(userID string) (*db.UserModel, error) {
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Exec(ctx)
	return user, err
}

// Get User data with dweets of user
func GetUserDweets(userID string) (*db.UserModel, error) {
	var user *db.UserModel
	var err error
	user, err = client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Dweets.Fetch(),
		db.User.Redweets.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, err
}

// Get User data with dweets that user liked
func GetUserLikes(userID string) (*db.UserModel, error) {
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.LikedDweets.Fetch(),
	).Exec(ctx)
	return user, err
}

// Get User data with followers of user
func GetFollowers(userID string) (*db.UserModel, error) {
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Followers.Fetch(),
	).Exec(ctx)
	return user, err
}

// Get User data with users that user follows
func GetFollowing(userID string) (*db.UserModel, error) {
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Following.Fetch(),
	).Exec(ctx)
	return user, err
}

// Get Replies to post
func GetPostBasic(postID string) (*db.DweetModel, error) {
	post, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Exec(ctx)
	return post, err
}

// Get Replies to post
func GetPostReplies(postID string, replies_to_fetch int) (*db.DweetModel, error) {
	var post *db.DweetModel
	var err error
	if replies_to_fetch < 0 {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch(),
			db.Dweet.ReplyTo.Fetch(),
			db.Dweet.LikeUsers.Fetch(),
		).Exec(ctx)
	} else {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.Author.Fetch(),
			db.Dweet.ReplyDweets.Fetch().Take(replies_to_fetch),
			db.Dweet.ReplyTo.Fetch(),
			db.Dweet.LikeUsers.Fetch(),
		).Exec(ctx)
	}
	return post, err
}

// Get redweets of post
func GetPostRedweets(postID string, redweets_to_fetch int) (*db.DweetModel, error) {
	var post *db.DweetModel
	var err error
	if redweets_to_fetch < 0 {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.RedweetDweets.Fetch(),
		).Exec(ctx)
	} else {
		post, err = client.Dweet.FindUnique(
			db.Dweet.ID.Equals(postID),
		).With(
			db.Dweet.RedweetDweets.Fetch().Take(redweets_to_fetch),
		).Exec(ctx)
	}
	return post, err
}

// Create a User
func NewUser(username, passwordHash, firstName, lastName, email, bio string) (*db.UserModel, error) {
	createdUser, err := client.User.CreateOne(
		db.User.Username.Set(username),
		db.User.PasswordHash.Set(passwordHash),
		db.User.FirstName.Set(firstName),
		db.User.Email.Set(email),
		db.User.Bio.Set(bio),
		db.User.CreatedAt.Set(time.Now()),
		db.User.LastName.Set(lastName),
	).Exec(ctx)
	return createdUser, err
}

// Create a Post
func NewDweet(body, authorID string, mediaLinks []string) (*db.DweetModel, error) {
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
	).Exec(ctx)
	return createdPost, err
}

// Add a like to a dweet
func NewLike(likedPostID, userID string) (*db.DweetModel, error) {
	// Check if user already liked this dweet
	likedPost, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(likedPostID),
	).With(
		db.Dweet.Author.Fetch(),
		db.Dweet.LikeUsers.Fetch(
			db.User.Username.Equals(userID),
		).With(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// If yes, then skip liking the dweet
	if len(likedPost.LikeUsers()) > 0 {
		return likedPost, err
	}

	// Else, if not already liked,
	// Create a Like on the post if not created already
	like, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(likedPostID),
	).With(
		db.Dweet.Author.Fetch(),
	).Update(
		db.Dweet.LikeCount.Increment(1),
		db.Dweet.LikeUsers.Link(
			db.User.Username.Equals(userID),
		),
	).Exec(ctx)
	if err != nil {
		return like, err
	}

	// Add post to user's liked dweets
	_, err = client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Update(
		db.User.LikedDweets.Link(
			db.Dweet.ID.Equals(like.ID),
		),
	).Exec(ctx)

	return like, err
}

// Create a reply to a post
func NewReply(originalPostID, userID, body string, mediaLinks []string) (*db.DweetModel, error) {
	now := time.Now()
	// Create a Reply
	createdReply, err := client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.ID.Set(genID(10)),
		db.Dweet.Author.Link(db.User.Username.Equals(userID)),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.IsReply.Set(true),
		db.Dweet.ReplyTo.Link(
			db.Dweet.ID.Equals(originalPostID),
		),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).With(
		db.Dweet.Author.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Update original Dweet to show reply
	_, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(originalPostID),
	).Update(
		db.Dweet.ReplyDweets.Link(
			db.Dweet.ID.Equals(createdReply.ID),
		),
		db.Dweet.ReplyCount.Increment(1),
	).Exec(ctx)

	return createdReply, err
}

// Create a new Redweet of a Dweet
func NewRedweet(originalPostID, userID string) (*db.RedweetModel, error) {
	// Create a Redweet
	createdRedweet, err := client.Redweet.CreateOne(
		db.Redweet.Author.Link(
			db.User.Username.Contains(userID),
		),
		db.Redweet.RedweetOf.Link(
			db.Dweet.ID.Equals(originalPostID),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Update original Dweet to show redweet
	_, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(originalPostID),
	).Update(
		db.Dweet.RedweetDweets.Link(
			db.Redweet.DbID.Equals(createdRedweet.DbID),
		),
		db.Dweet.RedweetCount.Increment(1),
	).Exec(ctx)

	return createdRedweet, err
}

// Create a follower relation
func NewFollower(followedID string, followerID string) (*db.UserModel, error) {
	// Check if user already followed this user
	myUser, err := client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).With(
		db.User.Followers.Fetch(
			db.User.Username.Equals(followerID),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// If yes, then skip following the user
	if len(myUser.Followers()) > 0 {
		return myUser, err
	}

	// Add follower to followed's follower list
	user, err := client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).Update(
		db.User.FollowerCount.Increment(1),
		db.User.Followers.Link(
			db.User.Username.Equals(followerID),
		),
	).Exec(ctx)
	if err != nil {
		return user, err
	}

	// Add followed to follower's following list
	_, err = client.User.FindUnique(
		db.User.Username.Equals(followerID),
	).Update(
		db.User.FollowingCount.Increment(1),
		db.User.Following.Link(
			db.User.Username.Equals(followedID),
		),
	).Exec(ctx)

	return user, err
}

// Update a user
func UpdateUser(userID, username, firstName, lastName, email, bio string) (*db.UserModel, error) {
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Update(
		db.User.FirstName.Set(firstName),
		db.User.LastName.Set(lastName),
		db.User.Email.Set(email),
		db.User.Bio.Set(bio),
	).Exec(ctx)

	return user, err
}

// Update a dweet
func UpdateDweet(postID, body string, mediaLinks []string) (*db.DweetModel, error) {
	post, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Update(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.LastUpdatedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return post, err
}

// Delete a follower relation
func DeleteFollower(followedID string, followerID string) (*db.UserModel, error) {
	// Decrement the follower and following counts
	user, err := client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).Update(
		db.User.FollowerCount.Decrement(1),
		db.User.Followers.Unlink(
			db.User.Username.Equals(followerID),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = client.User.FindUnique(
		db.User.Username.Equals(followerID),
	).Update(
		db.User.FollowingCount.Decrement(1),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, err
}

// Remove a like from a post
func DeleteLike(postID string, userID string) (*db.DweetModel, error) {
	// Find the post and decrease its likes by 1
	post, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.Author.Fetch(),
	).Update(
		db.Dweet.LikeCount.Decrement(1),
		db.Dweet.LikeUsers.Unlink(
			db.User.Username.Equals(userID),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return post, err
}

// Delete a User
func DeleteUser(userID string) (*db.UserModel, error) {
	// Get all the user's Dweets (we must delete these first since they depend on the User, and deleting the User first will render the DB invalid)
	user, err := GetUserDweets(userID)
	if err != nil {
		return nil, err
	}

	// Get all the user's Likes (we must delete these first since they depend on the User as well)
	userLikes, err := GetUserLikes(userID)
	if err != nil {
		return nil, err
	}

	// Delete all user's Dweets
	dweets := user.Dweets()
	for i := 0; i < len(dweets); i++ {
		DeleteDweet(dweets[i].ID)
	}

	// Delete all user's Redweets
	redweets := user.Redweets()
	for i := 0; i < len(redweets); i++ {
		DeleteRedweet(redweets[i].OriginalRedweetID, user.Username)
	}

	// Remove all likes of the User
	likes := userLikes.LikedDweets()
	for i := 0; i < len(likes); i++ {
		DeleteLike(likes[i].ID, userID)
	}

	// Delete the user
	user, err = client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Delete().Exec(ctx)

	return user, err
}

// Delete a Dweet
func DeleteDweet(postID string) (*db.DweetModel, error) {
	// Get all the replies to the post (these need to be deleted first since they depend on the root Dweet)
	post, err := GetPostBasic(postID)
	if err != nil {
		return nil, err
	}

	// If the Dweet itself is a reply, remove the reply from the original post
	if post.IsReply {
		// Find the dweet that was replied to
		id, exist := post.OriginalReplyID()
		if !exist {
			return nil, errors.New("original Dweet not found")
		}

		// Remove the Reply from the post
		_, err := client.Dweet.FindUnique(
			db.Dweet.ID.Equals(id),
		).Update(
			db.Dweet.ReplyCount.Decrement(1),
			db.Dweet.ReplyDweets.Unlink(
				db.Dweet.ID.Equals(postID),
			),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	dweet, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.RedweetDweets.Fetch().With(
			db.Redweet.Author.Fetch(),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	for _, redweet := range dweet.RedweetDweets() {
		DeleteRedweet(redweet.OriginalRedweetID, redweet.Author().Username)
	}

	dweet, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.ReplyDweets.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	for _, daughterDweet := range dweet.ReplyDweets() {
		DeleteDweet(daughterDweet.ID)
	}

	_, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Delete().Exec(ctx)
	if err != nil {
		return nil, err
	}

	// The following comment block is kept as an homage to the great recursive SQL function that once resided here.
	// May the soul of this legendary query rest in peace. It was a honor to use you.

	// Delete all the dependent posts (this includes redweets and replies to the post) recursively using RAW SQL
	// We use RAW SQL here because prisma-go-client doesn't support cascade deletes yet:
	// Link: https://github.com/prisma/prisma-client-go/issues/201
	// Recursive SQL function with modifications from: https://stackoverflow.com/q/10381243
	// delQuery := `WITH RECURSIVE all_posts (id, parentid1, root_id) AS (SELECT t1.db_id, t1.original_reply_id AS parentid1, t1.db_id AS root_id FROM public."Dweet" t1 UNION ALL SELECT c1.db_id, c1.original_reply_id AS parentid1, p.root_id FROM public."Dweet" c1 JOIN all_posts p ON (p.id = c1.original_reply_id) ) DELETE FROM public."Dweet" WHERE db_id IN ( SELECT id FROM all_posts WHERE root_id = $1);`
	// _, err = client.Prisma.ExecuteRaw(delQuery, post.DbID).Exec(ctx)

	return post, err
}

// Remove a Redweet
func DeleteRedweet(postID string, userID string) (*db.RedweetModel, error) {
	// Get all the replies to the redweet (these need to be deleted first since they depend on the root Redweet)
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Redweets.Fetch(
			db.Redweet.OriginalRedweetID.Equals(postID),
		).With(
			db.Redweet.Author.Fetch(),
			db.Redweet.RedweetOf.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// If no such redweet exists, return
	if len(user.Redweets()) == 0 {
		return nil, nil
	}

	// Remove the Redweet from the post
	_, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Update(
		db.Dweet.RedweetCount.Decrement(1),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	_, err = client.Redweet.FindUnique(
		db.Redweet.DbID.Equals(user.Redweets()[0].DbID),
	).Delete().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &user.Redweets()[0], err
}
