package acceptance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
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
	if myQueue.Status != "running" {
		return fmt.Errorf("queue status is not running. Actual: '%s'", myQueue.Status)
	}
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

func (q QueueResource) CheckDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := TestAcc.Provider.Meta().(*rabbithole.Client)
		vhost, err := rmqc.GetQueue("/", q.Name)
		if err != nil && err.(rabbithole.ErrorResponse).StatusCode != 404 {
			return fmt.Errorf("error retrieving queue '%s@%s': %#v", q.Name, q.Vhost, err)
		}

		if vhost != nil {
			return fmt.Errorf("queue still exists: %s@%s", q.Name, q.Vhost)
		}

		return nil
	}
}
