/*
 * SPDX-FileCopyrightText: The migadu-client.go Authors
 * SPDX-License-Identifier: 0BSD
 */

package idn

import (
	"reflect"
	"testing"
)

func TestConvertEmailsToASCII(t *testing.T) {
	testCases := map[string]struct {
		emails []string
		want   []string
	}{
		"single": {
			emails: []string{
				"test@example.com",
			},
			want: []string{
				"test@example.com",
			},
		},
		"multiple": {
			emails: []string{
				"test@example.com",
				"test@another.com",
			},
			want: []string{
				"test@example.com",
				"test@another.com",
			},
		},
		"idna": {
			emails: []string{
				"test@hoß.de",
			},
			want: []string{
				"test@xn--ho-hia.de",
			},
		},
		"mixed": {
			emails: []string{
				"test@hoß.de",
				"test@example.com",
			},
			want: []string{
				"test@xn--ho-hia.de",
				"test@example.com",
			},
		},
		"punycode": {
			emails: []string{
				"test@xn--ho-hia.de",
			},
			want: []string{
				"test@xn--ho-hia.de",
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := ConvertEmailsToASCII(testCase.emails)
			if err != nil {
				t.Errorf("ConvertEmailsToASCII() error = %v", err)
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("ConvertEmailsToASCII() got = %v, want %v", got, testCase.want)
			}
		})
	}
}
