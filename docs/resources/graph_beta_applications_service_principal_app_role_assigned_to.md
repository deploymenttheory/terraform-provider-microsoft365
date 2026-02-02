---
page_title: "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to Resource - terraform-provider-microsoft365"
subcategory: "Applications"

description: |-
  Manages app role assignments granted for a service principal using the /servicePrincipals/{id}/appRoleAssignedTo endpoint. This resource is used to enables assigning app roles defined by a resource service principal to users, groups, or client service principals. App roles assigned to service principals are also known as application permissions. These can be granted directly with app role assignments or through a consent experience.
  To grant an app role assignment, you need three identifiers:
  target_service_principal_object_id: The Object ID of the user, group, or client service principal to which you are assigning the app roleresource_object_id: The Object ID of the resource service principal which has defined the app roleapp_role_id: The ID of the appRole (defined on the resource service principal) to assign
  For more information, see the Microsoft Graph API documentation https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-approleassignedto?view=graph-rest-beta..
---

# microsoft365_graph_beta_applications_service_principal_app_role_assigned_to (Resource)

Manages app role assignments granted for a service principal using the `/servicePrincipals/{id}/appRoleAssignedTo` endpoint. This resource is used to enables assigning app roles defined by a resource service principal to users, groups, or client service principals. App roles assigned to service principals are also known as **application permissions**. These can be granted directly with app role assignments or through a consent experience.

To grant an app role assignment, you need three identifiers:
- `target_service_principal_object_id`: The Object ID of the user, group, or client service principal to which you are assigning the app role
- `resource_object_id`: The Object ID of the resource service principal which has defined the app role
- `app_role_id`: The ID of the appRole (defined on the resource service principal) to assign

For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-approleassignedto?view=graph-rest-beta)..

## Microsoft Documentation

