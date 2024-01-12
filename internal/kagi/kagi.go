package kagi

import (
	"log"
	"os"

	"code.vin047.com/focus-browser-server/internal/search"
	kagi "github.com/httpjamesm/kagigo"
)

type Client struct {
	kagiClient *kagi.Client
}

var loggerDebug = log.New(os.Stdout, "[DEBUG]: ", log.LstdFlags)

func NewClient(apiKey string) *Client {
	return &Client{
		kagiClient: kagi.NewClient(&kagi.ClientConfig{
			APIKey:     apiKey,
			APIVersion: "v0",
		}),
	}
}

func (c *Client) Search(query string) (*search.Result, error) {
	response, err := c.kagiClient.FastGPTCompletion(kagi.FastGPTCompletionParams{
		Query:     query,
		WebSearch: true,
		Cache:     true,
	})
	if err != nil {
		return &search.Result{}, err
	}

	loggerDebug.Printf("Tokens: %d, API balance: %.2f\n", response.Data.Tokens, response.Meta.APIBalance)

	return &search.Result{
		Response: response.Data.Output,
		References: func() []search.Reference {
			var references []search.Reference
			for _, reference := range response.Data.References {
				references = append(references, search.Reference{
					Title: reference.Title,
					Url:   reference.URL,
				})
			}
			return references
		}(),
	}, err
}
