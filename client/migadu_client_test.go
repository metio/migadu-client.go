/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package client_test

import (
	"github.com/metio/migadu-client.go/client"
	"time"
)

func newTestClient(endpoint string) *client.MigaduClient {
	username := "username"
	token := "token"
	c, err := client.New(&endpoint, &username, &token, 10*time.Second)
	if err != nil {
		panic(err)
	}
	return c
}
