package main

import (
	"context"
	"errors"
	"time"

	"dwitter_go_graphql/prisma/db"
)

var client *db.PrismaClient

var ctx context.Context

func ConnectDB() {
	// Connect to the database using prisma
	client = db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	ctx = context.Background()
}

func DisconnectDB() {
	// Disconnect from DB
	if err := client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}

func GetUser(userID string) (*db.UserModel, error) {
	// Get basic User data
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Exec(ctx)
	return user, err
}

func GetUserDweets(userID string, dweets_to_fetch int) (*db.UserModel, error) {
	// Get User data with dweets of user
	var user *db.UserModel
	var err error
	if dweets_to_fetch < 0 {
		user, err = client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Dweets.Fetch(),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		user, err = client.User.FindUnique(
			db.User.Username.Equals(userID),
		).With(
			db.User.Dweets.Fetch(),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func GetUserLikes(userID string) (*db.UserModel, error) {
	// Get User data with dweets that user liked
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.LikedDweets.Fetch(),
	).Exec(ctx)
	return user, err
}

func GetFollowers(userID string) (*db.UserModel, error) {
	// Get User data with followers of user
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Followers.Fetch(),
	).Exec(ctx)
	return user, err
}

func GetFollowing(userID string) (*db.UserModel, error) {
	// Get User data with users that user follows
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.Following.Fetch(),
	).Exec(ctx)
	return user, err
}

func GetPostBasic(postID string) (*db.DweetModel, error) {
	// Get Replies to post
	post, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Exec(ctx)
	return post, err
}

