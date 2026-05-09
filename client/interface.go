/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
	"context"
	"github.com/metio/migadu-client.go/model"
)

// ClientInterface defines the surface area of MigaduClient. Consumers that
// want to mock the client in tests should depend on this interface rather
// than the concrete *MigaduClient type.
type ClientInterface interface {
	GetAliases(ctx context.Context, domain string) (*model.Aliases, error)
	GetAlias(ctx context.Context, domain string, localPart string) (*model.Alias, error)
	CreateAlias(ctx context.Context, domain string, alias *model.Alias) (*model.Alias, error)
	UpdateAlias(ctx context.Context, domain string, localPart string, alias *model.Alias) (*model.Alias, error)
	DeleteAlias(ctx context.Context, domain string, localPart string) (*model.Alias, error)

	GetIdentities(ctx context.Context, domain string, localPart string) (*model.Identities, error)
	GetIdentity(ctx context.Context, domain string, localPart string, id string) (*model.Identity, error)
	CreateIdentity(ctx context.Context, domain string, localPart string, identity *model.Identity) (*model.Identity, error)
	UpdateIdentity(ctx context.Context, domain string, localPart string, id string, identity *model.Identity) (*model.Identity, error)
	DeleteIdentity(ctx context.Context, domain string, localPart string, id string) (*model.Identity, error)

	GetMailboxes(ctx context.Context, domain string) (*model.Mailboxes, error)
	GetMailbox(ctx context.Context, domain string, localPart string) (*model.Mailbox, error)
	CreateMailbox(ctx context.Context, domain string, mailbox *model.Mailbox) (*model.Mailbox, error)
	UpdateMailbox(ctx context.Context, domain string, localPart string, mailbox *model.Mailbox) (*model.Mailbox, error)
	DeleteMailbox(ctx context.Context, domain string, localPart string) (*model.Mailbox, error)

	GetRewriteRules(ctx context.Context, domain string) (*model.RewriteRules, error)
	GetRewriteRule(ctx context.Context, domain string, slug string) (*model.RewriteRule, error)
	CreateRewriteRule(ctx context.Context, domain string, rewrite *model.RewriteRule) (*model.RewriteRule, error)
	UpdateRewriteRule(ctx context.Context, domain string, slug string, rewrite *model.RewriteRule) (*model.RewriteRule, error)
	DeleteRewriteRule(ctx context.Context, domain string, slug string) (*model.RewriteRule, error)
}

var _ ClientInterface = (*MigaduClient)(nil)
