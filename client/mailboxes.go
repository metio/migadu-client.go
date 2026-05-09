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

// GetMailboxes returns mailboxes for a single domain
func (c *MigaduClient) GetMailboxes(ctx context.Context, domain string) (*model.Mailboxes, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetMailboxes: %w", err)
	}
	return do[model.Mailboxes](ctx, c, "GetMailboxes", http.MethodGet, nil,
		"domains", ascii, "mailboxes")
}

// GetMailbox returns specific mailbox
func (c *MigaduClient) GetMailbox(ctx context.Context, domain string, localPart string) (*model.Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetMailbox: %w", err)
	}
	return do[model.Mailbox](ctx, c, "GetMailbox", http.MethodGet, nil,
		"domains", ascii, "mailboxes", url.PathEscape(localPart))
}

// CreateMailbox creates a new mailbox
func (c *MigaduClient) CreateMailbox(ctx context.Context, domain string, mailbox *model.Mailbox) (*model.Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}
	if err := mailboxListsToASCII(mailbox); err != nil {
		return nil, fmt.Errorf("CreateMailbox: %w", err)
	}
	return do[model.Mailbox](ctx, c, "CreateMailbox", http.MethodPost, mailbox,
		"domains", ascii, "mailboxes")
}

// UpdateMailbox updates an existing mailbox
func (c *MigaduClient) UpdateMailbox(ctx context.Context, domain string, localPart string, mailbox *model.Mailbox) (*model.Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}
	if err := mailboxListsToASCII(mailbox); err != nil {
		return nil, fmt.Errorf("UpdateMailbox: %w", err)
	}
	return do[model.Mailbox](ctx, c, "UpdateMailbox", http.MethodPut, mailbox,
		"domains", ascii, "mailboxes", url.PathEscape(localPart))
}

// DeleteMailbox deletes an existing mailbox
func (c *MigaduClient) DeleteMailbox(ctx context.Context, domain string, localPart string) (*model.Mailbox, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteMailbox: %w", err)
	}
	return do[model.Mailbox](ctx, c, "DeleteMailbox", http.MethodDelete, nil,
		"domains", ascii, "mailboxes", url.PathEscape(localPart))
}

func mailboxListsToASCII(mailbox *model.Mailbox) error {
	if mailbox == nil {
		return nil
	}
	var err error
	if mailbox.SenderDenyList, err = idn.ConvertEmailsToASCII(mailbox.SenderDenyList); err != nil {
		return err
	}
	if mailbox.SenderAllowList, err = idn.ConvertEmailsToASCII(mailbox.SenderAllowList); err != nil {
		return err
	}
	if mailbox.RecipientDenyList, err = idn.ConvertEmailsToASCII(mailbox.RecipientDenyList); err != nil {
		return err
	}
	return nil
}
