// Package gql provides useful graphql API functionality
package gql

import (
	"errors"

	"dwitter_go_graphql/auth"
	"dwitter_go_graphql/database"
	"dwitter_go_graphql/schema"

	"github.com/graphql-go/graphql"
)

// Create a handler that handles graphql queries
var queryHandler = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"dweet": &graphql.Field{
				Type:        schema.DweetSchema,
				Description: "Get dweet by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						id, idPresent := params.Args["id"].(string)
						numReplies, numPresent := params.Args["repliesToFetch"].(int)
						replyOffset, offsetPresent := params.Args["repliesOffset"].(int)
						if idPresent && numPresent && offsetPresent {
							post, err := database.GetPost(id, numReplies, replyOffset, data["username"].(string))
							return post, err
						}
					} else {
						id, idPresent := params.Args["id"].(string)
						numReplies, numPresent := params.Args["repliesToFetch"].(int)
						replyOffset, offsetPresent := params.Args["repliesOffset"].(int)
						if idPresent && numPresent && offsetPresent {
							post, err := database.GetPostUnauth(id, numReplies, replyOffset)
							return post, err
						}
					}

					return nil, errors.New("param \"id\" or missing")
				},
			},
			// TODO: Advanced search
			"dweets": &graphql.Field{
				Type:        graphql.NewList(schema.DweetSchema),
				Description: "Search dweets by content",
				Args: graphql.FieldConfigArgument{
					"text": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						txt, txtPresent := params.Args["text"].(string)
						num, numPresent := params.Args["dweetsToFetch"].(int)
						numOffset, numOffsetPresent := params.Args["dweetsOffset"].(int)
						numReplies, numRepliesPresent := params.Args["repliesToFetch"].(int)
						replyOffset, replyOffsetPresent := params.Args["repliesOffset"].(int)
						if txtPresent && numPresent && numOffsetPresent && numRepliesPresent && replyOffsetPresent {
							posts, err := database.SearchPosts(txt, num, numOffset, numReplies, replyOffset, data["username"].(string))
							return posts, err
						}
					} else {
						txt, txtPresent := params.Args["text"].(string)
						num, numPresent := params.Args["dweetsToFetch"].(int)
						numOffset, numOffsetPresent := params.Args["dweetsOffset"].(int)
						numReplies, numRepliesPresent := params.Args["repliesToFetch"].(int)
						replyOffset, replyOffsetPresent := params.Args["repliesOffset"].(int)
						if txtPresent && numPresent && numOffsetPresent && numRepliesPresent && replyOffsetPresent {
							posts, err := database.SearchPostsUnauth(txt, num, numOffset, numReplies, replyOffset)
							return posts, err
						}
					}

					return nil, errors.New("param \"text\" missing")
				},
			},
			"user": &graphql.Field{
				Type:        schema.UserSchema,
				Description: "Get user by username",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						username, userPresent := params.Args["username"].(string)
						numDweets, numPresent := params.Args["dweetsToFetch"].(int)
						dweetOffset, dweetOffsetPresent := params.Args["dweetsOffset"].(int)
						if userPresent && numPresent && dweetOffsetPresent {
							user, err := database.GetUser(username, numDweets, dweetOffset, data["username"].(string))
							return user, err
						}
					} else {
						username, userPresent := params.Args["username"].(string)
						numDweets, numPresent := params.Args["dweetsToFetch"].(int)
						dweetOffset, dweetOffsetPresent := params.Args["dweetsOffset"].(int)
						if userPresent && numPresent && dweetOffsetPresent {
							user, err := database.GetUserUnauth(username, numDweets, dweetOffset)
							return user, err
						}
					}

					return nil, errors.New("param \"username\" missing")
				},
			},
			// TODO: Advanced search
			"users": &graphql.Field{
				Type:        graphql.NewList(schema.UserSchema),
				Description: "Search users by username",
				Args: graphql.FieldConfigArgument{
					"text": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"numberOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						txt, txtPresent := params.Args["text"].(string)
						num, numPresent := params.Args["numberToFetch"].(int)
						numOffset, numOffsetPresent := params.Args["numberOffset"].(int)
						numDweets, numDweetsPresent := params.Args["dweetsToFetch"].(int)
						dweetOffset, dweetOffsetPresent := params.Args["dweetsOffset"].(int)
						if txtPresent && numPresent && numOffsetPresent && numDweetsPresent && dweetOffsetPresent {
							posts, err := database.SearchUsers(txt, num, numOffset, numDweets, dweetOffset, data["username"].(string))
							return posts, err
						}
					} else {
						txt, txtPresent := params.Args["text"].(string)
						num, numPresent := params.Args["numberToFetch"].(int)
						numOffset, numOffsetPresent := params.Args["numberOffset"].(int)
						numDweets, numDweetsPresent := params.Args["dweetsToFetch"].(int)
						dweetOffset, dweetOffsetPresent := params.Args["dweetsOffset"].(int)
						if txtPresent && numPresent && numOffsetPresent && numDweetsPresent && dweetOffsetPresent {
							posts, err := database.SearchUsersUnauth(txt, num, numOffset, numDweets, dweetOffset)
							return posts, err
						}
					}

					return nil, errors.New("param \"text\" missing")
				},
			},
			"likedDweets": &graphql.Field{
				Type:        graphql.NewList(schema.DweetSchema),
				Description: "Get liked dweets of authenticated user",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"numberOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						numDweets, dweetPresent := params.Args["numberToFetch"].(int)
						numOffset, numOffsetPresent := params.Args["numberOffset"].(int)
						numReplies, repliesPresent := params.Args["repliesToFetch"].(int)
						replyOffset, replyOffsetPresent := params.Args["repliesOffset"].(int)
						if dweetPresent && repliesPresent && numOffsetPresent && replyOffsetPresent {
							post, err := database.GetLikedDweets(data["username"].(string), numDweets, numOffset, numReplies, replyOffset)
							return post, err
						}
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"followers": &graphql.Field{
				Type:        graphql.NewList(schema.UserSchema),
				Description: "Get followers of authenticated user",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"numberOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						numUsers, usersPresent := params.Args["numberToFetch"].(int)
						numOffset, usersOffsetPresent := params.Args["numberOffset"].(int)
						numDweets, dweetsPresent := params.Args["dweetsToFetch"].(int)
						dweetsOffset, dweetsOffsetPresent := params.Args["dweetsOffset"].(int)
						if usersPresent && dweetsPresent && usersOffsetPresent && dweetsOffsetPresent {
							post, err := database.GetFollowers(data["username"].(string), numUsers, numOffset, numDweets, dweetsOffset)
							return post, err
						}
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"following": &graphql.Field{
				Type:        graphql.NewList(schema.UserSchema),
				Description: "Get users that authenticated user follows",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"numberOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}
					if isAuth {
						numUsers, usersPresent := params.Args["numberToFetch"].(int)
						numOffset, usersOffsetPresent := params.Args["numberOffset"].(int)
						numDweets, dweetsPresent := params.Args["dweetsToFetch"].(int)
						dweetsOffset, dweetsOffsetPresent := params.Args["dweetsOffset"].(int)
						if usersPresent && dweetsPresent && usersOffsetPresent && dweetsOffsetPresent {
							post, err := database.GetFollowing(data["username"].(string), numUsers, numOffset, numDweets, dweetsOffset)
							return post, err
						}
					}

					return nil, errors.New("Unauthorized")
				},
			},
		},
	},
)

