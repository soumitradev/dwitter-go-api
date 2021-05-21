package gql

import (
	"dwitter_go_graphql/consts"
	"time"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

func InitSubscriptions() {
	go func() {
		for {
			// Every 5 mins, update the subscriptions
			time.Sleep(5 * time.Minute)
			subscriptions := consts.SubscriptionManager.Subscriptions()

			for conn := range subscriptions {
				// Things you have access to here:
				conn.ID()   // The connection ID
				conn.User() // The user returned from the Authenticate function

				for _, subscription := range subscriptions[conn] {
					// Things you have access to here:
					// subscription.ID            // The subscription ID (unique per conn)
					// subscription.OperationName // The name of the operation
					// subscription.Query         // The subscription query/queries string
					// subscription.Variables     // The subscription variables
					// subscription.Document      // The GraphQL AST for the subscription
					// subscription.Fields        // The names of top-level queries
					// subscription.Connection    // The GraphQL WS connection

					// Re-execute the subscription query
					params := graphql.Params{
						Schema:         Schema, // The GraphQL schema
						RequestString:  subscription.Query,
						VariableValues: subscription.Variables,
						OperationName:  subscription.OperationName,
						Context:        consts.BaseCtx,
					}
					result := graphql.Do(params)

					// Send query results back to the subscriber at any point
					data := graphqlws.DataMessagePayload{
						// Data can be anything (interface{})
						Data: result.Data,
						// Errors is optional ([]error)
						Errors: graphqlws.ErrorsFromGraphQLErrors(result.Errors),
					}
					subscription.SendData(&data)
				}
			}
		}
	}()
}
