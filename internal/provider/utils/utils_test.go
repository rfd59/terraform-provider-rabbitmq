package utils_test

import (
	"testing"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"
	"github.com/stretchr/testify/assert"
)

func TestUtils_PercentEncodeSlashes(t *testing.T) {
	assert := assert.New(t)

	type testCaseStruct struct {
		input    string
		expected string
	}

	for id, testCase := range map[string]testCaseStruct{
		"#1": {input: "xx/zz", expected: "xx%2Fzz"},
		"#2": {input: "xxzz", expected: "xxzz"},
		"#3": {input: "/", expected: "%2F"},
	} {
		t.Run(id, func(t *testing.T) {
			data := utils.PercentEncodeSlashes(testCase.input)

			assert.Equal(testCase.expected, data)
		})
	}
}

func TestUtils_PercentDecodeSlashes(t *testing.T) {
	assert := assert.New(t)

	type testCaseStruct struct {
		input    string
		expected string
	}

	for id, testCase := range map[string]testCaseStruct{
		"#1": {input: "xx%2Fzz", expected: "xx/zz"},
		"#2": {input: "xxzz", expected: "xxzz"},
		"#3": {input: "%2F", expected: "/"},
	} {
		t.Run(id, func(t *testing.T) {
			data := utils.PercentDecodeSlashes(testCase.input)

			assert.Equal(testCase.expected, data)
		})
	}
}
