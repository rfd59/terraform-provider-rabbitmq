package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider/utils"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description: "User --- The `rabbitmq_user` resource creates and manages a user.",
		Create:      CreateUser,
		Update:      UpdateUser,
		Read:        ReadUser,
		Delete:      DeleteUser,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the user.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"password": {
				Description: "The password of the user.\n~> **Note:** The value of this argument is plain-text so make sure to secure where this is defined.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"tags": {
				Description: "Which permission model to apply to the user. Valid options are: `management`, `policymaker`, `monitoring`, and `administrator`.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"max_connections": {
				Description: "To limit how many connection a user can open.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},
			"max_channels": {
				Description: "To limit how many channels, in total, a user can open.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
			},
		},
	}
}

func CreateUser(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	name := d.Get("name").(string)

	// Check if already exists
	_, not_found := rmqc.GetUser(name)
	if not_found == nil {
		return fmt.Errorf("error creating RabbitMQ user '%s': user already exists", name)
	}

	userSettings := rabbithole.UserSettings{
		Password: d.Get("password").(string),
		Tags:     userTagsToString(d),
	}

	limits := make(rabbithole.UserLimitsValues)

	if v, ok := d.GetOk("max_connections"); ok {
		if (len(v.(string))) > 0 {
			v_int, err := strconv.Atoi(v.(string))
			if err != nil {
				return fmt.Errorf("error converting 'max_connections' to int: %#v", v)
			}
			limits["max-connections"] = v_int
		}
	}

	if v, ok := d.GetOk("max_channels"); ok {
		if (len(v.(string))) > 0 {
			v_int, err := strconv.Atoi(v.(string))
			if err != nil {
				return fmt.Errorf("error converting 'max_channels' to int: %#v", v)
			}
			limits["max-channels"] = v_int
		}
	}

	resp, err := rmqc.PutUser(name, userSettings)
	if err != nil || resp.StatusCode >= 400 {
		return utils.FailApiResponse(err, resp, "creating", "user")
	}

	if len(limits) > 0 {
		resp, err = rmqc.PutUserLimits(name, limits)
		if err != nil || resp.StatusCode >= 400 {
			return utils.FailApiResponse(err, resp, "creating", "user limits")
		}
	}

	d.SetId(name)
	return ReadUser(d, meta)
}

func ReadUser(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)

	user, err := rmqc.GetUser(d.Id())
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}
	d.Set("name", user.Name)

	if len(user.Tags) > 0 {
		var tagList []string
		for _, v := range user.Tags {
			if v != "" {
				tagList = append(tagList, v)
			}
		}
		if len(tagList) > 0 {
			d.Set("tags", tagList)
		}
	}

	myUserLimits, err := rmqc.GetUserLimits(d.Id())
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}

	if len(myUserLimits) > 0 {
		if val, ok := myUserLimits[0].Value["max-connections"]; ok {
			d.Set("max_connections", strconv.Itoa(val))
		} else {
			d.Set("max_connections", nil)
		}

		if val, ok := myUserLimits[0].Value["max-channels"]; ok {
			d.Set("max_channels", strconv.Itoa(val))
		} else {
			d.Set("max_channels", nil)
		}
	}

	return nil
}

func UpdateUser(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)
	name := d.Id()

	userSettings := rabbithole.UserSettings{
		Password: d.Get("password").(string),
		Tags:     userTagsToString(d),
	}
	myUserLimits, err := rmqc.GetUserLimits(d.Id())
	if err != nil {
		return utils.CheckDeletedResource(d, err)
	}

	limits := make(rabbithole.UserLimitsValues)

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
			limits["max-connections"] = myUserLimits[0].Value["max-connections"]
		}
	}

	if _, ok := d.GetOk("max_channels"); ok {
		if d.HasChange("max_channels") {
			_, newMaxQueues := d.GetChange("max_channels")

			if v, ok := newMaxQueues.(string); ok {
				limits["max-channels"], err = strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("error converting 'max_channels' to int: %#v", v)
				}
			}
		} else {
			limits["max-channels"] = myUserLimits[0].Value["max-channels"]
		}
	}

	resp, err := rmqc.PutUser(name, userSettings)
	if err != nil || resp.StatusCode >= 400 {
		return utils.FailApiResponse(err, resp, "updating", "user")
	}

	resp, err = rmqc.DeleteUserLimits(name, rabbithole.UserLimits{"max-connections", "max-channels"})
	if err != nil || resp.StatusCode >= 400 {
		return utils.FailApiResponse(err, resp, "updating", "user limits")
	}

	if len(limits) > 0 {
		resp, err = rmqc.PutUserLimits(name, limits)
		if err != nil || resp.StatusCode >= 400 {
			return utils.FailApiResponse(err, resp, "updating", "user limits")
		}
	}

	return ReadUser(d, meta)
}

func DeleteUser(d *schema.ResourceData, meta interface{}) error {
	rmqc := meta.(*rabbithole.Client)
	name := d.Id()

	resp, err := rmqc.DeleteUserLimits(name, rabbithole.UserLimits{"max-connections", "max-channels"})
	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode != 404) {
		return utils.FailApiResponse(err, resp, "deleting", "user limits")
	}

	resp, err = rmqc.DeleteUser(name)
	if err != nil || (resp.StatusCode >= 400 && resp.StatusCode != 404) {
		return utils.FailApiResponse(err, resp, "deleting", "user")
	}

	return nil
}

func userTagsToString(d *schema.ResourceData) rabbithole.UserTags {
	tagList := rabbithole.UserTags{}

	for _, v := range d.Get("tags").([]interface{}) {
		if tag, ok := v.(string); ok {
			tagList = append(tagList, tag)
		}
	}

	return tagList
}
