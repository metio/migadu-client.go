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
)

// GetAliases returns all aliases for a single domain
func (c *MigaduClient) GetAliases(ctx context.Context, domain string) (*model.Aliases, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}
	return do[model.Aliases](ctx, c, "GetAliases", http.MethodGet, nil,
		"domains", ascii, "aliases")
}

// GetAlias returns specific alias
func (c *MigaduClient) GetAlias(ctx context.Context, domain string, localPart string) (*model.Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}
	return do[model.Alias](ctx, c, "GetAlias", http.MethodGet, nil,
		"domains", ascii, "aliases", url.PathEscape(localPart))
}

// CreateAlias creates a new alias
func (c *MigaduClient) CreateAlias(ctx context.Context, domain string, alias *model.Alias) (*model.Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}
	asciiEmails, err := idn.ConvertEmailsToASCII(alias.Destinations)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}
	alias.Destinations = asciiEmails
	return do[model.Alias](ctx, c, "CreateAlias", http.MethodPost, alias,
		"domains", ascii, "aliases")
}

// UpdateAlias updates an existing alias
func (c *MigaduClient) UpdateAlias(ctx context.Context, domain string, localPart string, alias *model.Alias) (*model.Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}
	asciiEmails, err := idn.ConvertEmailsToASCII(alias.Destinations)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}
	alias.Destinations = asciiEmails
	return do[model.Alias](ctx, c, "UpdateAlias", http.MethodPut, alias,
		"domains", ascii, "aliases", url.PathEscape(localPart))
}

// DeleteAlias deletes an existing alias
func (c *MigaduClient) DeleteAlias(ctx context.Context, domain string, localPart string) (*model.Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}
	return do[model.Alias](ctx, c, "DeleteAlias", http.MethodDelete, nil,
		"domains", ascii, "aliases", url.PathEscape(localPart))
}
