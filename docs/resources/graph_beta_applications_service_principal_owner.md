---
page_title: "microsoft365_graph_beta_applications_service_principal_owner Resource - terraform-provider-microsoft365"
subcategory: "Applications"
description: |-
  Manages an owner assignment for a Microsoft Entra Service Principal using the /servicePrincipals/{id}/owners endpoint. Owners are users or service principals who are allowed to modify the service principal object. As a recommended best practice, service principals should have at least two owners.
  For more information, see the Microsoft Graph API documentation https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-owners?view=graph-rest-beta.
---

# microsoft365_graph_beta_applications_service_principal_owner (Resource)

Manages an owner assignment for a Microsoft Entra Service Principal using the `/servicePrincipals/{id}/owners` endpoint. Owners are users or service principals who are allowed to modify the service principal object. As a recommended best practice, service principals should have at least two owners.

For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-owners?view=graph-rest-beta).

## Microsoft Documentation

- [servicePrincipal: List owners](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-list-owners?view=graph-rest-beta&tabs=http)
- [servicePrincipal: Add owner](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-owners?view=graph-rest-beta&tabs=http)
- [servicePrincipal: Remove owner](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-delete-owners?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Application.Read.All`
- `Directory.Read.All`
- `Application.ReadWrite.All`
- `Directory.ReadWrite.All`

**Optional:**
- `Application.ReadWrite.OwnedBy` (if managing only applications owned by the service principal)
- `User.Read.All` (when assigning user owners)
- `Application.Read.All` (when assigning service principal owners)

Find out more about the permissions required for managing service principals at Microsoft Learn [here](https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-beta).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.43.0 | Experimental | Initial release |

## Important Notes

- **Owner Types**: Owners can be either Users or Service Principals
- **Minimum One Owner**: Service principals should have at least one owner for proper management
- **Owner Object Type**: The `owner_object_type` attribute must be set correctly based on the owner type being assigned
- **Enterprise Applications**: Service principal owners manage the tenant-specific instance of an application

## Example Usage

### User as Service Principal Owner

```terraform
resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-application"
  description  = "Application with service principal"
}

resource "microsoft365_graph_beta_applications_service_principal" "example" {
  app_id = microsoft365_graph_beta_applications_application.example.app_id
}

resource "microsoft365_graph_beta_users_user" "sp_owner" {
  display_name        = "Service Principal Owner"
  user_principal_name = "sp.owner@mycompany.com"
  mail_nickname       = "sp.owner"
  account_enabled     = true
  password_profile = {
    password                           = "TempP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Assign user as service principal owner
resource "microsoft365_graph_beta_applications_service_principal_owner" "user_owner" {
  service_principal_id = microsoft365_graph_beta_applications_service_principal.example.id
  owner_id             = microsoft365_graph_beta_users_user.sp_owner.id
  owner_object_type    = "User"
}
```

### Service Principal as Service Principal Owner

```terraform
resource "microsoft365_graph_beta_applications_application" "managed_sp_app" {
  display_name = "my-managed-service-principal"
  description  = "Service principal managed by another service principal"
}

resource "microsoft365_graph_beta_applications_service_principal" "managed_sp" {
  app_id = microsoft365_graph_beta_applications_application.managed_sp_app.app_id
}

resource "microsoft365_graph_beta_applications_application" "manager_app" {
  display_name = "sp-manager"
  description  = "Service principal that manages other service principals"
}

resource "microsoft365_graph_beta_applications_service_principal" "manager_sp" {
  app_id = microsoft365_graph_beta_applications_application.manager_app.app_id
}

# Assign service principal as another service principal's owner
resource "microsoft365_graph_beta_applications_service_principal_owner" "sp_owner" {
  service_principal_id = microsoft365_graph_beta_applications_service_principal.managed_sp.id
  owner_id             = microsoft365_graph_beta_applications_service_principal.manager_sp.id
  owner_object_type    = "ServicePrincipal"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `owner_id` (String) The unique identifier (UUID) for the owner to be added to the service principal. This can be a user or service principal.
- `owner_object_type` (String) The type of object being added as an owner. This determines the correct Microsoft Graph API endpoint to use. Valid values: 'User', 'ServicePrincipal'.
- `service_principal_id` (String) The unique identifier (UUID) for the service principal.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for this service principal owner assignment. This is a composite ID formed by combining the service principal ID and owner ID.
- `owner_display_name` (String) The display name of the owner. Read-only.
- `owner_type` (String) The type of the owner object as returned by Microsoft Graph (e.g., 'User', 'ServicePrincipal'). Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

```shell
# Import a service principal owner by composite ID: service_principal_id/owner_id
terraform import microsoft365_graph_beta_applications_service_principal_owner.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"
```
