package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"golang.org/x/mod/semver"
)

type QueueResource struct {
	Name       string
	Vhost      string
	AutoDelete bool
	Durable    bool
	Arguments  map[string]interface{}
}

func (q *QueueResource) RequiredCreate(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		settings {}
	}`, data.ResourceType, data.ResourceLabel, q.Name)
}

func (q *QueueResource) OptionalUpdate(data TestData) string {
	q.AutoDelete = true
	q.Durable = true

	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		settings {
			auto_delete = %t
			durable = %t
		}
	}`, data.ResourceType, data.ResourceLabel, q.Name, q.AutoDelete, q.Durable)
}

func (q *QueueResource) OptionalUpdateArgument(data TestData) string {
	arg := "myKey"
	q.Arguments = map[string]interface{}{arg: "myValue"}

	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		settings {
			auto_delete = %t
			durable = %t
			arguments = {
				"%s" = "%s"
			}
		}
	}`, data.ResourceType, data.ResourceLabel, q.Name, q.AutoDelete, q.Durable, arg, q.Arguments[arg])
}

func (q *QueueResource) OptionalUpdateArgumentJson(data TestData) string {
	arg := "myKey"
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		settings {
			auto_delete = %t
			durable = %t
			arguments_json = jsonencode({
      			"%s" = "%s"
    		})
		}
	}`, data.ResourceType, data.ResourceLabel, q.Name, q.AutoDelete, q.Durable, arg, q.Arguments[arg])
}

func (q *QueueResource) XQueueTypeArgument(data TestData) string {
	arg := "x-queue-type"
	q.Arguments = map[string]interface{}{arg: "classic"}

	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		settings {
			arguments = {
				"%s" = "%s"
			}
		}
	}`, data.ResourceType, data.ResourceLabel, q.Name, arg, q.Arguments[arg])
}

func (q *QueueResource) ErrorBothArgumentsType(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		settings {
			arguments = {
				"k1" = "v1"
			}
			durable = true
			arguments_json = jsonencode({
      			"k2" = "v2"
    		})
		}
	}`, data.ResourceType, data.ResourceLabel, q.Name)
}

func (q *QueueResource) VhostDefaultQueueType_Step1(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		vhost = rabbitmq_vhost.test.name
		settings {
			durable = %t
		}
	}
	
	resource "rabbitmq_vhost" "test" {
 		name = "%s"
		default_queue_type = "quorum"
	}`, data.ResourceType, data.ResourceLabel, q.Name, q.Durable, q.Vhost)
}

func (q *QueueResource) VhostDefaultQueueType_Step2(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		vhost = rabbitmq_vhost.test.name
		settings {
			auto_delete = %t
			durable = %t
			arguments = {
				"x-queue-type" = "stream"
			}
		}
	}

	resource "rabbitmq_vhost" "test" {
 		name = "%s"
		default_queue_type = "quorum"
	}`, data.ResourceType, data.ResourceLabel, q.Name, q.AutoDelete, q.Durable, q.Vhost)
}

func (q *QueueResource) AlredayExist(data TestData) string {
	return fmt.Sprintf(`
	resource "%s" "%s" {
		name = "%s"
		settings {}
	}
	
	resource "%s" "already_exist" {
		name = "%s"
		settings {}
	}
	
	`, data.ResourceType, data.ResourceLabel, q.Name, data.ResourceType, q.Name)
}

func (q *QueueResource) DataSource(data TestData) string {
	return fmt.Sprintf(`
	data "%s" "%s" {
		name = "%s"
	}`, data.ResourceType, data.ResourceLabel, q.Name)
}

func (q QueueResource) ExistsInRabbitMQ() error {
	rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
	myQueue, err := rmqc.GetQueue(q.Vhost, q.Name)

	if err != nil {
		return fmt.Errorf("error retrieving queue '%s@%s': %#v", q.Name, q.Vhost, err)
	}
	if myQueue.Name != q.Name {
		return fmt.Errorf("queue name is not equal. Actual: '%s' Expected: %s", myQueue.Name, q.Name)
	}
	if myQueue.Vhost != q.Vhost {
		return fmt.Errorf("queue vhost is not equal. Actual: '%s' Expected: %s", myQueue.Vhost, q.Vhost)
	}
	if myQueue.AutoDelete != rabbithole.AutoDelete(q.AutoDelete) {
		return fmt.Errorf("queue autodelete is not equal. Actual: '%t' Expected: %t", myQueue.AutoDelete, q.AutoDelete)
	}
	if myQueue.Durable != q.Durable {
		return fmt.Errorf("queue durable is not equal. Actual: '%t' Expected: %t", myQueue.Durable, q.Durable)
	}
	// Comment this check because for RabbitMQ < 3.13, the state is set some seconds after to have created the queue.
	// if myQueue.Status != "running" {
	// 	return fmt.Errorf("queue status is not running. Actual: '%s'", myQueue.Status)
	// }
	if len(q.Arguments) != 0 {
		if len(myQueue.Arguments) != len(q.Arguments) {
			return fmt.Errorf("queue arguments length is not equal. Actual: '%d' Expected: %d [%v]", len(myQueue.Arguments), len(q.Arguments), myQueue.Arguments)
		}

		for k := range myQueue.Arguments {
			if myQueue.Arguments[k] != q.Arguments[k] {
				return fmt.Errorf("queue argument value for '%s' is not equal. Actual: '%v' Expected: %v", k, myQueue.Arguments[k], q.Arguments[k])
			}
		}
	}

	return nil
}

func (q QueueResource) CheckQueueTypeInRabbitMQ(queue_type string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
		myQueue, err := rmqc.GetQueue(q.Vhost, q.Name)
		if err != nil {
			return fmt.Errorf("error retrieving queue '%s@%s': %#v", q.Name, q.Vhost, err)
		}
		if myQueue.Type != queue_type {
			return fmt.Errorf("[%s@%s] queue type is not correct. Actual: '%s' Expected: %s", q.Name, q.Vhost, myQueue.Type, queue_type)
		}
		return nil
	}
}

func (q QueueResource) CheckDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
		vhost, err := rmqc.GetQueue(q.Vhost, q.Name)
		if err != nil && err.(rabbithole.ErrorResponse).StatusCode != 404 {
			return fmt.Errorf("error retrieving queue '%s@%s': %#v", q.Name, q.Vhost, err)
		}

		if vhost != nil {
			return fmt.Errorf("queue still exists: %s@%s", q.Name, q.Vhost)
		}

		return nil
	}
}

func (q QueueResource) SetDataSourceQueue(t *testing.T) {
	info := rabbithole.QueueSettings{}
	rmqc := TestAcc.Client(t)

	resp, err := rmqc.DeclareQueue(q.Vhost, q.Name, info)
	if err != nil || resp.StatusCode >= 400 {
		t.Errorf("Failed to init the test!")
	}
}

func (q QueueResource) DelDataSourceQueue(t *testing.T) {
	rmqc := TestAcc.Client(t)

	resp, err := rmqc.DeleteQueue(q.Vhost, q.Name)
	if err != nil || resp.StatusCode >= 400 {
		t.Errorf("Failed to reset the test!")
	}
}

// The 'default_queue_type' settings, for a virtual host, is implemented since RabbitMQ 3.10.
func (q QueueResource) SkipTestVhostDefaultQueueType(t *testing.T) {
	if semver.Compare("v"+TestAcc.Version, "v3.10") < 0 {
		t.Skip("Skipping testing: 'default_queue_type' settings (for a virtual host) is implemented since RabbitMQ 3.10!")
	}
}
