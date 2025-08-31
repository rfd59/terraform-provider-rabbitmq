package utils_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"
	"github.com/stretchr/testify/assert"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
)

func TestProvider_BuildResourceId(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("myName@myVhost", utils.BuildResourceId("myName", "myVhost"))
}

func TestProvider_ParseResourceId(t *testing.T) {
	assert := assert.New(t)

	type testExpectedStruct struct {
		name  string
		vhost string
		err   bool
	}
	type testCaseStruct struct {
		input    string
		expected testExpectedStruct
	}

	for id, testCase := range map[string]testCaseStruct{
		"Ok":       {input: "myName@myVhost", expected: testExpectedStruct{name: "myName", vhost: "myVhost", err: false}},
		"Error #1": {input: "myNamemyVhost", expected: testExpectedStruct{name: "", vhost: "", err: true}},
		"Error #2": {input: "my@Namemy@Vhost", expected: testExpectedStruct{name: "", vhost: "", err: true}},
		"Error #3": {input: "@", expected: testExpectedStruct{name: "", vhost: "", err: false}},
	} {
		t.Run(id, func(t *testing.T) {
			name, vhost, err := utils.ParseResourceId(testCase.input)

			assert.Equal(testCase.expected.name, name)
			assert.Equal(testCase.expected.vhost, vhost)
			assert.Equal(testCase.expected.err, err != nil)
		})
	}
}

func TestProvider_FailApiResponse(t *testing.T) {
	assert := assert.New(t)

	type testCaseStruct struct {
		err      error
		resp     *http.Response
		action   string
		name     string
		expected error
	}

	for id, testCase := range map[string]testCaseStruct{
		"#1": {err: errors.New("test error"), resp: &http.Response{}, action: "myAction", name: "myName", expected: errors.New("error myAction RabbitMQ myName: test error")},
		"#2": {err: nil, resp: &http.Response{Status: "500 Internal Server"}, action: "myAction", name: "myName", expected: errors.New("error myAction RabbitMQ myName: 500 Internal Server")},
	} {
		t.Run(id, func(t *testing.T) {
			data := utils.FailApiResponse(testCase.err, testCase.resp, testCase.action, testCase.name)

			assert.Equal(testCase.expected, data)
		})
	}
}

func TestProvider_CheckDeletedResource(t *testing.T) {
	assert := assert.New(t)

	type testExpectedStruct struct {
		id  string
		err bool
	}
	type testCaseStruct struct {
		id       string
		err      error
		expected testExpectedStruct
	}

	for id, testCase := range map[string]testCaseStruct{
		"#1": {id: "myId", err: errors.New("test error"), expected: testExpectedStruct{id: "myId", err: true}},
		"#2": {id: "myId", err: rabbithole.ErrorResponse{StatusCode: 200}, expected: testExpectedStruct{id: "myId", err: true}},
		"#3": {id: "myId", err: rabbithole.ErrorResponse{StatusCode: 404}, expected: testExpectedStruct{id: "", err: false}},
	} {
		t.Run(id, func(t *testing.T) {
			d := &schema.ResourceData{}
			d.SetId(testCase.id)

			data := utils.CheckDeletedResource(d, testCase.err)

			assert.Equal(testCase.expected.id, d.Id())
			assert.Equal(testCase.expected.err, data != nil)
		})
	}
}

func TestProvider_GetArgumentValue(t *testing.T) {
	assert := assert.New(t)

	type testExpectedStruct struct {
		value interface{}
		err   bool
	}
	type testCaseStruct struct {
		arg      map[string]interface{}
		expected testExpectedStruct
	}

	for id, testCase := range map[string]testCaseStruct{
		"string":        {arg: map[string]interface{}{"type": "string", "value": "myString"}, expected: testExpectedStruct{value: "myString", err: false}},
		"boolean":       {arg: map[string]interface{}{"type": "boolean", "value": "true"}, expected: testExpectedStruct{value: true, err: false}},
		"numeric int":   {arg: map[string]interface{}{"type": "numeric", "value": "12345"}, expected: testExpectedStruct{value: float64(12345), err: false}},
		"numeric float": {arg: map[string]interface{}{"type": "numeric", "value": "123.45"}, expected: testExpectedStruct{value: 123.45, err: false}},
		"list":          {arg: map[string]interface{}{"type": "list", "value": "myString, true, 12345"}, expected: testExpectedStruct{value: "myString, true, 12345", err: false}},
		"other":         {arg: map[string]interface{}{"type": "OtherType", "value": "12345"}, expected: testExpectedStruct{value: "12345", err: false}},
		"error boolean": {arg: map[string]interface{}{"type": "boolean", "value": "NotBooleanValue"}, expected: testExpectedStruct{value: nil, err: true}},
		"error numeric": {arg: map[string]interface{}{"type": "numeric", "value": "NotNumericValue"}, expected: testExpectedStruct{value: nil, err: true}},
	} {
		t.Run(id, func(t *testing.T) {

			data, err := utils.GetArgumentValue(testCase.arg)

			assert.Equal(testCase.expected.value, data)
			assert.Equal(testCase.expected.err, err != nil)
		})
	}
}

func TestProvider_GetArgumentType(t *testing.T) {
	assert := assert.New(t)

	type testCaseStruct struct {
		value    interface{}
		expected string
	}

	for id, testCase := range map[string]testCaseStruct{
		"string":        {value: "myString", expected: "string"},
		"boolean":       {value: true, expected: "boolean"},
		"numeric int":   {value: 12345, expected: "numeric"},
		"numeric float": {value: 123.45, expected: "numeric"},
		"other":         {value: errors.New("test"), expected: "string"},
	} {
		t.Run(id, func(t *testing.T) {

			data := utils.GetArgumentType(testCase.value)

			assert.Equal(testCase.expected, data)
		})
	}
}
