package database

import (
	"fmt"

	"github.com/soumitradev/Dwitter/backend/common"
	"github.com/soumitradev/Dwitter/backend/prisma/db"
	"github.com/soumitradev/Dwitter/backend/schema"
	"github.com/soumitradev/Dwitter/backend/util"
)

// Get User's liked dweets
func GetLikedDweets(userID string, numberToFetch int, numOffset int, repliesToFetch int, replyOffset int) ([]schema.DweetType, error) {
	// Validate params
	err := common.Validate.Var(userID, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	err = common.Validate.Var(numOffset, "gte=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	err = common.Validate.Var(replyOffset, "gte=0")
	if err != nil {
		return []schema.DweetType{}, err
	}

	var user *db.UserModel

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
					).OrderBy(
						db.Dweet.LikeCount.Order(db.DESC),
					),
					db.Dweet.LikeUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
					db.Dweet.RedweetUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
				).OrderBy(
					db.Dweet.PostedAt.Order(db.DESC),
				),
				db.User.Following.Fetch().OrderBy(
					db.User.FollowerCount.Order(db.DESC),
				),
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
					).OrderBy(
						db.Dweet.LikeCount.Order(db.DESC),
					).Take(repliesToFetch).Skip(replyOffset),
					db.Dweet.LikeUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
					db.Dweet.RedweetUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
				).OrderBy(
					db.Dweet.PostedAt.Order(db.DESC),
				),
				db.User.Following.Fetch().OrderBy(
					db.User.FollowerCount.Order(db.DESC),
				),
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
					).OrderBy(
						db.Dweet.LikeCount.Order(db.DESC),
					),
					db.Dweet.LikeUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
					db.Dweet.RedweetUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
				).OrderBy(
					db.Dweet.PostedAt.Order(db.DESC),
				).Take(numberToFetch).Skip(numOffset),
				db.User.Following.Fetch().OrderBy(
					db.User.FollowerCount.Order(db.DESC),
				),
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
					).OrderBy(
						db.Dweet.LikeCount.Order(db.DESC),
					).Take(repliesToFetch).Skip(replyOffset),
					db.Dweet.LikeUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
					db.Dweet.RedweetUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
				).OrderBy(
					db.Dweet.PostedAt.Order(db.DESC),
				).Take(numberToFetch).Skip(numOffset),
				db.User.Following.Fetch().OrderBy(
					db.User.FollowerCount.Order(db.DESC),
				),
			).Exec(common.BaseCtx)
		}
	}
	if err == db.ErrNotFound {
		return []schema.DweetType{}, fmt.Errorf("user not found: %v", err)
	}
	if err != nil {
		return []schema.DweetType{}, fmt.Errorf("internal server error: %v", err)
	}

	knownUsers := user.Following()
	knownUsers = append(knownUsers, *user)

	// Add common likes and return formatted
	var liked []schema.DweetType
	for _, dweet := range user.LikedDweets() {
		likes := dweet.LikeUsers()
		redweetUsers := dweet.RedweetUsers()

		// Find known people that liked thw dweet
		mutualLikes := util.HashIntersectUsers(likes, knownUsers)
		mutualRedweets := util.HashIntersectUsers(redweetUsers, knownUsers)

		liked = append(liked, schema.FormatAsDweetType(&dweet, mutualLikes, mutualRedweets))
	}
	return liked, err
}

// TODO: GetDweets, GetRedweets, GetRedweetedDweets, GetFeedObjects

// Get feed for authenticated user
func GetFeed(username string) ([]interface{}, error) {
	// Validate params
	err := common.Validate.Var(username, "required,alphanum,lte=20,gt=0")
	if err != nil {
		return []interface{}{}, err
	}

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
				).OrderBy(
					db.Dweet.LikeCount.Order(db.DESC),
				),
				db.Dweet.ReplyTo.Fetch().With(
					db.Dweet.Author.Fetch(),
				),
				db.Dweet.LikeUsers.Fetch().OrderBy(
					db.User.FollowerCount.Order(db.DESC),
				),
				db.Dweet.RedweetUsers.Fetch().OrderBy(
					db.User.FollowerCount.Order(db.DESC),
				),
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
					db.Dweet.LikeUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
					db.Dweet.RedweetUsers.Fetch().OrderBy(
						db.User.FollowerCount.Order(db.DESC),
					),
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

	knownUsers := user.Following()
	knownUsers = append(knownUsers, *user)

	var formatted []interface{}
	for _, post := range merged {
		var npost interface{}
		if dweet, ok := post.(db.DweetModel); ok {
			likes := util.HashIntersectUsers(dweet.LikeUsers(), knownUsers)
			redweets := util.HashIntersectUsers(dweet.RedweetUsers(), knownUsers)
			npost = schema.FormatAsDweetType(&dweet, likes, redweets)
		}
		if redweet, ok := post.(db.RedweetModel); ok {
			npost = schema.FormatAsRedweetType(&redweet)
		}
		formatted = append(formatted, npost)
	}
	return formatted, err
}
