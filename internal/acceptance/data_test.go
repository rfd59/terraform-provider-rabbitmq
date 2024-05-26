// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance_test

import (
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"
)

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
