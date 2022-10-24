package identityzone

import (
	"github.com/foundcloudry/terraform-provider-uaa/uaa/api"
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
)

// Mapper methods for mapping API objects to TF resources

func MapIdentityZoneToResource(identityZone *api.IdentityZone, data *schema.ResourceData) {

	data.SetId(identityZone.Id)
	data.Set(fields.AccountChooserEnabled.String(), identityZone.Config.AccountChooserEnabled)
	data.Set(fields.IsActive.String(), identityZone.IsActive)
	data.Set(fields.Name.String(), identityZone.Name)
	data.Set(fields.SubDomain.String(), identityZone.SubDomain)

	if identityZone.Config != nil {
		data.Set(fields.ClientSecretPolicy.String(), mapIdentityZoneClientSecretPolicyToInterface(identityZone.Config.ClientSecretPolicy))
		data.Set(fields.CorsPolicy.String(), mapIdentityZoneCorsPolicyToInterface(identityZone.Config.CorsPolicy))
		data.Set(fields.IdpDiscoveryEnabled.String(), &identityZone.Config.IdpDiscoveryEnabled)
		data.Set(fields.InputPrompts.String(), mapIdentityZoneInputPromptsToInterface(identityZone.Config.InputPrompts))
		data.Set(fields.IssuerUrl.String(), &identityZone.Config.IssuerUrl)
		data.Set(fields.SamlConfig.String(), mapIdentityZoneSamlConfigToInterface(identityZone.Config.Saml))
		data.Set(fields.TokenPolicy.String(), mapIdentityZoneTokenPolicyToInterface(identityZone.Config.TokenPolicy))

		if identityZone.Config.MfaConfig != nil {
			data.Set(fields.MfaEnabled.String(), identityZone.Config.MfaConfig.IsEnabled)
			data.Set(fields.MfaIdentityProviders.String(), identityZone.Config.MfaConfig.IdentityProviders)
		}

		if identityZone.Config.Links != nil {
			if identityZone.Config.Links != nil {
				data.Set(fields.HomeRedirectUrl.String(), identityZone.Config.Links.HomeRedirect)

				if identityZone.Config.Links.Logout != nil {
					data.Set(fields.LogoutRedirectParam.String(), identityZone.Config.Links.Logout.RedirectParameterName)
					data.Set(fields.LogoutRedirectUrl.String(), identityZone.Config.Links.Logout.RedirectUrl)
					data.Set(fields.LogoutAllowedRedirectUrls.String(), identityZone.Config.Links.Logout.AllowedRedirectUrls)
				}

				if identityZone.Config.Links.SelfService != nil {
					data.Set(fields.SelfServeEnabled.String(), identityZone.Config.Links.SelfService.Enabled)
					data.Set(fields.SelfServeSignupUrl.String(), identityZone.Config.Links.SelfService.SignupUrl)
					data.Set(fields.SelfServePasswordResetUrl.String(), identityZone.Config.Links.SelfService.PasswordResetUrl)
				}
			}
		}

		if identityZone.Config.UserConfig != nil {
			data.Set(fields.DefaultUserGroups.String(), &identityZone.Config.UserConfig.DefaultGroups)
		}

		if identityZone.Config.Branding != nil {
			data.Set(fields.Branding.String(), mapIdentityZoneBrandingToInterface(identityZone.Config.Branding))
		}
	}
}

func mapIdentityZoneBrandingToInterface(data *api.IdentityZoneBrandingConfig) []map[string]interface{} {
	if data == nil {
		return nil
	}

	return []map[string]interface{}{{
		brandingfields.BannerBackgroundColor.String(): data.Banner.BackgroundColor,
		brandingfields.BannerLogo.String():            data.Banner.Logo,
		brandingfields.BannerText.String():            data.Banner.Text,
		brandingfields.BannerTextColor.String():       data.Banner.TextColor,
		brandingfields.BannerUrl.String():             data.Banner.Url,
		brandingfields.CompanyName.String():           data.CompanyName,
		brandingfields.CompanyLogo.String():           data.CompanyLogo,
		brandingfields.Favicon.String():               data.Favicon,
		brandingfields.FooterText.String():            data.FooterText,
		brandingfields.FooterLinks.String():           data.Favicon,
		brandingfields.FooterLinks.String():           mapIdentityZoneBrandingFooterLinksToInterface(data),
	}}
}

