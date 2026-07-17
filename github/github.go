package github

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Client manages communication with the GitHub API.
type Client struct {
	client     *http.Client
	MaxRetries int
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		client:     httpClient,
		MaxRetries: 3,
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.MaxRetries; i++ {
		resp, err = c.client.Do(req.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if (resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests) {
			retryAfter := resp.Header.Get("Retry-After")
			seconds, parseErr := strconv.Atoi(retryAfter)
			if parseErr == nil && seconds > 0 {
				select {
				case <-time.After(time.Duration(seconds) * time.Second):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
		}
		break
	}
	return resp, err
}