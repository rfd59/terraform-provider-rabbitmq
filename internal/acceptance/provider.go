package acceptance

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/rabbitmq"
)

var TestAcc TestProvider

type TestProvider struct {
	Providers map[string]*schema.Provider
	Provider  *schema.Provider
}

func init() {
	TestAcc.Provider = rabbitmq.Provider()
	TestAcc.Providers = map[string]*schema.Provider{
		"rabbitmq": TestAcc.Provider,
	}
}

func (TestProvider) PreCheck(t *testing.T) {
	for _, name := range []string{"RABBITMQ_ENDPOINT", "RABBITMQ_USERNAME", "RABBITMQ_PASSWORD"} {
		if v := os.Getenv(name); v == "" {
			t.Fatal("RABBITMQ_ENDPOINT, RABBITMQ_USERNAME and RABBITMQ_PASSWORD must be set for acceptance tests")
		}
	}
}
