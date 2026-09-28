/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package model

// Mailboxes is the data model that wraps multiple mailboxes
type Mailboxes struct {
	Mailboxes []Mailbox `json:"mailboxes"`
}

// Mailbox is the data model for a single mailbox.
//
// The Migadu API accepts writes to most fields, but silently ignores writes to
// some of them: it answers 200 and echoes the unchanged value. Those fields are
// effectively read-only and should be treated as computed by callers:
//
//   - IsActive, ActivatedAt, ChangedAt are managed by the server.
//   - MonthlyOutgoingLimit mirrors DailyOutgoingLimit. Setting DailyOutgoingLimit
//     changes both; setting MonthlyOutgoingLimit alone has no effect.
//
// The remaining limits (DailyIncomingLimit, WeeklyIncomingLimit,
// MonthlyIncomingLimit, DailyOutgoingLimit, WeeklyOutgoingLimit) and
// WildcardSender are writable.
type Mailbox struct {
	LocalPart             string     `json:"local_part"`
	DomainName            string     `json:"domain_name"`
	Address               string     `json:"address"`
	Name                  string     `json:"name"`
	IsInternal            bool       `json:"is_internal"`
	IsActive              bool       `json:"is_active"`
	MaySend               bool       `json:"may_send"`
	MayReceive            bool       `json:"may_receive"`
	MayAccessImap         bool       `json:"may_access_imap"`
	MayAccessPop3         bool       `json:"may_access_pop3"`
	MayAccessManageSieve  bool       `json:"may_access_managesieve"`
	WildcardSender        bool       `json:"wildcard_sender"`
	PasswordMethod        string     `json:"password_method"`
	Password              string     `json:"password"`
	PasswordRecoveryEmail string     `json:"password_recovery_email"`
	SpamAction            string     `json:"spam_action"`
	SpamAggressiveness    string     `json:"spam_aggressiveness"`
	Expirable             bool       `json:"expireable"`
	ExpiresOn             string     `json:"expires_on"`
	RemoveUponExpiry      bool       `json:"remove_upon_expiry"`
	SenderDenyList        []string   `json:"sender_denylist"`
	SenderAllowList       []string   `json:"sender_allowlist"`
	RecipientDenyList     []string   `json:"recipient_denylist"`
	DailyIncomingLimit    int        `json:"daily_incoming_limit"`
	WeeklyIncomingLimit   int        `json:"weekly_incoming_limit"`
	MonthlyIncomingLimit  int        `json:"monthly_incoming_limit"`
	DailyOutgoingLimit    int        `json:"daily_outgoing_limit"`
	WeeklyOutgoingLimit   int        `json:"weekly_outgoing_limit"`
	MonthlyOutgoingLimit  int        `json:"monthly_outgoing_limit"`
	AutoRespondActive     bool       `json:"autorespond_active"`
	AutoRespondSubject    string     `json:"autorespond_subject"`
	AutoRespondBody       string     `json:"autorespond_body"`
	AutoRespondExpiresOn  string     `json:"autorespond_expires_on"`
	FooterActive          bool       `json:"footer_active"`
	FooterPlainBody       string     `json:"footer_plain_body"`
	FooterHtmlBody        string     `json:"footer_html_body"`
	StorageUsage          float64    `json:"storage_usage"`
	ActivatedAt           string     `json:"activated_at"`
	ChangedAt             string     `json:"changed_at"`
	Delegations           []string   `json:"delegations"`
	Identities            []Identity `json:"identities"`
}
