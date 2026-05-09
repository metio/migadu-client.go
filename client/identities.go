/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
	"context"
	"fmt"
	"github.com/metio/migadu-client.go/model"
	"golang.org/x/net/idna"
	"net/http"
	"net/url"
)

// GetIdentities returns identities for a single mailbox
func (c *MigaduClient) GetIdentities(ctx context.Context, domain string, localPart string) (*model.Identities, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetIdentities: %w", err)
	}
	return do[model.Identities](ctx, c, "GetIdentities", http.MethodGet, nil,
		"domains", ascii, "mailboxes", url.PathEscape(localPart), "identities")
}

// GetIdentity returns a specific identity
func (c *MigaduClient) GetIdentity(ctx context.Context, domain string, localPart string, id string) (*model.Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetIdentity: %w", err)
	}
	return do[model.Identity](ctx, c, "GetIdentity", http.MethodGet, nil,
		"domains", ascii, "mailboxes", url.PathEscape(localPart), "identities", url.PathEscape(id))
}

// CreateIdentity creates a new identity
func (c *MigaduClient) CreateIdentity(ctx context.Context, domain string, localPart string, identity *model.Identity) (*model.Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateIdentity: %w", err)
	}
	return do[model.Identity](ctx, c, "CreateIdentity", http.MethodPost, identity,
		"domains", ascii, "mailboxes", url.PathEscape(localPart), "identities")
}

// UpdateIdentity updates an existing identity
func (c *MigaduClient) UpdateIdentity(ctx context.Context, domain string, localPart string, id string, identity *model.Identity) (*model.Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateIdentity: %w", err)
	}
	return do[model.Identity](ctx, c, "UpdateIdentity", http.MethodPut, identity,
		"domains", ascii, "mailboxes", url.PathEscape(localPart), "identities", url.PathEscape(id))
}

// DeleteIdentity deletes an existing identity
func (c *MigaduClient) DeleteIdentity(ctx context.Context, domain string, localPart string, id string) (*model.Identity, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteIdentity: %w", err)
	}
	return do[model.Identity](ctx, c, "DeleteIdentity", http.MethodDelete, nil,
		"domains", ascii, "mailboxes", url.PathEscape(localPart), "identities", url.PathEscape(id))
}
