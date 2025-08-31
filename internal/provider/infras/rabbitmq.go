package infras

import (
	"net/http"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
)

type IRabbitMQInfra interface {
	GetExchange(vhost, exchange string) (rec *rabbithole.DetailedExchangeInfo, err error)
	DeclareExchange(vhost, exchange string, info rabbithole.ExchangeSettings) (res *http.Response, err error)
	DeleteExchange(vhost, exchange string) (res *http.Response, err error)
}

type RabbitMQInfra struct {
	cli *rabbithole.Client
}

func NewRabbitMQInfra(rmqc *rabbithole.Client) *RabbitMQInfra {
	return &RabbitMQInfra{
		cli: rmqc,
	}
}

func (i *RabbitMQInfra) GetExchange(vhost, exchange string) (rec *rabbithole.DetailedExchangeInfo, err error) {
	return i.cli.GetExchange(vhost, exchange)
}

func (i *RabbitMQInfra) DeclareExchange(vhost, exchange string, info rabbithole.ExchangeSettings) (res *http.Response, err error) {
	return i.cli.DeclareExchange(vhost, exchange, info)
}

func (i *RabbitMQInfra) DeleteExchange(vhost, exchange string) (res *http.Response, err error) {
	return i.cli.DeleteExchange(vhost, exchange)
}
