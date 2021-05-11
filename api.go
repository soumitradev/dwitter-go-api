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
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						id, present := params.Args["id"].(string)
						if present {
							post, err := AuthGetPost(id, 10, data["username"].(string))
							return post, err
						}
					} else {

						id, present := params.Args["id"].(string)
						if present {
							post, err := NoAuthGetPost(id, 10)
							return post, err
						}
					}

					return nil, errors.New("param \"username\" missing")
				},
			},
			"user": &graphql.Field{
				Type:        userSchema,
				Description: "Get user by username",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, isAuth, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}

					if isAuth {
						username, present := params.Args["username"].(string)
						if present {
							user, err := AuthGetUser(username, 10, data["username"].(string))
							return user, err
						}
					} else {

						username, present := params.Args["username"].(string)
						if present {
							user, err := NoAuthGetUser(username, 10)
							return user, err
						}
					}

					return nil, errors.New("param \"username\" missing")
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
