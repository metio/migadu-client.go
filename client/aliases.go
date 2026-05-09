/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package client

import (
	"bytes"
	"context"
	"encoding/json"
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

	reqURL, err := url.JoinPath(c.Endpoint, "domains", ascii, "aliases")
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}

	response := model.Aliases{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetAliases: %w", err)
	}

	return &response, nil
}

// GetAlias returns specific alias
func (c *MigaduClient) GetAlias(ctx context.Context, domain string, localPart string) (*model.Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	reqURL, err := url.JoinPath(c.Endpoint, "domains", ascii, "aliases", url.PathEscape(localPart))
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	response := model.Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("GetAlias: %w", err)
	}

	return &response, nil
}

// CreateAlias creates a new alias
func (c *MigaduClient) CreateAlias(ctx context.Context, domain string, alias *model.Alias) (*model.Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	reqURL, err := url.JoinPath(c.Endpoint, "domains", ascii, "aliases")
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	asciiEmails, err := idn.ConvertEmailsToASCII(alias.Destinations)
	if err != nil {
		return nil, err
	}
	alias.Destinations = asciiEmails

	requestBody, err := json.Marshal(alias)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	response := model.Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("CreateAlias: %w", err)
	}

	return &response, nil
}

// UpdateAlias updates an existing alias
func (c *MigaduClient) UpdateAlias(ctx context.Context, domain string, localPart string, alias *model.Alias) (*model.Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	reqURL, err := url.JoinPath(c.Endpoint, "domains", ascii, "aliases", url.PathEscape(localPart))
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	asciiEmails, err := idn.ConvertEmailsToASCII(alias.Destinations)
	if err != nil {
		return nil, err
	}
	alias.Destinations = asciiEmails

	requestBody, err := json.Marshal(alias)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	response := model.Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("UpdateAlias: %w", err)
	}

	return &response, nil
}

// DeleteAlias deletes an existing alias
func (c *MigaduClient) DeleteAlias(ctx context.Context, domain string, localPart string) (*model.Alias, error) {
	ascii, err := idna.ToASCII(domain)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	reqURL, err := url.JoinPath(c.Endpoint, "domains", ascii, "aliases", url.PathEscape(localPart))
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	responseBody, err := c.doRequest(request)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	response := model.Alias{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("DeleteAlias: %w", err)
	}

	return &response, nil
}
