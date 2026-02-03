---
page_title: "microsoft365_graph_beta_applications_application Data Source - terraform-provider-microsoft365"
subcategory: "Applications"

description: |-
  Retrieves information about a Microsoft Entra ID (Azure AD) application using the /applications endpoint. This data source is used to query application details by ID, app ID, display name, or advanced OData filtering.
---

# microsoft365_graph_beta_applications_application (Data Source)

Retrieves information about a Microsoft Entra ID (Azure AD) application using the `/applications` endpoint. This data source is used to query application details by ID, app ID, display name, or advanced OData filtering.

## Microsoft Documentation

- [application resource type](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta)
- [Get application](https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta)
- [List applications](https://learn.microsoft.com/en-us/graph/api/application-list?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `Application.Read.All`
- `Directory.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.43.0-alpha | Experimental | Initial release |

## Important Notes

### OData Query Support

This method supports the `$count`, `$expand`, `$filter`, `$orderby`, `$search`, `$select`, and `$top` OData query parameters to help customize the response. Some queries are supported only when you use the ConsistencyLevel header set to eventual and `$count`. For more information, see [Advanced query capabilities on directory objects](https://learn.microsoft.com/en-us/graph/aad-advanced-queries).

**Important:** When using the `odata_query` attribute, provide only the filter expression itself without the `$filter=` prefix. For example:
- **Correct:** `odata_query = "displayName eq 'My Application'"`
- **Incorrect:** `odata_query = "$filter=displayName eq 'My Application'"`

### Key Credentials

By default, this API doesn't return the public key value of the key in the `keyCredentials` property unless `keyCredentials` is specified in a `$select` query. For example, `$select=id,appId,keyCredentials`.

The use of `$select` to get `keyCredentials` for applications has a throttling limit of 150 requests per minute for every tenant.

### Retry and Consistency

This datasource implements automatic retry logic with up to 6 attempts over 60 seconds to handle eventual consistency in Microsoft Entra ID. This is particularly important when querying for applications that were recently created.

## Lookup Methods

This data source provides four mutually exclusive methods to query for an application:

1. **object_id**: Look up an application by its unique object ID (GUID). This is the most direct and efficient lookup method.
2. **app_id**: Look up an application by its application (client) ID (GUID). The system uses an OData filter internally.
3. **display_name**: Look up an application by its exact display name (case-sensitive). The system uses an OData filter internally.
4. **odata_query**: Use a custom OData `$filter` query for advanced filtering. Must return exactly one application.

You must specify exactly one of these attributes. Using multiple lookup methods simultaneously will result in a validation error.

### OData Query Examples

When using `odata_query`, you can use any valid OData filter expression without the `$filter=` prefix. The provider automatically applies the necessary OData parameters and headers (including ConsistencyLevel: eventual) when executing the query.

**Note:** Do not include `$filter=` in the `odata_query` value. Provide only the filter expression itself.

## Example Usage

### Lookup by Object ID

```terraform
# Look up an application by its object ID (most direct and efficient method)
data "microsoft365_graph_beta_applications_application" "by_object_id" {
  object_id = "00000000-0000-0000-0000-000000000000" # Replace with actual object ID
}

# Output the application details
output "app_by_object_id" {
  value = {
    id           = data.microsoft365_graph_beta_applications_application.by_object_id.id
    app_id       = data.microsoft365_graph_beta_applications_application.by_object_id.app_id
    display_name = data.microsoft365_graph_beta_applications_application.by_object_id.display_name
  }
}
```

### Lookup by App ID

```terraform
# Look up an application by its application (client) ID
# This is useful when you know the app ID but not the object ID
data "microsoft365_graph_beta_applications_application" "by_app_id" {
  app_id = "00000003-0000-0000-c000-000000000000" # Example: Microsoft Graph app ID
}

# Output the application details
output "app_by_app_id" {
  value = {
    id               = data.microsoft365_graph_beta_applications_application.by_app_id.id
    display_name     = data.microsoft365_graph_beta_applications_application.by_app_id.display_name
    sign_in_audience = data.microsoft365_graph_beta_applications_application.by_app_id.sign_in_audience
    identifier_uris  = data.microsoft365_graph_beta_applications_application.by_app_id.identifier_uris
    publisher_domain = data.microsoft365_graph_beta_applications_application.by_app_id.publisher_domain
  }
}
```

### Lookup by Display Name

```terraform
# Look up an application by its exact display name (case-sensitive)
data "microsoft365_graph_beta_applications_application" "by_display_name" {
  display_name = "My Application"
}

# Output the application details
output "app_by_display_name" {
  value = {
    id          = data.microsoft365_graph_beta_applications_application.by_display_name.id
    app_id      = data.microsoft365_graph_beta_applications_application.by_display_name.app_id
    description = data.microsoft365_graph_beta_applications_application.by_display_name.description
    tags        = data.microsoft365_graph_beta_applications_application.by_display_name.tags
  }
}
```

### OData Query - Simple

```terraform
# Look up an application using a simple OData query filter
# This example finds an application that starts with a specific prefix
data "microsoft365_graph_beta_applications_application" "by_odata_simple" {
  odata_query = "startswith(displayName, 'Contoso')"
}

# Output the application details
output "app_by_odata_simple" {
  value = {
    id           = data.microsoft365_graph_beta_applications_application.by_odata_simple.id
    app_id       = data.microsoft365_graph_beta_applications_application.by_odata_simple.app_id
    display_name = data.microsoft365_graph_beta_applications_application.by_odata_simple.display_name
  }
}
```

### OData Query - Complex

```terraform
# Look up an application using a complex OData query with multiple conditions
# This example finds a single-tenant application with a specific name pattern
data "microsoft365_graph_beta_applications_application" "by_odata_complex" {
  odata_query = "displayName eq 'Production API' and signInAudience eq 'AzureADMyOrg'"
}

# Output comprehensive application details
output "app_by_odata_complex" {
  value = {
    id                            = data.microsoft365_graph_beta_applications_application.by_odata_complex.id
    app_id                        = data.microsoft365_graph_beta_applications_application.by_odata_complex.app_id
    display_name                  = data.microsoft365_graph_beta_applications_application.by_odata_complex.display_name
    sign_in_audience              = data.microsoft365_graph_beta_applications_application.by_odata_complex.sign_in_audience
    identifier_uris               = data.microsoft365_graph_beta_applications_application.by_odata_complex.identifier_uris
    is_fallback_public_client     = data.microsoft365_graph_beta_applications_application.by_odata_complex.is_fallback_public_client
    is_device_only_auth_supported = data.microsoft365_graph_beta_applications_application.by_odata_complex.is_device_only_auth_supported
  }
}
```

### OData Query - Filter by Tags

```terraform
# Look up an application using OData query filtering by tags
# This example finds an application with a specific tag
data "microsoft365_graph_beta_applications_application" "by_tags" {
  odata_query = "tags/any(t:t eq 'Production')"
}

# Output the application details with tags
output "app_by_tags" {
  value = {
    id           = data.microsoft365_graph_beta_applications_application.by_tags.id
    app_id       = data.microsoft365_graph_beta_applications_application.by_tags.app_id
    display_name = data.microsoft365_graph_beta_applications_application.by_tags.display_name
    tags         = data.microsoft365_graph_beta_applications_application.by_tags.tags
  }
}
```

### Retrieving API Permissions

```terraform
# Look up an application and output its API permissions configuration
data "microsoft365_graph_beta_applications_application" "with_api_permissions" {
  display_name = "My API Application"
}

# Output API configuration details
output "api_configuration" {
  value = {
    id                       = data.microsoft365_graph_beta_applications_application.with_api_permissions.id
    display_name             = data.microsoft365_graph_beta_applications_application.with_api_permissions.display_name
    identifier_uris          = data.microsoft365_graph_beta_applications_application.with_api_permissions.identifier_uris
    api                      = data.microsoft365_graph_beta_applications_application.with_api_permissions.api
    app_roles                = data.microsoft365_graph_beta_applications_application.with_api_permissions.app_roles
    required_resource_access = data.microsoft365_graph_beta_applications_application.with_api_permissions.required_resource_access
  }
  sensitive = true # API configuration may contain sensitive scope information
}
```

### Retrieving Authentication Configuration

```terraform
# Look up an application and output its authentication configuration
data "microsoft365_graph_beta_applications_application" "with_auth_config" {
  app_id = "00000000-0000-0000-0000-000000000000" # Replace with actual app ID
}

# Output web application authentication configuration
output "web_auth_config" {
  value = {
    id                        = data.microsoft365_graph_beta_applications_application.with_auth_config.id
    display_name              = data.microsoft365_graph_beta_applications_application.with_auth_config.display_name
    sign_in_audience          = data.microsoft365_graph_beta_applications_application.with_auth_config.sign_in_audience
    is_fallback_public_client = data.microsoft365_graph_beta_applications_application.with_auth_config.is_fallback_public_client
    web                       = data.microsoft365_graph_beta_applications_application.with_auth_config.web
    spa                       = data.microsoft365_graph_beta_applications_application.with_auth_config.spa
    public_client             = data.microsoft365_graph_beta_applications_application.with_auth_config.public_client
  }
}
```

### Retrieving Credentials Metadata

```terraform
# Look up an application and output its credentials information
# Note: Actual secret values are not returned by the API for security reasons
data "microsoft365_graph_beta_applications_application" "with_credentials" {
  display_name = "My Application with Credentials"
}

# Output credentials metadata (not the actual secrets)
output "credentials_info" {
  value = {
    id                   = data.microsoft365_graph_beta_applications_application.with_credentials.id
    display_name         = data.microsoft365_graph_beta_applications_application.with_credentials.display_name
    key_credentials      = data.microsoft365_graph_beta_applications_application.with_credentials.key_credentials
    password_credentials = data.microsoft365_graph_beta_applications_application.with_credentials.password_credentials
  }
  sensitive = true # Credentials information should be marked sensitive
}

# Note: To retrieve the public key value in key_credentials, you must use $select=keyCredentials
# in an OData query. The key value is only returned when explicitly requested.
```

### Using with Service Principal Resource

```terraform
# Example: Using the application datasource to reference an existing application
# and create a service principal for it

# Look up an existing application
data "microsoft365_graph_beta_applications_application" "existing" {
  display_name = "My Existing Application"
}

# Create a service principal for the application
resource "microsoft365_graph_beta_applications_service_principal" "sp" {
  app_id                       = data.microsoft365_graph_beta_applications_application.existing.app_id
  account_enabled              = true
  app_role_assignment_required = true

  tags = [
    "WindowsAzureActiveDirectoryIntegratedApp",
    "Production"
  ]
}

# Output the relationship
output "application_and_sp" {
  value = {
    application_id           = data.microsoft365_graph_beta_applications_application.existing.id
    application_app_id       = data.microsoft365_graph_beta_applications_application.existing.app_id
    application_name         = data.microsoft365_graph_beta_applications_application.existing.display_name
    service_principal_id     = microsoft365_graph_beta_applications_service_principal.sp.id
    service_principal_app_id = microsoft365_graph_beta_applications_service_principal.sp.app_id
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `app_id` (String) The unique identifier for the application that is assigned by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`). One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.
- `display_name` (String) The display name for the application. Maximum length is 256 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on null values), `$search`, and `$orderby`. One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.
- `object_id` (String) The object ID of the application. One of `object_id`, `app_id`, `display_name`, or `odata_query` must be specified.
- `odata_query` (String) Custom OData filter query. Use this for advanced filtering when the standard lookup attributes don't meet your needs. Cannot be combined with `object_id`, `app_id`, or `display_name`. Example: `displayName eq 'My Application' and signInAudience eq 'AzureADMyOrg'`
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `api` (Attributes) Specifies settings for an application that implements a web API. (see [below for nested schema](#nestedatt--api))
- `app_roles` (Attributes Set) The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable. (see [below for nested schema](#nestedatt--app_roles))
- `created_date_time` (String) The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, and `eq` on null values) and `$orderby`.
- `deleted_date_time` (String) The date and time the application was deleted. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
- `description` (String) Free text field to provide a description of the application object to end users. The maximum allowed size is 1,024 characters. Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `startsWith`) and `$search`.
- `disabled_by_microsoft_status` (String) Specifies whether Microsoft has disabled the registered application. The possible values are: null (default value), `NotDisabled`, and `DisabledDueToViolationOfServicesAgreement` (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement). Supports `$filter` (`eq`, `ne`, `not`). Read-only.
- `group_membership_claims` (Set of String) Configures the groups claim issued in a user or OAuth 2.0 access token that the application expects. To set this attribute, use one of the following string values: `None`, `SecurityGroup` (for security groups and Microsoft Entra roles), `All` (this gets all security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of).
- `id` (String) Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).
- `identifier_uris` (Set of String) Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique across Microsoft Entra ID. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).
- `info` (Attributes) Basic profile information of the application, such as it's marketing, support, terms of service, and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more information, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, and `eq` on null values). (see [below for nested schema](#nestedatt--info))
- `is_device_only_auth_supported` (Boolean) Specifies whether this application supports device authentication without a user. The default is false.
- `is_fallback_public_client` (Boolean) Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false, which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID can't determine the client application type. For example, the ROPC flow where the application is configured without specifying a redirect URI. In those cases Microsoft Entra ID interprets the application type based on the value of this property.
- `key_credentials` (Attributes Set) The collection of key credentials associated with the application. This is a read-only attribute. To manage certificate credentials, use the `microsoft365_graph_beta_applications_application_certificate_credential` resource instead. (see [below for nested schema](#nestedatt--key_credentials))
- `notes` (String) Notes relevant for the management of the application.
- `optional_claims` (Attributes) Application developers can configure optional claims in their Microsoft Entra applications to specify the claims that are sent to their application by the Microsoft security token service. For more information, see How to: Provide optional claims to your app. (see [below for nested schema](#nestedatt--optional_claims))
- `owner_user_ids` (Set of String) The user IDs of the owners for the application. At least one owner is typically required when creating an application. Owners are a set of non-admin users or service principals allowed to modify this object.
- `parental_control_settings` (Attributes) Specifies parental control settings for an application. (see [below for nested schema](#nestedatt--parental_control_settings))
- `password_credentials` (Attributes Set) The collection of password credentials associated with the application. Not nullable. (see [below for nested schema](#nestedatt--password_credentials))
- `public_client` (Attributes) Specifies settings for installed clients such as desktop or mobile devices. (see [below for nested schema](#nestedatt--public_client))
- `publisher_domain` (String) The verified publisher domain for the application. Read-only. Supports `$filter` (`eq`, `ne`, `ge`, `le`, `startsWith`).
- `required_resource_access` (Attributes Set) Specifies the resources that the application needs to access. This property also specifies the set of delegated permissions and application roles that it needs for each of those resources. This configuration of access to the required resources drives the consent experience. No more than 50 resource services (APIs) can be configured. Beginning mid-October 2021, the total number of required permissions must not exceed 400. For more information, see Limits on requested permissions per app. Not nullable. Supports `$filter` (`eq`, `not`, `ge`, `le`). (see [below for nested schema](#nestedatt--required_resource_access))
- `service_management_reference` (String) References application or service contact information from a Service or Asset Management database. Nullable.
- `sign_in_audience` (String) Specifies the Microsoft accounts that are supported for the current application. The possible values are: `AzureADMyOrg` (default), `AzureADMultipleOrgs`, `AzureADandPersonalMicrosoftAccount`, and `PersonalMicrosoftAccount`. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you may need to change other properties first. For more information, see Validation differences for signInAudience. Supports `$filter` (`eq`, `ne`, `not`).
- `sign_in_audience_restrictions` (Attributes) Specifies restrictions on the supported account types specified in signInAudience. The value type determines the restrictions that can be applied: unrestrictedAudience (There are no additional restrictions on the supported account types allowed by signInAudience) or allowedTenantsAudience (The application can only be used in the specified Entra tenants. Only supported when signInAudience is AzureADMultipleOrgs). Default is a value of type unrestrictedAudience. Returned only on `$select`. (see [below for nested schema](#nestedatt--sign_in_audience_restrictions))
- `spa` (Attributes) Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens. (see [below for nested schema](#nestedatt--spa))
- `tags` (Set of String) Custom strings that can be used to categorize and identify the application. Not nullable. Strings added here also appear in the tags property of any associated service principals. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`) and `$search`.
- `web` (Attributes) Specifies settings for a web application. (see [below for nested schema](#nestedatt--web))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--api"></a>
### Nested Schema for `api`

Read-Only:

- `accept_mapped_claims` (Boolean) Allows an application to use claims mapping without specifying a custom signing key.
- `known_client_applications` (Set of String) Used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app. If you set the appID of the client app to this value, the user only consents once to the client app. Microsoft Entra ID knows that consenting to the client means implicitly consenting to the web API and automatically provisions service principals for both APIs at the same time. Both the client and the web API app must be registered in the same tenant.
- `oauth2_permission_scopes` (Attributes Set) The definition of the delegated permissions exposed by the web API represented by this application registration. These delegated permissions may be requested by a client application, and may be granted by users or administrators during consent. Delegated permissions are sometimes referred to as OAuth 2.0 scopes. (see [below for nested schema](#nestedatt--api--oauth2_permission_scopes))
- `pre_authorized_applications` (Attributes Set) Lists the client applications that are preauthorized with the specified delegated permissions to access this application's APIs. Users aren't required to consent to any preauthorized application (for the permissions specified). However, any other permissions not listed in preAuthorizedApplications (requested through incremental consent for example) will require user consent. (see [below for nested schema](#nestedatt--api--pre_authorized_applications))
- `requested_access_token_version` (Number) Specifies the access token version expected by this resource. This changes the version and format of the JWT produced independent of the endpoint or client used to request the access token. The endpoint used, v1.0 or v2.0, is chosen by the client and only impacts the version of id_tokens. Resources need to explicitly configure requestedAccessTokenVersion to indicate the supported access token format. Possible values for requestedAccessTokenVersion are 1, 2, or null. If the value is null, this defaults to 1, which corresponds to the v1.0 endpoint. If signInAudience on the application is configured as AzureADandPersonalMicrosoftAccount or PersonalMicrosoftAccount, the value for this property must be 2.

<a id="nestedatt--api--oauth2_permission_scopes"></a>
### Nested Schema for `api.oauth2_permission_scopes`

Read-Only:

- `admin_consent_description` (String) A description of the delegated permissions, intended to be read by an administrator granting the permission on behalf of all users. This text appears in tenant-wide admin consent experiences.
- `admin_consent_display_name` (String) The permission's title, intended to be read by an administrator granting the permission on behalf of all users.
- `id` (String) Unique scope permission identifier inside the oauth2PermissionScopes collection. Required.
- `is_enabled` (Boolean) When you create or update a permission, this property must be set to true (which is the default). To delete a permission, this property must first be set to false. At that point, in a subsequent call, the permission may be removed.
- `type` (String) The possible values are: `User` and `Admin`. Specifies whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator consent should be required for the permissions. While Microsoft Graph defines the default consent requirement for each permission, the tenant administrator may override the behavior in their organization (by allowing, restricting, or limiting user consent to this delegated permission). For more information, see Configure how users consent to applications.
- `user_consent_description` (String) A description of the delegated permissions, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.
- `user_consent_display_name` (String) A title for the permission, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.
- `value` (String) Specifies the value to include in the scp (scope) claim in access tokens. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with `.`.


<a id="nestedatt--api--pre_authorized_applications"></a>
### Nested Schema for `api.pre_authorized_applications`

Read-Only:

- `app_id` (String) The unique identifier for the client application.
- `delegated_permission_ids` (Set of String) The unique identifier for the scopes the client application is granted.



<a id="nestedatt--app_roles"></a>
### Nested Schema for `app_roles`

Read-Only:

- `allowed_member_types` (Set of String) Specifies whether this app role can be assigned to users and groups (by setting to `['User']`), to other application's (by setting to `['Application']`, or both (by setting to `['User', 'Application']`). App roles supporting assignment to other applications' service principals are also known as application permissions. The 'Application' value is only supported for app roles defined on application entities. Required.
- `description` (String) The description for the app role. This is displayed when the app role is being assigned and, if the app role functions as an application permission, during consent experiences. Required.
- `display_name` (String) Display name for the permission that appears in the app role assignment and consent experiences. Required.
- `id` (String) Unique role identifier inside the appRoles collection. When creating a new app role, a new GUID identifier must be provided. Required.
- `is_enabled` (Boolean) Defines whether the application's app role is enabled or disabled. Required.
- `origin` (String) Specifies if the app role is defined on the application object or on the servicePrincipal entity. Must not be included in any POST or PATCH requests. Read-only.
- `value` (String) Specifies the value to include in the roles claim in ID tokens and access tokens authenticating an assigned user or service principal. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with `.`. Nullable.


<a id="nestedatt--info"></a>
### Nested Schema for `info`

Read-Only:

- `logo_url` (String) CDN URL to the application's logo. Read-only.
- `marketing_url` (String) Link to the application's marketing page. For example, https://www.contoso.com/app/marketing.
- `privacy_statement_url` (String) Link to the application's privacy statement. For example, https://www.contoso.com/app/privacy.
- `support_url` (String) Link to the application's support page. For example, https://www.contoso.com/app/support.
- `terms_of_service_url` (String) Link to the application's terms of service statement. For example, https://www.contoso.com/app/termsofservice.


<a id="nestedatt--key_credentials"></a>
### Nested Schema for `key_credentials`

Read-Only:

- `custom_key_identifier` (String) A 40-character binary type that can be used to identify the credential. Optional. When not provided in the payload, defaults to the thumbprint of the certificate.
- `display_name` (String) Friendly name for the key. Optional.
- `end_date_time` (String) The date and time at which the credential expires. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
- `key` (String, Sensitive) Value for the key credential. Should be a Base64 encoded value. Returned only on $select for a single object, that is, GET applications/{applicationId}?$select=keyCredentials or GET servicePrincipals/{servicePrincipalId}?$select=keyCredentials; otherwise, it's always null. From a .cer certificate, you can read the key using the Convert.ToBase64String() method. For more information, see Get the certificate key.
- `key_id` (String) The unique identifier (GUID) for the key.
- `start_date_time` (String) The date and time at which the credential becomes valid. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
- `type` (String) The type of key credential; for example, `Symmetric`, `AsymmetricX509Cert`, or `X509CertAndPassword`.
- `usage` (String) A string that describes the purpose for which the key can be used; for example, `None​`, `Verify`, `PairwiseIdentifier`, `Sign`.


<a id="nestedatt--optional_claims"></a>
### Nested Schema for `optional_claims`

Read-Only:

- `access_token` (Attributes Set) The optional claims returned in the JWT access token. (see [below for nested schema](#nestedatt--optional_claims--access_token))
- `id_token` (Attributes Set) The optional claims returned in the JWT ID token. (see [below for nested schema](#nestedatt--optional_claims--id_token))
- `saml2_token` (Attributes Set) The optional claims returned in the SAML token. (see [below for nested schema](#nestedatt--optional_claims--saml2_token))

<a id="nestedatt--optional_claims--access_token"></a>
### Nested Schema for `optional_claims.access_token`

Read-Only:

- `additional_properties` (Set of String) Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
- `essential` (Boolean) If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
- `name` (String) The name of the optional claim. Required.
- `source` (String) The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.


<a id="nestedatt--optional_claims--id_token"></a>
### Nested Schema for `optional_claims.id_token`

Read-Only:

- `additional_properties` (Set of String) Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
- `essential` (Boolean) If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
- `name` (String) The name of the optional claim. Required.
- `source` (String) The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.


<a id="nestedatt--optional_claims--saml2_token"></a>
### Nested Schema for `optional_claims.saml2_token`

Read-Only:

- `additional_properties` (Set of String) Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.
- `essential` (Boolean) If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.
- `name` (String) The name of the optional claim. Required.
- `source` (String) The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.



<a id="nestedatt--parental_control_settings"></a>
### Nested Schema for `parental_control_settings`

Read-Only:

- `countries_blocked_for_minors` (Set of String) Specifies the two-letter ISO country codes. Access to the application will be blocked for minors from the countries specified in this list.
- `legal_age_group_rule` (String) Specifies the legal age group rule that applies to users of the app. Can be set to one of the following values: `Allow`, `RequireConsentForPrivacyServices`, `RequireConsentForMinors`, `RequireConsentForKids`, `BlockMinors`.


<a id="nestedatt--password_credentials"></a>
### Nested Schema for `password_credentials`

Read-Only:

- `custom_key_identifier` (String) Do not use.
- `display_name` (String) Friendly name for the password. Optional.
- `end_date_time` (String) The date and time at which the password expires represented using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.
- `hint` (String) Contains the first three characters of the password. Read-only.
- `key_id` (String) The unique identifier for the password. Required.
- `secret_text` (String, Sensitive) Read-only; Contains the strong passwords generated by Microsoft Entra ID that are 16-64 characters in length. The generated password value is only returned during the initial POST request to addPassword. There is no way to retrieve this password in the future.
- `start_date_time` (String) The date and time at which the password becomes valid. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.


<a id="nestedatt--public_client"></a>
### Nested Schema for `public_client`

Read-Only:

- `redirect_uris` (Set of String) Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.


<a id="nestedatt--required_resource_access"></a>
### Nested Schema for `required_resource_access`

Read-Only:

- `resource_access` (Attributes Set) The list of OAuth2.0 permission scopes and app roles that the application requires from the specified resource. Required. (see [below for nested schema](#nestedatt--required_resource_access--resource_access))
- `resource_app_id` (String) The unique identifier for the resource that the application requires access to. This should be equal to the appId declared on the target resource application. Required.

<a id="nestedatt--required_resource_access--resource_access"></a>
### Nested Schema for `required_resource_access.resource_access`

Read-Only:

- `id` (String) The unique identifier of an app role or delegated permission exposed by the resource application. For delegated permissions, this should match the id property of one of the delegated permissions in the oauth2PermissionScopes collection of the resource application's service principal. For app roles (application permissions), this should match the id property of an app role in the appRoles collection of the resource application's service principal. Required.
- `type` (String) Specifies whether the id property references a delegated permission or an app role (application permission). The possible values are: `Scope` (for delegated permissions) or `Role` (for app roles). Required.



<a id="nestedatt--sign_in_audience_restrictions"></a>
### Nested Schema for `sign_in_audience_restrictions`

Read-Only:

- `allowed_tenant_ids` (Set of String) The list of allowed tenant IDs. Only applicable when odata_type is `#microsoft.graph.allowedTenantsAudience`.
- `is_home_tenant_allowed` (Boolean) Indicates whether the home tenant is allowed. Only applicable when odata_type is `#microsoft.graph.allowedTenantsAudience`.
- `odata_type` (String) The OData type. Must be `#microsoft.graph.allowedTenantsAudience` or `#microsoft.graph.unrestrictedAudience`.


<a id="nestedatt--spa"></a>
### Nested Schema for `spa`

Read-Only:

- `redirect_uris` (Set of String) Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.


<a id="nestedatt--web"></a>
### Nested Schema for `web`

Read-Only:

- `home_page_url` (String) Home page or landing page of the application.
- `implicit_grant_settings` (Attributes) Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow. (see [below for nested schema](#nestedatt--web--implicit_grant_settings))
- `logout_url` (String) Specifies the URL that is used by Microsoft's authorization service to log out a user using front-channel, back-channel or SAML logout protocols.
- `redirect_uri_settings` (Attributes Set) Specifies the index of the URLs where user tokens are sent for sign-in. This is only valid for applications using SAML. Note: If not specified, the API may auto-generate settings based on redirect_uris. To manage this field, you must provide at least one entry; empty arrays are not supported as the API auto-generates values. (see [below for nested schema](#nestedatt--web--redirect_uri_settings))
- `redirect_uris` (Set of String) Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.

<a id="nestedatt--web--implicit_grant_settings"></a>
### Nested Schema for `web.implicit_grant_settings`

Read-Only:

- `enable_access_token_issuance` (Boolean) Specifies whether this web application can request an access token using the OAuth 2.0 implicit flow.
- `enable_id_token_issuance` (Boolean) Specifies whether this web application can request an ID token using the OAuth 2.0 implicit flow.


<a id="nestedatt--web--redirect_uri_settings"></a>
### Nested Schema for `web.redirect_uri_settings`

Read-Only:

- `index` (Number) Identifies the specific URI within the redirectURIs collection in SAML SSO flows. Defaults to null. The index is unique across all the redirectUris for the application.
- `uri` (String) Specifies the URI that tokens are sent to.
