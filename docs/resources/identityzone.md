---
page_title: "Cloud Foundry UAA: uaa_identity_zone"
---

# Group Data Source

Gets information on a Cloud Foundry UAA identity zone.

## Example Usage

The following example looks up an identity zone named 'my-zone'.

```
resource uaa_identity_zone "myzone" {
    name = "my-zone"
    sub_domain = "my-zone"
    
    branding {
        company_name = "My Company, LLC."
        company_logo = "<base64 encoded logo>"
        favicon = "<base64 encoded favicon>"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the identity zone to look up
* `sub_domain` - (Required) Unique subdomain for the running instance. May only contain legal characters for a subdomain name.
* `account_chooser_enabled` - (Optional) This flag enables the account choosing functionality. If set to true in the config the IDP is chosen by discovery. Otherwise, the user can enter the IDP by providing the origin.
* [`branding`](#branding) - (Optional) Branding customization details.  Documented below.
* [`client_secret_policy`](#client_secret_policy) - (Optional) The rules that are enforced when creating/updating client secrets. Documented below.
* [`cors_policy`](#cors_policy) - (Optional) The CORS policy defined for the identity zone. Documented Below.
* `default_user_groups` - (Optional) Default groups each user in the zone inherits.
* `home_redirect_url` - (Optional) Overrides the UAA home page and issues a redirect to this URL when the browser requests `/` and `/home`.
* [`input_prompt`](#input_prompt) - (Optional) List of fields that users are prompted for to login. Defaults to username, password, and passcode. Documented Below.
* `idp_discovery_enabled` - (Optional) IDP Discovery should be set to true if you have configured more than one identity provider for UAA. The discovery relies on email domain being set for each additional provider
* `issuer_url` - (Optional) Issuer of this zone. Must be a valid URL.
* `is_active` - (Optional) Indicates whether the identity zone is active. Defaults to true.
* `logout_redirect_url` - (Optional) Logout redirect url
* `logout_redirect_param` - (Optional) The name of the redirect parameter
* `logout_allowed_redirect_urls` - (Optional) Allowed logout redirect urls
* `mfa_enabled` - `true` (Optional) if Multi-factor Authentication (MFA) is enabled for the identity zone. Defaults to false
* `mfa_identity_providers` - (Optional) Only trigger MFA when user is using an identity provider whose origin key matches one of these values
* `name` - (Optional) Human-readable zone name
* [`saml_config`](#saml_config) - (Optional) SAML configuration for the identity zone. Documented Below.
* `self_serve_enabled` - (Optional) Whether users are allowed to sign up or reset their passwords via the UI
* `self_serve_signup_url` - (Optional) Where users are directed upon clicking the account creation link
* `self_serve_pw_reset_url` - (Optional) Where users are directed upon clicking the password reset link
* [`token_policy`](#token_policy) - (Optional) Various fields pertaining to the JWT access and refresh tokens.  Documented below.

### client_secret_policy

* `max_length` - Maximum number of characters required for secret to be considered valid (defaults to 255).
* `min_digits` - Minimum number of digits required for secret to be considered valid (defaults to 0).
* `min_length` - Minimum number of characters required for secret to be considered valid (defaults to 0).
* `min_lower_case_chars` - Minimum number of lowercase characters required for secret to be considered valid (defaults to 0).
* `min_special_chars` - Minimum number of special characters required for secret to be considered valid (defaults to 0).
* `min_upper_case_chars` - Minimum number of uppercase characters required for secret to be considered valid (defaults to 0).

### cors_policy

* `name` - The effected CORS resource (allowed values: default, xhr)
* `allowed_origins` - `Access-Control-Allow-Origin` header. Indicates whether a resource can be shared based by returning the value of the Origin request header, "*", or "null" in the response.
* `allowed_origin_patterns` - Indicates whether a resource can be shared based by returning the value of the Origin patterns
* `allowed_uris` - The list of allowed URIs.
* `allowed_uri_patterns` - The list of allowed URI patterns.
* `allowed_headers` - `Access-Control-Allow-Headers` header. Indicates which header field names can be used during the actual response
* `allowed_methods` - `Access-Control-Allow-Methods` header. Indicates which method will be used in the actual request as part of the preflight request.
* `allowed_credentials` - `Access-Control-Allow-Credentials` header. Indicates whether the response to request can be exposed when the omit credentials flag is unset. When part of the response to a preflight request it indicates that the actual request can include user credentials..
* `max_age` - `Access-Control-Max-Age` header. Indicates how long the results of a preflight request can be cached in a preflight result cache

### input_prompt

* `name` - Name of field
* `type` - What kind of field this is (e.g. text or password)
* `value` - Actual text displayed on prompt for field

### saml_config

* `active_key_id` - The ID of the key that should be used for signing metadata and assertions.
* `assertion_ttl_seconds` - The ID of the key that should be used for signing metadata and assertions.
* `disable_in_response_to_check` - If `true`, this zone will not validate the InResponseToField part of an incoming IDP assertion. Please see` https://docs.spring.io/spring-security-saml/docs/current/reference/html/chapter-troubleshooting.html`
* `entity_id` - Unique ID of the SAML2 entity
* `is_assertion_signed` - If `true`, the SAML provider will sign all assertions
* `is_request_signed` - Exposed SAML metadata property. If `true`, the service provider will sign all outgoing authentication requests. Defaults to `true`.
* [`key`](#saml_configkey) - A list of the SAML provider's private keys. Documented below.
* `want_assertion_signed` - 	Exposed SAML metadata property. If `true`, all assertions received by the SAML provider must be signed. Defaults to `true`.
* `want_authn_request_signed` - If `true`, the authentication request from the partner service provider must be signed.

### saml_config.key
* `name` - The name of the SAML key
* `certificate` - The 

### token_policy
* `access_token_ttl` - Time in seconds between when a access token is issued and when it expires. Defaults to global `accessTokenValidity`
* `active_key_id` - The name of the key that is being used to sign tokens
* `is_jwt_revocable` - Set to true if JWT tokens should be stored in the token store, and thus made individually revocable. Opaque tokens are always stored and revocable.
* `is_refresh_token_unique` - If true, uaa will only issue one refresh token per client_id/user_id combination. Defaults to `false`.
* `refresh_token_format` - The format for the refresh token. Allowed values are `jwt`, `opaque`. Defaults to `jwt`.
* `refresh_token_ttl` - Time in seconds between when a refresh token is issued and when it expires. Defaults to global `refreshTokenValidity`

### branding
* `banner_bg_color` - Hexadecimal color code for banner background color, does not allow color namesThis name is used on the UAA Pages and in account management related communication in UAA
* `banner_logo` - This is base64 encoded PNG data displayed in a banner at the top of the UAA login page, overrides banner text
* `banner_text` - This is text displayed in a banner at the top of the UAA login page
* `banner_text_color` - Hexadecimal color code for banner text color, does not allow color names
* `banner_url` - The UAA login banner will be a link pointing to this url
* `company_name` - This name is used on the UAA Pages and in account management related communication in UAA
* `company_logo` - This is a base64Url encoded PNG image which will be used as the logo on all UAA pages like Login, Sign Up etc.
* `favicon` - This is a base64 encoded PNG image which will be used as the favicon for the UAA pages
* `footer_text` - 	This text appears on the footer of all UAA pages
* [`footer_links`](#footer_links) - These links appear on the footer of all UAA pages. You may choose to add multiple urls for things like Support, Terms of Service etc. Documented below.

### footer_links
* `name` - The link text to be displayed.
* `url` - The url for the href of the link displayed.
