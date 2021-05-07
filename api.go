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
						post, err := APIGetPost(id, -1)
						return post, err
					}
					return nil, nil
				},
			},
			"user": &graphql.Field{
				Type:        userSchema,
				Description: "Get user by mention",
				Args: graphql.FieldConfigArgument{
					"mention": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					mention, success := params.Args["mention"].(string)
					if success {
						post, err := APIGetUser(mention, -1, -1)
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
