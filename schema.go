package main

import (
	"time"

	"github.com/graphql-go/graphql"
)

// Create Go structs and GraphQL objects for types

type BasicUserType struct {
	Username       string    `json:"username"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Bio            string    `json:"bio"`
	FollowerCount  int       `json:"follower_count"`
	FollowingCount int       `json:"following_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserType struct {
	Username       string           `json:"username"`
	FirstName      string           `json:"first_name"`
	LastName       string           `json:"last_name"`
	Email          string           `json:"email"`
	Bio            string           `json:"bio"`
	Dweets         []BasicDweetType `json:"dweets"`
	LikedDweets    []BasicDweetType `json:"liked_dweets"`
	FollowerCount  int              `json:"follower_count"`
	Followers      []BasicUserType  `json:"followers"`
	FollowingCount int              `json:"following_count"`
	Following      []BasicUserType  `json:"following"`
	CreatedAt      time.Time        `json:"created_at"`
}

type BasicDweetType struct {
	DweetBody         string        `json:"dweet_body"`
	ID                string        `json:"id"`
	Author            BasicUserType `json:"author"`
	AuthorID          string        `json:"author_id"`
	PostedAt          time.Time     `json:"posted_at"`
	LastUpdatedAt     time.Time     `json:"last_updated_at"`
	LikeCount         int           `json:"like_count"`
	IsReply           bool          `json:"is_reply"`
	OriginalReplyID   string        `json:"original_reply_id"`
	ReplyCount        int           `json:"reply_count"`
	IsRedweet         bool          `json:"is_redweet"`
	OriginalRedweetID string        `json:"original_redweet_id"`
	RedweetCount      int           `json:"redweet_count"`
	Media             []string      `json:"media"`
}

type DweetType struct {
	DweetBody         string           `json:"dweet_body"`
	ID                string           `json:"id"`
	Author            BasicUserType    `json:"author"`
	AuthorID          string           `json:"author_id"`
	PostedAt          time.Time        `json:"posted_at"`
	LastUpdatedAt     time.Time        `json:"last_updated_at"`
	LikeCount         int              `json:"like_count"`
	LikeUsers         []BasicUserType  `json:"like_users"`
	IsReply           bool             `json:"is_reply"`
	OriginalReplyID   string           `json:"original_reply_id"`
	ReplyTo           BasicDweetType   `json:"reply_to"`
	ReplyCount        int              `json:"reply_count"`
	ReplyDweets       []BasicDweetType `json:"reply_dweets"`
	IsRedweet         bool             `json:"is_redweet"`
	OriginalRedweetID string           `json:"original_redweet_id"`
	RedweetOf         BasicDweetType   `json:"redweet_of"`
	RedweetCount      int              `json:"redweet_count"`
	RedweetDweets     []BasicDweetType `json:"redweet_dweets"`
	Media             []string         `json:"media"`
}

var basicUserSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BasicUser",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"bio": &graphql.Field{
				Type: graphql.String,
			},
			"follower_count": &graphql.Field{
				Type: graphql.Int,
			},
			"following_count": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var userSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"bio": &graphql.Field{
				Type: graphql.String,
			},
			"dweets": &graphql.Field{
				Type: graphql.NewList(basicDweetSchema),
			},
			"liked_dweets": &graphql.Field{
				Type: graphql.NewList(basicDweetSchema),
			},
			"follower_count": &graphql.Field{
				Type: graphql.Int,
			},
			"followers": &graphql.Field{
				Type: graphql.NewList(basicUserSchema),
			},
			"following_count": &graphql.Field{
				Type: graphql.Int,
			},
			"following": &graphql.Field{
				Type: graphql.NewList(basicUserSchema),
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

var basicDweetSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BasicDweet",
		Fields: graphql.Fields{
			"dweet_body": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: basicUserSchema,
			},
			"author_id": &graphql.Field{
				Type: graphql.String,
			},
			"posted_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"last_updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"like_count": &graphql.Field{
				Type: graphql.Int,
			},
			"is_reply": &graphql.Field{
				Type: graphql.Boolean,
			},
			"original_reply_id": &graphql.Field{
				Type: graphql.String,
			},
			"reply_count": &graphql.Field{
				Type: graphql.Int,
			},
			"is_redweet": &graphql.Field{
				Type: graphql.Boolean,
			},
			"original_redweet_id": &graphql.Field{
				Type: graphql.String,
			},
			"redweet_count": &graphql.Field{
				Type: graphql.Int,
			},
			"media": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

var dweetSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Dweet",
		Fields: graphql.Fields{
			"dweet_body": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: basicUserSchema,
			},
			"author_id": &graphql.Field{
				Type: graphql.String,
			},
			"posted_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"last_updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"like_count": &graphql.Field{
				Type: graphql.Int,
			},
			"like_users": &graphql.Field{
				Type: graphql.NewList(basicUserSchema),
			},
			"is_reply": &graphql.Field{
				Type: graphql.Boolean,
			},
			"original_reply_id": &graphql.Field{
				Type: graphql.String,
			},
			"reply_to": &graphql.Field{
				Type: basicDweetSchema,
			},
			"reply_count": &graphql.Field{
				Type: graphql.Int,
			},
			"reply_dweets": &graphql.Field{
				Type: graphql.NewList(basicDweetSchema),
			},
			"is_redweet": &graphql.Field{
				Type: graphql.Boolean,
			},
			"original_redweet_id": &graphql.Field{
				Type: graphql.String,
			},
			"redweet_of": &graphql.Field{
				Type: basicDweetSchema,
			},
			"redweet_count": &graphql.Field{
				Type: graphql.Int,
			},
			"redweet_dweets": &graphql.Field{
				Type: graphql.NewList(basicDweetSchema),
			},
			"media": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)
