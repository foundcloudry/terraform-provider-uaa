package util

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/envvars"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/provider"
	"os"
	"strings"
	"testing"
)

var uaaProvider = provider.Provider()

var ProviderFactories = map[string]func() (*schema.Provider, error){

	"uaa": func() (*schema.Provider, error) {
		return uaaProvider, nil
	},
}

func UaaSession() *api.Session {
	return uaaProvider.Meta().(*api.Session)
}

func IntegrationTestPreCheck(t *testing.T) {

	if !testAccEnvironmentSet() {
		t.Fatal("Acceptance environment has not been set.")
	}
}

func testAccEnvironmentSet() bool {

	loginEndpoint := os.Getenv(envvars.UaaLoginUrl.String())
	authEndpoint := os.Getenv(envvars.UaaAuthUrl.String())
	clientID := os.Getenv(envvars.UaaClientId.String())
	clientSecret := os.Getenv(envvars.UaaClientSecret.String())

	if len(loginEndpoint) == 0 || len(authEndpoint) == 0 || len(clientID) == 0 || len(clientSecret) == 0 {
		envVars := strings.Join([]string{
			envvars.UaaLoginUrl.String(),
			envvars.UaaAuthUrl.String(),
			envvars.UaaClientId.String(),
			envvars.UaaClientSecret.String(),
		}, ", ")
		fmt.Println(envVars + " must be set when running tests.")
		return false
	}
	return true
}
