/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultUserAgent = "migadu-client.go"

// MigaduClient holds all necessary information to interact with the Migadu API
type MigaduClient struct {
	httpClient  *http.Client
	endpoint    string
	username    string
	token       string
	userAgent   string
	rateLimiter *rate.Limiter
}

// Option is a functional option for configuring a MigaduClient.
// Pass one or more options to New to customize client behavior.
type Option func(*MigaduClient)

// WithRateLimit configures a rate limiter that allows at most rateLimit requests
// per rateInterval. The limiter blocks each request until a token is available,
// so callers do not need to implement back-off themselves.
//
// Example: WithRateLimit(60, 2*time.Minute) permits 60 requests every 2 minutes.
func WithRateLimit(rateLimit int, rateInterval time.Duration) Option {
	return func(c *MigaduClient) {
		c.rateLimiter = rate.NewLimiter(rate.Every(rateInterval/time.Duration(rateLimit)), rateLimit)
	}
}

// WithDefaultRateLimit applies the rate limit recommended by the Migadu API:
// 60 requests per 2 minutes (~1 request per 2 seconds).
func WithDefaultRateLimit() Option {
	return WithRateLimit(60, 2*time.Minute)
}

// WithUserAgent overrides the User-Agent header sent on every request.
func WithUserAgent(userAgent string) Option {
	return func(c *MigaduClient) {
		c.userAgent = userAgent
	}
}

// New creates a new MigaduClient
func New(endpoint, username, token *string, timeout time.Duration, opts ...Option) (*MigaduClient, error) {
	if username == nil {
		return nil, errors.New("no username supplied")
	}
	if token == nil {
		return nil, errors.New("no token supplied")
	}

	c := MigaduClient{
		httpClient: &http.Client{Timeout: timeout},
		endpoint:   "https://api.migadu.com/v1/",
		username:   *username,
		token:      *token,
		userAgent:  defaultUserAgent,
	}

	if endpoint != nil {
		c.endpoint = *endpoint
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c, nil
}

func (c *MigaduClient) doRequest(req *http.Request) ([]byte, error) {
	if c.rateLimiter != nil {
		err := c.rateLimiter.Wait(req.Context())
		if err != nil {
			return nil, err
		}
	}

	req.SetBasicAuth(c.username, c.token)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, &RequestError{StatusCode: res.StatusCode, ResponseBody: body}
	}

	return body, err
}

type RequestError struct {
	StatusCode   int
	ResponseBody []byte
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("status: %d, body: %s", e.StatusCode, e.ResponseBody)
}

// do executes an authenticated request and decodes a JSON response of type T.
// label is prefixed onto any returned error. If body is nil, no request body
// is sent; otherwise body is JSON-encoded. pathSegments are joined onto the
// endpoint via url.JoinPath; callers are responsible for url.PathEscape on
// any user-controlled segments.
func do[T any](ctx context.Context, c *MigaduClient, label, method string, body any, pathSegments ...string) (*T, error) {
	reqURL, err := url.JoinPath(c.endpoint, pathSegments...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}

	var reader io.Reader = http.NoBody
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		reader = bytes.NewBuffer(encoded)
	}

	request, err := http.NewRequestWithContext(ctx, method, reqURL, reader)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}

	var result T
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}

	return &result, nil
}
