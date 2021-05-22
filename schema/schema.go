package schema

import (
	"time"

	"github.com/graphql-go/graphql"
)

// Create Go structs and GraphQL objects for types

type BasicUserType struct {
	Username       string    `json:"username"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	Bio            string    `json:"bio"`
	PfpUrl         string    `json:"pfpURL"`
	FollowerCount  int       `json:"followerCount"`
	FollowingCount int       `json:"followingCount"`
	CreatedAt      time.Time `json:"createdAt"`
}

type UserType struct {
	Username       string           `json:"username"`
	FirstName      string           `json:"firstName"`
	LastName       string           `json:"lastName"`
	Email          string           `json:"email"`
	Bio            string           `json:"bio"`
	PfpUrl         string           `json:"pfpURL"`
	Dweets         []BasicDweetType `json:"dweets"`
	LikedDweets    []BasicDweetType `json:"likedDweets"`
	FollowerCount  int              `json:"followerCount"`
	Followers      []BasicUserType  `json:"followers"`
	FollowingCount int              `json:"followingCount"`
	Following      []BasicUserType  `json:"following"`
	CreatedAt      time.Time        `json:"createdAt"`
}

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
	Media           []string         `json:"media"`
}

type RedweetType struct {
	Author            BasicUserType  `json:"author"`
	AuthorID          string         `json:"authorID"`
	RedweetOf         BasicDweetType `json:"redweetOf"`
	OriginalRedweetID string         `json:"originalRedweetID"`
	RedweetTime       time.Time      `json:"redweetTime"`
}

var BasicUserSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BasicUser",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
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

var UserSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
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
			"media": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	},
)

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
