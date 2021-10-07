// Package database provides some functions to interface with the posstgresql database
package database

import (
	"errors"

	"dwitter_go_graphql/common"
	"dwitter_go_graphql/prisma/db"
)

func init() {
	ConnectDB()
}

// Connect to the database using prisma
func ConnectDB() {
	common.Client = db.NewClient()
	if err := common.Client.Prisma.Connect(); err != nil {
		panic(err)
	}
}

// Disconnect from DB
func DisconnectDB() {
	if err := common.Client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}

// Delete a Dweet
func deleteDweet(postID string) (*db.DweetModel, error) {
	// Get all the replies to the post (these need to be deleted first since they depend on the root Dweet)
	post, err := common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Exec(common.BaseCtx)
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
		_, err := common.Client.Dweet.FindUnique(
			db.Dweet.ID.Equals(id),
		).Update(
			db.Dweet.ReplyCount.Decrement(1),
			db.Dweet.ReplyDweets.Unlink(
				db.Dweet.ID.Equals(postID),
			),
		).Exec(common.BaseCtx)
		if err != nil {
			return nil, err
		}
	}

	dweet, err := common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.RedweetDweets.Fetch().With(
			db.Redweet.Author.Fetch(),
		),
	).Exec(common.BaseCtx)
	if err != nil {
		return nil, err
	}

	for _, redweet := range dweet.RedweetDweets() {
		deleteRedweet(redweet.OriginalRedweetID, redweet.Author().Username)
	}

	dweet, err = common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).With(
		db.Dweet.ReplyDweets.Fetch(),
	).Exec(common.BaseCtx)
	if err != nil {
		return nil, err
	}
	for _, daughterDweet := range dweet.ReplyDweets() {
		deleteDweet(daughterDweet.ID)
	}

	_, err = common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Delete().Exec(common.BaseCtx)
	if err != nil {
		return nil, err
	}

	// The following comment block is kept as a homage to the great recursive SQL function that once resided here.
	// May the soul of this legendary query rest in peace. It was a honor to use you.

	// Delete all the dependent posts (this includes redweets and replies to the post) recursively using RAW SQL
	// We use RAW SQL here because prisma-go-client doesn't support cascade deletes yet:
	// Link: https://github.com/prisma/prisma-client-go/issues/201
	// Recursive SQL function with modifications from: https://stackoverflow.com/q/10381243
	// delQuery := `WITH RECURSIVE all_posts (id, parentid1, root_id) AS (SELECT t1.db_id, t1.original_reply_id AS parentid1, t1.db_id AS root_id FROM public."Dweet" t1 UNION ALL SELECT c1.db_id, c1.original_reply_id AS parentid1, p.root_id FROM public."Dweet" c1 JOIN all_posts p ON (p.id = c1.original_reply_id) ) DELETE FROM public."Dweet" WHERE db_id IN ( SELECT id FROM all_posts WHERE root_id = $1);`
	// _, err = common.Client.Prisma.ExecuteRaw(delQuery, post.DbID).Exec(common.BaseCtx)

	return post, err
}

// Remove a Redweet
func deleteRedweet(postID string, username string) (*db.RedweetModel, error) {
	// Get all the replies to the redweet (these need to be deleted first since they depend on the root Redweet)
	user, err := common.Client.User.FindUnique(
		db.User.Username.Equals(username),
	).With(
		db.User.Redweets.Fetch(
			db.Redweet.OriginalRedweetID.Equals(postID),
		).With(
			db.Redweet.Author.Fetch(),
			db.Redweet.RedweetOf.Fetch().With(
				db.Dweet.Author.Fetch(),
			),
		),
	).Exec(common.BaseCtx)
	if err != nil {
		return nil, err
	}

	// If no such redweet exists, return
	if len(user.Redweets()) == 0 {
		return nil, nil
	}

	// Remove the Redweet from the post
	_, err = common.Client.Dweet.FindUnique(
		db.Dweet.ID.Equals(postID),
	).Update(
		db.Dweet.RedweetCount.Decrement(1),
	).Exec(common.BaseCtx)
	if err != nil {
		return nil, err
	}

	_, err = common.Client.Redweet.FindUnique(
		db.Redweet.DbID.Equals(user.Redweets()[0].DbID),
	).Delete().Exec(common.BaseCtx)
	if err != nil {
		return nil, err
	}

	return &user.Redweets()[0], err
}
