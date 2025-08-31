package resources_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core/resources"
	mock_test "github.com/rfd59/terraform-provider-rabbitmq/test/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExchange_CreateExchange_AlreadyExist(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{Read: mock_test.RabbitMQInfraMock_Exchange{Err: nil, Rec: nil}}

	// Test
	d := getResourseDataExchange_Basic(t)
	err := resources.CreateExchange(d, mock)

	// Assert the expected behavior
	require.Error(err)
	require.ErrorContains(err, "exchange already exists")
	require.ErrorContains(err, "myName")
	assert.Empty(d.Id())
}

func TestExchange_CreateExchange_DataError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{Read: mock_test.RabbitMQInfraMock_Exchange{Err: errors.New("exchange not found!"), Rec: nil}}

	// Test
	d := getResourseDataExchange_ArgumentError(t)
	err := resources.CreateExchange(d, mock)

	// Assert the expected behavior
	require.Error(err)
	require.ErrorContains(err, "failed to parse number")
	require.ErrorContains(err, "myValue")
	assert.Empty(d.Id())
}

func TestExchange_CreateExchange_ErrorDeclare(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{
		Read:   mock_test.RabbitMQInfraMock_Exchange{Err: errors.New("exchange not found!"), Rec: nil},
		Create: mock_test.RabbitMQInfraMock_Response{Err: errors.New("exchange not created!"), Res: nil},
	}

	// Test
	d := getResourseDataExchange_Basic(t)
	err := resources.CreateExchange(d, mock)

	// Assert the expected behavior
	require.Error(err)
	require.ErrorContains(err, "exchange not created")
	assert.Empty(d.Id())
}

func TestExchange_CreateExchange_Success(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{
		Read:   mock_test.RabbitMQInfraMock_Exchange{Err: errors.New("exchange not found!"), Rec: nil},
		Create: mock_test.RabbitMQInfraMock_Response{Err: nil, Res: &http.Response{StatusCode: 200}},
	}

	// Test
	d := getResourseDataExchange_Full(t)
	err := resources.CreateExchange(d, mock)

	// Assert the expected behavior
	require.NoError(err)
	assert.Equal("myName@myVhost", d.Id())
}

func TestExchange_ReadExchange_FailedId(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test
	d := getResourseDataExchange_Basic(t)
	err := resources.ReadExchange(d, nil)

	// Assert the expected behavior
	require.Error(err)
	require.ErrorContains(err, "unable to parse resource id")
	assert.Empty(d.Id())
}

func TestExchange_ReadExchange_ErrorGet(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{Read: mock_test.RabbitMQInfraMock_Exchange{Err: errors.New("mock error"), Rec: nil}}

	// Test
	d := getResourseDataExchange_Empty(t)
	d.SetId("myName@myVhost")
	err := resources.ReadExchange(d, mock)

	// Assert the expected behavior
	require.Error(err)
	require.ErrorContains(err, "mock error")
	assert.Empty(d.Get("name"))
}

func TestExchange_ReadExchange_Success(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{Read: mock_test.RabbitMQInfraMock_Exchange{Err: nil, Rec: &rabbithole.DetailedExchangeInfo{
		Name:       "myName",
		Vhost:      "MyVhost",
		Type:       "topic",
		Durable:    true,
		AutoDelete: true,
		Internal:   false,
		Arguments:  map[string]interface{}{"alternate-exchange": "myAlternateExchange", "myStringKey": "myStringValue", "myNumericKey": 12345, "myBooleanKey": true},
	}}}

	// Test
	d := getResourseDataExchange_Empty(t)
	d.SetId("myName@myVhost")
	err := resources.ReadExchange(d, mock)

	// Assert the expected behavior
	require.NoError(err)
	assert.Equal("myName", d.Get("name"))
	assert.Equal("MyVhost", d.Get("vhost"))
	assert.Equal("topic", d.Get("type"))
	assert.True(d.Get("durable").(bool))
	assert.True(d.Get("auto_delete").(bool))
	assert.False(d.Get("internal").(bool))
	assert.Equal("myAlternateExchange", d.Get("alternate_exchange"))
	set := d.Get("argument").(*schema.Set)
	assert.Len(set.List(), 3)
}

func TestExchange_DeleteExchange_FailedId(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test
	d := getResourseDataExchange_Basic(t)
	err := resources.DeleteExchange(d, nil)

	// Assert the expected behavior
	require.Error(err)
	require.ErrorContains(err, "unable to parse resource id")
	assert.Empty(d.Id())
}

func TestExchange_DeleteExchange_ErrorDelete(t *testing.T) {
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{Delete: mock_test.RabbitMQInfraMock_Response{Err: errors.New("mock error"), Res: nil}}

	// Test
	d := getResourseDataExchange_Empty(t)
	d.SetId("myName@myVhost")
	err := resources.DeleteExchange(d, mock)

	// Assert the expected behavior
	require.Error(err)
	require.ErrorContains(err, "mock error")
}

func TestExchange_DeleteExchange_Success(t *testing.T) {
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{Delete: mock_test.RabbitMQInfraMock_Response{Err: nil, Res: &http.Response{StatusCode: 200}}}

	// Test
	d := getResourseDataExchange_Empty(t)
	d.SetId("myName@myVhost")
	err := resources.DeleteExchange(d, mock)

	// Assert the expected behavior
	require.NoError(err)
}

func getResourseDataExchange_Basic(t *testing.T) *schema.ResourceData {
	raw := map[string]interface{}{
		"name":  "myName",
		"vhost": "myVhost",
		"type":  "direct",
	}

	return schema.TestResourceDataRaw(t, resources.Exchange(), raw)
}

func getResourseDataExchange_Full(t *testing.T) *schema.ResourceData {
	raw := map[string]interface{}{
		"name":               "myName",
		"vhost":              "myVhost",
		"type":               "direct",
		"durable":            false,
		"auto_delete":        true,
		"internal":           true,
		"alternate_exchange": "myAlternateExchange",
		"argument": []interface{}{map[string]interface{}{
			"key":   "myKey",
			"value": "myValue",
			"type":  "string",
		}},
	}

	return schema.TestResourceDataRaw(t, resources.Exchange(), raw)
}

func getResourseDataExchange_ArgumentError(t *testing.T) *schema.ResourceData {
	raw := map[string]interface{}{
		"name":  "myName",
		"vhost": "myVhost",
		"type":  "direct",
		"argument": []interface{}{map[string]interface{}{
			"key":   "myKey",
			"value": "myValue",
			"type":  "numeric",
		}},
	}

	return schema.TestResourceDataRaw(t, resources.Exchange(), raw)
}

func getResourseDataExchange_Empty(t *testing.T) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, resources.Exchange(), map[string]interface{}{})
}
