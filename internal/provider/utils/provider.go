package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
)

func BuildResourceId(name, vhost string) string {
	return fmt.Sprintf("%s@%s", name, vhost)
}

// get the resource name and rabbitmq vhost from the resource id
func ParseResourceId(resourceId string) (name, vhost string, err error) {
	parts := strings.Split(resourceId, "@")
	if len(parts) != 2 {
		err = fmt.Errorf("unable to parse resource id: %s", resourceId)
		return
	}
	name = parts[0]
	vhost = parts[1]
	return
}

func FailApiResponse(err error, resp *http.Response, action string, name string) error {
	if err != nil {
		return fmt.Errorf("error %s RabbitMQ %s: %v", action, name, err)
	} else {
		return fmt.Errorf("error %s RabbitMQ %s: %s", action, name, resp.Status)
	}
}

func CheckDeletedResource(d *schema.ResourceData, err error) error {
	var errorResponse rabbithole.ErrorResponse
	if errors.As(err, &errorResponse) {
		if errorResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
	}
	return err
}

func GetArgumentValue(arg map[string]interface{}) (interface{}, error) {
	switch arg["type"].(string) {
	case "numeric":
		if value, err := strconv.ParseFloat(arg["value"].(string), 64); err != nil {
			return nil, fmt.Errorf("failed to parse number %q", arg["value"].(string))
		} else {
			return value, nil
		}
	case "boolean":
		if value, err := strconv.ParseBool(arg["value"].(string)); err != nil {
			return nil, fmt.Errorf("failed to parse boolean %q", arg["value"].(string))
		} else {
			return value, nil
		}
	case "list":
		return arg["value"].(string), nil
	default:
		return arg["value"].(string), nil
	}
}

func GetArgumentType(value interface{}) string {
	switch value.(type) {
	case int:
		return "numeric"
	case float64:
		return "numeric"
	case bool:
		return "boolean"
	default:
		return "string"
	}
}
