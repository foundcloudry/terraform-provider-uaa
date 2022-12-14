package group

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
	"github.com/foundcloudry/terraform-provider-uaa/test"
	"github.com/foundcloudry/terraform-provider-uaa/test/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"
)

const ref = "uaa_group.new-group"
const originalDisplayName = "new.group.for.tests"
const originalDescription = "A group used for testing group resource functionality"
const updatedDisplayName = "updated.display.name"
const updatedDescription = "An updated description for the group resource"

func createTestGroupResourceAttr(attribute, value string) string {
	if attribute == "" || value == "" {
		return ""
	}
	return `	` + attribute + ` = "` + value + `"`
}

func createTestGroupResource(displayName, description, zoneId string) string {
	return `resource uaa_group "new-group" {
		` + createTestGroupResourceAttr("display_name", displayName) + `
		` + createTestGroupResourceAttr("description", description) + `
		` + createTestGroupResourceAttr("zone_id", zoneId) + `
	}`
}

func TestGroupResource_normal(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.VerifyEnvironmentVariablesAreSet(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testGroupDestroyed(originalDisplayName),
			Steps: []resource.TestStep{
				{
					Config: createTestGroupResource(originalDisplayName, originalDescription, ""),
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", originalDisplayName),
						resource.TestCheckResourceAttr(ref, "description", originalDescription),
						resource.TestCheckResourceAttr(ref, "zone_id", test.DefaultZoneId),
					),
				},
				{
					Config: createTestGroupResource(originalDisplayName, updatedDescription, ""),
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", originalDisplayName),
						resource.TestCheckResourceAttr(ref, "description", updatedDescription),
						resource.TestCheckResourceAttr(ref, "zone_id", test.DefaultZoneId),
					),
				},
				{
					Config: createTestGroupResource(updatedDisplayName, updatedDescription, ""),
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", updatedDisplayName),
						resource.TestCheckResourceAttr(ref, "description", updatedDescription),
						resource.TestCheckResourceAttr(ref, "zone_id", test.DefaultZoneId),
					),
				},
				{
					Config: createTestGroupResource(updatedDisplayName, updatedDescription, test.UpdatedZoneId),
					Check: resource.ComposeTestCheckFunc(
						checkDataSourceGroupExists(ref),
						resource.TestCheckResourceAttrSet(ref, "id"),
						resource.TestCheckResourceAttr(ref, "display_name", updatedDisplayName),
						resource.TestCheckResourceAttr(ref, "description", updatedDescription),
						resource.TestCheckResourceAttr(ref, "zone_id", test.UpdatedZoneId),
					),
				},
			},
		})
}

func TestGroupResource_createError(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.VerifyEnvironmentVariablesAreSet(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testGroupDestroyed(ref),
			Steps: []resource.TestStep{
				{
					Config:      createTestGroupResource("", originalDescription, test.DefaultZoneId),
					ExpectError: regexp.MustCompile("The argument \"display_name\" is required, but no definition was found."),
				},
			},
		})
}

func testGroupDestroyed(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, zoneId := range []string{test.DefaultZoneId, test.UpdatedZoneId} {
			if err := testGroupDestroyedInZone(id, zoneId); err != nil {
				return err
			}
		}
		return nil
	}
}

func testGroupDestroyedInZone(id, zoneId string) error {

	gm := util.UaaSession().GroupManager()
	if _, err := gm.FindByDisplayName(id, zoneId); err != nil {
		switch err.(type) {
		case *errors.ModelNotFoundError:
			return nil
		default:
			return err
		}
	}
	return fmt.Errorf("group with id '%s' still exists in cloud foundry", id)
}
