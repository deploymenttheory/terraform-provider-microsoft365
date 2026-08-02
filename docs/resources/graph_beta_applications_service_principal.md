---
page_title: "microsoft365_graph_beta_applications_service_principal Resource - terraform-provider-microsoft365"
subcategory: "Applications"
description: |-
  Manages a Service Principal in Microsoft Entra ID. Service principals are the local representation of an application object in a specific tenant. They define what the app can do in the specific tenant, who can access the app, and what resources the app can access.
  This resource models only the properties the service principal itself owns. Properties that Microsoft Entra reflects from the backing application — such as its display name, home page, logout URL, reply URLs, app roles and permission scopes — are deliberately not exposed here and will not be added: they change whenever the application changes. Read them from the microsoft365_graph_beta_applications_application resource, or from the service principal data source for an application this configuration does not manage.
  For more information, see the Microsoft Graph API documentation https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-serviceprincipals?view=graph-rest-beta.
---

# microsoft365_graph_beta_applications_service_principal (Resource)

Manages a Service Principal in Microsoft Entra ID. Service principals are the local representation of an application object in a specific tenant. They define what the app can do in the specific tenant, who can access the app, and what resources the app can access.

This resource models only the properties the service principal itself owns. Properties that Microsoft Entra reflects from the backing application — such as its display name, home page, logout URL, reply URLs, app roles and permission scopes — are deliberately not exposed here and will not be added: they change whenever the application changes. Read them from the `microsoft365_graph_beta_applications_application` resource, or from the service principal data source for an application this configuration does not manage.

For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-serviceprincipals?view=graph-rest-beta).

## Microsoft Documentation

- [servicePrincipal resource type](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-beta)
- [Create servicePrincipal](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-serviceprincipals?view=graph-rest-beta&tabs=http)
- [Get servicePrincipal](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-get?view=graph-rest-beta&tabs=http)
- [Update servicePrincipal](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-update?view=graph-rest-beta&tabs=http)
- [Delete servicePrincipal](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Application.Read.All`
- `Application.ReadWrite.All`
- `Directory.Read.All`
- `Directory.ReadWrite.All`

**Optional:**
- `Application.ReadWrite.OwnedBy` (if managing only applications owned by the service principal)

Find out more about the permissions required for managing service principals at Microsoft Learn [here](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-beta).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.43.0 | Experimental | Initial release |

## Important Notes

- **Enterprise Applications**: Service principals are the tenant-specific representation of applications (also called Enterprise Applications in the portal)
- **App ID Required**: You must provide the `app_id` (client ID) of an existing application registration
- **Single Sign-On**: Configure `preferred_single_sign_on_mode` for SSO (options: `saml`, `password`, `notSupported`)
- **App Role Assignment**: Set `app_role_assignment_required = true` to require users/groups be assigned before accessing the application
- **Tags**: Special tags control visibility and behavior:
  - `HideApp`: Hides the application from My Apps portal
  - `WindowsAzureActiveDirectoryIntegratedApp`: Marks as integrated app

## Example Usage

### Minimal Service Principal

```terraform
resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-application"
  description  = "Application for service principal"
}

# Create service principal for the application
resource "microsoft365_graph_beta_applications_service_principal" "example" {
  app_id = microsoft365_graph_beta_applications_application.example.app_id
}
```

### Service Principal with Full Configuration

```terraform
resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-saml-application"
  description  = "SAML-based enterprise application"
}