func mapIdentityZoneBrandingFooterLinksToInterface(data *api.IdentityZoneBrandingConfig) (footerLinks []map[string]interface{}) {

	for name, url := range data.FooterLinks {
		footerLinks = append(footerLinks, map[string]interface{}{
			footerlinkfields.Name.String(): name,
			footerlinkfields.Url.String():  url,
		})
	}

	return footerLinks
}

func mapIdentityZoneCorsPolicyToInterface(data *api.IdentityZoneCorsPolicy) []map[string]interface{} {

	if data == nil {
		return nil
	}

	return []map[string]interface{}{
		mapIdentityZoneCorsConfigurationToInterface(corsconfignames.Default, data.DefaultConfiguration),
		mapIdentityZoneCorsConfigurationToInterface(corsconfignames.Xhr, data.XhrConfiguration),
	}
}

func mapIdentityZoneCorsConfigurationToInterface(name corsconfignames.CorsConfigName, data *api.IdentityZoneCorsConfig) map[string]interface{} {

	if data == nil {
		return nil
	}

	return map[string]interface{}{
		corsconfigfields.AllowedOrigins.String():        data.AllowedOrigins,
		corsconfigfields.AllowedOriginPatterns.String(): data.AllowedOriginPatterns,
		corsconfigfields.AllowedUris.String():           data.AllowedUris,
		corsconfigfields.AllowedUriPatterns.String():    data.AllowedUriPatterns,
		corsconfigfields.AllowedHeaders.String():        data.AllowedHeaders,
		corsconfigfields.AllowedMethods.String():        data.AllowedMethods,
		corsconfigfields.AllowedCredentials.String():    data.AllowedCredentials,
		corsconfigfields.Name.String():                  name.String(),
		corsconfigfields.MaxAge.String():                data.MaxAge,
	}
}

func mapIdentityZoneSamlConfigToInterface(data *api.IdentityZoneSamlConfig) []map[string]interface{} {

	if data == nil {
		return nil
	}

	return []map[string]interface{}{{
		samlconfigfields.ActiveKeyId.String():              data.ActiveKeyId,
		samlconfigfields.AssertionTtlSeconds.String():      data.AssertionTtlSeconds,
		samlconfigfields.DisableInResponseToCheck.String(): data.DisableInResponseToCheck,
		samlconfigfields.EntityId.String():                 data.EntityId,
		samlconfigfields.IsAssertionSigned.String():        data.IsAssertionSigned,
		samlconfigfields.IsRequestSigned.String():          data.IsRequestSigned,
		samlconfigfields.Key.String():                      mapIdentityZoneSamlKeysToInterface(data.Keys),
		samlconfigfields.WantAssertionSigned.String():      data.WantAssertionSigned,
		samlconfigfields.WantAuthRequestSigned.String():    data.WantAuthnRequestSigned,
	}}
}

func mapIdentityZoneSamlKeysToInterface(data map[string]*api.IdentityZoneSamlKey) (keys []map[string]interface{}) {

	for name, key := range data {
		keys = append(keys, map[string]interface{}{
			samlkeyfields.Certificate.String(): key.Certificate,
			samlkeyfields.Name.String():        name,
		})
	}

	return keys
}

func mapIdentityZoneClientSecretPolicyToInterface(data *api.IdentityZoneClientSecretPolicy) []map[string]interface{} {

	if data == nil {
		return nil
	}

	return []map[string]interface{}{{
		clientsecretpolicyfields.MaxLength.String():         data.MaxLength,
		clientsecretpolicyfields.MinDigits.String():         data.MinDigit,
		clientsecretpolicyfields.MinLength.String():         data.MinLength,
		clientsecretpolicyfields.MinLowerCaseChars.String(): data.MinLowerCaseCharacter,
		clientsecretpolicyfields.MinSpecialChars.String():   data.MinSpecialCharacter,
		clientsecretpolicyfields.MinUpperCaseChars.String(): data.MinUpperCaseCharacter,
	}}
}

