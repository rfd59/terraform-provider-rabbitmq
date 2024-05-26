package provider_test

import (
	"fmt"
	"strings"
	"testing"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"github.com/rfd59/terraform-provider-rabbitmq/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPermissions(t *testing.T) {
	var permissionInfo rabbithole.PermissionInfo
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: testAccPermissionsCheckDestroy(&permissionInfo),
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionsConfig_basic,
				Check: testAccPermissionsCheck(
					"rabbitmq_permissions.test", &permissionInfo,
				),
			},
			{
				Config: testAccPermissionsConfig_update,
				Check: testAccPermissionsCheck(
					"rabbitmq_permissions.test", &permissionInfo,
				),
			},
		},
	})
}

func TestAccPermissions_Empty(t *testing.T) {

	configCreate := `
resource "rabbitmq_vhost" "test" {
    name = "test"
}

resource "rabbitmq_user" "test" {
    name = "mctest"
    password = "foobar"
}

resource "rabbitmq_permissions" "test" {
    user = rabbitmq_user.test.name
    vhost = rabbitmq_vhost.test.name
    permissions {
        configure = ""
        write = ""
        read = ""
    }
}`
	var permissionInfo rabbithole.PermissionInfo
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAcc.PreCheck(t) },
		Providers:    acceptance.TestAcc.Providers,
		CheckDestroy: testAccPermissionsCheckDestroy(&permissionInfo),
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: testAccPermissionsCheck(
					"rabbitmq_permissions.test", &permissionInfo,
				),
			},
		},
	})
}

func testAccPermissionsCheck(rn string, permissionInfo *rabbithole.PermissionInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("permission id not set")
		}

		rmqc := acceptance.TestAcc.Provider.Meta().(*rabbithole.Client)
		perms, err := rmqc.ListPermissions()
		if err != nil {
			return fmt.Errorf("Error retrieving permissions: %s", err)
		}

		userParts := strings.Split(rs.Primary.ID, "@")
		for _, perm := range perms {
			if perm.User == userParts[0] && perm.Vhost == userParts[1] {
				permissionInfo = &perm
				return nil
			}
		}

		return fmt.Errorf("Unable to find permissions for user %s", rn)
	}
}

func testAccPermissionsCheckDestroy(permissionInfo *rabbithole.PermissionInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rmqc := acceptance.TestAcc.Provider.Meta().(*rabbithole.Client)
		perms, err := rmqc.ListPermissions()
		if err != nil {
			return fmt.Errorf("Error retrieving permissions: %s", err)
		}

		for _, perm := range perms {
			if perm.User == permissionInfo.User && perm.Vhost == permissionInfo.Vhost {
				return fmt.Errorf("Permissions still exist for user %s@%s", permissionInfo.User, permissionInfo.Vhost)
			}
		}

		return nil
	}
}

const testAccPermissionsConfig_basic = `
resource "rabbitmq_vhost" "test" {
    name = "test"
}

resource "rabbitmq_user" "test" {
    name = "mctest"
    password = "foobar"
    tags = ["administrator"]
}

resource "rabbitmq_permissions" "test" {
    user = "${rabbitmq_user.test.name}"
    vhost = "${rabbitmq_vhost.test.name}"
    permissions {
        configure = ".*"
        write = ".*"
        read = ".*"
    }
}`

const testAccPermissionsConfig_update = `
resource "rabbitmq_vhost" "test" {
    name = "test"
}

resource "rabbitmq_user" "test" {
    name = "mctest"
    password = "foobar"
    tags = ["administrator"]
}

resource "rabbitmq_permissions" "test" {
    user = "${rabbitmq_user.test.name}"
    vhost = "${rabbitmq_vhost.test.name}"
    permissions {
        configure = ".*"
        write = ".*"
        read = ""
    }
}`
