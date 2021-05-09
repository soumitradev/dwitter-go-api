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
						Type: graphql.String,
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
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					username, success := params.Args["username"].(string)
					if success {
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
		Name:   "Mutation",
		Fields: graphql.Field{},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryHandler,
		Mutation: mutationHandler,
	},
)
