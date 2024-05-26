package provider_test

import (
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFederationUpstream_importBasic(t *testing.T) {
	var upstream rabbithole.FederationUpstream
	resourceName := "rabbitmq_federation_upstream.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: testAccFederationUpstreamCheckDestroy(&upstream),
		Steps: []resource.TestStep{
			{
				Config: testAccFederationUpstream_create(),
				Check: testAccFederationUpstreamCheck(
					resourceName, &upstream,
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
