package identityzone

import (
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/brandingfields"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/clientsecretpolicyfields"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/corsconfigfields"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/corsconfignames"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/fields"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/footerlinkfields"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/inputpromptfields"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/samlconfigfields"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/samlkeyfields"
	"github.com/foundcloudry/terraform-provider-uaa/uaa/identityzone/tokenpolicyfields"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var identityZoneSchema = map[string]*schema.Schema{
	fields.AccountChooserEnabled.String(): {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	fields.DefaultUserGroups.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	fields.HomeRedirectUrl.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.InputPrompts.String(): {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: InputPromptSchema,
		},
	},
	fields.IdpDiscoveryEnabled.String(): {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	fields.IsActive.String(): {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	fields.IssuerUrl.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.LogoutRedirectParam.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	fields.LogoutRedirectUrl.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	fields.LogoutAllowedRedirectUrls.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	fields.MfaEnabled.String(): {
		Type:     schema.TypeBool,
		Optional: true,
		Computed: true,
	},
	fields.MfaIdentityProviders.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	fields.Name.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.SelfServeEnabled.String(): {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	fields.SelfServeSignupUrl.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.SelfServePasswordResetUrl.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	fields.SubDomain.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	fields.Branding.String(): {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: BrandingSchema,
		},
	},
	fields.ClientSecretPolicy.String(): {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: ClientSecretPolicySchema,
		},
	},
	fields.CorsPolicy.String(): {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 2,
		Elem: &schema.Resource{
			Schema: CorsPolicySchema,
		},
	},
	fields.SamlConfig.String(): {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: SamlConfigSchema,
		},
	},
	fields.TokenPolicy.String(): {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: TokenPolicySchema,
		},
	},
}

var BrandingSchema = map[string]*schema.Schema{
	brandingfields.CompanyName.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	brandingfields.CompanyLogo.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	brandingfields.Favicon.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	brandingfields.FooterText.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	brandingfields.FooterLinks.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: FooterLinkSchema,
		},
	},
}

var FooterLinkSchema = map[string]*schema.Schema{
	footerlinkfields.Name.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	footerlinkfields.Url.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
}

var ClientSecretPolicySchema = map[string]*schema.Schema{
	clientsecretpolicyfields.MaxLength.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	clientsecretpolicyfields.MinLength.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	clientsecretpolicyfields.MinUpperCaseChars.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	clientsecretpolicyfields.MinLowerCaseChars.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	clientsecretpolicyfields.MinDigits.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	clientsecretpolicyfields.MinSpecialChars.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
}

var TokenPolicySchema = map[string]*schema.Schema{
	tokenpolicyfields.AccessTokenTtl.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	tokenpolicyfields.RefreshTokenTtl.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	tokenpolicyfields.IsJwtRevocable.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	tokenpolicyfields.IsRefreshTokenUnique.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	tokenpolicyfields.RefreshTokenFormat.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "jwt",
	},
	tokenpolicyfields.ActiveKeyId.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
}

var SamlConfigSchema = map[string]*schema.Schema{
	samlconfigfields.ActiveKeyId.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	samlconfigfields.AssertionTtlSeconds.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
	samlconfigfields.DisableInResponseToCheck.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	samlconfigfields.EntityId.String(): {
		Type:     schema.TypeString,
		Optional: true,
	},
	samlconfigfields.IsAssertionSigned.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	samlconfigfields.IsRequestSigned.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	samlconfigfields.Key.String(): {
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: SamlConfigKeySchema,
		},
	},
	samlconfigfields.WantAssertionSigned.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	samlconfigfields.WantAuthRequestSigned.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
}

var SamlConfigKeySchema = map[string]*schema.Schema{
	samlkeyfields.Certificate.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	samlkeyfields.Name.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
}

var CorsPolicySchema = map[string]*schema.Schema{
	corsconfigfields.AllowedOrigins.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	corsconfigfields.AllowedOriginPatterns.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	corsconfigfields.AllowedUris.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	corsconfigfields.AllowedUriPatterns.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	corsconfigfields.AllowedHeaders.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	corsconfigfields.AllowedMethods.String(): {
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	corsconfigfields.AllowedCredentials.String(): {
		Type:     schema.TypeBool,
		Optional: true,
	},
	corsconfigfields.Name.String(): {
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(corsconfignames.CorsConfigNames, false),
	},
	corsconfigfields.MaxAge.String(): {
		Type:     schema.TypeInt,
		Optional: true,
	},
}

var InputPromptSchema = map[string]*schema.Schema{
	inputpromptfields.Name.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
	inputpromptfields.Type.String(): {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "text",
	},
	inputpromptfields.Value.String(): {
		Type:     schema.TypeString,
		Required: true,
	},
}

// The only required field for looking up an existing identity zone is the `name`.  All other fields should be optional
// and computed.  We can iterate over the resource schema and change those properties to avoid managing two schemas
// that are otherwise identical.
var dataSourceSchema = mapSchemaForDataSource(identityZoneSchema)

func mapSchemaForDataSource(originalSchema map[string]*schema.Schema) map[string]*schema.Schema {

	dsSchema := map[string]*schema.Schema{}

	for k, v := range originalSchema {
		isName := k == fields.Name.String()
		dsSchema[k] = &schema.Schema{
			Type:     v.Type,
			Required: isName,
			Computed: !isName,
			Elem:     v.Elem,
		}
		if v.Type == schema.TypeList {
			if elem, ok := v.Elem.(*schema.Resource); ok {
				dsSchema[k].Elem = &schema.Resource{
					Schema: mapSchemaForDataSource(elem.Schema),
				}
			}
		}
	}

	return dsSchema
}
