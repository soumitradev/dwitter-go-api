package main

import (
	"github.com/graphql-go/graphql"
)

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
			"login": &graphql.Field{
				Type:        loginResponseSchema,
				Description: "Log into Dwitter",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					user, err := LoginUser(
						params.Args["username"].(string),
						params.Args["password"].(string),
					)

					if err != nil {
						return nil, err
					}

					return user, nil
				},
			},
		},
	},
)

var schema, SchemaError = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryHandler,
		Mutation: mutationHandler,
	},
)
