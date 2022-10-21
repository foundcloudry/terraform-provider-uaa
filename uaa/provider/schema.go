package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/envvars"
	"github.com/jlpospisil/terraform-provider-uaa/uaa/provider/fields"
)

var Schema = map[string]*schema.Schema{
	fields.LoginEndpoint.String(): {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc(envvars.UaaLoginUrl.String(), ""),
	},
	fields.AuthEndpoint.String(): {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc(envvars.UaaAuthUrl.String(), ""),
	},
	fields.ClientId.String(): {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc(envvars.UaaClientId.String(), ""),
	},
	fields.ClientSecret.String(): {
		Type:        schema.TypeString,
		Optional:    true,
		DefaultFunc: schema.EnvDefaultFunc(envvars.UaaClientSecret.String(), ""),
	},
	fields.CaCert.String(): {
		Type:        schema.TypeString,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_CA_CERT", ""),
	},
	fields.SkipSslValidation.String(): {
		Type:        schema.TypeBool,
		Required:    true,
		DefaultFunc: schema.EnvDefaultFunc("UAA_SKIP_SSL_VALIDATION", "true"),
	},
}