- [appRoleAssignment resource type](https://learn.microsoft.com/en-us/graph/api/resources/approleassignment?view=graph-rest-beta)
- [List appRoleAssignments](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-list-approleassignedto?view=graph-rest-beta)
- [Create appRoleAssignment](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-approleassignedto?view=graph-rest-beta)
- [Delete appRoleAssignment](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-delete-approleassignedto?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Application.Read.All`
- `Directory.Read.All`
- `Application.ReadWrite.All`
- `AppRoleAssignment.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.38.0 | Experimental | Initial release |

## Important Notes

- This resource creates an `appRoleAssignment` that grants an app role to a principal (user, group, or service principal)
- The `resource_id` is the service principal ID of the application that defines the app roles (e.g., Microsoft Graph)
- The `principal_id` is the ID of the user, group, or service principal receiving the permission
- The `app_role_id` must be a valid app role ID from the resource application's `appRoles` collection
- Use `00000000-0000-0000-0000-000000000000` as the `app_role_id` for the default app role when no specific role is required

## Common App Role IDs

### Microsoft Graph (appId: 00000003-0000-0000-c000-000000000000)

| Permission | App Role ID | Description |
|------------|-------------|-------------|
| User.Read.All | `df021288-bdef-4463-88db-98f22de89214` | Read all users' full profiles |
| Directory.Read.All | `7ab1d382-f21e-4acd-a863-ba3e13f7da61` | Read directory data |
| Application.Read.All | `9a5d68dd-52b0-4cc2-bd40-abcf44ac3a30` | Read all applications |
| Application.ReadWrite.All | `1bfefb4e-e0b5-418b-a88f-73c46d2cc8e9` | Read and write all applications |

### Microsoft Entra Agent ID Permissions

For managing AI agent identities, see the [Agent ID permissions reference](https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities).

| Permission | App Role ID | Description |
|------------|-------------|-------------|
| AgentIdentity.Read.All | `b2b8f011-2898-4234-9092-5059f6c1ebfa` | Read all agent identities |
| AgentIdentity.ReadWrite.All | `dcf7150a-88d4-4fe6-9be1-c2744c455397` | Read and write all agent identities |
| AgentIdentity.DeleteRestore.All | `5b016f9b-18eb-41d4-869a-66931914d1c8` | Delete and restore agent identities |
| AgentCardManifest.Read.All | `3ee18438-e6e5-4858-8f1c-d7b723b45213` | Read agent card manifests |
| AgentCollection.Read.All | `e65ee1da-d1d5-467b-bdd0-3e9bb94e6e0c` | Read all agent collections |
| AgentInstance.Read.All | `799a4732-85b8-4c67-b048-75f0e88a232b` | Read all agent instances |
| AgentInstance.ReadWrite.All | `07abdd95-78dc-4353-bd32-09f880ea43d0` | Read and write all agent instances |

## Example Usage

### Basic Example

```terraform
# Example: Assign an app role to a service principal
# This grants the "User.Read.All" permission from Microsoft Graph to a client service principal

resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "example" {
  # The Object ID of the service principal that exposes the app role (e.g., Microsoft Graph)
  resource_object_id = "00000003-0000-0000-c000-000000000000" # Microsoft Graph service principal Object ID

  # The app role ID to assign (e.g., User.Read.All = df021288-bdef-4463-88db-98f22de89214)
  app_role_id = "df021288-bdef-4463-88db-98f22de89214"

  # The Object ID of the service principal being granted the app role
  target_service_principal_object_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # Your application's service principal Object ID

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Example: Assign app role to a group
# This grants application permissions to all members of the security group

resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "group_assignment" {
  # The Object ID of the service principal that exposes the app role
  resource_object_id = var.resource_service_principal_object_id

  # The app role ID from the resource's appRoles collection
  app_role_id = var.app_role_id

  # The Object ID of the security group being granted the app role
  target_service_principal_object_id = var.security_group_id
}

# Example: Default app role assignment (no specific role)
# Use this when the application doesn't define specific app roles

resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "default_role" {
  resource_object_id = var.resource_service_principal_object_id

  # Default app role ID when no specific roles are defined
  app_role_id = "00000000-0000-0000-0000-000000000000"

  target_service_principal_object_id = var.client_service_principal_object_id
}
```

### Agent Identity Permissions Example

```terraform
# Example: Assign Microsoft Entra Agent ID permissions to a service principal
# These permissions are required for managing AI agent identities in Microsoft Entra
# Reference: https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities

# Get the Microsoft Graph service principal (resource that defines the permissions)
data "microsoft365_graph_beta_applications_service_principal" "msgraph" {
  filter_type  = "display_name"
  filter_value = "Microsoft Graph"
}

# Example: Agent Identity Read permissions
# AgentIdentity.Read.All - Read all agent identities
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_identity_read" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "b2b8f011-2898-4234-9092-5059f6c1ebfa" # AgentIdentity.Read.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Identity ReadWrite permissions
# AgentIdentity.ReadWrite.All - Read and write all agent identities
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_identity_readwrite" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "dcf7150a-88d4-4fe6-9be1-c2744c455397" # AgentIdentity.ReadWrite.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Identity Delete/Restore permissions
# AgentIdentity.DeleteRestore.All - Delete and restore agent identities
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_identity_delete_restore" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "5b016f9b-18eb-41d4-869a-66931914d1c8" # AgentIdentity.DeleteRestore.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Registry permissions
# AgentCardManifest.Read.All - Read agent card manifests
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_card_manifest_read" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "3ee18438-e6e5-4858-8f1c-d7b723b45213" # AgentCardManifest.Read.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Collection permissions
# AgentCollection.Read.All - Read all agent collections
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_collection_read" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "e65ee1da-d1d5-467b-bdd0-3e9bb94e6e0c" # AgentCollection.Read.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Instance permissions
# AgentInstance.Read.All - Read all agent instances
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_instance_read" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "799a4732-85b8-4c67-b048-75f0e88a232b" # AgentInstance.Read.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Instance ReadWrite permissions
# AgentInstance.ReadWrite.All - Read and write all agent instances
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_instance_readwrite" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "07abdd95-78dc-4353-bd32-09f880ea43d0" # AgentInstance.ReadWrite.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Variable for the target service principal
variable "client_service_principal_object_id" {
  description = "The Object ID of the service principal to grant agent identity permissions to"
  type        = string
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_role_id` (String) The identifier (ID) for the app role which is assigned to the principal. This app role must be exposed in the `appRoles` property on the resource application's service principal (`resource_object_id`). If the resource application has not declared any app roles, a default app role ID of `00000000-0000-0000-0000-000000000000` can be specified to signal that the principal is assigned to the resource app without any specific app roles.
- `resource_object_id` (String) The Object ID of the service principal that exposes the app roles (permissions). This is the API whose permissions you are granting. For Microsoft 365 permissions, this is typically the Microsoft Graph service principal (appId: 00000003-0000-0000-c000-000000000000). Other examples include SharePoint Online, Exchange Online, or your own custom APIs.
- `target_service_principal_object_id` (String) The Object ID of the service principal being granted the app role. This is the enterprise app (service principal) that will receive the permission.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The time when the app role assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. Read-only.
- `id` (String) The unique identifier for the app role assignment.
- `principal_display_name` (String) The display name of the user, group, or service principal that was granted the app role assignment. Read-only.
- `principal_type` (String) The type of the assigned principal. This can be either `User`, `Group`, or `ServicePrincipal`. Read-only.
- `resource_display_name` (String) The display name of the resource app's service principal to which the assignment is made. Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
# Import an existing app role assignment
# The import ID is the app role assignment ID returned by Microsoft Graph

# {app_role_assignment_id} - The unique identifier of the app role assignment
terraform import microsoft365_graph_beta_applications_service_principal_app_role_assigned_to.example {app_role_assignment_id}
```

