package datasources_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/core/datasources"
	mock_test "github.com/rfd59/terraform-provider-rabbitmq/test/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExchange_ReadExchange_Error(t *testing.T) {
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{Read: mock_test.RabbitMQInfraMock_Exchange{Err: errors.New("mock error"), Rec: nil}}

	// Test
	d := getResourseDataExchange_Basic(t)
	diag := datasources.ReadExchange(d, mock)

	// Assert the expected behavior
	require.True(diag.HasError())
}

func TestExchange_ReadExchange_Success(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Mock RabbitMQ Infrastructure
	mock := &mock_test.RabbitMQInfraMock{Read: mock_test.RabbitMQInfraMock_Exchange{Err: nil, Rec: &rabbithole.DetailedExchangeInfo{
		Name:       "myName",
		Vhost:      "myVhost",
		Type:       "topic",
		Durable:    true,
		AutoDelete: true,
		Internal:   false,
		Arguments:  map[string]interface{}{"alternate-exchange": "myAlternateExchange", "myStringKey": "myStringValue", "myNumericKey": 12345, "myBooleanKey": true},
	}}}

	// Test
	d := getResourseDataExchange_Basic(t)
	diag := datasources.ReadExchange(d, mock)

	// Assert the expected behavior
	require.False(diag.HasError())
	assert.Equal("myName@myVhost", d.Id())
	assert.Equal("topic", d.Get("type"))
	assert.True(d.Get("durable").(bool))
	assert.True(d.Get("auto_delete").(bool))
	assert.False(d.Get("internal").(bool))
	assert.Equal("myAlternateExchange", d.Get("alternate_exchange"))
	set := d.Get("argument").(*schema.Set)
	assert.Len(set.List(), 3)
}

func getResourseDataExchange_Basic(t *testing.T) *schema.ResourceData {
	raw := map[string]interface{}{
		"name":  "myName",
		"vhost": "myVhost",
	}

	return schema.TestResourceDataRaw(t, datasources.Exchange(), raw)
}
