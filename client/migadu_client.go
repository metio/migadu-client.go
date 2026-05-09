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
	"strconv"
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
	retry       *retryConfig
}

type retryConfig struct {
	maxRetries int
	maxBackoff time.Duration
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

// WithRetry enables automatic retries for transient failures: network errors
// and HTTP 429, 502, 503, and 504 responses. Up to maxRetries additional
// attempts are made with exponential backoff (starting at 500ms, doubling
// each attempt) capped at maxBackoff. When the server returns a Retry-After
// header on a 429, that value is honored instead of the computed backoff.
func WithRetry(maxRetries int, maxBackoff time.Duration) Option {
	return func(c *MigaduClient) {
		c.retry = &retryConfig{maxRetries: maxRetries, maxBackoff: maxBackoff}
	}
}

// WithDefaultRetry enables retries with sensible defaults: 3 retries and a
// 30-second cap on backoff.
func WithDefaultRetry() Option {
	return WithRetry(3, 30*time.Second)
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
	req.SetBasicAuth(c.username, c.token)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Add("Content-Type", "application/json")

	maxAttempts := 1
	if c.retry != nil {
		maxAttempts += c.retry.maxRetries
	}

	var (
		lastBody       []byte
		lastStatus     int
		lastErr        error
		lastRetryAfter time.Duration
	)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			wait := computeBackoff(attempt, lastRetryAfter, c.retry.maxBackoff)
			select {
			case <-time.After(wait):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
			if req.GetBody != nil {
				fresh, err := req.GetBody()
				if err != nil {
					return nil, err
				}
				req.Body = fresh
			}
		}

		if c.rateLimiter != nil {
			if err := c.rateLimiter.Wait(req.Context()); err != nil {
				return nil, err
			}
		}

		body, status, retryAfter, err := c.attemptOnce(req)
		if err != nil {
			if c.retry == nil || !isNetworkRetryable(err) {
				return nil, err
			}
			lastErr, lastBody, lastStatus, lastRetryAfter = err, nil, 0, 0
			continue
		}
		if status == http.StatusOK {
			return body, nil
		}
		if c.retry == nil || !isStatusRetryable(status) {
			return nil, &RequestError{StatusCode: status, ResponseBody: body}
		}
		lastErr, lastBody, lastStatus, lastRetryAfter = nil, body, status, retryAfter
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, &RequestError{StatusCode: lastStatus, ResponseBody: lastBody}
}

func (c *MigaduClient) attemptOnce(req *http.Request) ([]byte, int, time.Duration, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, 0, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, 0, err
	}
	return body, res.StatusCode, parseRetryAfter(res.Header.Get("Retry-After")), nil
}

func isNetworkRetryable(err error) bool {
	return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
}

func isStatusRetryable(status int) bool {
	switch status {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

func parseRetryAfter(header string) time.Duration {
	if header == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(header); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second
	}
	if t, err := http.ParseTime(header); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

func computeBackoff(attempt int, retryAfter, max time.Duration) time.Duration {
	if retryAfter > 0 {
		if retryAfter > max {
			return max
		}
		return retryAfter
	}
	shift := attempt - 1
	if shift > 30 {
		shift = 30
	}
	backoff := 500 * time.Millisecond << shift
	if backoff > max {
		return max
	}
	return backoff
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