// Create a handler that handles graphql mutations
var mutationHandler = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type:        schema.UserSchema,
				Description: "Create a user",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"firstName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"lastName": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"bio": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					username, usernamePresent := params.Args["username"].(string)
					password, passwordPresent := params.Args["password"].(string)
					firstName, firstPresent := params.Args["firstName"].(string)
					lastName, lastPresent := params.Args["lastName"].(string)
					bio, bioPresent := params.Args["bio"].(string)
					email, emailPresent := params.Args["email"].(string)
					if usernamePresent && passwordPresent && firstPresent && lastPresent && bioPresent && emailPresent {
						user, err := database.SignUpUser(username, password, firstName, lastName, bio, email)
						return user, err
					}
					return nil, errors.New("invalid request: missing argument")
				},
			},
			"createDweet": &graphql.Field{
				Type:        schema.DweetSchema,
				Description: "Create a dweet authored by authenticated user",
				Args: graphql.FieldConfigArgument{
					"body": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"media": &graphql.ArgumentConfig{
						Type:         graphql.NewList(graphql.String),
						DefaultValue: []interface{}{},
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create dweet, and return formatted
						body, bodyPresent := params.Args["body"].(string)
						media, mediaPresent := params.Args["media"].([]interface{})
						if bodyPresent && mediaPresent {
							mediaList := []string{}
							for _, link := range media {
								mediaList = append(mediaList, link.(string))
							}
							dweet, err := database.NewDweet(body, data["username"].(string), mediaList)
							return dweet, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"createReply": &graphql.Field{
				Type:        schema.DweetSchema,
				Description: "Create a reply to a dweet by authenticated user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"body": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"media": &graphql.ArgumentConfig{
						Type:         graphql.NewList(graphql.String),
						DefaultValue: []interface{}{},
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create a reply to a dweet, and return formatted
						originalID, idPresent := params.Args["id"].(string)
						body, bodyPresent := params.Args["body"].(string)
						media, mediaPresent := params.Args["media"].([]interface{})
						if bodyPresent && mediaPresent && idPresent {
							mediaList := []string{}
							for _, link := range media {
								mediaList = append(mediaList, link.(string))
							}
							dweet, err := database.NewReply(originalID, body, data["username"].(string), mediaList)
							return dweet, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"redweet": &graphql.Field{
				Type:        schema.RedweetSchema,
				Description: "Create a redweet of a dweet by authenticated user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create a redweet, and return formatted
						originalID, idPresent := params.Args["id"].(string)
						if idPresent {
							redweet, err := database.Redweet(originalID, data["username"].(string))
							return redweet, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"follow": &graphql.Field{
				Type:        schema.UserSchema,
				Description: "Make authenticated user follow another user",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						username, userPresent := params.Args["username"].(string)
						dweetsToFetch, dweetsPresent := params.Args["dweetsToFetch"].(int)
						dweetOffset, offsetPresent := params.Args["dweetsOffset"].(int)

						if username == data["username"].(string) {
							return nil, errors.New("can't follow self")
						}

						if userPresent && dweetsPresent && offsetPresent {
							user, err := database.Follow(username, data["username"].(string), dweetsToFetch, dweetOffset)
							return user, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"like": &graphql.Field{
				Type:        schema.DweetSchema,
				Description: "Add authenticated user like to a dweet",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user like dweet, and return formatted
						id, idPresent := params.Args["id"].(string)
						repliesToFetch, repliesPresent := params.Args["repliesToFetch"].(int)
						replyOffset, offsetPresent := params.Args["repliesOffset"].(int)
						if idPresent && repliesPresent && offsetPresent {
							dweet, err := database.Like(id, data["username"].(string), repliesToFetch, replyOffset)
							return dweet, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"unlike": &graphql.Field{
				Type:        schema.DweetSchema,
				Description: "Remove authenticated user's like from dweet",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user unlike dweet, and return formatted
						id, idPresent := params.Args["id"].(string)
						repliesToFetch, repliesPresent := params.Args["repliesToFetch"].(int)
						replyOffset, offsetPresent := params.Args["repliesOffset"].(int)
						if idPresent && repliesPresent && offsetPresent {
							dweet, err := database.Unlike(id, data["username"].(string), repliesToFetch, replyOffset)
							return dweet, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"unfollow": &graphql.Field{
				Type:        schema.UserSchema,
				Description: "Make authenticated user unfollow another user",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user unfollow the other user, and return formatted
						username, userPresent := params.Args["username"].(string)
						dweetsToFetch, numPresent := params.Args["dweetsToFetch"].(int)
						dweetOffset, offsetPresent := params.Args["dweetsOffset"].(int)

						if username == data["username"].(string) {
							return nil, errors.New("can't unfollow self")
						}

						if userPresent && numPresent && offsetPresent {
							user, err := database.Unfollow(username, data["username"].(string), dweetsToFetch, dweetOffset)
							return user, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"editDweet": &graphql.Field{
				Type:        schema.DweetSchema,
				Description: "Edit a dweet authored by authenticated user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"body": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"media": &graphql.ArgumentConfig{
						Type:         graphql.NewList(graphql.String),
						DefaultValue: []interface{}{},
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Edit dweet, and return formatted
						id, idPresent := params.Args["id"].(string)
						body, bodyPresent := params.Args["body"].(string)
						media, mediaPresent := params.Args["media"].([]interface{})
						repliesToFetch, numPresent := params.Args["repliesToFetch"].(int)
						replyOffset, offsetPresent := params.Args["repliesOffset"].(int)
						if bodyPresent && mediaPresent && idPresent && numPresent && offsetPresent {
							mediaList := []string{}
							for _, link := range media {
								mediaList = append(mediaList, link.(string))
							}
							dweet, err := database.UpdateDweet(id, data["username"].(string), body, mediaList, repliesToFetch, replyOffset)
							return dweet, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"editUser": &graphql.Field{
				Type:        schema.UserSchema,
				Description: "Edit authenticated user",
				Args: graphql.FieldConfigArgument{
					"firstName": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"lastName": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"email": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"bio": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"pfpURL": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"dweetsOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"followersToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"followersOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"followingToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"followingOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Edit user, and return formatted
						firstName, firstPresent := params.Args["firstName"].(string)
						lastName, lastPresent := params.Args["lastName"].(string)
						email, emailPresent := params.Args["email"].(string)
						bio, bioPresent := params.Args["email"].(string)
						PfpUrl, pfpPresent := params.Args["pfpURL"].(string)
						dweetsToFetch, dweetsPresent := params.Args["dweetsToFetch"].(int)
						dweetOffset, dweetOffsetPresent := params.Args["dweetsOffset"].(int)
						followersToFetch, followersPresent := params.Args["followersToFetch"].(int)
						followersOffset, followersOffsetPresent := params.Args["followersOffset"].(int)
						followingToFetch, followingPresent := params.Args["followingToFetch"].(int)
						followingOffset, followingOffsetPresent := params.Args["followingOffset"].(int)
						if firstPresent && lastPresent && emailPresent && bioPresent && pfpPresent && dweetsPresent && dweetOffsetPresent && followersPresent && followersOffsetPresent && followingPresent && followingOffsetPresent {
							user, err := database.UpdateUser(data["username"].(string), firstName, lastName, email, bio, PfpUrl, dweetsToFetch, dweetOffset, followersToFetch, followersOffset, followingToFetch, followingOffset)
							return user, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"deleteDweet": &graphql.Field{
				Type:        schema.DweetSchema,
				Description: "Delete dweet authored by user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesOffset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Delete dweet, and return formatted
						id, idPresent := params.Args["id"].(string)
						repliesToFetch, repliesPresent := params.Args["repliesToFetch"].(int)
						replyOffset, offsetPresent := params.Args["repliesOffset"].(int)
						if idPresent && repliesPresent && offsetPresent {
							dweet, err := database.DeleteDweet(id, data["username"].(string), repliesToFetch, replyOffset)
							return dweet, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"unredweet": &graphql.Field{
				Type:        schema.RedweetSchema,
				Description: "Unredweet a dweet",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user unredweet the dweet, and return formatted
						id, present := params.Args["id"].(string)
						if present {
							redweet, err := database.DeleteRedweet(id, data["username"].(string))
							return redweet, err
						}
						return nil, errors.New("invalid request: missing argument")
					}

					return nil, errors.New("Unauthorized")
				},
			},
		},
	},
)

// Create a handler that handles graphql mutations
var subscriptionHandler = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"feed": &graphql.Field{
				Type: schema.FeedObjectSchema,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := auth.VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						obj, err := database.GetFeed(data["username"].(string))
						return obj, err
					}

					return nil, errors.New("Unauthorized")
				},
			},
		},
	},
)

// Create schema from handlers
var Schema, SchemaError = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:        queryHandler,
		Mutation:     mutationHandler,
		Subscription: subscriptionHandler,
	},
)
