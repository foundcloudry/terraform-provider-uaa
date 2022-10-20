package identityzone

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jlpospisil/terraform-provider-uaa/test/util"
	"regexp"
	"testing"
)

const dataSource = `
data uaa_identity_zone "uaa" {
	name = "uaa"
}
`

const dataSourceNotFound = `
data uaa_identity_zone "not-found" {
	name = "not-found"
}
`

func TestDataSource_normal(t *testing.T) {
	ref := "data.uaa_identity_zone.uaa"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: dataSource,
					Check: resource.ComposeTestCheckFunc(
						checkIdentityZoneExists(ref),
						resource.TestCheckResourceAttr(ref, "id", "uaa"),
						resource.TestCheckResourceAttr(ref, "name", "uaa"),
						resource.TestCheckResourceAttr(ref, "is_active", "true"),
					),
				},
			},
		})
}

func TestGroupDataSourceClient_notFound(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:      dataSourceNotFound,
					ExpectError: regexp.MustCompile(".*Identity Zone not-found not found.*"),
				},
			},
		})
}

func checkIdentityZoneExists(resource string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}

		util.UaaSession().Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID

		identityZone, err := util.UaaSession().IdentityZoneManager().FindById(id)
		if err != nil {
			return err
		}
		if err := util.AssertSame(identityZone.Id, id); err != nil {
			return err
		}

		return nil
	}
}
