package client

import (
	"code.cloudfoundry.org/cli/cf/errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jlpospisil/terraform-provider-uaa/test"
	"github.com/jlpospisil/terraform-provider-uaa/test/util"
	"regexp"
	"testing"
)

// TODO: when dns is figured out for containerized tests, remove this constant; just use `test.UpdatedZoneId` directly
const updatedZoneId = test.DefaultZoneId

const clientResource = `
resource "uaa_client" "client1" {
    client_id = "my-name"
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
    client_secret = "mysecret"
}
`

const clientResourceUpdateSecret = `
resource "uaa_client" "client1" {
    client_id = "my-name"
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
    client_secret = "newsecret"
}
`

const clientResourceUpdateZone = `
resource "uaa_client" "client1" {
    client_id = "my-name"
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
    client_secret = "newsecret"
	zone_id = "` + updatedZoneId + `"
}
`

const clientResourceWithoutSecret = `
resource "uaa_client" "client2" {
    client_id = "my-name2"
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
}
`

const clientResourceWithScope = `
resource "uaa_client" "client3" {
    client_id = "my-name-scope"
    scope = ["uaa.admin", "openid"]
    authorized_grant_types = [ "client_credentials" ]
    redirect_uri = [ "https://uaa.local.pcfdev.io/login" ]
    client_secret = "mysecret"
}
`

func TestAccClient_normal(t *testing.T) {
	ref := "uaa_client.client1"
	clientid := "my-name"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testClientDestroyed(clientid),
			Steps: []resource.TestStep{
				{
					Config: clientResource,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref, test.DefaultZoneId),
						testAccCheckValidSecret(ref, "mysecret", test.DefaultZoneId),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						util.TestCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						util.TestCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
					),
				},
				{
					Config: clientResourceUpdateSecret,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref, test.DefaultZoneId),
						testAccCheckValidSecret(ref, "newsecret", test.DefaultZoneId),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						util.TestCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						util.TestCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
					),
				},
				{
					Config: clientResourceUpdateZone,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref, updatedZoneId),
						testAccCheckValidSecret(ref, "newsecret", updatedZoneId),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						util.TestCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						util.TestCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
					),
				},
			},
		})
}

func TestAccClient_scope(t *testing.T) {
	ref := "uaa_client.client3"
	clientid := "my-name-scope"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testClientDestroyed(clientid),
			Steps: []resource.TestStep{
				{
					Config: clientResourceWithScope,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckClientExists(ref, test.DefaultZoneId),
						testAccCheckValidSecret(ref, "mysecret", test.DefaultZoneId),
						resource.TestCheckResourceAttr(ref, "client_id", clientid),
						util.TestCheckResourceSet(ref, "authorized_grant_types", []string{"client_credentials"}),
						util.TestCheckResourceSet(ref, "redirect_uri", []string{"https://uaa.local.pcfdev.io/login"}),
						util.TestCheckResourceSet(ref, "scope", []string{"openid", "uaa.admin"}),
					),
				},
			},
		})
}

func TestAccClient_createError(t *testing.T) {
	clientId := "my-name2"

	resource.Test(t,
		resource.TestCase{
			PreCheck:          func() { util.IntegrationTestPreCheck(t) },
			ProviderFactories: util.ProviderFactories,
			CheckDestroy:      testClientDestroyed(clientId),
			Steps: []resource.TestStep{
				{
					Config:      clientResourceWithoutSecret,
					ExpectError: regexp.MustCompile(".*Client secret is required for client_credentials.*"),
				},
			},
		})
}

func testAccCheckValidSecret(resource, secret, zoneId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}

		id := rs.Primary.ID
		subDomain := ""
		if zoneId != "uaa" {
			subDomain = zoneId
		}

		auth := util.UaaSession().AuthManager()
		if _, err := auth.GetClientToken(id, secret, subDomain); err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckClientExists(resource, zoneId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("client '%s' not found in terraform state", resource)
		}
		util.UaaSession().Log.DebugMessage("terraform state for resource '%s': %# v", resource, rs)

		id := rs.Primary.ID
		um := util.UaaSession().ClientManager()

		// check client exists
		_, err := um.GetClient(id, zoneId)
		if err != nil {
			return err
		}
		return nil
	}
}

func testClientDestroyed(id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, zoneId := range []string{test.DefaultZoneId, test.UpdatedZoneId} {
			if err := testClientDestroyedInZone(id, zoneId); err != nil {
				return err
			}
		}
		return nil
	}
}

func testClientDestroyedInZone(id, zoneId string) error {
	um := util.UaaSession().ClientManager()
	if _, err := um.FindByClientID(id, zoneId); err != nil {
		switch err.(type) {
		case *errors.ModelNotFoundError:
			return nil
		default:
			return err
		}
	}
	return fmt.Errorf("client with id '%s' still exists in cloud foundry", id)
}
