/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package simulator

import (
	"github.com/metio/migadu-client.go/model"
	"net/http"
	"testing"
)

// State is the optional state of the Migadu API. Use this to populate the simulator before a test.
type State struct {
	Mailboxes  []model.Mailbox
	Aliases    []model.Alias
	Identities []model.Identity
	Rewrites   []model.RewriteRule
	StatusCode int
}

// MigaduAPI returns a handler function that simulates the Migadu API.
func MigaduAPI(t *testing.T, state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if aliasesUrlPattern.MatchString(r.URL.Path) {
			handleAliases(t, &state.Aliases, state.StatusCode).ServeHTTP(w, r)
		} else if identitiesUrlPattern.MatchString(r.URL.Path) {
			handleIdentities(t, &state.Identities, state.StatusCode).ServeHTTP(w, r)
		} else if rewritesPattern.MatchString(r.URL.Path) {
			handleRewriteRules(t, &state.Rewrites, state.StatusCode).ServeHTTP(w, r)
		} else if mailboxesPattern.MatchString(r.URL.Path) {
			handleMailboxes(t, &state.Mailboxes, state.StatusCode).ServeHTTP(w, r)
		} else {
			t.Errorf("No Handler for URL: %s", r.URL.Path)
		}
	}
}
