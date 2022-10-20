package identityzone

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jlpospisil/terraform-provider-uaa/test/util"
	"regexp"
	"testing"
)

const ref = "uaa_identity_zone.new-test-zone"
const originalName = "New Test Zone"
const originalSubDomain = "original-int-test-zone"
const updatedName = "Updated Test Zone"
const updatedSubdomain = "updated-int-test-zone"

func createTestResourceAttr(attribute, value string) string {
	if attribute == "" || value == "" {
		return ""
	}
	return `	` + attribute + ` = "` + value + `"`
}

func createTestResource(name, subDomain string) string {
	return `resource uaa_identity_zone "new-test-zone" {
		` + createTestResourceAttr("name", name) + `
		` + createTestResourceAttr("sub_domain", subDomain) + `
		default_user_groups = ["openid"]
	}`
}

func TestResource_normal(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testCheckDestroyed(),
			Steps: []resource.TestStep{
				{
					Config: createTestResource(originalName, originalSubDomain),
					Check: resource.ComposeTestCheckFunc(
						checkIdentityZoneExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "name", originalName),
						resource.TestCheckResourceAttr(ref, "sub_domain", originalSubDomain),
					),
				},
				{
					Config: createTestResource(originalName, updatedSubdomain),
					Check: resource.ComposeTestCheckFunc(
						checkIdentityZoneExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "name", originalName),
						resource.TestCheckResourceAttr(ref, "sub_domain", updatedSubdomain),
					),
				},
				{
					Config: createTestResource(updatedName, updatedSubdomain),
					Check: resource.ComposeTestCheckFunc(
						checkIdentityZoneExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "name", updatedName),
						resource.TestCheckResourceAttr(ref, "sub_domain", updatedSubdomain),
					),
				},
			},
		})
}

func TestResource_createError(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testCheckDestroyed(),
			Steps: []resource.TestStep{
				{
					Config:      createTestResource("", originalSubDomain),
					ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
				},
				{
					Config:      createTestResource(originalName, ""),
					ExpectError: regexp.MustCompile("The argument \"sub_domain\" is required, but no definition was found."),
				},
			},
		})
}

func testCheckDestroyed() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, name := range []string{originalName, updatedName} {
			izm := util.UaaSession().IdentityZoneManager()
			zone, err := izm.FindByName(name)
			if zone != nil {
				return fmt.Errorf("identity zone with name '%s' still exists in cloud foundry", name)
			}
			if err != nil {
				switch err.(type) {
				case *errors.ModelNotFoundError:
				default:
					return err
				}
			}
		}

		return nil
	}
}
