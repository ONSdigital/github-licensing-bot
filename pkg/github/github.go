package github

import (
	"context"

	"github.com/ONSdigital/graphql"
)

type (

	// Client wraps a GraphQL client for communicating with the GitHub API.
	Client struct {
		apiBase string
		token   string
		client  *graphql.Client
	}
)

// NewClient instantiates a new GraphQL client.
func NewClient(apiBase, token string) *Client {
	return &Client{
		apiBase: apiBase,
		token:   token,
		client:  graphql.NewClient(apiBase + "/graphql"),
	}
}

// Run wraps the underlying graphql.Run function, authomatically adding an authentication header and background context.
func (c Client) Run(request *graphql.Request, response interface{}) error {
	request.Header.Set("Authorization", "Bearer "+c.token)
	return c.client.Run(context.Background(), request, response)
}
