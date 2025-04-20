package acceptance

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
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
	if v.HasDefaultQueueTypeUpdateFeature() {
		v.DefaultQueueType = "quorum"
	} else {
		//Set the default value that was set during Create step
		v.DefaultQueueType = "classic"
	}
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
	if v.HasDescriptionUpdateFeature() {
		v.Description = data.RandomString()
	}
	if v.HasDefaultQueueTypeUpdateFeature() {
		v.DefaultQueueType = "stream"
	}
	v.Tracing = !v.Tracing
	return v.OptionalCreate(data)
}

func (v *VhostResource) OptionalUpdateLimits(data TestData) string {
	v.MaxConnections = data.RandomIntegerString()
	v.MaxQueues = data.RandomIntegerString()
	return v.OptionalCreate(data)
}

func (v *VhostResource) OptionalRemove(data TestData) string {
	if v.HasDescriptionUpdateFeature() {
		v.Description = ""
	}
	if v.HasDefaultQueueTypeUpdateFeature() {
		v.DefaultQueueType = ""
	}
	v.MaxConnections = ""
	v.MaxQueues = ""
	v.Tracing = false
	if v.HasDescriptionUpdateFeature() && v.HasDefaultQueueTypeUpdateFeature() {
		return fmt.Sprintf(`
		resource "%s" "%s" {
			name = "%s"
		}`, data.ResourceType, data.ResourceLabel, v.Name)
	} else {
		if !v.HasDescriptionUpdateFeature() {
			return fmt.Sprintf(`
			resource "%s" "%s" {
				name = "%s"
				description = "%s"
			}`, data.ResourceType, data.ResourceLabel, v.Name, v.Description)
		} else {
			//So No HasDefaultQueueTypeUpdateFeature
			return fmt.Sprintf(`
			resource "%s" "%s" {
				name = "%s"
				default_queue_type = "%s"
			}`, data.ResourceType, data.ResourceLabel, v.Name, v.DefaultQueueType)
		}

	}

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

func (v *VhostResource) ErrorDefaultQueueTypeAttribute(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		default_queue_type = "%s"
	}`, data.ResourceType, data.ResourceLabel, v.Name, v.DefaultQueueType)
}

func (v *VhostResource) DataSource(data TestData) string {
	return fmt.Sprintf(`
	data "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, v.Name)
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
		if v.HasDescriptionUpdateFeature() {
			return fmt.Errorf("vhost description is not equal. Actual: '%s' Expected: %s", myVhost.Description, v.Description)
		}
	}

	if hasDefaultQueueTypeFeature() {
		if len(v.DefaultQueueType) == 0 {
			if myVhost.DefaultQueueType != "classic" {
				if v.HasDefaultQueueTypeUpdateFeature() {
					return fmt.Errorf("vhost default_queue_type is not set to the default value. Actual: '%s' Expected: classic", myVhost.DefaultQueueType)
				}
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
	if hasDefaultQueueTypeFeature() {
		return []string{}
	} else {
		return []string{"default_queue_type"}
	}
}

func (v VhostResource) SetDataSourceVhost(t *testing.T) {
	settings := rabbithole.VhostSettings{}
	rmqc := TestAcc.Client(t)

	resp, err := rmqc.PutVhost(v.Name, settings)
	if err != nil || resp.StatusCode >= 400 {
		t.Errorf("Failed to init the test!")
	}
}

func (v VhostResource) DelDataSourceVhost(t *testing.T) {
	rmqc := TestAcc.Client(t)

	resp, err := rmqc.DeleteVhost(v.Name)
	if err != nil || resp.StatusCode >= 400 {
		t.Errorf("Failed to reset the test!")
	}
}

// 'Description' field can't be updated in 3.8. It was fixed in 3.9 and later
func (v VhostResource) HasDescriptionUpdateFeature() bool {
	return TestAcc.ValidFeature("3.9")
}

// 'DefaultQueueType' field can't be updated in 3.10. It was fixed in 3.11 and later
func (v VhostResource) HasDefaultQueueTypeUpdateFeature() bool {
	return TestAcc.ValidFeature("3.11")
}

// 'DefaultQueueType' field is present into RappbitMQ 3.10 and latter
func hasDefaultQueueTypeFeature() bool {
	return TestAcc.ValidFeature("3.10")
}