func mapIdentityZoneTokenPolicyToInterface(data *api.IdentityZoneTokenPolicy) []map[string]interface{} {

	if data == nil {
		return nil
	}

	return []map[string]interface{}{{
		tokenpolicyfields.AccessTokenTtl.String():       data.AccessTokenTtl,
		tokenpolicyfields.ActiveKeyId.String():          data.ActiveKeyId,
		tokenpolicyfields.IsJwtRevocable.String():       data.IsJwtRevocable,
		tokenpolicyfields.IsRefreshTokenUnique.String(): data.IsRefreshTokenUnique,
		tokenpolicyfields.RefreshTokenFormat.String():   data.RefreshTokenFormat,
		tokenpolicyfields.RefreshTokenTtl.String():      data.RefreshTokenTtl,
	}}
}

func mapIdentityZoneInputPromptsToInterface(data []*api.InputPrompt) (prompts []map[string]interface{}) {

	for _, prompt := range data {
		prompts = append(prompts, map[string]interface{}{
			inputpromptfields.Name.String():  prompt.Name,
			inputpromptfields.Type.String():  prompt.Type,
			inputpromptfields.Value.String(): prompt.Value,
		})
	}

	return prompts
}

// Mapper methods for mapping TF resources to API objects

func MapResourceToIdentityZone(data *schema.ResourceData) *api.IdentityZone {

	return &api.IdentityZone{
		Id:        data.Id(),
		IsActive:  data.Get(fields.IsActive.String()).(bool),
		Name:      data.Get(fields.Name.String()).(string),
		SubDomain: data.Get(fields.SubDomain.String()).(string),
		Config:    mapResourceToIdentityZoneConfig(data),
	}
}

func mapResourceToIdentityZoneConfig(data *schema.ResourceData) *api.IdentityZoneConfig {

	config := &api.IdentityZoneConfig{
		AccountChooserEnabled: data.Get(fields.AccountChooserEnabled.String()).(bool),
		IdpDiscoveryEnabled:   data.Get(fields.IdpDiscoveryEnabled.String()).(bool),
		InputPrompts:          mapResourceToIdentityZoneInputPrompts(data),
		IssuerUrl:             data.Get(fields.IssuerUrl.String()).(string),
		MfaConfig:             mapResourceToIdentityZoneMfaConfig(data),
		UserConfig:            mapResourceToIdentityZoneUserConfig(data),
	}

	if clientSecretPolicy := mapResourceToIdentityZoneClientSecretPolicy(data); clientSecretPolicy != nil {
		config.ClientSecretPolicy = clientSecretPolicy
	}
	if corsPolicy := mapResourceToIdentityZoneCorsPolicy(data); corsPolicy != nil {
		config.CorsPolicy = corsPolicy
	}
	if tokenPolicy := mapResourceToIdentityZoneTokenPolicy(data); tokenPolicy != nil {
		config.TokenPolicy = tokenPolicy
	}
	if samlConfig := mapResourceToIdentityZoneSamlConfig(data); samlConfig != nil {
		config.Saml = samlConfig
	}
	if branding := mapResourceToIdentityZoneBrandingConfig(data); branding != nil {
		config.Branding = branding
	}

	return config
}

func mapResourceToIdentityZoneClientSecretPolicy(data *schema.ResourceData) *api.IdentityZoneClientSecretPolicy {

	if list := getFieldAsList(data, fields.ClientSecretPolicy.String()); len(list) == 1 {
		clientSecretPolicy := list[0]
		return &api.IdentityZoneClientSecretPolicy{
			MaxLength:             clientSecretPolicy[clientsecretpolicyfields.MaxLength.String()].(*int64),
			MinLength:             clientSecretPolicy[clientsecretpolicyfields.MinLength.String()].(*int64),
			MinUpperCaseCharacter: clientSecretPolicy[clientsecretpolicyfields.MinUpperCaseChars.String()].(*int64),
			MinLowerCaseCharacter: clientSecretPolicy[clientsecretpolicyfields.MinLowerCaseChars.String()].(*int64),
			MinDigit:              clientSecretPolicy[clientsecretpolicyfields.MinDigits.String()].(*int64),
			MinSpecialCharacter:   clientSecretPolicy[clientsecretpolicyfields.MinSpecialChars.String()].(*int64),
		}
	}

	return nil
}

