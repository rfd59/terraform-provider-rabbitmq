package core

import (
	"fmt"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"

	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CreateExchange(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)

	// Check if already exists
	_, not_found := rmqc.GetExchange(vhost, name)
	if not_found == nil {
		return fmt.Errorf("error creating RabbitMQ exchange '%s': exchange already exists", name)
	}

	// Build exchange info
	info, err := makeInfoExchange(d)
	if err != nil {
		return fmt.Errorf("error creating RabbitMQ exchange '%s': %v", name, err)
	}

	// Declare the exchange
	resp, err := rmqc.DeclareExchange(vhost, name, info)
	if err != nil || resp.StatusCode >= 400 {
		return utils.FailApiResponse(err, resp, "creating", "exchange")
	}

	//Save the id
	d.SetId(utils.BuildResourceId(name, vhost))

	return nil
}

func ReadExchange(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := utils.ParseResourceId(d.Id())
	if err != nil {
		return err
	}

	exchange, err := rmqc.GetExchange(vhost, name)
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}

	d.Set("name", exchange.Name)
	d.Set("vhost", exchange.Vhost)
	d.Set("type", exchange.Type)
	d.Set("durable", exchange.Durable)
	d.Set("auto_delete", exchange.AutoDelete)
	d.Set("internal", exchange.Internal)

	if len(exchange.Arguments) > 0 {
		if val := exchange.Arguments["alternate-exchange"]; val != nil {
			d.Set("alternate_exchange", val)
			delete(exchange.Arguments, "alternate-exchange")
		}

		var args []interface{}
		for key, value := range exchange.Arguments {
			args = append(args, map[string]interface{}{"key": key, "value": fmt.Sprintf("%v", value), "type": utils.GetArgumentType(value)})
		}
		d.Set("argument", args)
	}

	return nil
}

func DeleteExchange(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := utils.ParseResourceId(d.Id())
	if err != nil {
		return err
	}

	resp, err := rmqc.DeleteExchange(vhost, name)
	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode != 404) {
		return utils.FailApiResponse(err, resp, "deleting", "exchange")
	}

	return nil
}

func makeInfoExchange(d *schema.ResourceData) (info rabbithole.ExchangeSettings, err error) {
	info.Type = d.Get("type").(string)
	info.Durable = d.Get("durable").(bool)
	info.AutoDelete = d.Get("auto_delete").(bool)
	info.Internal = d.Get("internal").(bool)

	info.Arguments = make(map[string]interface{})
	if v := d.Get("alternate_exchange").(string); len(v) > 0 {
		info.Arguments["alternate-exchange"] = v
	}

	args := d.Get("argument").(*schema.Set)
	for _, v := range args.List() {
		arg := v.(map[string]interface{})
		if value, err := utils.GetArgumentValue(arg); err != nil {
			return rabbithole.ExchangeSettings{}, err
		} else {
			info.Arguments[arg["key"].(string)] = value
		}
	}

	return
}
