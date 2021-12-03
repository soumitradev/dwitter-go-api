// Package schema provides useful custom types and functions to format database objects into these types
package schema

import (
	"time"

	"github.com/graphql-go/graphql"
)

// Create Go structs and GraphQL objects for types

// A User object without any relation fields
type BasicUserType struct {
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Bio            string    `json:"bio"`
	PfpURL         string    `json:"pfpURL"`
	FollowerCount  int       `json:"followerCount"`
	FollowingCount int       `json:"followingCount"`
	CreatedAt      time.Time `json:"createdAt"`
}

// A User object
type UserType struct {
	Username        string           `json:"username"`
	Name            string           `json:"name"`
	Email           string           `json:"email"`
	Bio             string           `json:"bio"`
	PfpURL          string           `json:"pfpURL"`
	Dweets          []BasicDweetType `json:"dweets"`
	Redweets        []RedweetType    `json:"redweets"`
	FeedObjects     []interface{}    `json:"feedObjects"`
	RedweetedDweets []BasicDweetType `json:"redweetedDweets"`
	LikedDweets     []BasicDweetType `json:"likedDweets"`
	FollowerCount   int              `json:"followerCount"`
	Followers       []BasicUserType  `json:"followers"`
	FollowingCount  int              `json:"followingCount"`
	Following       []BasicUserType  `json:"following"`
	CreatedAt       time.Time        `json:"createdAt"`
}

// A Dweet object without any relation fields except for Author (a necessary relation field)
type BasicDweetType struct {
	DweetBody       string        `json:"dweetBody"`
	ID              string        `json:"id"`
	Author          BasicUserType `json:"author"`
	AuthorID        string        `json:"authorID"`
	PostedAt        time.Time     `json:"postedAt"`
	LastUpdatedAt   time.Time     `json:"lastUpdatedAt"`
	LikeCount       int           `json:"likeCount"`
	IsReply         bool          `json:"isReply"`
	OriginalReplyID string        `json:"originalReplyID"`
	ReplyCount      int           `json:"replyCount"`
	RedweetCount    int           `json:"redweetCount"`
	Media           []string      `json:"media"`
}

// A Dweet object
type DweetType struct {
	DweetBody       string           `json:"dweetBody"`
	ID              string           `json:"id"`
	Author          BasicUserType    `json:"author"`
	AuthorID        string           `json:"authorID"`
	PostedAt        time.Time        `json:"postedAt"`
	LastUpdatedAt   time.Time        `json:"lastUpdatedAt"`
	LikeCount       int              `json:"likeCount"`
	LikeUsers       []BasicUserType  `json:"likeUsers"`
	IsReply         bool             `json:"isReply"`
	OriginalReplyID string           `json:"originalReplyID"`
	ReplyTo         BasicDweetType   `json:"replyTo"`
	ReplyCount      int              `json:"replyCount"`
	ReplyDweets     []BasicDweetType `json:"replyDweets"`
	RedweetCount    int              `json:"redweetCount"`
	RedweetUsers    []BasicUserType  `json:"redweetUsers"`
	Media           []string         `json:"media"`
}

// A Redweet Object
type RedweetType struct {
	Author            BasicUserType  `json:"author"`
	AuthorID          string         `json:"authorID"`
	RedweetOf         BasicDweetType `json:"redweetOf"`
	OriginalRedweetID string         `json:"originalRedweetID"`
	RedweetTime       time.Time      `json:"redweetTime"`
}

// GraphQL schema for basic user
var BasicUserSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BasicUser",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"bio": &graphql.Field{
				Type: graphql.String,
			},
			"pfpURL": &graphql.Field{
				Type: graphql.String,
			},
			"followerCount": &graphql.Field{
				Type: graphql.Int,
			},
			"followingCount": &graphql.Field{
				Type: graphql.Int,
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

// GraphQL schema for user
var UserSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"bio": &graphql.Field{
				Type: graphql.String,
			},
			"pfpURL": &graphql.Field{
				Type: graphql.String,
			},
			"dweets": &graphql.Field{
				Type: graphql.NewList(BasicDweetSchema),
			},
			"redweets": &graphql.Field{
				Type: graphql.NewList(RedweetSchema),
			},
			"redweetedDweets": &graphql.Field{
				Type: graphql.NewList(BasicDweetSchema),
			},
			"feedObjects": &graphql.Field{
				Type: FeedObjectSchema,
			},
			"likedDweets": &graphql.Field{
				Type: graphql.NewList(BasicDweetSchema),
			},
			"followerCount": &graphql.Field{
				Type: graphql.Int,
			},
			"followers": &graphql.Field{
				Type: graphql.NewList(BasicUserSchema),
			},
			"followingCount": &graphql.Field{
				Type: graphql.Int,
			},
			"following": &graphql.Field{
				Type: graphql.NewList(BasicUserSchema),
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

// GraphQL schema for basic dweet
var BasicDweetSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BasicDweet",
		Fields: graphql.Fields{
			"dweetBody": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: BasicUserSchema,
			},
			"authorID": &graphql.Field{
				Type: graphql.String,
			},
			"postedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"lastUpdatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"likeCount": &graphql.Field{
				Type: graphql.Int,
			},
			"isReply": &graphql.Field{
				Type: graphql.Boolean,
			},
			"originalReplyID": &graphql.Field{
				Type: graphql.String,
			},
			"replyCount": &graphql.Field{
				Type: graphql.Int,
			},
			"redweetCount": &graphql.Field{
				Type: graphql.Int,
			},
			"media": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

// GraphQL schema for dweet
var DweetSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Dweet",
		Fields: graphql.Fields{
			"dweetBody": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: BasicUserSchema,
			},
			"authorID": &graphql.Field{
				Type: graphql.String,
			},
			"postedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"lastUpdatedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"likeCount": &graphql.Field{
				Type: graphql.Int,
			},
			"likeUsers": &graphql.Field{
				Type: graphql.NewList(BasicUserSchema),
			},
			"isReply": &graphql.Field{
				Type: graphql.Boolean,
			},
			"originalReplyID": &graphql.Field{
				Type: graphql.String,
			},
			"replyTo": &graphql.Field{
				Type: BasicDweetSchema,
			},
			"replyCount": &graphql.Field{
				Type: graphql.Int,
			},
			"replyDweets": &graphql.Field{
				Type: graphql.NewList(BasicDweetSchema),
			},
			"redweetCount": &graphql.Field{
				Type: graphql.Int,
			},
			"redweetUsers": &graphql.Field{
				Type: graphql.NewList(BasicUserSchema),
			},
			"media": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

// GraphQL schema for redweet
var RedweetSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Redweet",
		Fields: graphql.Fields{
			"author": &graphql.Field{
				Type: BasicUserSchema,
			},
			"authorID": &graphql.Field{
				Type: graphql.String,
			},
			"redweetOf": &graphql.Field{
				Type: BasicDweetSchema,
			},
			"originalRedweetID": &graphql.Field{
				Type: graphql.String,
			},
			"redweetTime": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

// A GraphQL union type for objects that may appear on a feed. i.e. Dweets and Redweets
var FeedObjectSchema = graphql.NewList(graphql.NewUnion(graphql.UnionConfig{
	Name:        "FeedObject",
	Types:       []*graphql.Object{DweetSchema, RedweetSchema},
	Description: "An object representing either a dweet or a redweet object.",
	ResolveType: func(params graphql.ResolveTypeParams) *graphql.Object {
		if _, ok := params.Value.(DweetType); ok {
			return DweetSchema
		} else {
			return RedweetSchema
		}
	},
}))
