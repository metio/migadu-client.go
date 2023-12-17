/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package simulator

import (
	"encoding/json"
	"fmt"
	"github.com/metio/migadu-client.go/idn"
	"github.com/metio/migadu-client.go/model"
	"golang.org/x/net/idna"
	"io"
	"net/http"
	"regexp"
	"testing"
)

var aliasesUrlPattern = regexp.MustCompile("/domains/(.*)/aliases/?(.*)?")

func handleAliases(t *testing.T, aliases *[]model.Alias, forcedStatusCode int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := aliasesUrlPattern.FindStringSubmatch(r.URL.Path)
		if matches == nil {
			t.Errorf("Expected to request to match %s, got: %s", aliasesUrlPattern, r.URL.Path)
		}

		domain, err := idna.ToASCII(matches[1])
		if err != nil {
			t.Errorf("Could not convert %s to ASCII because of: %v", matches[1], err)
		}
		localPart := matches[2]

		if forcedStatusCode > 0 {
			w.WriteHeader(forcedStatusCode)
			return
		}

		if r.Method == http.MethodPost {
			handleCreateAlias(w, r, t, aliases, domain)
		}
		if r.Method == http.MethodPut {
			handleUpdateAlias(w, r, t, aliases, domain, localPart)
		}
		if r.Method == http.MethodDelete {
			handleDeleteAlias(w, r, t, aliases, domain, localPart)
		}
		if r.Method == http.MethodGet {
			if localPart == "" {
				handleGetAliases(w, r, t, aliases, domain)
			} else {
				handleGetAlias(w, r, t, aliases, domain, localPart)
			}
		}
	}
}

func handleGetAliases(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]model.Alias, domain string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases", domain) {
		t.Errorf("Expected to request '/domains/%s/aliases', got: %s", domain, r.URL.Path)
	}

	var found []model.Alias
	for _, alias := range *aliases {
		if alias.DomainName == domain {
			found = append(found, alias)
		}
	}
	w.WriteHeader(http.StatusOK)
	writeJsonResponse(t, w, model.Aliases{Aliases: found})
}

func handleGetAlias(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]model.Alias, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/aliases/%s', got: %s", domain, localPart, r.URL.Path)
	}

	missing := true
	for _, alias := range *aliases {
		if alias.DomainName == domain && alias.LocalPart == localPart {
			missing = false
			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, alias)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleDeleteAlias(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]model.Alias, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/aliases/%s', got: %s", domain, localPart, r.URL.Path)
	}

	missing := true
	for index, alias := range *aliases {
		if alias.DomainName == domain && alias.LocalPart == localPart {
			missing = false
			c := *aliases
			c[index] = c[len(c)-1]
			*aliases = c[:len(c)-1]

			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, alias)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleUpdateAlias(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]model.Alias, domain string, localPart string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases/%s", domain, localPart) {
		t.Errorf("Expected to request '/domains/%s/aliases/%s', got: %s", domain, localPart, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	requestAlias := model.Alias{}
	err = json.Unmarshal(requestBody, &requestAlias)
	if err != nil {
		t.Errorf("Could not unmarshall alias")
	}

	requestAlias.DomainName = domain
	requestAlias.LocalPart = localPart
	requestAlias.Address = fmt.Sprintf("%s@%s", requestAlias.LocalPart, domain)
	ascii, err := idn.ConvertEmailsToASCII(requestAlias.Destinations)
	if err != nil {
		t.Errorf("Could not convert to punycode")
	}
	requestAlias.Destinations = ascii

	missing := true
	for index, alias := range *aliases {
		if alias.DomainName == domain && alias.LocalPart == localPart {
			missing = false
			c := *aliases
			c[index] = requestAlias
			*aliases = c

			w.WriteHeader(http.StatusOK)
			writeJsonResponse(t, w, requestAlias)
		}
	}
	if missing {
		w.WriteHeader(http.StatusNotFound)
	}
}

func handleCreateAlias(w http.ResponseWriter, r *http.Request, t *testing.T, aliases *[]model.Alias, domain string) {
	if r.URL.Path != fmt.Sprintf("/domains/%s/aliases", domain) {
		t.Errorf("Expected to request '/domains/%s/aliases', got: %s", domain, r.URL.Path)
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Could not read body")
	}

	alias := model.Alias{}
	err = json.Unmarshal(requestBody, &alias)
	if err != nil {
		t.Errorf("Could not unmarshall alias")
	}

	alias.DomainName = domain
	alias.Address = fmt.Sprintf("%s@%s", alias.LocalPart, domain)
	ascii, err := idn.ConvertEmailsToASCII(alias.Destinations)
	if err != nil {
		t.Errorf("Could not convert to punycode")
	}
	alias.Destinations = ascii

	for _, existingAlias := range *aliases {
		if existingAlias.DomainName == alias.DomainName && existingAlias.LocalPart == alias.LocalPart {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	*aliases = append(*aliases, alias)

	w.WriteHeader(http.StatusOK)
	writeJsonResponse(t, w, alias)
}
