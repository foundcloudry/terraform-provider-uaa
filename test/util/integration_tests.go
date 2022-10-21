package util

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/test"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/api"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/envvars"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/provider"
	"net"
	"os"
	"regexp"
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

func VerifyEnvironmentVariablesAreSet(t *testing.T) {

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
		t.Fatal("Acceptance environment has not been set.")
	}
}

func WarnIfTestZoneSubDomainDoesNotResolve(t *testing.T) {

	uaaHost := os.Getenv(envvars.UaaAuthUrl.String())
	var re = regexp.MustCompile(`(https?://)?([^:]*)(:\d+)?(/.*)?`)
	testZoneHost := re.ReplaceAllString(uaaHost, test.UpdatedZoneId+`.${2}`)

	if _, err := net.LookupHost(testZoneHost); err != nil {
		warningMessageLines := []string{
			"",
			"*********************************************************************************************************",
			"**\t" + fmt.Sprintf("Error: %s", err),
			"**\t",
			"**\t" + fmt.Sprintf("Could not reolve test zone sub-domain `%s`.  This test will likely fail.", testZoneHost),
			"**\t",
			"**\t" + "If you're running UAA locally, try adding the following entry to `/etc/hosts` to resolve this issue:",
			"**\t" + fmt.Sprintf("\t127.0.0.1 %s", testZoneHost),
			"**\t",
			"**\t" + "If you're testing against a cloud or other non-local deployment, ensure that an identity zone",
			"**\t" + "with an `id` and `sub-domain` of `" + test.UpdatedZoneId + "` exists and that DNS is properly configured.",
			"*********************************************************************************************************",
		}
		t.Logf(strings.Join(warningMessageLines, "\n"))
	}
}
