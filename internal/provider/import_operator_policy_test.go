package provider_test

import (
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOperatorPolicy_importBasic(t *testing.T) {
	resourceName := "rabbitmq_operator_policy.test"
	var operatorPolicy rabbithole.OperatorPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: testAccOperatorPolicyCheckDestroy(&operatorPolicy),
		Steps: []resource.TestStep{
			{
				Config: testAccOperatorPolicyConfig_basic,
				Check: testAccOperatorPolicyCheck(
					resourceName, &operatorPolicy,
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
