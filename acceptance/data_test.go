// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance_test

import (
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/acceptance"
)

func TestRandomIntOfLength(t *testing.T) {
	td := acceptance.TestData{
		//RandomInteger: 112233445566779999,
	}

	cases := []struct {
		len      int
		expected int
	}{
		{
			len:      18,
			expected: 112233445566779999,
		},
		{
			len:      17,
			expected: 11223344556677999,
		},
		{
			len:      16,
			expected: 1122334455667799,
		},
		{
			len:      15,
			expected: 112233445566799,
		},
		{
			len:      14,
			expected: 11223344556699,
		},
		{
			len:      10,
			expected: 1122334499,
		},
		{
			len:      9,
			expected: 112233499,
		},
		{
			len:      8,
			expected: 11223399,
		},
	}

	for _, c := range cases {
		result := td.RandomIntOfLength(c.len)
		if result != c.expected {
			t.Fatalf("For length %d expected %d but got %d", c.len, c.expected, result)
		}
	}
}

func TestBuildArrayString(t *testing.T) {
	td := acceptance.TestData{}

	cases := []struct {
		input    []string
		expected string
	}{
		{
			input:    []string{},
			expected: "[]",
		},
		{
			input:    []string{"azerty"},
			expected: "[\"azerty\"]",
		},
		{
			input:    []string{"azerty", "ytreza", "qwerty", "ytrewq"},
			expected: "[\"azerty\",\"ytreza\",\"qwerty\",\"ytrewq\"]",
		},
	}

	for _, c := range cases {
		result := td.BuildArrayString(c.input)
		if result != c.expected {
			t.Fatalf("for array %#v, expected %s but got %s", c.input, c.expected, result)
		}
	}
}
