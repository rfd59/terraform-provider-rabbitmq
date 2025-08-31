package provider

import (
	"fmt"
	"log"
	"strconv"

	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVhost() *schema.Resource {
	return &schema.Resource{
		Description: "Virtual Host --- The `rabbitmq_vhost` resource creates and manages a vhost.",
		Create:      CreateVhost,
		Read:        ReadVhost,
		Delete:      DeleteVhost,
		Update:      UpdateVhost,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the vhost.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "A friendly description.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},
			"default_queue_type": {
				Description:  "Default queue type for new queues. The available values are `classic`, `quorum` or `stream`. Defaults to `classic`.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Default:      "classic",
				ValidateFunc: validateDefaultQueueTypeAttribute,
			},
			"tracing": {
				Description: "To enable/disable tracing. Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     false,
			},
			"max_connections": {
				Description: "To limit the total number of concurrent client connections in vhost.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},
			"max_queues": {
				Description: "To limit the total number of queues in vhost.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},
		},
	}
}

func validateDefaultQueueTypeAttribute(val interface{}, key string) (warns []string, errs []error) {
	value := val.(string)

	// Define the allowed values
	allowedValues := map[string]struct{}{
		"classic": {},
		"quorum":  {},
		"stream":  {},
	}

	// Check if the value is in the allowed values
	if _, ok := allowedValues[value]; !ok {
		errs = append(errs, fmt.Errorf("%q must be one of [classic, quorum, stream], got: %s", key, value))
	}

	return warns, errs
}

func CreateVhost(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	vhost := d.Get("name").(string)

	// Check if already exists
	_, not_found := rmqc.GetVhost(vhost)
	if not_found == nil {
		return fmt.Errorf("error creating RabbitMQ vhost '%s': vhost already exists", vhost)
	}

	var settings rabbithole.VhostSettings

	if v, ok := d.Get("default_queue_type").(string); ok && v != "" {
		settings.DefaultQueueType = v
	}

	if v, ok := d.Get("description").(string); ok && v != "" {
		settings.Description = v
	}

	if v, ok := d.Get("tracing").(bool); ok {
		settings.Tracing = v
	}

	limits := make(rabbithole.VhostLimitsValues)

	if v, ok := d.GetOk("max_connections"); ok {
		if (len(v.(string))) > 0 {
			v_int, err := strconv.Atoi(v.(string))
			if err != nil {
				return fmt.Errorf("error converting 'max_connections' to int: %#v", v)
			}
			limits["max-connections"] = v_int
		}
	}

	if v, ok := d.GetOk("max_queues"); ok {
		if (len(v.(string))) > 0 {
			v_int, err := strconv.Atoi(v.(string))
			if err != nil {
				return fmt.Errorf("error converting 'max_queues' to int: %#v", v)
			}
			limits["max-queues"] = v_int
		}
	}

	resp, err := rmqc.PutVhost(vhost, settings)
	if err != nil || resp.StatusCode >= 400 {
		return utils.FailApiResponse(err, resp, "creating", "vhost")
	}

	if len(limits) > 0 {
		resp, err = rmqc.PutVhostLimits(vhost, limits)
		if err != nil || resp.StatusCode >= 400 {
			return utils.FailApiResponse(err, resp, "creating", "vhost limits")
		}
	}

	d.SetId(vhost)
	return ReadVhost(d, meta)
}

func ReadVhost(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	vhost, err := rmqc.GetVhost(d.Id())
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}
	d.Set("name", vhost.Name)

	if len(vhost.DefaultQueueType) > 0 && vhost.DefaultQueueType != "undefined" {
		d.Set("default_queue_type", vhost.DefaultQueueType)
	}

	if len(vhost.Description) > 0 {
		d.Set("description", vhost.Description)
	}

	myVhostLimits, err := rmqc.GetVhostLimits(d.Id())
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}

	if len(myVhostLimits) > 0 {
		if val, ok := myVhostLimits[0].Value["max-connections"]; ok {
			d.Set("max_connections", strconv.Itoa(val))
		} else {
			d.Set("max_connections", "") // set as unlimited
		}

		if val, ok := myVhostLimits[0].Value["max-queues"]; ok {
			d.Set("max_queues", strconv.Itoa(val))
		} else {
			d.Set("max_queues", "") // set as unlimited
		}
	}

	d.Set("tracing", vhost.Tracing)

	return nil
}

func UpdateVhost(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	vhost, err := rmqc.GetVhost(d.Id())
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}

	myVhostLimits, err := rmqc.GetVhostLimits(d.Id())
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}

	var settings rabbithole.VhostSettings
	limits := make(rabbithole.VhostLimitsValues)

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")

		if v, ok := newDescription.(string); ok && v != "" {
			settings.Description = v
		}
	} else {
		settings.Description = vhost.Description
	}

	if d.HasChange("default_queue_type") {
		_, newDefaultQueueType := d.GetChange("default_queue_type")

		if v, ok := newDefaultQueueType.(string); ok && v != "" {
			settings.DefaultQueueType = v
		}
	} else {
		settings.DefaultQueueType = vhost.DefaultQueueType
	}

	if d.HasChange("tracing") {
		_, newTracing := d.GetChange("tracing")

		if v, ok := newTracing.(bool); ok {
			settings.Tracing = v
		}
	} else {
		settings.Tracing = vhost.Tracing
	}

	if _, ok := d.GetOk("max_connections"); ok {
		if d.HasChange("max_connections") {
			_, newMaxConnections := d.GetChange("max_connections")

			if v, ok := newMaxConnections.(string); ok {
				limits["max-connections"], err = strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("error converting 'max_connections' to int: %#v", v)
				}
			}
		} else {
			limits["max-connections"] = myVhostLimits[0].Value["max-connections"]
		}
	}

	if _, ok := d.GetOk("max_queues"); ok {
		if d.HasChange("max_queues") {
			_, newMaxQueues := d.GetChange("max_queues")

			if v, ok := newMaxQueues.(string); ok {
				limits["max-queues"], err = strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("error converting 'max_queues' to int: %#v", v)
				}
			}
		} else {
			limits["max-queues"] = myVhostLimits[0].Value["max-queues"]
		}
	}

	resp, err := rmqc.PutVhost(vhost.Name, settings)
	log.Printf("[DEBUG] RabbitMQ: vhost creation response: %#v", resp)
	if err != nil || resp.StatusCode >= 400 {
		return utils.FailApiResponse(err, resp, "updating", "vhost")
	}

	resp, err = rmqc.DeleteVhostLimits(vhost.Name, rabbithole.VhostLimits{"max-connections", "max-queues"})
	if err != nil || resp.StatusCode >= 400 {
		return utils.FailApiResponse(err, resp, "updating", "vhost limits")
	}

	if len(limits) > 0 {
		resp, err = rmqc.PutVhostLimits(vhost.Name, limits)
		if err != nil || resp.StatusCode >= 400 {
			return utils.FailApiResponse(err, resp, "updating", "vhost limits")
		}
	}

	return ReadVhost(d, meta)
}

func DeleteVhost(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	resp, err := rmqc.DeleteVhost(d.Id())
	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode != 404) {
		return utils.FailApiResponse(err, resp, "deleting", "vhost limits")
	}

	return nil
}
