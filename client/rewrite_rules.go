/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
	"context"
	"fmt"
	"github.com/metio/migadu-client.go/idn"
	"github.com/metio/migadu-client.go/model"
	"golang.org/x/net/idna"
	"net/http"
	"net/url"
	"strings"
)

type rewriteRuleJson struct {
	*model.RewriteRule
	Destinations string `json:"destinations"`
}

// GetRewriteRules returns rewrite rules for a single domain
func (c *MigaduClient) GetRewriteRules(ctx context.Context, domain string) (*model.RewriteRules, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRules: %w", err)
	}
	return do[model.RewriteRules](ctx, c, "GetRewriteRules", http.MethodGet, nil,
		"domains", ascii, "rewrites")
}

// GetRewriteRule returns a specific rewrite rule
func (c *MigaduClient) GetRewriteRule(ctx context.Context, domain string, slug string) (*model.RewriteRule, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetRewriteRule: %w", err)
	}
	return do[model.RewriteRule](ctx, c, "GetRewriteRule", http.MethodGet, nil,
		"domains", ascii, "rewrites", url.PathEscape(slug))
}

// CreateRewriteRule creates a new rewrite rule
func (c *MigaduClient) CreateRewriteRule(ctx context.Context, domain string, rewrite *model.RewriteRule) (*model.RewriteRule, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateRewriteRule: %w", err)
	}
	body, err := encodeRewriteRule(rewrite)
	if err != nil {
		return nil, fmt.Errorf("CreateRewriteRule: %w", err)
	}
	return do[model.RewriteRule](ctx, c, "CreateRewriteRule", http.MethodPost, body,
		"domains", ascii, "rewrites")
}

// UpdateRewriteRule updates an existing rewrite rule
func (c *MigaduClient) UpdateRewriteRule(ctx context.Context, domain string, slug string, rewrite *model.RewriteRule) (*model.RewriteRule, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateRewriteRule: %w", err)
	}
	body, err := encodeRewriteRule(rewrite)
	if err != nil {
		return nil, fmt.Errorf("UpdateRewriteRule: %w", err)
	}
	return do[model.RewriteRule](ctx, c, "UpdateRewriteRule", http.MethodPut, body,
		"domains", ascii, "rewrites", url.PathEscape(slug))
}

// DeleteRewriteRule deletes an existing rewrite rule
func (c *MigaduClient) DeleteRewriteRule(ctx context.Context, domain string, slug string) (*model.RewriteRule, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteRewriteRule: %w", err)
	}
	return do[model.RewriteRule](ctx, c, "DeleteRewriteRule", http.MethodDelete, nil,
		"domains", ascii, "rewrites", url.PathEscape(slug))
}

func encodeRewriteRule(rewrite *model.RewriteRule) (any, error) {
	if rewrite == nil {
		return nil, nil
	}
	asciiEmails, err := idn.ConvertEmailsToASCII(rewrite.Destinations)
	if err != nil {
		return nil, err
	}
	rewrite.Destinations = asciiEmails
	return rewriteRuleJson{RewriteRule: rewrite, Destinations: strings.Join(rewrite.Destinations, ",")}, nil
}
