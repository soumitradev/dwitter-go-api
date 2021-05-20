package main

import (
	"errors"

	"github.com/graphql-go/graphql"
)

// Create a handler that handles graphql queries
var queryHandler = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"dweet": &graphql.Field{
				Type:        dweetSchema,
				Description: "Get dweet by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},

					// TODO: Pagination
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						id, idPresent := params.Args["id"].(string)
						numReplies, numPresent := params.Args["repliesToFetch"].(int)
						if idPresent && numPresent {
							post, err := AuthGetPost(id, numReplies, data["username"].(string))
							return post, err
						}
					} else {
						id, idPresent := params.Args["id"].(string)
						numReplies, numPresent := params.Args["repliesToFetch"].(int)
						if idPresent && numPresent {
							post, err := NoAuthGetPost(id, numReplies)
							return post, err
						}
					}

					return nil, errors.New("param \"id\" or missing")
				},
			},
			// TODO: Advanced search
			"dweets": &graphql.Field{
				Type:        graphql.NewList(dweetSchema),
				Description: "Search dweets by content",
				Args: graphql.FieldConfigArgument{
					"text": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},

					// TODO: Pagination
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						txt, txtPresent := params.Args["text"].(string)
						num, numPresent := params.Args["dweetsToFetch"].(int)
						numReplies, numRepliesPresent := params.Args["repliesToFetch"].(int)
						if txtPresent && numPresent && numRepliesPresent {
							posts, err := AuthSearchPosts(txt, num, numReplies, data["username"].(string))
							return posts, err
						}
					} else {
						txt, txtPresent := params.Args["text"].(string)
						num, numPresent := params.Args["dweetsToFetch"].(int)
						numReplies, numRepliesPresent := params.Args["repliesToFetch"].(int)
						if txtPresent && numPresent && numRepliesPresent {
							posts, err := NoAuthSearchPosts(txt, num, numReplies)
							return posts, err
						}
					}

					return nil, errors.New("param \"text\" or missing")
				},
			},
			"user": &graphql.Field{
				Type:        userSchema,
				Description: "Get user by username",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					// TODO: Pagination
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						username, userPresent := params.Args["username"].(string)
						numReplies, numPresent := params.Args["dweetsToFetch"].(int)
						if userPresent && numPresent {
							user, err := AuthGetUser(username, numReplies, data["username"].(string))
							return user, err
						}
					} else {

						username, userPresent := params.Args["username"].(string)
						numDweets, numPresent := params.Args["dweetsToFetch"].(int)
						if userPresent && numPresent {
							user, err := NoAuthGetUser(username, numDweets)
							return user, err
						}
					}

					return nil, errors.New("param \"username\" missing")
				},
			},
			// TODO: Pagination
			"likedDweets": &graphql.Field{
				Type:        graphql.NewList(dweetSchema),
				Description: "Get liked dweets of authenticated user",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					// TODO: Pagination
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						numDweets, dweetPresent := params.Args["numberToFetch"].(int)
						numReplies, repliesPresent := params.Args["repliesToFetch"].(int)
						if dweetPresent && repliesPresent {
							post, err := FetchLikedDweets(data["username"].(string), numDweets, numReplies)
							return post, err
						}
					}

					return nil, errors.New("Unauthorized")
				},
			},
			// TODO: Pagination
			"followers": &graphql.Field{
				Type:        graphql.NewList(userSchema),
				Description: "Get followers of authenticated user",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					// TODO: Pagination
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						numUsers, usersPresent := params.Args["numberToFetch"].(int)
						numDweets, dweetsPresent := params.Args["dweetsToFetch"].(int)
						if usersPresent && dweetsPresent {
							post, err := FetchFollowers(data["username"].(string), numUsers, numDweets)
							return post, err
						}
					}

					return nil, errors.New("Unauthorized")
				},
			},
			// TODO: Pagination
			"following": &graphql.Field{
				Type:        graphql.NewList(userSchema),
				Description: "Get users that authenticated user follows",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					// TODO: Pagination
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}
					if isAuth {
						numUsers, usersPresent := params.Args["numberToFetch"].(int)
						numDweets, dweetsPresent := params.Args["dweetsToFetch"].(int)
						if usersPresent && dweetsPresent {
							post, err := FetchFollowing(data["username"].(string), numUsers, numDweets)
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
				Type:        userSchema,
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
					user, err := SignUpUser(
						params.Args["username"].(string),
						params.Args["password"].(string),
						params.Args["firstName"].(string),
						params.Args["lastName"].(string),
						params.Args["bio"].(string),
						params.Args["email"].(string),
					)
					if err != nil {
						return nil, err
					}

					return user, nil
				},
			},
			"createDweet": &graphql.Field{
				Type:        dweetSchema,
				Description: "Create a dweet authored by authenticated user",
				Args: graphql.FieldConfigArgument{
					"body": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"media": &graphql.ArgumentConfig{
						Type:         graphql.NewList(graphql.String),
						DefaultValue: []string{},
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create dweet, and return formatted
						body, bodyPresent := params.Args["body"].(string)
						media, mediaPresent := params.Args["media"].([]string)
						if bodyPresent && mediaPresent {
							dweet, err := AuthCreateDweet(body, data["username"].(string), media)
							return dweet, err
						}
						return nil, errors.New("invalid request, \"body\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"createReply": &graphql.Field{
				Type:        dweetSchema,
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
						DefaultValue: []string{},
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create a reply to a dweet, and return formatted
						originalID, idPresent := params.Args["id"].(string)
						body, bodyPresent := params.Args["body"].(string)
						media, mediaPresent := params.Args["media"].([]string)
						if bodyPresent && mediaPresent && idPresent {
							dweet, err := AuthCreateReply(originalID, body, data["username"].(string), media)
							return dweet, err
						}
						return nil, errors.New("invalid request, \"id\", or \"body\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"redweet": &graphql.Field{
				Type:        redweetSchema,
				Description: "Create a redweet of a dweet by authenticated user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create a reply to a dweet, and return formatted
						originalID, idPresent := params.Args["id"].(string)
						if idPresent {
							redweet, err := AuthCreateRedweet(originalID, data["username"].(string))
							return redweet, err
						}
						return nil, errors.New("invalid request, \"id\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"follow": &graphql.Field{
				Type:        userSchema,
				Description: "Make authenticated user follow another user",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					// TODO: Pagination
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						username, userPresent := params.Args["username"].(string)
						dweetsToFetch, dweetsPresent := params.Args["dweetsToFetch"].(int)
						if userPresent && dweetsPresent {
							user, err := AuthFollow(username, data["username"].(string), dweetsToFetch)
							return user, err
						}
						return nil, errors.New("invalid request, \"username\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"like": &graphql.Field{
				Type:        dweetSchema,
				Description: "Add authenticated user like to a dweet",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					// TODO: Pagination
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user like dweet, and return formatted
						id, idPresent := params.Args["id"].(string)
						repliesToFetch, repliesPresent := params.Args["repliesToFetch"].(int)
						if idPresent && repliesPresent {
							dweet, err := AuthLike(id, data["username"].(string), repliesToFetch)
							return dweet, err
						}
						return nil, errors.New("invalid request, \"id\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"unlike": &graphql.Field{
				Type:        dweetSchema,
				Description: "Remove authenticated user's like from dweet",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					// TODO: Pagination
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user like dweet, and return formatted
						id, idPresent := params.Args["id"].(string)
						repliesToFetch, repliesPresent := params.Args["repliesToFetch"].(int)
						if idPresent && repliesPresent {
							dweet, err := AuthUnlike(id, data["username"].(string), repliesToFetch)
							return dweet, err
						}
						return nil, errors.New("invalid request, \"id\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"unfollow": &graphql.Field{
				Type:        userSchema,
				Description: "Make authenticated user unfollow another user",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					// TODO: Pagination
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						username, userPresent := params.Args["username"].(string)
						dweetsToFetch, numPresent := params.Args["dweetsToFetch"].(int)
						if userPresent && numPresent {
							user, err := AuthUnfollow(username, data["username"].(string), dweetsToFetch)
							return user, err
						}
						return nil, errors.New("invalid request, \"username\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"editDweet": &graphql.Field{
				Type:        dweetSchema,
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
					// TODO: Pagination
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Edit dweet, and return formatted
						id, idPresent := params.Args["id"].(string)
						body, bodyPresent := params.Args["body"].(string)
						media, mediaPresent := params.Args["media"].([]interface{})
						repliesToFetch, numPresent := params.Args["repliesToFetch"].(int)
						if bodyPresent && mediaPresent && idPresent && numPresent {
							mediaList := []string{}
							for _, link := range media {
								mediaList = append(mediaList, link.(string))
							}
							dweet, err := AuthUpdateDweet(id, data["username"].(string), body, mediaList, repliesToFetch)
							return dweet, err
						}
						return nil, errors.New("invalid request, \"body\" or \"media\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"editUser": &graphql.Field{
				Type:        userSchema,
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
					// TODO: Pagination
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					// TODO: Pagination
					"followersToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
					// TODO: Pagination
					"followingToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Edit dweet, and return formatted
						firstName, firstPresent := params.Args["firstName"].(string)
						lastName, lastPresent := params.Args["lastName"].(string)
						email, emailPresent := params.Args["email"].(string)
						bio, bioPresent := params.Args["email"].(string)
						dweetsToFetch, dweetsPresent := params.Args["dweetsToFetch"].(int)
						followersToFetch, followersPresent := params.Args["followersToFetch"].(int)
						followingToFetch, followingPresent := params.Args["followingToFetch"].(int)
						if firstPresent && lastPresent && emailPresent && bioPresent && dweetsPresent && followersPresent && followingPresent {
							user, err := AuthUpdateUser(data["username"].(string), firstName, lastName, email, bio, dweetsToFetch, followersToFetch, followingToFetch)
							return user, err
						}
						return nil, errors.New("invalid request, \"body\" or \"media\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"deleteDweet": &graphql.Field{
				Type:        dweetSchema,
				Description: "Delete dweet authored by user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					// TODO: Pagination
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						id, idPresent := params.Args["id"].(string)
						repliesToFetch, repliesPresent := params.Args["repliesToFetch"].(int)
						if idPresent && repliesPresent {
							dweet, err := AuthDeleteDweet(id, data["username"].(string), repliesToFetch)
							return dweet, err
						}
						return nil, errors.New("invalid request, \"id\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"unredweet": &graphql.Field{
				Type:        redweetSchema,
				Description: "Unredweet a dweet",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						id, present := params.Args["id"].(string)
						if present {
							redweet, err := AuthDeleteRedweet(id, data["username"].(string))
							return redweet, err
						}
						return nil, errors.New("invalid request, \"id\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			// "authTest": &graphql.Field{
			// 	Type:        graphql.String,
			// 	Description: "Log into Dwitter",
			// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// 		// tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
			// 		// data, _, err := VerifyToken(tokenString)
			// 		// if err != nil {
			// 		// 	return nil, err
			// 		// }
			// 		// return fmt.Sprintf("Username: %v", data["username"]), err
			// 		return nil, nil
			// 	},
			// },
		},
	},
)

// Create a handler that handles graphql mutations
var subscriptionHandler = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"feed": &graphql.Field{
				Type: feedObjectSchema,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyAccessToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						obj, err := GetFeed(data["username"].(string))
						return obj, err
					}

					return nil, errors.New("Unauthorized")
				},
			},
		},
	},
)

// Create schema from handlers
var schema, SchemaError = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:        queryHandler,
		Mutation:     mutationHandler,
		Subscription: subscriptionHandler,
	},
)
