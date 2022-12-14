package api

import (
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/cf/errors"
	"code.cloudfoundry.org/cli/cf/net"
	"fmt"
	"net/url"
)

type IdentityZoneManager struct {
	log *Logger
	api *UaaApi
}

func newIdentityZoneManager(config coreconfig.Reader, gateway net.Gateway, logger *Logger) (izm *IdentityZoneManager, err error) {

	api, err := newUaaApi(config, gateway)
	if err != nil {
		return
	}

	izm = &IdentityZoneManager{
		log: logger,
		api: api,
	}
	return
}

// CRUD methods

func (manager *IdentityZoneManager) Create(identityZone *IdentityZone) (*IdentityZone, error) {

	if err := manager.api.Post("/identity-zones", identityZone, &identityZone); err != nil {
		return nil, err
	}

	return identityZone, nil
}

func (manager *IdentityZoneManager) FindById(id string) (*IdentityZone, error) {

	path := fmt.Sprintf("/identity-zones/%s", id)
	identityZone := &IdentityZone{}
	err := manager.api.Get(path, identityZone)
	if err != nil {
		return nil, err
	}

	return identityZone, nil
}

func (manager *IdentityZoneManager) FindByName(name string) (*IdentityZone, error) {

	displayNameFilter := url.QueryEscape(fmt.Sprintf(`name Eq "%s"`, name))
	path := fmt.Sprintf("/identity-zones?filter=%s", displayNameFilter)

	identityZones := &[]IdentityZone{}
	err := manager.api.Get(path, identityZones)
	if err != nil {
		return nil, err
	}

	for _, identityZone := range *identityZones {
		if identityZone.Name == name {
			return &identityZone, nil
		}
	}

	return nil, errors.NewModelNotFoundError("Identity Zone", name)
}

func (manager *IdentityZoneManager) Update(id string, identityZone *IdentityZone) (*IdentityZone, error) {

	path := fmt.Sprintf("/identity-zones/%s", id)
	if err := manager.api.Put(path, identityZone, &identityZone); err != nil {
		return nil, err
	}

	return identityZone, nil
}

func (manager *IdentityZoneManager) Delete(id string) error {

	return manager.api.Delete(fmt.Sprintf("/identity-zones/%s", id))
}

// DTOs

type IdentityZone struct {
	Id        string              `json:"id"`
	IsActive  bool                `json:"active"`
	Name      string              `json:"name,omitempty"`
	SubDomain string              `json:"subdomain,omitempty"`
	Config    *IdentityZoneConfig `json:"config,omitempty"`
}

type IdentityZoneConfig struct {
	AccountChooserEnabled bool                            `json:"accountChooserEnabled"`
	Branding              *IdentityZoneBrandingConfig     `json:"branding,omitempty"`
	ClientSecretPolicy    *IdentityZoneClientSecretPolicy `json:"clientSecretPolicy,omitempty"`
	CorsPolicy            *IdentityZoneCorsPolicy         `json:"corsPolicy,omitempty"`
	IdpDiscoveryEnabled   bool                            `json:"idpDiscoveryEnabled"`
	InputPrompts          []*InputPrompt                  `json:"prompts,omitempty"`
	IssuerUrl             string                          `json:"issuer,omitempty"`
	Links                 *IdentityZoneLinks              `json:"links,omitempty"`
	MfaConfig             *MfaConfig                      `json:"MfaConfig,omitempty"`
	TokenPolicy           *IdentityZoneTokenPolicy        `json:"tokenPolicy,omitempty"`
	Saml                  *IdentityZoneSamlConfig         `json:"samlConfig,omitempty"`
	UserConfig            *UserConfig                     `json:"userConfig,omitempty"`
}

type IdentityZoneBrandingConfig struct {
	Banner      *IdentityZoneBrandingBanner `json:"banner,omitempty"`
	CompanyName string                      `json:"companyName,omitempty"`
	CompanyLogo string                      `json:"productLogo,omitempty"`
	Favicon     string                      `json:"squareLogo,omitempty"`
	FooterText  string                      `json:"footerLegalText,omitempty"`
	FooterLinks map[string]string           `json:"footerLinks,omitempty"`
}

