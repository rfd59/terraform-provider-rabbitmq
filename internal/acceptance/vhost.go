package acceptance

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"golang.org/x/mod/semver"
)

type VhostResource struct {
	Name             string
	Description      string
	DefaultQueueType string
	MaxConnections   string
	MaxQueues        string
	Tracing          bool
}

func (v *VhostResource) RequiredCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, v.Name)
}

func (v *VhostResource) RequiredUpdate(data TestData) string {
	v.Description = data.RandomString()
	v.DefaultQueueType = "quorum"
	v.Tracing = true

	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		description = "%s"
		default_queue_type = "%s"
		tracing = %t
	}`, data.ResourceType, data.ResourceLabel, v.Name, v.Description, v.DefaultQueueType, v.Tracing)
}

func (v *VhostResource) OptionalCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		description = "%s"
		default_queue_type = "%s"
		max_connections = "%s"
		max_queues = "%s"
		tracing = %t
	}`, data.ResourceType, data.ResourceLabel, v.Name, v.Description, v.DefaultQueueType, v.MaxConnections, v.MaxQueues, v.Tracing)
}

func (v *VhostResource) OptionalUpdate(data TestData) string {
	v.Description = data.RandomString()
	v.DefaultQueueType = "stream"
	v.Tracing = !v.Tracing
	return v.OptionalCreate(data)
}

func (v *VhostResource) OptionalUpdateLimits(data TestData) string {
	v.MaxConnections = data.RandomIntegerString()
	v.MaxQueues = data.RandomIntegerString()
	return v.OptionalCreate(data)
}

func (v *VhostResource) OptionalRemove(data TestData) string {
	v.Description = ""
	v.DefaultQueueType = ""
	v.MaxConnections = ""
	v.MaxQueues = ""
	v.Tracing = false
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, v.Name)
}
func (v *VhostResource) ErrorConvertingCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		max_connections = "%s"
		max_queues = "%s"
	}`, data.ResourceType, data.ResourceLabel, v.Name, v.MaxConnections, v.MaxQueues)
}

func (u *VhostResource) ErrorConvertingUpdate(data TestData, connections string, queues string) string {
	u.MaxConnections = connections
	u.MaxQueues = queues
	return u.ErrorConvertingCreate(data)
}

func (v VhostResource) ExistsInRabbitMQ() error {

	rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
	myVhost, err := rmqc.GetVhost(v.Name)
	if err != nil {
		return fmt.Errorf("error retrieving vhost '%s': %#v", v.Name, err)
	}
	if myVhost.Name != v.Name {
		return fmt.Errorf("vhost name is not equal. Actual: '%s' Expected: %s", myVhost.Name, v.Name)
	}
	if myVhost.Description != v.Description {
		return fmt.Errorf("vhost description is not equal. Actual: '%s' Expected: %s", myVhost.Description, v.Description)
	}

	if hasDefaultQueueTypeFeature(rmqc.Overview()) {
		if len(v.DefaultQueueType) == 0 {
			if myVhost.DefaultQueueType != "classic" {
				return fmt.Errorf("vhost default_queue_type is not set to the default value. Actual: '%s' Expected: classic", myVhost.DefaultQueueType)
			}
		} else {
			if myVhost.DefaultQueueType != v.DefaultQueueType {
				return fmt.Errorf("vhost default_queue_type is not equal. Actual: '%s' Expected: %s", myVhost.DefaultQueueType, v.DefaultQueueType)
			}
		}
	}

	myVhostLimits, err := rmqc.GetVhostLimits(v.Name)
	if err != nil {
		return fmt.Errorf("error retrieving vhost limit '%s': %#v", v.Name, err)
	}
	if len(myVhostLimits) == 0 {
		if v.MaxConnections != "" && v.MaxQueues != "" {
			return fmt.Errorf("vhost limit is not empty: %#v", myVhostLimits)
		}
	} else {
		if myVhostLimits[0].Value["max-connections"] == 0 && v.MaxConnections == "" {
			//It's OK (specific case).
		} else if strconv.Itoa(myVhostLimits[0].Value["max-connections"]) != v.MaxConnections {
			return fmt.Errorf("vhost limit 'max-connections' is not equal. Actual: '%d' Expected: %s", myVhostLimits[0].Value["max-connections"], v.MaxConnections)
		}
		if myVhostLimits[0].Value["max-queues"] == 0 && v.MaxQueues == "" {
			//It's OK (specific case).
		} else if strconv.Itoa(myVhostLimits[0].Value["max-queues"]) != v.MaxQueues {
			return fmt.Errorf("vhost limit 'max-queues' is not equal. Actual: '%d' Expected: %s", myVhostLimits[0].Value["max-queues"], v.MaxQueues)
		}
	}

	return nil
}

func (v VhostResource) CheckDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
		vhost, err := rmqc.GetVhost(v.Name)
		if err != nil && err.(rabbithole.ErrorResponse).StatusCode != 404 {
			return fmt.Errorf("error retrieving vhost '%s': %#v", v.Name, err)
		}

		if vhost != nil {
			return fmt.Errorf("vhost still exists: %s", v.Name)
		}

		return nil
	}
}

func (v VhostResource) ImportStateVerifyIgnore() []string {
	rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
	if hasDefaultQueueTypeFeature(rmqc.Overview()) {
		return []string{}
	} else {
		return []string{"default_queue_type"}
	}
}

func hasDefaultQueueTypeFeature(overview *rabbithole.Overview, err error) bool {
	// DefaultQueueType is into RappbitMQ 3.10 and latter
	if err != nil {
		return true
	} else {
		return semver.Compare("v"+overview.RabbitMQVersion, "v3.10") >= 0
	}
}
