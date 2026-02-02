---
page_title: "microsoft365_graph_beta_applications_application Resource - terraform-provider-microsoft365"
subcategory: "Applications"
description: |-
  Manages an application in Microsoft Entra ID using the /applications endpoint. Any application that outsources authentication to Microsoft Entra ID must be registered in the Microsoft identity platform. Application registration involves telling Microsoft Entra ID about your application, including the URL where it's located, the URL to send replies after authentication, the URI to identify your application, and more.
---

# microsoft365_graph_beta_applications_application (Resource)

Manages an application in Microsoft Entra ID using the `/applications` endpoint. Any application that outsources authentication to Microsoft Entra ID must be registered in the Microsoft identity platform. Application registration involves telling Microsoft Entra ID about your application, including the URL where it's located, the URL to send replies after authentication, the URI to identify your application, and more.

## Microsoft Documentation

- [application resource type](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta)
- [Create application](https://learn.microsoft.com/en-us/graph/api/application-post-applications?view=graph-rest-beta&tabs=http)
- [Get application](https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta&tabs=http)
- [Update application](https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta&tabs=http)
- [Delete application](https://learn.microsoft.com/en-us/graph/api/application-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Application.Read.All`
- `Directory.Read.All`
- `Application.ReadWrite.All`
- `Directory.ReadWrite.All`

**Optional:**
- `Application.ReadWrite.OwnedBy` (if managing only applications owned by the service principal)

Find out more about the permissions required for managing applications at Microsoft Learn [here](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.43.0 | Experimental | Initial release |

## Important Notes

- **Application Types**: Choose the appropriate `sign_in_audience` based on your needs:
  - `AzureADMyOrg`: Single tenant (organization users only)
  - `AzureADMultipleOrgs`: Multitenant (any Azure AD organization)
  - `AzureADandPersonalMicrosoftAccount`: Multitenant + Personal Microsoft accounts
  - `PersonalMicrosoftAccount`: Personal Microsoft accounts only
- **Identifier URIs**: Must use verified domains for your tenant or the `api://` format
- **Owners**: Applications require at least one owner for proper management
- **Hard Delete**: The `hard_delete` attribute controls whether the application is permanently deleted or soft-deleted (recoverable for 30 days)

## Example Usage

### Minimal Application

```terraform
resource "microsoft365_graph_beta_applications_application" "minimal" {
  display_name = "my-minimal-app"
  description  = "A minimal application with only required fields"
}
```

### Web Application

```terraform
resource "microsoft365_graph_beta_applications_application" "web_app" {
  display_name     = "my-web-application"
  description      = "Web application with OIDC authentication"
  sign_in_audience = "AzureADMyOrg"

  identifier_uris = [
    "https://mycompany.com/my-web-app"
  ]

  web = {
    home_page_url = "https://my-web-app.azurewebsites.net"
    logout_url    = "https://my-web-app.azurewebsites.net/signout"
    redirect_uris = [
      "https://my-web-app.azurewebsites.net/signin-oidc"
    ]
    implicit_grant_settings = {
      enable_access_token_issuance = false
      enable_id_token_issuance     = true
    }
  }

  required_resource_access = []
}
```

### Single Page Application (SPA)

```terraform
resource "microsoft365_graph_beta_applications_application" "spa" {
  display_name     = "my-single-page-app"
  description      = "Single Page Application (React, Angular, Vue)"
  sign_in_audience = "AzureADMultipleOrgs"

  identifier_uris = [
    "https://mycompany.com/my-spa"
  ]

  spa = {
    redirect_uris = [
      "http://localhost:3000",
      "https://my-spa.azurestaticapps.net"
    ]
  }

  required_resource_access = []
}
```

### Public Client (Mobile/Desktop)

```terraform
resource "microsoft365_graph_beta_applications_application" "mobile_app" {
  display_name              = "my-mobile-application"
  description               = "Mobile or desktop application (public client)"
  sign_in_audience          = "AzureADMyOrg"
  is_fallback_public_client = true

  public_client = {
    redirect_uris = [
      "http://localhost",
      "ms-appx-web://microsoft.aad.brokerplugin/my-mobile-app"
    ]
  }

  required_resource_access = []
}
```

### Multitenant Application

```terraform
resource "microsoft365_graph_beta_applications_application" "multitenant" {
  display_name     = "my-multitenant-app"
  description      = "Multitenant application with personal Microsoft account support"
  sign_in_audience = "AzureADandPersonalMicrosoftAccount"

  identifier_uris = [
    "https://mycompany.com/my-multitenant-app"
  ]

  web = {
    home_page_url = "https://contoso.com"
    redirect_uris = [
      "https://contoso.com/signin-oidc"
    ]
    implicit_grant_settings = {
      enable_access_token_issuance = false
      enable_id_token_issuance     = true
    }
  }

  spa = {
    redirect_uris = [
      "https://contoso.com/spa"
    ]
  }

  required_resource_access = []

  tags = [
    "multitenant",
    "production"
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name for the application. Maximum length is 256 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on null values), `$search`, and `$orderby`.

### Optional

- `api` (Attributes) Specifies settings for an application that implements a web API. (see [below for nested schema](#nestedatt--api))
- `app_roles` (Attributes Set) The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable. (see [below for nested schema](#nestedatt--app_roles))
- `description` (String) Free text field to provide a description of the application object to end users. The maximum allowed size is 1,024 characters. Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `startsWith`) and `$search`.
- `group_membership_claims` (Set of String) Configures the groups claim issued in a user or OAuth 2.0 access token that the application expects. To set this attribute, use one of the following string values: `None`, `SecurityGroup` (for security groups and Microsoft Entra roles), `All` (this gets all security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of).
- `hard_delete` (Boolean) When `true`, the application will be permanently deleted (hard delete) during destroy. When `false` (default), the application will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. Note: This field defaults to `false` on import since the API does not return this value.
- `identifier_uris` (Set of String) Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique across Microsoft Entra ID. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).
- `info` (Attributes) Basic profile information of the application, such as it's marketing, support, terms of service, and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more information, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, and `eq` on null values). (see [below for nested schema](#nestedatt--info))
- `is_device_only_auth_supported` (Boolean) Specifies whether this application supports device authentication without a user. The default is false.
- `is_fallback_public_client` (Boolean) Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false, which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID can't determine the client application type. For example, the ROPC flow where the application is configured without specifying a redirect URI. In those cases Microsoft Entra ID interprets the application type based on the value of this property.
- `notes` (String) Notes relevant for the management of the application.
- `optional_claims` (Attributes) Application developers can configure optional claims in their Microsoft Entra applications to specify the claims that are sent to their application by the Microsoft security token service. For more information, see How to: Provide optional claims to your app. (see [below for nested schema](#nestedatt--optional_claims))
- `owner_user_ids` (Set of String) The user IDs of the owners for the application. At least one owner is typically required when creating an application. Owners are a set of non-admin users or service principals allowed to modify this object.
- `parental_control_settings` (Attributes) Specifies parental control settings for an application. (see [below for nested schema](#nestedatt--parental_control_settings))
- `prevent_duplicate_names` (Boolean) If set to `true`, the provider will check for existing applications with the same display name and return an error if one is found. This helps prevent accidentally creating duplicate applications. Note: This field defaults to `false` on import since the API does not return this value.
- `public_client` (Attributes) Specifies settings for installed clients such as desktop or mobile devices. (see [below for nested schema](#nestedatt--public_client))
- `required_resource_access` (Attributes Set) Specifies the resources that the application needs to access. This property also specifies the set of delegated permissions and application roles that it needs for each of those resources. This configuration of access to the required resources drives the consent experience. No more than 50 resource services (APIs) can be configured. Beginning mid-October 2021, the total number of required permissions must not exceed 400. For more information, see Limits on requested permissions per app. Not nullable. Supports `$filter` (`eq`, `not`, `ge`, `le`). (see [below for nested schema](#nestedatt--required_resource_access))
- `service_management_reference` (String) References application or service contact information from a Service or Asset Management database. Nullable.
- `sign_in_audience` (String) Specifies the Microsoft accounts that are supported for the current application. The possible values are: `AzureADMyOrg` (default), `AzureADMultipleOrgs`, `AzureADandPersonalMicrosoftAccount`, and `PersonalMicrosoftAccount`. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you may need to change other properties first. For more information, see Validation differences for signInAudience. Supports `$filter` (`eq`, `ne`, `not`).
- `sign_in_audience_restrictions` (Attributes) Specifies restrictions on the supported account types specified in signInAudience. The value type determines the restrictions that can be applied: unrestrictedAudience (There are no additional restrictions on the supported account types allowed by signInAudience) or allowedTenantsAudience (The application can only be used in the specified Entra tenants. Only supported when signInAudience is AzureADMultipleOrgs). Default is a value of type unrestrictedAudience. Returned only on `$select`. (see [below for nested schema](#nestedatt--sign_in_audience_restrictions))
- `spa` (Attributes) Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens. (see [below for nested schema](#nestedatt--spa))
- `tags` (Set of String) Custom strings that can be used to categorize and identify the application. Not nullable. Strings added here also appear in the tags property of any associated service principals. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`) and `$search`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `web` (Attributes) Specifies settings for a web application. (see [below for nested schema](#nestedatt--web))

### Read-Only

- `app_id` (String) The unique identifier for the application that is assigned by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`).
- `created_date_time` (String) The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, and `eq` on null values) and `$orderby`.
- `deleted_date_time` (String) The date and time the application was deleted. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
- `disabled_by_microsoft_status` (String) Specifies whether Microsoft has disabled the registered application. The possible values are: null (default value), `NotDisabled`, and `DisabledDueToViolationOfServicesAgreement` (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement). Supports `$filter` (`eq`, `ne`, `not`). Read-only.
- `id` (String) Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).
- `publisher_domain` (String) The verified publisher domain for the application. Read-only. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).

<a id="nestedatt--api"></a>
### Nested Schema for `api`

Optional:

- `accept_mapped_claims` (Boolean) Allows an application to use claims mapping without specifying a custom signing key.
- `known_client_applications` (Set of String) Used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app. If you set the appID of the client app to this value, the user only consents once to the client app. Microsoft Entra ID knows that consenting to the client means implicitly consenting to the web API and automatically provisions service principals for both APIs at the same time. Both the client and the web API app must be registered in the same tenant.
- `oauth2_permission_scopes` (Attributes Set) The definition of the delegated permissions exposed by the web API represented by this application registration. These delegated permissions may be requested by a client application, and may be granted by users or administrators during consent. Delegated permissions are sometimes referred to as OAuth 2.0 scopes. (see [below for nested schema](#nestedatt--api--oauth2_permission_scopes))
- `pre_authorized_applications` (Attributes Set) Lists the client applications that are preauthorized with the specified delegated permissions to access this application's APIs. Users aren't required to consent to any preauthorized application (for the permissions specified). However, any other permissions not listed in preAuthorizedApplications (requested through incremental consent for example) will require user consent. (see [below for nested schema](#nestedatt--api--pre_authorized_applications))
- `requested_access_token_version` (Number) Specifies the access token version expected by this resource. This changes the version and format of the JWT produced independent of the endpoint or client used to request the access token. The endpoint used, v1.0 or v2.0, is chosen by the client and only impacts the version of id_tokens. Resources need to explicitly configure requestedAccessTokenVersion to indicate the supported access token format. Possible values for requestedAccessTokenVersion are 1, 2, or null. If the value is null, this defaults to 1, which corresponds to the v1.0 endpoint. If signInAudience on the application is configured as AzureADandPersonalMicrosoftAccount or PersonalMicrosoftAccount, the value for this property must be 2.

<a id="nestedatt--api--oauth2_permission_scopes"></a>
### Nested Schema for `api.oauth2_permission_scopes`

Required:

- `id` (String) Unique scope permission identifier inside the oauth2PermissionScopes collection. Required.

Optional:

- `admin_consent_description` (String) A description of the delegated permissions, intended to be read by an administrator granting the permission on behalf of all users. This text appears in tenant-wide admin consent experiences.
- `admin_consent_display_name` (String) The permission's title, intended to be read by an administrator granting the permission on behalf of all users.
- `is_enabled` (Boolean) When you create or update a permission, this property must be set to true (which is the default). To delete a permission, this property must first be set to false. At that point, in a subsequent call, the permission may be removed.
- `type` (String) The possible values are: `User` and `Admin`. Specifies whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator consent should be required for the permissions. While Microsoft Graph defines the default consent requirement for each permission, the tenant administrator may override the behavior in their organization (by allowing, restricting, or limiting user consent to this delegated permission). For more information, see Configure how users consent to applications.
- `user_consent_description` (String) A description of the delegated permissions, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.
- `user_consent_display_name` (String) A title for the permission, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.
- `value` (String) Specifies the value to include in the scp (scope) claim in access tokens. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with `.`.


<a id="nestedatt--api--pre_authorized_applications"></a>
### Nested Schema for `api.pre_authorized_applications`

Required:

- `app_id` (String) The unique identifier for the client application.
- `delegated_permission_ids` (Set of String) The unique identifier for the scopes the client application is granted.



<a id="nestedatt--app_roles"></a>
### Nested Schema for `app_roles`

Required:

- `allowed_member_types` (Set of String) Specifies whether this app role can be assigned to users and groups (by setting to `['User']`), to other application's (by setting to `['Application']`, or both (by setting to `['User', 'Application']`). App roles supporting assignment to other applications' service principals are also known as application permissions. The 'Application' value is only supported for app roles defined on application entities. Required.
- `description` (String) The description for the app role. This is displayed when the app role is being assigned and, if the app role functions as an application permission, during consent experiences. Required.
- `display_name` (String) Display name for the permission that appears in the app role assignment and consent experiences. Required.
- `id` (String) Unique role identifier inside the appRoles collection. When creating a new app role, a new GUID identifier must be provided. Required.
- `is_enabled` (Boolean) Defines whether the application's app role is enabled or disabled. Required.

Optional:

- `value` (String) Specifies the value to include in the roles claim in ID tokens and access tokens authenticating an assigned user or service principal. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with `.`. Nullable.

Read-Only:

- `origin` (String) Specifies if the app role is defined on the application object or on the servicePrincipal entity. Must not be included in any POST or PATCH requests. Read-only.


<a id="nestedatt--info"></a>
### Nested Schema for `info`

Optional:

- `marketing_url` (String) Link to the application's marketing page. For example, https://www.contoso.com/app/marketing.
- `privacy_statement_url` (String) Link to the application's privacy statement. For example, https://www.contoso.com/app/privacy.
- `support_url` (String) Link to the application's support page. For example, https://www.contoso.com/app/support.
- `terms_of_service_url` (String) Link to the application's terms of service statement. For example, https://www.contoso.com/app/termsofservice.

Read-Only:

- `logo_url` (String) CDN URL to the application's logo. Read-only.


<a id="nestedatt--optional_claims"></a>
### Nested Schema for `optional_claims`

Optional:

- `access_token` (Attributes Set) The optional claims returned in the JWT access token. (see [below for nested schema](#nestedatt--optional_claims--access_token))
- `id_token` (Attributes Set) The optional claims returned in the JWT ID token. (see [below for nested schema](#nestedatt--optional_claims--id_token))
- `saml2_token` (Attributes Set) The optional claims returned in the SAML token. (see [below for nested schema](#nestedatt--optional_claims--saml2_token))

<a id="nestedatt--optional_claims--access_token"></a>
### Nested Schema for `optional_claims.access_token`

Required:

- `name` (String) The name of the optional claim. Required.

Optional:

- `additional_properties` (Set of String) Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
- `essential` (Boolean) If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
- `source` (String) The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.


<a id="nestedatt--optional_claims--id_token"></a>
### Nested Schema for `optional_claims.id_token`

Required:

- `name` (String) The name of the optional claim. Required.

Optional:

- `additional_properties` (Set of String) Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
- `essential` (Boolean) If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
- `source` (String) The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.


<a id="nestedatt--optional_claims--saml2_token"></a>
### Nested Schema for `optional_claims.saml2_token`

Required:

- `name` (String) The name of the optional claim. Required.

Optional:

- `additional_properties` (Set of String) Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
- `essential` (Boolean) If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
- `source` (String) The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.



<a id="nestedatt--parental_control_settings"></a>
### Nested Schema for `parental_control_settings`

Optional:

- `countries_blocked_for_minors` (Set of String) Specifies the two-letter ISO country codes. Access to the application will be blocked for minors from the countries specified in this list.
- `legal_age_group_rule` (String) Specifies the legal age group rule that applies to users of the app. Can be set to one of the following values: `Allow`, `RequireConsentForPrivacyServices`, `RequireConsentForMinors`, `RequireConsentForKids`, `BlockMinors`.


<a id="nestedatt--public_client"></a>
### Nested Schema for `public_client`

Optional:

- `redirect_uris` (Set of String) Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.


<a id="nestedatt--required_resource_access"></a>
### Nested Schema for `required_resource_access`

Required:

- `resource_access` (Attributes Set) The list of OAuth2.0 permission scopes and app roles that the application requires from the specified resource. Required. (see [below for nested schema](#nestedatt--required_resource_access--resource_access))
- `resource_app_id` (String) The unique identifier for the resource that the application requires access to. This should be equal to the appId declared on the target resource application. Required.

<a id="nestedatt--required_resource_access--resource_access"></a>
### Nested Schema for `required_resource_access.resource_access`

Required:

- `id` (String) The unique identifier of an app role or delegated permission exposed by the resource application. For delegated permissions, this should match the id property of one of the delegated permissions in the oauth2PermissionScopes collection of the resource application's service principal. For app roles (application permissions), this should match the id property of an app role in the appRoles collection of the resource application's service principal. Required.
- `type` (String) Specifies whether the id property references a delegated permission or an app role (application permission). The possible values are: `Scope` (for delegated permissions) or `Role` (for app roles). Required.



<a id="nestedatt--sign_in_audience_restrictions"></a>
### Nested Schema for `sign_in_audience_restrictions`

Optional:

- `allowed_tenant_ids` (Set of String) The list of allowed tenant IDs. Only applicable when odata_type is `#microsoft.graph.allowedTenantsAudience`.
- `is_home_tenant_allowed` (Boolean) Indicates whether the home tenant is allowed. Only applicable when odata_type is `#microsoft.graph.allowedTenantsAudience`.
- `odata_type` (String) The OData type. Must be `#microsoft.graph.allowedTenantsAudience` or `#microsoft.graph.unrestrictedAudience`.


<a id="nestedatt--spa"></a>
### Nested Schema for `spa`

Optional:

- `redirect_uris` (Set of String) Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--web"></a>
### Nested Schema for `web`

Optional:

- `home_page_url` (String) Home page or landing page of the application.
- `implicit_grant_settings` (Attributes) Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow. (see [below for nested schema](#nestedatt--web--implicit_grant_settings))
- `logout_url` (String) Specifies the URL that is used by Microsoft's authorization service to log out a user using front-channel, back-channel or SAML logout protocols.
- `redirect_uri_settings` (Attributes Set) Specifies the index of the URLs where user tokens are sent for sign-in. This is only valid for applications using SAML. Note: If not specified, the API may auto-generate settings based on redirect_uris. To manage this field, you must provide at least one entry; empty arrays are not supported as the API auto-generates values. (see [below for nested schema](#nestedatt--web--redirect_uri_settings))
- `redirect_uris` (Set of String) Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.

<a id="nestedatt--web--implicit_grant_settings"></a>
### Nested Schema for `web.implicit_grant_settings`

Optional:

- `enable_access_token_issuance` (Boolean) Specifies whether this web application can request an access token using the OAuth 2.0 implicit flow.
- `enable_id_token_issuance` (Boolean) Specifies whether this web application can request an ID token using the OAuth 2.0 implicit flow.


<a id="nestedatt--web--redirect_uri_settings"></a>
### Nested Schema for `web.redirect_uri_settings`

Optional:

- `index` (Number) Identifies the specific URI within the redirectURIs collection in SAML SSO flows. Defaults to null. The index is unique across all the redirectUris for the application.
- `uri` (String) Specifies the URI that tokens are sent to.

## Import

```shell
# Simple import - defaults to prevent_duplicate_names=false and hard_delete=false
terraform import microsoft365_graph_beta_applications_application.example "00000000-0000-0000-0000-000000000000"

# Extended import - with hard_delete enabled for permanent deletion
terraform import microsoft365_graph_beta_applications_application.example "00000000-0000-0000-0000-000000000000:hard_delete=true"

# Extended import - with both prevent_duplicate_names and hard_delete enabled
terraform import microsoft365_graph_beta_applications_application.example "00000000-0000-0000-0000-000000000000:prevent_duplicate_names=true:hard_delete=true"
```