# Create service principal with full configuration
resource "microsoft365_graph_beta_applications_service_principal" "example" {
  app_id                        = microsoft365_graph_beta_applications_application.example.app_id
  account_enabled               = true
  app_role_assignment_required  = true
  description                   = "Enterprise application for SSO access"
  login_url                     = "https://login.mycompany.com"
  notes                         = "Managed by Terraform - Production environment"
  notification_email_addresses  = ["admin@mycompany.com", "security@mycompany.com"]
  preferred_single_sign_on_mode = "saml"

  tags = [
    "HideApp",
    "WindowsAzureActiveDirectoryIntegratedApp"
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_id` (String) The application (client) ID of the application for which to create the service principal. Required.

### Optional

- `account_enabled` (Boolean) True if the service principal account is enabled; otherwise, false. If set to false, then no users are able to sign in to this app, even if they're assigned to it. Supports `$filter` (`eq`, `ne`, `not`, `in`).
- `alternative_names` (Set of String) Used to retrieve service principals by subscription, identify resource group and full resource IDs for managed identities. Set to `[]` to clear previously configured values. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`).
- `app_role_assignment_required` (Boolean) Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports `$filter` (`eq`, `ne`, `NOT`).
- `description` (String) Free text field to provide an internal end-user facing description of the service principal. End-user portals such MyApps displays the application description in this field. The maximum allowed size is 1,024 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `startsWith`) and `$search`.
- `hard_delete` (Boolean) When `true`, the service principal will be permanently deleted (hard delete) during destroy. When `false` (default), the service principal will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. Note: This field defaults to `false` on import since the API does not return this value.
- `login_url` (String) Specifies the URL where the service provider redirects the user to Microsoft Entra ID to authenticate. Microsoft Entra ID uses the URL to launch the application from Microsoft 365 or the Microsoft Entra My Apps. When blank, Microsoft Entra ID performs IdP-initiated sign-on for applications configured with SAML-based single sign-on.
- `notes` (String) Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1,024 characters.
- `notification_email_addresses` (Set of String) Specifies the list of email addresses where Microsoft Entra ID sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Microsoft Entra Gallery applications.
- `preferred_single_sign_on_mode` (String) Specifies the single sign-on mode configured for this application. Microsoft Entra ID uses the preferred single sign-on mode to launch the application from Microsoft 365 or the Microsoft Entra My Apps. The supported values are `password`, `saml`, `notSupported`, and `oidc`.
- `saml_single_sign_on_settings` (Attributes) The collection for settings related to SAML single sign-on. When this block is absent from the configuration, updates actively clear the property on the service principal. (see [below for nested schema](#nestedatt--saml_single_sign_on_settings))
- `tags` (Set of String) Custom strings that can be used to categorize and identify the service principal. Note: Microsoft may automatically add system-managed tags in addition to the tags you specify.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `token_encryption_key_id` (String) Specifies the keyId of a public key from the key credentials collection. When configured, Microsoft Entra ID issues tokens for this application encrypted using the key specified by this property. The referenced key credential must already exist on the service principal. When this attribute is absent from the configuration, updates actively clear the property on the service principal.

### Read-Only

- `app_owner_organization_id` (String) Contains the tenant ID where the application is registered. Applicable only to service principals backed by applications. Read-only. Equivalent to `application_tenant_id` in the azuread provider.
- `application_template_id` (String) Unique identifier of the application template that the associated application was created from. Read-only. `null` if the app wasn't created from an application template.
- `created_by_app_id` (String) The appId of the application that created this service principal. Set internally by Microsoft Entra ID. Read-only.
- `id` (String) The unique identifier (object ID) for the service principal. Read-only.
- `is_disabled` (Boolean) Specifies whether the service principal is deactivated so the app can't obtain new access tokens or access protected resources. Read-only; the API rejects writes to this property.
- `key_credentials` (Attributes Set) The collection of key credentials associated with the service principal. Read-only on this resource; certificates are added through dedicated credential resources or the addTokenSigningCertificate API. Private key material is never returned by the API, and this resource does not expose the public `key` field. (see [below for nested schema](#nestedatt--key_credentials))
- `password_credentials` (Attributes Set) The collection of password credentials associated with the service principal. Read-only. The secret itself is never returned by the API. (see [below for nested schema](#nestedatt--password_credentials))
- `preferred_token_signing_key_end_date_time` (String) Specifies the expiration date of the key credential used for token signing, marked by `preferred_token_signing_key_thumbprint`. Read-only.
- `preferred_token_signing_key_thumbprint` (String) The thumbprint of the certificate used to sign SAML responses for applications with `preferred_single_sign_on_mode` set to `saml`. Read-only on this resource; it is set when a token signing certificate is activated on the service principal.
- `service_principal_type` (String) Identifies if the service principal represents an application or a managed identity. Read-only.

<a id="nestedatt--saml_single_sign_on_settings"></a>
### Nested Schema for `saml_single_sign_on_settings`

Required:

- `relay_state` (String) The relative URI the service provider would redirect to after completion of the single sign-on flow.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--key_credentials"></a>
### Nested Schema for `key_credentials`

Read-Only:

- `custom_key_identifier` (String) A base64-encoded custom key identifier.
- `display_name` (String) The friendly name for the key.
- `end_date_time` (String) The date and time at which the credential expires.
- `key_id` (String) The unique identifier for the key.
- `start_date_time` (String) The date and time at which the credential becomes valid.
- `type` (String) The type of key credential, for example `Symmetric` or `AsymmetricX509Cert`.
- `usage` (String) A string that describes the purpose for which the key can be used, for example `Verify` or `Sign`.


<a id="nestedatt--password_credentials"></a>
### Nested Schema for `password_credentials`

Read-Only:

- `custom_key_identifier` (String) A base64-encoded custom key identifier.
- `display_name` (String) The friendly name for the password.
- `end_date_time` (String) The date and time at which the password expires.
- `hint` (String) Contains the first three characters of the password. Read-only.
- `key_id` (String) The unique identifier for the password.
- `start_date_time` (String) The date and time at which the password becomes valid.

## Import

```shell
# Simple import - defaults to hard_delete=false (soft delete with 30-day recovery)
terraform import microsoft365_graph_beta_applications_service_principal.example "00000000-0000-0000-0000-000000000000"

# Extended import - with hard_delete enabled for permanent deletion
terraform import microsoft365_graph_beta_applications_service_principal.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
```
