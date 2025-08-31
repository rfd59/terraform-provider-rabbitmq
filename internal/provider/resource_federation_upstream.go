package provider

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"
)

func resourceFederationUpstream() *schema.Resource {
	return &schema.Resource{
		Description: "Federation Upstream --- The `rabbitmq_federation_upstream` resource creates and manages a federation upstream parameter.",
		Create:      CreateFederationUpstream,
		Read:        ReadFederationUpstream,
		Update:      UpdateFederationUpstream,
		Delete:      DeleteFederationUpstream,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the federation upstream.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"vhost": {
				Description: "The vhost to create the resource in.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			// "federation-upstream"
			"component": {
				Description: "Set to _federation-upstream_ by the underlying RabbitMQ provider. You do not set this attribute but will see it in state and plan output.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"definition": {
				Description: "The configuration of the federation upstream. The structure is described below.",
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// applicable to both federated exchanges and queues
						"uri": {
							Description: "The AMQP Uri for the upstream.\n~> **Note:** The Uri may contain sensitive information, such as a password.",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},

						"prefetch_count": {
							Description: "Maximum number of unacknowledged messages that may be in flight over a federation link at one time. Defaults to `1000`.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1000,
						},

						"reconnect_delay": {
							Description: "Time in seconds to wait after a network link goes down before attempting reconnection. Defaults to `5`.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     5,
						},

						"ack_mode": {
							Description: "Determines how the link should acknowledge messages. Possible values are `on-confirm`, `on-publish` and `no-ack`. Defaults to `on-confirm`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "on-confirm",
							ValidateFunc: validation.StringInSlice([]string{
								"on-confirm",
								"on-publish",
								"no-ack",
							}, false),
						},

						"trust_user_id": {
							Description: "Determines how federation should interact with the validated user-id feature. Default is `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						// applicable to federated exchanges only
						"exchange": {
							Description: "**Federated Exchanges Only**: The name of the upstream exchange.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"max_hops": {
							Description: "**Federated Exchanges Only**: Maximum number of federation links that messages can traverse before being dropped. Defaults to `1`.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
						},
						"expires": {
							Description: "**Federated Exchanges Only**: The expiry time (in milliseconds) after which an upstream queue for a federated exchange may be deleted if a connection to the upstream is lost.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"message_ttl": {
							Description: "**Federated Exchanges Only**: The expiry time (in milliseconds) for messages in the upstream queue for a federated exchange (see `expires`).",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						// applicable to federated queues only
						"queue": {
							Description: "**Federated Queues Only**: The name of the upstream queue.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func CreateFederationUpstream(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)
	defList := d.Get("definition").([]interface{})

	defMap, ok := defList[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to parse federation upstream definition")
	}

	if err := putFederationUpstream(rmqc, vhost, name, defMap); err != nil {
		return err
	}

	id := fmt.Sprintf("%s@%s", name, vhost)
	d.SetId(id)

	return ReadFederationUpstream(d, meta)
}

func ReadFederationUpstream(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := utils.ParseResourceId(d.Id())
	if err != nil {
		return err
	}

	upstream, err := rmqc.GetFederationUpstream(vhost, name)
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}

	log.Printf("[DEBUG] RabbitMQ: Federation upstream retrieved for %s: %#v", d.Id(), upstream)

	d.Set("name", upstream.Name)
	d.Set("vhost", upstream.Vhost)
	d.Set("component", upstream.Component)

	var uri string
	if len(upstream.Definition.Uri) > 0 {
		uri = upstream.Definition.Uri[0]
	}
	defMap := map[string]interface{}{
		"uri":             uri,
		"prefetch_count":  upstream.Definition.PrefetchCount,
		"reconnect_delay": upstream.Definition.ReconnectDelay,
		"ack_mode":        upstream.Definition.AckMode,
		"trust_user_id":   upstream.Definition.TrustUserId,
		"exchange":        upstream.Definition.Exchange,
		"max_hops":        upstream.Definition.MaxHops,
		"expires":         upstream.Definition.Expires,
		"message_ttl":     upstream.Definition.MessageTTL,
		"queue":           upstream.Definition.Queue,
	}

	defList := [1]map[string]interface{}{defMap}
	d.Set("definition", defList)

	return nil
}

func UpdateFederationUpstream(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := utils.ParseResourceId(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange("definition") {
		_, newDef := d.GetChange("definition")

		defList := newDef.([]interface{})
		defMap, ok := defList[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("unable to parse federation definition")
		}

		if err := putFederationUpstream(rmqc, vhost, name, defMap); err != nil {
			return err
		}
	}

	return ReadFederationUpstream(d, meta)
}

func DeleteFederationUpstream(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := utils.ParseResourceId(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] RabbitMQ: Attempting to delete federation upstream for %s", d.Id())

	resp, err := rmqc.DeleteFederationUpstream(vhost, name)
	log.Printf("[DEBUG] RabbitMQ: Federation upstream delete response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		// the upstream was automatically deleted
		return nil
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error deleting RabbitMQ federation upstream: %s", resp.Status)
	}

	return nil
}

func putFederationUpstream(rmqc *rabbithole.Client, vhost string, name string, defMap map[string]interface{}) error {
	definition := rabbithole.FederationDefinition{}

	log.Printf("[DEBUG] RabbitMQ: Attempting to put federation definition for %s@%s: %#v", name, vhost, defMap)

	if v, ok := defMap["uri"].(string); ok {
		definition.Uri = []string{v}
	}

	if v, ok := defMap["expires"].(int); ok {
		definition.Expires = v
	}

	if v, ok := defMap["message_ttl"].(int); ok {
		definition.MessageTTL = int32(v)
	}

	if v, ok := defMap["max_hops"].(int); ok {
		definition.MaxHops = v
	}

	if v, ok := defMap["prefetch_count"].(int); ok {
		definition.PrefetchCount = v
	}

	if v, ok := defMap["reconnect_delay"].(int); ok {
		definition.ReconnectDelay = v
	}

	if v, ok := defMap["ack_mode"].(string); ok {
		definition.AckMode = v
	}

	if v, ok := defMap["trust_user_id"].(bool); ok {
		definition.TrustUserId = v
	}

	if v, ok := defMap["exchange"].(string); ok {
		definition.Exchange = v
	}

	if v, ok := defMap["queue"].(string); ok {
		definition.Queue = v
	}

	log.Printf("[DEBUG] RabbitMQ: Attempting to declare federation upstream for %s@%s: %#v", name, vhost, definition)

	resp, err := rmqc.PutFederationUpstream(vhost, name, definition)
	log.Printf("[DEBUG] RabbitMQ: Federation upstream declare response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error creating RabbitMQ federation upstream: %s", resp.Status)
	}

	return nil
}
