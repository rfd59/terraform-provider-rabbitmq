package mock_test

import (
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
)

type RabbitMQInfraMock struct {
	Read   RabbitMQInfraMock_Exchange
	Create RabbitMQInfraMock_Response
	Delete RabbitMQInfraMock_Response
}

type RabbitMQInfraMock_Exchange struct {
	Rec *rabbithole.DetailedExchangeInfo
	Err error
}

type RabbitMQInfraMock_Response struct {
	Res *http.Response
	Err error
}

func (i *RabbitMQInfraMock) GetExchange(vhost, exchange string) (rec *rabbithole.DetailedExchangeInfo, err error) {
	return i.Read.Rec, i.Read.Err
}

func (i *RabbitMQInfraMock) DeclareExchange(vhost, exchange string, info rabbithole.ExchangeSettings) (res *http.Response, err error) {
	return i.Create.Res, i.Create.Err
}

func (i *RabbitMQInfraMock) DeleteExchange(vhost, exchange string) (res *http.Response, err error) {
	return i.Delete.Res, i.Delete.Err
}
