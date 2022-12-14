package user

import (
	"fmt"
	"github.com/foundcloudry/terraform-provider-uaa/test"
	"github.com/foundcloudry/terraform-provider-uaa/test/util"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

const userDataResource = `

data "uaa_user" "admin-user" {
   name = "admin"
}
`

func TestUserDataSource(t *testing.T) {

	ref := "data.uaa_user.admin-user"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.VerifyEnvironmentVariablesAreSet(t) },
			ProviderFactories: util.ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: userDataResource,
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceUserExists(ref, test.DefaultZoneId),
						resource.TestCheckResourceAttr(
							ref, "name", "admin"),
					),
				},
			},
		})
}

func checkDataSourceUserExists(resource, zoneId string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("user '%s' not found in terraform state", resource)
		}

		util.UaaSession().Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID
		name := rs.Primary.Attributes["name"]

		var (
			err  error
			user api.UAAUser
		)

		user, err = util.UaaSession().UserManager().FindByUsername(name, zoneId)
		if err != nil {
			return err
		}
		if err := util.AssertSame(user.Id, id); err != nil {
			return err
		}

		return nil
	}
}