type IdentityZoneBrandingBanner struct {
	BackgroundColor string `json:"backgroundColor,omitempty"`
	Logo            string `json:"logo,omitempty"`
	Text            string `json:"text,omitempty"`
	TextColor       string `json:"textColor,omitempty"`
	Url             string `json:"link,omitempty"`
}

type IdentityZoneClientSecretPolicy struct {
	MaxLength             *int64 `json:"maxLength,omitempty"`
	MinLength             *int64 `json:"minLength,omitempty"`
	MinUpperCaseCharacter *int64 `json:"requireUpperCaseCharacter,omitempty"`
	MinLowerCaseCharacter *int64 `json:"requireLowerCaseCharacter,omitempty"`
	MinDigit              *int64 `json:"requireDigit,omitempty"`
	MinSpecialCharacter   *int64 `json:"requireSpecialCharacter,omitempty"`
}

type IdentityZoneCorsPolicy struct {
	DefaultConfiguration *IdentityZoneCorsConfig `json:"defaultConfiguration,omitempty"`
	XhrConfiguration     *IdentityZoneCorsConfig `json:"xhrConfiguration,omitempty"`
}

type IdentityZoneCorsConfig struct {
	AllowedOrigins        []string `json:"allowedOrigins,omitempty"`
	AllowedOriginPatterns []string `json:"allowedOriginPatterns,omitempty"`
	AllowedUris           []string `json:"allowedUris,omitempty"`
	AllowedUriPatterns    []string `json:"allowedUriPatterns,omitempty"`
	AllowedHeaders        []string `json:"allowedHeaders,omitempty"`
	AllowedMethods        []string `json:"allowedMethods,omitempty"`
	AllowedCredentials    bool     `json:"allowedCredentials"`
	MaxAge                *int64   `json:"maxAge,omitempty"`
}

type IdentityZoneTokenPolicy struct {
	AccessTokenTtl       *int64 `json:"accessTokenValidity,omitempty"`
	RefreshTokenTtl      *int64 `json:"refreshTokenValidity,omitempty"`
	IsJwtRevocable       bool   `json:"jwtRevocable"`
	IsRefreshTokenUnique bool   `json:"refreshTokenUnique"`
	RefreshTokenFormat   string `json:"refreshTokenFormat,omitempty"`
	ActiveKeyId          string `json:"activeKeyId,omitempty"`
}

type IdentityZoneSamlConfig struct {
	ActiveKeyId              string                          `json:"activeKeyId,omitempty"`
	AssertionTtlSeconds      *int64                          `json:"assertionTimeToLiveSeconds,omitempty"`
	Certificate              string                          `json:"certificate,omitempty"`
	DisableInResponseToCheck bool                            `json:"disableInResponseToCheck"`
	EntityId                 string                          `json:"entityID,omitempty"`
	IsAssertionSigned        bool                            `json:"assertionSigned"`
	IsRequestSigned          bool                            `json:"requestSigned"`
	Keys                     map[string]*IdentityZoneSamlKey `json:"keys,omitempty"`
	WantAssertionSigned      bool                            `json:"wantAssertionSigned"`
	WantAuthnRequestSigned   bool                            `json:"wantAuthnRequestSigned"`
}

type IdentityZoneSamlKey struct {
	Certificate string `json:"certificate,omitempty"`
}

type IdentityZoneLinks struct {
	HomeRedirect string                   `json:"homeRedirect,omitempty"`
	Logout       *IdentityZoneLogoutLinks `json:"logout,omitempty"`
	SelfService  *SelfServiceLinks        `json:"selfService,omitempty"`
}

type IdentityZoneLogoutLinks struct {
	RedirectUrl           string   `json:"redirectUrl,omitempty"`
	RedirectParameterName string   `json:"redirectParameterName,omitempty"`
	AllowedRedirectUrls   []string `json:"whitelist"`
}

type SelfServiceLinks struct {
	Enabled          bool   `json:"selfServiceLinksEnabled"`
	SignupUrl        string `json:"signup,omitempty"`
	PasswordResetUrl string `json:"passwd,omitempty"`
}

type InputPrompt struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"text,omitempty"`
}

type UserConfig struct {
	DefaultGroups []string `json:"defaultGroups,omitempty"`
}

type MfaConfig struct {
	IsEnabled         bool     `json:"enabled"`
	IdentityProviders []string `json:"identityProviders,omitempty"`
}
