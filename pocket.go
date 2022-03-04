package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/olinbg/go-pocket-sdk/internal/request"
	"github.com/olinbg/go-pocket-sdk/pkg/input"
	"github.com/olinbg/go-pocket-sdk/pkg/response"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	host         = "https://getpocket.com/v3"
	authorizeUrl = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"

	endpointAdd = "/add"
	endpointGet = "/get"

	endpointRequestToken = "/oauth/request"
	endpointAuthorize    = "/oauth/authorize"

	// xErrorHeader used to parse error message from Headers on non-2XX responses
	xErrorHeader = "X-Error"

	defaultTimeout = 5 * time.Second
)

// Client is a getpocket API client
type Client struct {
	client      *http.Client
	consumerKey string
}

// NewClient creates a new client instance with your app key (to generate key visit https://getpocket.com/developer/apps/)
func NewClient(consumerKey string) (*Client, error) {
	if consumerKey == "" {
		return nil, errors.New("consumer key is empty")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		consumerKey: consumerKey,
	}, nil
}

// GetRequestToken obtains the request token that is used to authorize user in your application
func (c *Client) GetRequestToken(ctx context.Context, redirectUrl string) (string, error) {
	inp := &request.RequestTokenRequest{
		ConsumerKey: c.consumerKey,
		RedirectURI: redirectUrl,
	}

	values, err := c.doHTTP(ctx, endpointRequestToken, inp)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = json.Unmarshal([]byte(values), &data)
	if err != nil {
		return "", err
	}

	code := data["code"].(string)

	return code, nil
}

// GetAuthorizationURL generates link to authorize user
func (c *Client) GetAuthorizationURL(requestToken, redirectUrl string) (string, error) {
	if requestToken == "" || redirectUrl == "" {
		return "", errors.New("empty params")
	}

	return fmt.Sprintf(authorizeUrl, requestToken, redirectUrl), nil
}

// Authorize generates access token for user, that authorized in your app via link
func (c *Client) Authorize(ctx context.Context, requestToken string) (*response.AuthorizeResponse, error) {
	if requestToken == "" {
		return nil, errors.New("empty request token")
	}

	inp := &request.AuthorizeRequest{
		Code:        requestToken,
		ConsumerKey: c.consumerKey,
	}

	values, err := c.doHTTP(ctx, endpointAuthorize, inp)
	if err != nil {
		return nil, err
	}

	var authResponse response.AuthorizeResponse

	err = json.Unmarshal([]byte(values), &authResponse)

	if err != nil {
		return nil, err
	}

	if authResponse.AccessToken == "" {
		return nil, errors.New("empty access token in API response")
	}

	return &authResponse, err
}

// Add creates new item in Pocket list
func (c *Client) Add(ctx context.Context, input input.AddInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	req := input.GenerateRequest(c.consumerKey)
	_, err := c.doHTTP(ctx, endpointAdd, req)

	return err
}

// Get request already existing items in Pocket list
func (c *Client) Get(ctx context.Context, input input.GetInput) (*response.GetResponse, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	req := input.GenerateRequest(c.consumerKey)
	values, err := c.doHTTP(ctx, endpointGet, req)

	if err != nil {
		return nil, err
	}

	var getResponse response.GetResponse

	err = json.Unmarshal([]byte(values), &getResponse)

	if err != nil {
		return nil, err
	}

	return &getResponse, err
}

func (c *Client) doHTTP(ctx context.Context, endpoint string, body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", errors.WithMessage(err, "failed to marshal input body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, host+endpoint, bytes.NewBuffer(b))
	if err != nil {
		return "", errors.WithMessage(err, "failed to create new request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF8")
	req.Header.Set("X-Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", errors.WithMessage(err, "failed to send http request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Sprintf("API Error: %s", resp.Header.Get(xErrorHeader))
		return "", errors.New(err)
	}

	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read request body")
	}

	return string(respB), nil
}