func mapResourceToIdentityZoneCorsPolicy(data *schema.ResourceData) (corsPolicy *api.IdentityZoneCorsPolicy) {

	for _, p := range getFieldAsList(data, fields.CorsPolicy.String()) {
		policy := &api.IdentityZoneCorsConfig{
			AllowedOrigins:        p[corsconfigfields.AllowedOrigins.String()].([]string),
			AllowedOriginPatterns: p[corsconfigfields.AllowedOriginPatterns.String()].([]string),
			AllowedUris:           p[corsconfigfields.AllowedUris.String()].([]string),
			AllowedUriPatterns:    p[corsconfigfields.AllowedOriginPatterns.String()].([]string),
			AllowedHeaders:        p[corsconfigfields.AllowedHeaders.String()].([]string),
			AllowedMethods:        p[corsconfigfields.AllowedMethods.String()].([]string),
			AllowedCredentials:    p[corsconfigfields.AllowedCredentials.String()].(bool),
			MaxAge:                p[corsconfigfields.MaxAge.String()].(*int64),
		}

		switch p[corsconfigfields.Name.String()].(string) {
		case corsconfignames.Default.String():
			corsPolicy.DefaultConfiguration = policy
		case corsconfignames.Xhr.String():
			corsPolicy.XhrConfiguration = policy
		}
	}

	return corsPolicy
}

func mapResourceToIdentityZoneMfaConfig(data *schema.ResourceData) *api.MfaConfig {

	providersList := data.Get(fields.MfaIdentityProviders.String()).(*schema.Set).List()
	providers := make([]string, len(providersList))
	for i, v := range providersList {
		providers[i] = v.(string)
	}

	return &api.MfaConfig{
		IsEnabled:         data.Get(fields.MfaEnabled.String()).(bool),
		IdentityProviders: providers,
	}
}

func mapResourceToIdentityZoneInputPrompts(data *schema.ResourceData) (inputPrompts []*api.InputPrompt) {

	for _, prompt := range getFieldAsList(data, fields.InputPrompts.String()) {
		inputPrompts = append(inputPrompts, &api.InputPrompt{
			Name:  prompt[inputpromptfields.Name.String()].(string),
			Type:  prompt[inputpromptfields.Type.String()].(string),
			Value: prompt[inputpromptfields.Value.String()].(string),
		})
	}

	return inputPrompts
}

func mapResourceToIdentityZoneTokenPolicy(data *schema.ResourceData) *api.IdentityZoneTokenPolicy {

	if list := getFieldAsList(data, fields.TokenPolicy.String()); len(list) == 1 {
		tokenPolicy := list[0]
		return &api.IdentityZoneTokenPolicy{
			AccessTokenTtl:       tokenPolicy[tokenpolicyfields.AccessTokenTtl.String()].(*int64),
			RefreshTokenTtl:      tokenPolicy[tokenpolicyfields.RefreshTokenTtl.String()].(*int64),
			IsJwtRevocable:       tokenPolicy[tokenpolicyfields.IsJwtRevocable.String()].(bool),
			IsRefreshTokenUnique: tokenPolicy[tokenpolicyfields.IsRefreshTokenUnique.String()].(bool),
			RefreshTokenFormat:   tokenPolicy[tokenpolicyfields.RefreshTokenFormat.String()].(string),
			ActiveKeyId:          tokenPolicy[tokenpolicyfields.ActiveKeyId.String()].(string),
		}
	}

	return nil
}

