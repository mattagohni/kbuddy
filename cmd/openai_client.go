package cmd

import (
	"context"
	"github.com/franciscoescher/goopenai"
)

type OpenAiClient interface {
	CreateCompletions(ctx context.Context, req goopenai.CreateCompletionsRequest) (goopenai.CreateCompletionsResponse, error)
}

type MyOpenAiClient struct {
	client *goopenai.Client
}

func NewMyOpenAiClient(apiKey, organization string) *MyOpenAiClient {
	client := goopenai.NewClient(apiKey, organization)
	return &MyOpenAiClient{
		client: client,
	}
}

func (c *MyOpenAiClient) CreateCompletions(ctx context.Context, req goopenai.CreateCompletionsRequest) (goopenai.CreateCompletionsResponse, error) {
	return c.client.CreateCompletions(ctx, req)
}
