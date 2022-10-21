package user

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jlpospisil/terraform-provider-uaa/test"
	"github.com/jlpospisil/terraform-provider-uaa/test/util"
	"testing"
)

const ldapUserResource = `

resource "uaa_user" "manager1" {
	name = "manager1@acme.com"
	origin = "ldap"
}
`

const userResourceWithGroups = `

resource "uaa_user" "admin-service-user" {
	name = "cf-admin"
	password = "qwerty"
	zone_id = "` + test.DefaultZoneId + `"
	given_name = "Build"
	family_name = "User"
	groups = [ "cloud_controller.admin", "scim.read", "scim.write" ]
}
`

const userResourceWithGroupsUpdate = `

resource "uaa_user" "admin-service-user" {
	name = "cf-admin"
	password = "asdfg"
	zone_id = "` + test.UpdatedZoneId + `"
	email = "cf-admin@acme.com"
	groups = [ "cloud_controller.admin", "clients.admin", "uaa.admin", "doppler.firehose" ]
}
`

func TestAccUser_LdapOrigin_normal(t *testing.T) {

	ref := "uaa_user.manager1"
	username := "manager1@acme.com"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testAccCheckUserDestroyed(username),
			Steps: []resource.TestStep{
				{
					Config: ldapUserResource,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckUserExists(ref, test.DefaultZoneId),
						resource.TestCheckResourceAttr(ref, "name", username),
						resource.TestCheckResourceAttr(ref, "origin", "ldap"),
						resource.TestCheckResourceAttr(ref, "email", username),
						resource.TestCheckResourceAttr(ref, "zone_id", test.DefaultZoneId),
					),
				},
			},
		})
}

func TestAccUser_WithGroups_normal(t *testing.T) {

	ref := "uaa_user.admin-service-user"
	username := "cf-admin"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testAccCheckUserDestroyed(username),
			Steps: []resource.TestStep{
				{
					Config: userResourceWithGroups,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckUserExists(ref, test.DefaultZoneId),
						resource.TestCheckResourceAttr(ref, "name", username),
						resource.TestCheckResourceAttr(ref, "password", "qwerty"),
						resource.TestCheckResourceAttr(ref, "email", username),
						resource.TestCheckResourceAttr(ref, "zone_id", test.DefaultZoneId),
						util.TestCheckResourceSet(ref, "groups", []string{
							"cloud_controller.admin",
							"scim.read",
							"scim.write",
						}),
					),
				},

				{
					Config: userResourceWithGroupsUpdate,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckUserExists(ref, test.UpdatedZoneId),
						resource.TestCheckResourceAttr(ref, "name", "cf-admin"),
						resource.TestCheckResourceAttr(ref, "password", "asdfg"),
						resource.TestCheckResourceAttr(ref, "email", "cf-admin@acme.com"),
						resource.TestCheckResourceAttr(ref, "zone_id", test.UpdatedZoneId),
						util.TestCheckResourceSet(ref, "groups", []string{
							"clients.admin",
							"cloud_controller.admin",
							"doppler.firehose",
							"uaa.admin",
						}),
					),
				},
			},
		})
}

func testAccCheckUserExists(resource, zoneId string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("user '%s' not found in terraform state", resource)
		}

		util.UaaSession().Log.DebugMessage(
			"terraform state for resource '%s': %# v",
			resource, rs)

		id := rs.Primary.ID
		attributes := rs.Primary.Attributes

		um := util.UaaSession().UserManager()
		user, err := um.GetUser(id, zoneId)
		if err != nil {
			return err
		}

		util.UaaSession().Log.DebugMessage(
			"retrieved user for resource '%s' with id '%s': %# v",
			resource, id, user)

		if err := util.AssertEquals(attributes, "name", user.Username); err != nil {
			return err
		}
		if err := util.AssertEquals(attributes, "origin", user.Origin); err != nil {
			return err
		}
		if err := util.AssertEquals(attributes, "given_name", user.Name.GivenName); err != nil {
			return err
		}
		if err := util.AssertEquals(attributes, "family_name", user.Name.FamilyName); err != nil {
			return err
		}
		if err := util.AssertEquals(attributes, "email", user.Emails[0].Value); err != nil {
			return err
		}

		var groups []interface{}
		for _, g := range user.Groups {
			isDefault, err := um.IsDefaultGroup(zoneId, g.Display)
			if err != nil {
				return err
			}
			if !isDefault {
				groups = append(groups, g.Display)
			}
		}
		if err := util.AssertSetEquals(attributes, "groups", groups); err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckUserDestroyed(username string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		for _, zoneId := range []string{test.DefaultZoneId, test.UpdatedZoneId} {
			if err := testCheckUserDoesNotExistInZone(username, zoneId); err != nil {
				return err
			}
		}
		return nil
	}
}

func testCheckUserDoesNotExistInZone(username, zoneId string) error {

	um := util.UaaSession().UserManager()
	if _, err := um.FindByUsername(username, zoneId); err != nil {
		switch err.(type) {
		case *errors.ModelNotFoundError:
			return nil
		default:
			return err
		}
	}
	return fmt.Errorf("user with username '%s' still exists in cloud foundry", username)
}