func mapResourceToIdentityZoneSamlConfig(data *schema.ResourceData) *api.IdentityZoneSamlConfig {

	if list := getFieldAsList(data, fields.SamlConfig.String()); len(list) == 1 {
		samlConfig := list[0]
		keys := samlConfig[samlconfigfields.Key.String()].([]map[string]interface{})
		return &api.IdentityZoneSamlConfig{
			ActiveKeyId:              samlConfig[samlconfigfields.ActiveKeyId.String()].(string),
			AssertionTtlSeconds:      samlConfig[samlconfigfields.AssertionTtlSeconds.String()].(*int64),
			DisableInResponseToCheck: samlConfig[samlconfigfields.DisableInResponseToCheck.String()].(bool),
			EntityId:                 samlConfig[samlconfigfields.ActiveKeyId.String()].(string),
			IsAssertionSigned:        samlConfig[samlconfigfields.IsAssertionSigned.String()].(bool),
			IsRequestSigned:          samlConfig[samlconfigfields.IsRequestSigned.String()].(bool),
			Keys:                     mapResourceToIdentityZoneSamlKeys(keys),
			WantAssertionSigned:      samlConfig[samlconfigfields.WantAssertionSigned.String()].(bool),
			WantAuthnRequestSigned:   samlConfig[samlconfigfields.WantAuthRequestSigned.String()].(bool),
		}
	}

	return nil
}

func mapResourceToIdentityZoneSamlKeys(resourceKeys []map[string]interface{}) (keys map[string]*api.IdentityZoneSamlKey) {

	for _, key := range resourceKeys {
		name := key[samlkeyfields.Name.String()].(string)
		keys[name] = &api.IdentityZoneSamlKey{
			Certificate: key[samlkeyfields.Certificate.String()].(string),
		}
	}

	return keys
}

func mapResourceToIdentityZoneBrandingConfig(data *schema.ResourceData) *api.IdentityZoneBrandingConfig {

	if list := getFieldAsList(data, fields.Branding.String()); len(list) == 1 {
		branding := list[0]
		brandingConfig := &api.IdentityZoneBrandingConfig{
			Banner: &api.IdentityZoneBrandingBanner{
				BackgroundColor: branding[brandingfields.BannerBackgroundColor.String()].(string),
				Logo:            branding[brandingfields.BannerLogo.String()].(string),
				Text:            branding[brandingfields.BannerText.String()].(string),
				TextColor:       branding[brandingfields.BannerTextColor.String()].(string),
				Url:             branding[brandingfields.BannerUrl.String()].(string),
			},
			CompanyName: branding[brandingfields.CompanyName.String()].(string),
			CompanyLogo: branding[brandingfields.CompanyLogo.String()].(string),
			Favicon:     branding[brandingfields.Favicon.String()].(string),
			FooterText:  branding[brandingfields.FooterText.String()].(string),
		}

		if links, ok := branding[brandingfields.FooterLinks.String()].(*schema.Set); ok {
			brandingConfig.FooterLinks = mapResourceToIdentityZoneBrandingFooterLinks(links)
		}

		return brandingConfig
	}

	return nil
}

func mapResourceToIdentityZoneBrandingFooterLinks(data *schema.Set) map[string]string {

	linksList := data.List()
	footerLinks := make(map[string]string, len(linksList))
	for _, l := range linksList {
		if link, ok := l.(map[string]interface{}); ok {
			name := link[footerlinkfields.Name.String()].(string)
			footerLinks[name] = link[footerlinkfields.Url.String()].(string)
		}
	}

	return footerLinks
}

func mapResourceToIdentityZoneUserConfig(data *schema.ResourceData) *api.UserConfig {

	groupsList := data.Get(fields.DefaultUserGroups.String()).(*schema.Set).List()
	groups := make([]string, len(groupsList))
	for i, v := range groupsList {
		groups[i] = v.(string)
	}

	return &api.UserConfig{
		DefaultGroups: groups,
	}
}

func getFieldAsList(data *schema.ResourceData, field string) []map[string]interface{} {

	if value, isSet := data.GetOk(field); isSet {
		if valueAsList, canBeCast := value.([]map[string]interface{}); canBeCast {
			return valueAsList
		} else if valueAsList, canBeCast := value.([]interface{}); canBeCast {
			list := make([]map[string]interface{}, len(valueAsList))
			for i, v := range valueAsList {
				if listItem, canCast := v.(map[string]interface{}); canCast {
					list[i] = listItem
				}
			}
			return list
		}
	}
	return []map[string]interface{}{}
}
