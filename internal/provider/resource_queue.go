package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
)

func resourceQueue() *schema.Resource {
	return &schema.Resource{
		Description: "The rabbitmq_queue resource creates and manages a queue.",
		Create:      CreateQueue,
		Read:        ReadQueue,
		Delete:      DeleteQueue,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the queue.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"vhost": {
				Description: "The vhost to create the resource in. Defaults to `/`.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/",
				ForceNew:    true,
			},

			"settings": {
				Description: "The settings of the queue. The structure is described below.",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"durable": {
							Description: "Whether the queue survives server restarts. Defaults to `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
						},

						"auto_delete": {
							Description: "Whether the queue will self-delete when all consumers have unsubscribed. Defaults to `false`.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
						},

						"arguments": {
							Description:   "Additional key/value settings for the queue. All values will be sent to RabbitMQ as a string. If you require non-string values, use `arguments_json`.\n~> **Note:** Either this or `arguments_json` must be specified but not both.",
							Type:          schema.TypeMap,
							Optional:      true,
							ConflictsWith: []string{"settings.0.arguments_json"},
							ForceNew:      true,
						},

						"arguments_json": {
							Description:      "A nested JSON string which contains additional settings for the queue. This is useful for when the arguments contain non-string values.\n~> **Note:** Either this or `arguments` must be specified but not both.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringIsJSON,
							ConflictsWith:    []string{"settings.0.arguments"},
							DiffSuppressFunc: structure.SuppressJsonDiff,
							ForceNew:         true,
						},
					},
				},
			},

			"type": {
				Description: "The queue type created. The value are `classic`, `quorum` or `stream`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func CreateQueue(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)
	vhost := d.Get("vhost").(string)

	// Check if already exists
	_, not_found := rmqc.GetQueue(vhost, name)
	if not_found == nil {
		return fmt.Errorf("Error creating RabbitMQ queue '%s': queue already exists", name)
	}

	settingsList := d.Get("settings").([]interface{})

	settingsMap, ok := settingsList[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Unable to parse settings")
	}

	// If arguments_json is used, unmarshal it into a generic interface
	// and use it as the "arguments" key for the queue.
	if v, ok := settingsMap["arguments_json"].(string); ok && v != "" {
		var arguments map[string]interface{}
		err := json.Unmarshal([]byte(v), &arguments)
		if err != nil {
			return err
		}

		delete(settingsMap, "arguments_json")
		settingsMap["arguments"] = arguments
	}

	if err := declareQueue(rmqc, vhost, name, settingsMap); err != nil {
		return err
	}

	id := fmt.Sprintf("%s@%s", name, vhost)
	d.SetId(id)

	return ReadQueue(d, meta)
}

func ReadQueue(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := parseResourceId(d)
	if err != nil {
		return err
	}

	queueSettings, err := rmqc.GetQueue(vhost, name)
	if err != nil {
		return checkDeleted(d, err)
	}

	log.Printf("[DEBUG] RabbitMQ: Queue retrieved for %s: %#v", d.Id(), queueSettings)

	d.Set("name", queueSettings.Name)
	d.Set("vhost", queueSettings.Vhost)
	d.Set("type", queueSettings.Type)

	e := make(map[string]interface{})
	e["durable"] = queueSettings.Durable
	e["auto_delete"] = queueSettings.AutoDelete

	// Check if "x-queue-type" was explicitly defined in Terraform configuration
	xQueueTypeDefined := false
	if args, ok := d.GetOk("settings.0.arguments"); ok {
		argsMap := args.(map[string]interface{})
		if _, exists := argsMap["x-queue-type"]; exists {
			xQueueTypeDefined = true
		}
	}
	if argsJson, ok := d.GetOk("settings.0.arguments_json"); ok {
		var jsonArgs map[string]interface{}
		if err := json.Unmarshal([]byte(argsJson.(string)), &jsonArgs); err == nil {
			if _, exists := jsonArgs["x-queue-type"]; exists {
				xQueueTypeDefined = true
			}
		}
	}

	// Delete "x-queue-type" if it was not defined in Terraform configuration to keep a clean terraform plan
	if !xQueueTypeDefined {
		delete(queueSettings.Arguments, "x-queue-type")
	}

	// The user may have used either `arguments` or `arguments_json` to populate this originally.
	// We need to preserve that decision here so that a subsequent Terraform plan for the
	// same configuration wouldn't produce an errant diff that moves the value from one
	// to the other without changing any values.
	// These two arguments are mutually exclusive due to ConflictsWith in the schema.
	// `arguments` cannot receive any values other than a string (d.Set will fail), therefore any drift
	// containing nonstring values AND the configuration originated from `arguments`,
	// will now be encoded to `arguments_json`.
	if _, ok := d.GetOk("settings.0.arguments_json"); ok || nonStringInArguments(queueSettings.Arguments) {
		bytes, err := json.Marshal(queueSettings.Arguments)
		if err != nil {
			return err
		}
		e["arguments_json"] = string(bytes)
	} else {
		e["arguments"] = queueSettings.Arguments
	}

	queue := make([]map[string]interface{}, 1)
	queue[0] = e

	return d.Set("settings", queue)
}

func DeleteQueue(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name, vhost, err := parseResourceId(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] RabbitMQ: Attempting to delete queue for %s", d.Id())

	resp, err := rmqc.DeleteQueue(vhost, name)
	log.Printf("[DEBUG] RabbitMQ: Queue delete response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		// the queue was automatically deleted
		return nil
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error deleting RabbitMQ queue '%s': %s", name, resp.Status)
	}

	return nil
}

func declareQueue(rmqc *rabbithole.Client, vhost string, name string, settingsMap map[string]interface{}) error {
	queueSettings := rabbithole.QueueSettings{}

	if v, ok := settingsMap["durable"].(bool); ok {
		queueSettings.Durable = v
	}

	if v, ok := settingsMap["auto_delete"].(bool); ok {
		queueSettings.AutoDelete = v
	}

	if v, ok := settingsMap["arguments"].(map[string]interface{}); ok {
		queueSettings.Arguments = v
	}

	log.Printf("[DEBUG] RabbitMQ: Attempting to declare queue for %s@%s: %#v", name, vhost, queueSettings)

	resp, err := rmqc.DeclareQueue(vhost, name, queueSettings)
	log.Printf("[DEBUG] RabbitMQ: Queue declare response: %#v", resp)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error declaring RabbitMQ queue '%s': %s", name, resp.Status)
	}

	return nil
}

func nonStringInArguments(args map[string]interface{}) bool {
	for _, val := range args {
		switch val.(type) {
		case string:
			continue
		default:
			return true
		}
	}
	return false
}