func GetPostReplies(postID string, replies_to_fetch int) (*db.DweetModel, error) {
	// Get Replies to post
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

func GetPostRedweets(postID string, redweets_to_fetch int) (*db.DweetModel, error) {
	// Get redweets of post
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

func NewUser(username, passwordHash, firstName, lastName, email, bio string) (*db.UserModel, error) {
	// Create a User
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

func NewDweet(body, authorID string, mediaLinks []string) (*db.DweetModel, error) {
	now := time.Now()
	// Create a Post
	createdPost, err := client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.ID.Set(genID(10)),
		db.Dweet.Author.Link(db.User.Username.Equals(authorID)),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).Exec(ctx)
	return createdPost, err
}

func NewLike(likedPostID, userID string) (*db.DweetModel, error) {
	// Check if user already liked this dweet
	myUser, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).With(
		db.User.LikedDweets.Fetch(
			db.Dweet.ID.Equals(likedPostID),
		).With(
			db.Dweet.Author.Fetch(),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	increment := 1

	// If yes, then skip liking the dweet
	if len(myUser.LikedDweets()) > 0 {
		increment = 0
	}

	// Create a Like on the post if not created already
	like, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(likedPostID),
	).Update(
		db.Dweet.LikeCount.Increment(increment),
		db.Dweet.LikeUsers.Link(
			db.User.Username.Equals(userID),
		),
	).Exec(ctx)
	if err != nil {
		return like, err
	}

	// Return dweet
	if len(myUser.LikedDweets()) > 0 {
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

func NewReply(originalPostID, userID, body string, mediaLinks []string) (*db.DweetModel, error) {
	// Create a reply to a post

	now := time.Now()
	// Get post and user
	post, err := GetPostBasic(originalPostID)
	if err != nil {
		return nil, err
	}

	user, err := GetUser(userID)
	if err != nil {
		return nil, err
	}

	// Create a Reply
	createdReply, err := client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.ID.Set(genID(10)),
		db.Dweet.Author.Link(db.User.Username.Equals(user.Username)),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.IsReply.Set(true),
		db.Dweet.ReplyTo.Link(
			db.Dweet.ID.Equals(post.ID),
		),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Update original Dweet to show reply
	_, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(post.ID),
	).Update(
		db.Dweet.ReplyDweets.Link(
			db.Dweet.ID.Equals(createdReply.ID),
		),
		db.Dweet.ReplyCount.Increment(1),
	).Exec(ctx)

	return createdReply, err
}

func NewRedweet(originalPostID, userID string) (*db.DweetModel, error) {
	// Create a new Redweet of a Dweet

	now := time.Now()
	// Get post and user
	post, err := GetPostBasic(originalPostID)
	if err != nil {
		return nil, err
	}

	user, err := GetUser(userID)
	if err != nil {
		return nil, err
	}

	// Create a Redweet
	createdRedweet, err := client.Dweet.CreateOne(
		db.Dweet.DweetBody.Set(post.DweetBody),
		db.Dweet.ID.Set(genID(10)),
		db.Dweet.Author.Link(db.User.Username.Equals(user.Username)),
		db.Dweet.Media.Set(post.Media),
		db.Dweet.IsRedweet.Set(true),
		db.Dweet.RedweetOf.Link(
			db.Dweet.ID.Equals(post.ID),
		),
		db.Dweet.PostedAt.Set(now),
		db.Dweet.LastUpdatedAt.Set(now),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Update original Dweet to show redweet
	_, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(post.ID),
	).Update(
		db.Dweet.RedweetDweets.Link(
			db.Dweet.ID.Equals(createdRedweet.ID),
		),
		db.Dweet.LastUpdatedAt.Set(now),
		db.Dweet.RedweetCount.Increment(1),
	).Exec(ctx)

	return createdRedweet, err
}

func NewFollower(followedID string, followerID string) (*db.UserModel, error) {
	// Create a follower relation

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

func UpdateUser(userID, username, firstName, lastName, email, bio string) (*db.UserModel, error) {
	// Update a user
	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Update(
		db.User.Username.Set(username),
		db.User.FirstName.Set(firstName),
		db.User.LastName.Set(lastName),
		db.User.Email.Set(email),
		db.User.Bio.Set(bio),
	).Exec(ctx)

	return user, err
}

func UpdateDweet(postID, body string, mediaLinks []string) (*db.DweetModel, error) {
	// Update a dweet
	post, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.RedweetDweets.Fetch(),
	).Update(
		db.Dweet.DweetBody.Set(body),
		db.Dweet.Media.Set(mediaLinks),
		db.Dweet.LastUpdatedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	redweets := post.RedweetDweets()
	for i := 0; i < len(redweets); i++ {
		_, err := client.Dweet.FindUnique(
			db.Dweet.ID.Equals(redweets[i].ID),
		).Update(
			db.Dweet.DweetBody.Set(body),
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Return updated post
	post, err = GetPostBasic(postID)
	if err != nil {
		return nil, err
	}

	return post, err
}

func DeleteFollower(followedID string, followerID string) (*db.UserModel, error) {
	// Delete a follower relation

	// Decrement the follower and following counts
	user, err := client.User.FindUnique(
		db.User.Username.Equals(followedID),
	).Update(
		db.User.FollowerCount.Decrement(1),
	).Exec(ctx)
	if err != nil {
		return user, err
	}

	userDBID1 := user.DbID

	other_user, err := client.User.FindUnique(
		db.User.Username.Equals(followerID),
	).Update(
		db.User.FollowingCount.Decrement(1),
	).Exec(ctx)
	if err != nil {
		return user, err
	}

	userDBID2 := other_user.DbID

	// Delete the follower-following relation using RAW SQL (prisma-client-go has almost no documentation on unlinking relations)
	delFollowerQuery := `DELETE FROM public."_Follow" WHERE "A"= $1 AND "B"= $2;`
	_, err = client.Prisma.ExecuteRaw(delFollowerQuery, userDBID1, userDBID2).Exec(ctx)

	return user, err
}

func DeleteLike(postID string, userID string) (*db.DweetModel, error) {
	// Remove a like from a post

	// Find the post and decrease its likes by 1
	post, err := client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Update(
		db.Dweet.LikeCount.Decrement(1),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	postDBID := post.DbID

	user, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	userDBID := user.DbID

	// Delete the like relation using RAW SQL (prisma-client-go has almost no documentation on unlinking relations)
	delLikeQuery := `DELETE FROM public."_Likes" WHERE "A"= $1 AND "B"= $2;`
	_, err = client.Prisma.ExecuteRaw(delLikeQuery, postDBID, userDBID).Exec(ctx)

	return post, err
}

func DeleteUser(userID string) (*db.UserModel, error) {
	// Delete a User

	// Get all the user's Dweets (we must delete these first since they depend on the User, and deleting the User first will render the DB invalid)
	user, err := GetUserDweets(userID, -1)
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
		if dweets[i].IsRedweet {
			DeleteRedweet(dweets[i].ID)
		} else {
			DeleteDweet(dweets[i].ID)
		}
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

func DeleteDweet(postID string) (*db.DweetModel, error) {
	// Delete a Dweet

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
		).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Delete all the dependent posts (this includes redweets and replies to the post) recursively using RAW SQL
	// We use RAW SQL here because prisma-go-client doesn't support cascade deletes yet:
	// Link: https://github.com/prisma/prisma-client-go/issues/201

	// Recursive SQL function with modifications from: https://stackoverflow.com/q/10381243
	delQuery := `with recursive all_posts (id, parentid1, parentid2, root_id) as (select t1.db_id, t1.original_reply_id as parentid1, t1.original_redweet_id as parentid2, t1.db_id as root_id from public."Dweet" t1 union all select c1.db_id, c1.original_reply_id as parentid1, c1.original_redweet_id as parentid2, p.root_id from public."Dweet" c1 join all_posts p on ((p.id = c1.original_reply_id) OR (p.id = c1.original_redweet_id)) ) DELETE FROM public."Dweet"  WHERE db_id IN (SELECT id FROM all_posts WHERE root_id = $1);`
	_, err = client.Prisma.ExecuteRaw(delQuery, post.DbID).Exec(ctx)

	return post, err
}

func DeleteRedweet(postID string) (*db.DweetModel, error) {
	// Remove a Redweet

	// Get all the replies to the redweet (these need to be deleted first since they depend on the root Redweet)
	post, err := GetPostBasic(postID)
	if err != nil {
		return nil, err
	}

	// Find the dweet that was redweeted
	id, exist := post.OriginalRedweetID()
	if !exist {
		return nil, errors.New("original Dweet not found")
	}

	// Remove the Redweet from the post
	_, err = client.Dweet.FindUnique(
		db.Dweet.ID.Equals(id),
	).Update(
		db.Dweet.RedweetCount.Decrement(1),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Delete all the dependent posts (this includes redweets and replies to the post) recursively using RAW SQL
	// We use RAW SQL here because prisma-go-client doesn't support cascade deletes yet:
	// Link: https://github.com/prisma/prisma-client-go/issues/201

	// Recursive SQL function with modifications from: https://stackoverflow.com/q/10381243
	delQuery := `with recursive all_posts (id, parentid1, parentid2, root_id) as (select t1.db_id, t1.original_reply_id as parentid1, t1.original_redweet_id as parentid2, t1.db_id as root_id from public."Dweet" t1 union all select c1.db_id, c1.original_reply_id as parentid1, c1.original_redweet_id as parentid2, p.root_id from public."Dweet" c1 join all_posts p on ((p.id = c1.original_reply_id) OR (p.id = c1.original_redweet_id)) ) DELETE FROM public."Dweet"  WHERE db_id IN (SELECT id FROM all_posts WHERE root_id = $1);`
	_, err = client.Prisma.ExecuteRaw(delQuery, post.DbID).Exec(ctx)

	return post, err
}
