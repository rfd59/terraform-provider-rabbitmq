package infras_test

import (
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/infras"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRabbitMQ_GetExchange(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	infra := infras.NewRabbitMQInfra(&rabbithole.Client{})
	rec, err := infra.GetExchange("myVhost", "myExchange")

	require.Error(err)
	assert.Nil(rec)
}

func TestRabbitMQ_DeclareExchange(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	infra := infras.NewRabbitMQInfra(&rabbithole.Client{})
	res, err := infra.DeclareExchange("myVhost", "myExchange", rabbithole.ExchangeSettings{})

	require.Error(err)
	assert.Nil(res)
}

func TestRabbitMQ_DeleteExchange(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	infra := infras.NewRabbitMQInfra(&rabbithole.Client{})
	res, err := infra.DeleteExchange("myVhost", "myExchange")

	require.Error(err)
	assert.Nil(res)
}
