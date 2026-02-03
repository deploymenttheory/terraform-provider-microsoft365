---
page_title: "microsoft365_graph_beta_applications_service_principal Data Source - terraform-provider-microsoft365"
subcategory: "Applications"

description: |-
  Retrieves information about a Microsoft Entra ID service principal using the /servicePrincipals endpoint. This data source is used to query enterprise applications and managed identities by ID, app ID, display name, or advanced OData filtering.
---

# microsoft365_graph_beta_applications_service_principal (Data Source)

Retrieves information about a Microsoft Entra ID service principal using the `/servicePrincipals` endpoint. This data source is used to query enterprise applications and managed identities by ID, app ID, display name, or advanced OData filtering.

## Microsoft Documentation

- [servicePrincipal resource type](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `Application.Read.All`
- `Application.ReadWrite.All`
- `Directory.Read.All`
- `Directory.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.31.0-alpha | Experimental | Initial release |
| v0.43.0-alpha | Preview | Changed Odata filter logic and added more examples |

## Important Notes

### OData Query Support

This method supports the `$count`, `$expand`, `$filter`, `$orderby`, `$search`, `$select`, and `$top` OData query parameters to help customize the response. Some queries are supported only when you use the ConsistencyLevel header set to eventual and `$count`. For more information, see [Advanced query capabilities on directory objects](https://learn.microsoft.com/en-us/graph/aad-advanced-queries).

**Important:** When using the `odata_query` attribute, provide only the filter expression itself without the `$filter=` prefix. For example:
- **Correct:** `odata_query = "displayName eq 'Microsoft Intune'"`
- **Incorrect:** `odata_query = "$filter=displayName eq 'Microsoft Intune'"`

### Key Credentials

By default, this API doesn't return the public key value of the key in the `keyCredentials` property unless `keyCredentials` is specified in a `$select` query. For example, `$select=id,appId,keyCredentials`.

The use of `$select` to get `keyCredentials` for service principals has a throttling limit of 150 requests per minute for every tenant.

### Retry and Consistency

This datasource implements automatic retry logic with up to 6 attempts over 60 seconds to handle eventual consistency in Microsoft Entra ID. This is particularly important when querying for service principals that were recently created.

## Lookup Methods

This data source provides four mutually exclusive methods to query for a service principal:

1. **object_id**: Look up a service principal by its unique object ID (GUID). This is the most direct and efficient lookup method.
2. **app_id**: Look up a service principal by its application (client) ID (GUID). The system uses an OData filter internally.
3. **display_name**: Look up a service principal by its exact display name (case-sensitive). The system uses an OData filter internally.
4. **odata_query**: Use a custom OData `$filter` query for advanced filtering. Must return exactly one service principal.

You must specify exactly one of these attributes. Using multiple lookup methods simultaneously will result in a validation error.

### OData Query Examples

When using `odata_query`, you can use any valid OData filter expression without the `$filter=` prefix. The provider automatically applies the necessary OData parameters and headers (including ConsistencyLevel: eventual) when executing the query.

**Note:** Do not include `$filter=` in the `odata_query` value. Provide only the filter expression itself.

## Example Usage

### Lookup by Object ID

```terraform
# Retrieve a service principal by its Object ID
data "microsoft365_graph_beta_applications_service_principal" "by_id" {
  object_id = "3b6f95b0-2064-4cc9-b5e5-1ab72af707b3"
}
```

### Lookup by App ID

```terraform
# Retrieve a service principal by its Application (Client) ID
data "microsoft365_graph_beta_applications_service_principal" "by_app_id" {
  app_id = "63e61dc2-f593-4a6f-92b9-92e4d2c03d4f"
}
```

### Lookup by Display Name

```terraform
# Retrieve a service principal by its display name
data "microsoft365_graph_beta_applications_service_principal" "by_display_name" {
  display_name = "Microsoft Intune SCCM Connector"
}
```

### OData Query - Filter with Conditions

```terraform
# Retrieve a service principal using OData query with filter conditions
data "microsoft365_graph_beta_applications_service_principal" "odata_filter" {
  odata_query = "preferredSingleSignOnMode ne 'notSupported' and displayName eq 'Microsoft Intune'"
}
```

### OData Query - Advanced Multiple Conditions

```terraform
# Retrieve a service principal using OData query with multiple filter conditions
data "microsoft365_graph_beta_applications_service_principal" "odata_advanced" {
  odata_query = "servicePrincipalType eq 'Application' and accountEnabled eq true"
}
```

### OData Query - Comprehensive SAML Filter

```terraform
# Retrieve a service principal using comprehensive OData query filters
# This example demonstrates filtering for SAML-based service principals
data "microsoft365_graph_beta_applications_service_principal" "odata_comprehensive" {
  odata_query = "preferredSingleSignOnMode eq 'saml' and servicePrincipalType eq 'Application'"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `app_id` (String) The unique identifier for the associated application (client ID). Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`). One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.
- `display_name` (String) The display name for the service principal. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on null values), `$search`, and `$orderby`. One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.
- `object_id` (String) The object ID of the service principal. One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.
- `odata_query` (String) Custom OData filter query. Use this for advanced filtering when the standard lookup attributes don't meet your needs. Cannot be combined with `object_id`, `app_id`, or `display_name`. Example: `displayName eq 'My Service Principal' and servicePrincipalType eq 'Application'`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `account_enabled` (Boolean) true if the service principal account is enabled; otherwise, false. If set to false, then no users are able to sign in to this app, even if they're assigned to it. Supports `$filter` (`eq`, `ne`, `not`, `in`).
- `app_display_name` (String) The display name exposed by the associated application.
- `app_owner_organization_id` (String) Contains the tenant ID where the application is registered. This is applicable only to service principals backed by applications. Supports `$filter` (`eq`, `ne`, `NOT`, `ge`, `le`).
- `app_role_assignment_required` (Boolean) Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports `$filter` (`eq`, `ne`, `NOT`).
- `application_template_id` (String) Unique identifier of the applicationTemplate that the servicePrincipal was created from. Read-only. Supports `$filter` (`eq`, `ne`, `NOT`, `startsWith`).
- `deleted_date_time` (String) The date and time the service principal was deleted. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
- `disabled_by_microsoft_status` (String) Specifies whether Microsoft has disabled the registered application. Possible values are: `null` (default value), `NotDisabled`, and `DisabledDueToViolationOfServicesAgreement` (reasons include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement). Supports `$filter` (`eq`, `ne`, `NOT`).
- `error_url` (String) Deprecated. Do not use.
- `homepage` (String) Home page or landing page of the application.
- `id` (String) Unique identifier for the service principal object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).
- `info` (Attributes) Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more information, see How to: Add Terms of service and privacy statement for registered Azure AD apps. (see [below for nested schema](#nestedatt--info))
- `login_url` (String) Specifies the URL where the service provider redirects the user to Azure AD to authenticate. Azure AD uses the URL to launch the application from Microsoft 365 or the Azure AD My Apps. When blank, Azure AD performs IdP-initiated sign-on for applications configured with SAML-based single sign-on.
- `logout_url` (String) Specifies the URL that the Microsoft's authorization service uses to sign out a user using OpenId Connect front-channel, back-channel, or SAML sign out protocols.
- `notes` (String) Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1,024 characters.
- `notification_email_addresses` (Set of String) Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery applications.
- `preferred_single_sign_on_mode` (String) Specifies the single sign-on mode configured for this application. Azure AD uses the preferred single sign-on mode to launch the application from Microsoft 365 or the Azure AD My Apps. The supported values are `password`, `saml`, `notSupported`, and `oidc`.
- `preferred_token_signing_key_end_date_time` (String) Specifies the expiration date of the keyCredential used for token signing, marked by preferredTokenSigningKeyThumbprint. Updating this attribute isn't currently supported. For details, see ServicePrincipal property differences. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.
- `preferred_token_signing_key_thumbprint` (String) This property can be used on SAML applications (apps that have preferredSingleSignOnMode set to saml) to control which certificate is used to sign the SAML responses. For applications that aren't SAML, don't write or otherwise rely on this property.
- `publisher_name` (String) The display name of the tenant in which the associated application is registered. Provided only when the application publisher is from a different tenant. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).
- `reply_urls` (Set of String) The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application. Not nullable.
- `saml_metadata_url` (String) The URL where the service exposes SAML metadata for federation.
- `saml_single_sign_on_settings` (Attributes) The collection for settings related to saml single sign-on. (see [below for nested schema](#nestedatt--saml_single_sign_on_settings))
- `service_principal_names` (Set of String) Contains the list of identifiersUris, copied over from the associated application. More values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Azure AD. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).
- `service_principal_type` (String) Identifies if the service principal represents an Application, a ManagedIdentity, or a legacy application (socialIdp). This is set by Azure AD internally. For a service principal that represents an Application this is set as Application. For a service principal that represent a managed identity this is set as ManagedIdentity. For a service principal representing a legacy app this is set as SocialIdp. Supports `$filter` (`eq`, `ne`, `NOT`, `in`).
- `sign_in_audience` (String) Specifies the Microsoft accounts that are supported for the current application. Supported values are `AzureADMyOrg`, `AzureADMultipleOrgs`, `AzureADandPersonalMicrosoftAccount`, `PersonalMicrosoftAccount`. Read-only. Supports `$filter` (`eq`, `ne`, `NOT`, `startsWith`).
- `tags` (Set of String) Custom strings that can be used to categorize and identify the service principal. Not nullable. Supports `$filter` (`eq`, `ne`, `NOT`, `ge`, `le`, `startsWith`).
- `verified_publisher` (Attributes) Specifies the verified publisher of the application which this service principal represents. (see [below for nested schema](#nestedatt--verified_publisher))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--info"></a>
### Nested Schema for `info`

Read-Only:

- `logo_url` (String) CDN URL to the application's logo.
- `marketing_url` (String) Link to the application's marketing page.
- `privacy_statement_url` (String) Link to the application's privacy statement.
- `support_url` (String) Link to the application's support page.
- `terms_of_service_url` (String) Link to the application's terms of service statement.


<a id="nestedatt--saml_single_sign_on_settings"></a>
### Nested Schema for `saml_single_sign_on_settings`

Read-Only:

- `relay_state` (String) The relative URI the service provider would redirect to after completion of the single sign-on flow.


<a id="nestedatt--verified_publisher"></a>
### Nested Schema for `verified_publisher`

Read-Only:

- `added_date_time` (String) The timestamp when the verified publisher was first added or most recently updated.
- `display_name` (String) The verified publisher name from the app publisher's Microsoft Partner Network (MPN) account.
- `verified_publisher_id` (String) The ID of the verified publisher from the app publisher's Partner Center account.