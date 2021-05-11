package main

import (
	"fmt"

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
					id, success := params.Args["id"].(string)
					if success {
						post, err := NoAuthGetPost(id, 10)
						return post, err
					}
					return nil, nil
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
					username, present := params.Args["username"].(string)
					if present {
						post, err := NoAuthGetUser(username, 10)
						return post, err
					}
					return nil, nil
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
			"authTest": &graphql.Field{
				Type:        graphql.String,
				Description: "Log into Dwitter",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					tokenString := params.Info.RootValue.(map[string]interface{})["token"].(string)
					data, err := VerifyToken(tokenString)
					if err != nil {
						return nil, err
					}
					return fmt.Sprintf("Username: %v", data["username"]), err
				},
			},
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
