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
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: -1,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
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

					return nil, errors.New("param \"id\" or \"repliesToFetch\" missing")
				},
			},
			"user": &graphql.Field{
				Type:        userSchema,
				Description: "Get user by username",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"dweetsToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: -1,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						username, userPresent := params.Args["username"].(string)
						numReplies, numPresent := params.Args["repliesToFetch"].(int)
						if userPresent && numPresent {
							user, err := AuthGetUser(username, numReplies, data["username"].(string))
							return user, err
						}
					} else {

						username, userPresent := params.Args["username"].(string)
						numDweets, numPresent := params.Args["repliesToFetch"].(int)
						if userPresent && numPresent {
							user, err := NoAuthGetUser(username, numDweets)
							return user, err
						}
					}

					return nil, errors.New("param \"username\" missing")
				},
			},
			"likedDweets": &graphql.Field{
				Type:        graphql.NewList(dweetSchema),
				Description: "Get liked dweets of authenticated user",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: -1,
					},
					"repliesToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: -1,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						numDweets, dweetPresent := params.Args["numberToFetch"].(int)
						numReplies, repliesPresent := params.Args["numberToFetch"].(int)
						if dweetPresent && repliesPresent {
							post, err := FetchLikedDweets(data["username"].(string), numDweets, numReplies)
							return post, err
						}
					}

					return nil, errors.New("Unauthorized")
				},
			},
			// TODO: FINISHED SO FAR, NEED TO FINISH MORE CHECKS FOR SUBFIELDS
			"followers": &graphql.Field{
				Type:        graphql.NewList(userSchema),
				Description: "Get followers of authenticated user",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: -1,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						num, present := params.Args["numberToFetch"].(int)
						if present {
							post, err := FetchFollowers(data["username"].(string), num)
							return post, err
						}
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"following": &graphql.Field{
				Type:        graphql.NewList(userSchema),
				Description: "Get users that authenticated user follows",
				Args: graphql.FieldConfigArgument{
					"numberToFetch": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: -1,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						num, present := params.Args["numberToFetch"].(int)
						if present {
							post, err := FetchFollowing(data["username"].(string), num)
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
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user, err := SignUpUser(
						params.Args["username"].(string),
						params.Args["password"].(string),
						params.Args["firstName"].(string),
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
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create dweet, and return formatted
						body, bodyPresent := params.Args["body"].(string)
						media, mediaPresent := params.Args["media"].([]string)
						if bodyPresent && mediaPresent {
							dweetObj, err := NewDweet(body, data["username"].(string), media)
							user := FormatAsDweetType(dweetObj)
							return user, err
						}
						return nil, errors.New("invalid request, \"body\" or \"media\" not present")
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
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create a reply to a dweet, and return formatted
						body, bodyPresent := params.Args["body"].(string)
						originalID, idPresent := params.Args["id"].(string)
						media, mediaPresent := params.Args["media"].([]string)
						if bodyPresent && mediaPresent && idPresent {
							dweetObj, err := NewReply(originalID, body, data["username"].(string), media)
							user := FormatAsDweetType(dweetObj)
							return user, err
						}
						return nil, errors.New("invalid request, \"id\", \"body\", or \"media\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"createRedweet": &graphql.Field{
				Type:        dweetSchema,
				Description: "Create a redweet of a dweet by authenticated user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Create a reply to a dweet, and return formatted
						originalID, idPresent := params.Args["id"].(string)
						if idPresent {
							dweetObj, err := NewRedweet(originalID, data["username"].(string))
							dweet := FormatAsDweetType(dweetObj)
							return dweet, err
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
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						username, present := params.Args["username"].(string)
						if present {
							userObj, err := NewFollower(username, data["username"].(string))
							user := FormatAsUserType(userObj)
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
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user like dweet, and return formatted
						id, present := params.Args["id"].(string)
						if present {
							dweetObj, err := NewLike(id, data["username"].(string))
							dweet := FormatAsDweetType(dweetObj)
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
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user like dweet, and return formatted
						id, present := params.Args["id"].(string)
						if present {
							dweetObj, err := DeleteLike(id, data["username"].(string))
							dweet := FormatAsDweetType(dweetObj)
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
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						username, present := params.Args["username"].(string)
						if present {
							userObj, err := DeleteFollower(username, data["username"].(string))
							user := FormatAsUserType(userObj)
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
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Edit dweet, and return formatted
						id, idPresent := params.Args["id"].(string)
						body, bodyPresent := params.Args["body"].(string)
						media, mediaPresent := params.Args["media"].([]string)
						if bodyPresent && mediaPresent && idPresent {
							dweet, err := AuthUpdateDweet(id, data["username"].(string), body, media)
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
						DefaultValue: nil,
					},
					"lastName": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: nil,
					},
					"email": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: nil,
					},
					"bio": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: nil,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Edit dweet, and return formatted
						firstName, firstPresent := params.Args["firstName"].(string)
						lastName, lastPresent := params.Args["lastName"].(string)
						email, emailPresent := params.Args["email"].(string)
						bio, bioPresent := params.Args["email"].(string)
						if firstPresent && lastPresent && emailPresent && bioPresent {
							user, err := AuthUpdateUser(data["username"].(string), firstName, lastName, email, bio)
							return user, err
						}
						return nil, errors.New("invalid request, \"body\" or \"media\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"deleteDweet": &graphql.Field{
				Type:        userSchema,
				Description: "Delete dweet authored by user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						id, present := params.Args["id"].(string)
						if present {
							dweet, err := AuthDeleteDweet(id, data["username"].(string))
							return dweet, err
						}
						return nil, errors.New("invalid request, \"id\" not present")
					}

					return nil, errors.New("Unauthorized")
				},
			},
			"unredweet": &graphql.Field{
				Type:        userSchema,
				Description: "Unredweet a dweet",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					// Check authentication
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						// Make user follow the other user, and return formatted
						id, present := params.Args["id"].(string)
						if present {
							dweet, err := AuthDeleteRedweet(id, data["username"].(string))
							return dweet, err
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

// Create schema from handlers
var schema, SchemaError = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryHandler,
		Mutation: mutationHandler,
	},
)
