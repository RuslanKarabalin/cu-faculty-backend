package cuclient

import (
	"context"
	"encoding/json"
	"faculty/internal/model"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func New(httpClient *http.Client, baseURL string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

func (c *Client) get(ctx context.Context, cookie string, segments ...string) ([]byte, error) {
	endpoint, err := url.JoinPath(c.baseURL, segments...)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.AddCookie(&http.Cookie{Name: CookieName, Value: cookie})

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	if httpResp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrUpstream, httpResp.StatusCode)
	}

	return io.ReadAll(io.LimitReader(httpResp.Body, 1<<20))
}

func (c *Client) Authorize(ctx context.Context, cookie string) (*model.CuUserResp, error) {
	body, err := c.get(ctx, cookie, "api", "student-hub", "students", "me")
	if err != nil {
		return nil, err
	}

	var user model.CuUserResp
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &user, nil
}

func (c *Client) StudentEduInfo(ctx context.Context, cookie string) ([]model.CuEduPlaceResp, error) {
	body, err := c.get(ctx, cookie, "api", "student-hub", "education-info", "me")
	if err != nil {
		return nil, err
	}

	var eduInfo []model.CuEduPlaceResp
	if err := json.Unmarshal(body, &eduInfo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return eduInfo, nil
}
