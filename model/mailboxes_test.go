/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package model_test

import (
	"encoding/json"
	"github.com/metio/migadu-client.go/model"
	"testing"
)

// mailboxResponse mirrors the payload the Migadu API returns for
// GET /v1/domains/{domain}/mailboxes/{local_part}, with every key the API is
// known to send. It exists to pin the JSON keys of model.Mailbox against the
// real wire format: a round trip through the struct alone cannot catch a
// misspelled tag, because it would serialize and deserialize symmetrically.
//
// Note that "forwardings" is deliberately not modelled. It is a subresource
// with its own endpoints, and encoding/json ignores the unknown key.
const mailboxResponse = `{
  "local_part": "test",
  "domain_name": "example.com",
  "address": "test@example.com",
  "name": "Some Name",
  "is_internal": false,
  "is_active": true,
  "may_send": true,
  "may_receive": true,
  "may_access_imap": true,
  "may_access_pop3": false,
  "may_access_managesieve": true,
  "wildcard_sender": true,
  "password_recovery_email": "recovery@example.com",
  "spam_action": "folder",
  "spam_aggressiveness": "default",
  "expireable": false,
  "expires_on": null,
  "remove_upon_expiry": false,
  "sender_denylist": [],
  "sender_allowlist": [],
  "recipient_denylist": [],
  "daily_incoming_limit": 100,
  "weekly_incoming_limit": 200,
  "monthly_incoming_limit": 300,
  "daily_outgoing_limit": 400,
  "weekly_outgoing_limit": 500,
  "monthly_outgoing_limit": 400,
  "footer_active": false,
  "storage_usage": 0.0,
  "activated_at": "09/07/26",
  "changed_at": "26/07/26",
  "forwardings": [],
  "delegations": [],
  "identities": []
}`

func TestMailbox_UnmarshalJSON(t *testing.T) {
	var mailbox model.Mailbox
	if err := json.Unmarshal([]byte(mailboxResponse), &mailbox); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{name: "is_active", got: mailbox.IsActive, want: true},
		{name: "wildcard_sender", got: mailbox.WildcardSender, want: true},
		{name: "daily_incoming_limit", got: mailbox.DailyIncomingLimit, want: 100},
		{name: "weekly_incoming_limit", got: mailbox.WeeklyIncomingLimit, want: 200},
		{name: "monthly_incoming_limit", got: mailbox.MonthlyIncomingLimit, want: 300},
		{name: "daily_outgoing_limit", got: mailbox.DailyOutgoingLimit, want: 400},
		{name: "weekly_outgoing_limit", got: mailbox.WeeklyOutgoingLimit, want: 500},
		{name: "monthly_outgoing_limit", got: mailbox.MonthlyOutgoingLimit, want: 400},
		{name: "activated_at", got: mailbox.ActivatedAt, want: "09/07/26"},
		{name: "changed_at", got: mailbox.ChangedAt, want: "26/07/26"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s got = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}
