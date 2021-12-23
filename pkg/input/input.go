package input

import (
	"github.com/pkg/errors"
	"github.com/tierko/go-pocket-sdk/internal/request"
	"strings"
)

type (
	// AddInput holds data necessary to create new item in Pocket list
	AddInput struct {
		URL         string
		Title       string
		Tags        []string
		AccessToken string
	}

	GetInput struct {
		AccessToken string
	}
)

func (i AddInput) Validate() error {
	if i.URL == "" {
		return errors.New("required URL values is empty")
	}

	if i.AccessToken == "" {
		return errors.New("access token is empty")
	}

	return nil
}

func (i GetInput) Validate() error {
	if i.AccessToken == "" {
		return errors.New("access token is empty")
	}

	return nil
}

func (i AddInput) GenerateRequest(consumerKey string) request.AddRequest {
	return request.AddRequest{
		URL:         i.URL,
		Tags:        strings.Join(i.Tags, ","),
		Title:       i.Title,
		AccessToken: i.AccessToken,
		ConsumerKey: consumerKey,
	}
}

func (i GetInput) generateRequest(consumerKey string) request.GetRequest {
	return request.GetRequest{
		ConsumerKey: consumerKey,
		AccessToken: i.AccessToken,
	}
}
