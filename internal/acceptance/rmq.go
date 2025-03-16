package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/provider"
	"golang.org/x/mod/semver"
)

var TestAcc RMQ

type RMQ struct {
	Providers map[string]*schema.Provider
	Provider  *schema.Provider
	Version   string
}

func init() {
	TestAcc.Provider = provider.New()
	TestAcc.Providers = map[string]*schema.Provider{
		"rabbitmq": TestAcc.Provider,
	}
}

func (RMQ) PreCheck(t *testing.T) {
	rmqc := TestAcc.Client(t)

	overview, _ := rmqc.Overview()
	TestAcc.Version = overview.RabbitMQVersion
}

func (RMQ) ValidFeature(miniVersion string) bool {
	return semver.Compare("v"+TestAcc.Version, "v"+miniVersion) >= 0
}

func (RMQ) Client(t *testing.T) *rabbithole.Client {
	for _, name := range []string{"RABBITMQ_ENDPOINT", "RABBITMQ_USERNAME", "RABBITMQ_PASSWORD"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("RABBITMQ_ENDPOINT, RABBITMQ_USERNAME and RABBITMQ_PASSWORD must be set for acceptance tests")
		}
	}

	rmqc, err := rabbithole.NewClient(os.Getenv("RABBITMQ_ENDPOINT"), os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"))
	if err != nil {
		t.Fatal("Can't connect to RabbitMQ!")
	}

	return rmqc
}
