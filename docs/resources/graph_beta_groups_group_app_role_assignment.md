---
page_title: "microsoft365_graph_beta_groups_group_app_role_assignment Resource - terraform-provider-microsoft365"
subcategory: "Groups"
description: |-
  Manages Azure AD/Entra group app role assignments using the /groups/{group-id}/appRoleAssignments endpoint. This resource is used to enables assigning app roles to security groups, allowing all direct members of the group to inherit the assigned permissions. Security groups with dynamic memberships are supported.
  Important Notes:
  All direct members of the assigned group will be considered as having the app roleAdditional licenses might be required to use a group to manage access to applicationsThe resource requires three key identifiers: principal ID (group), resource ID (service principal), and app role ID
  Required Permissions:
  AppRoleAssignment.ReadWrite.All + Group.Read.All (least privileged)Delegated scenarios: The signed-in user must be assigned one of the supported Microsoft Entra roles (Directory Readers, Directory Writers, Application Administrator, Cloud Application Administrator, etc.).
---

# microsoft365_graph_beta_groups_group_app_role_assignment (Resource)

Manages Azure AD/Entra group app role assignments using the `/groups/{group-id}/appRoleAssignments` endpoint. This resource is used to enables assigning app roles to security groups, allowing all direct members of the group to inherit the assigned permissions. Security groups with dynamic memberships are supported.

**Important Notes:**
- All direct members of the assigned group will be considered as having the app role
- Additional licenses might be required to use a group to manage access to applications
- The resource requires three key identifiers: principal ID (group), resource ID (service principal), and app role ID

**Required Permissions:**
- `AppRoleAssignment.ReadWrite.All` + `Group.Read.All` (least privileged)
- Delegated scenarios: The signed-in user must be assigned one of the supported Microsoft Entra roles (Directory Readers, Directory Writers, Application Administrator, Cloud Application Administrator, etc.).

## Microsoft Documentation

- [appRoleAssignment resource type](https://learn.microsoft.com/en-us/graph/api/resources/approleassignment?view=graph-rest-beta)
- [List appRoleAssignments granted to a group](https://learn.microsoft.com/en-us/graph/api/group-list-approleassignments?view=graph-rest-beta)
- [Grant an appRoleAssignment to a group](https://learn.microsoft.com/en-us/graph/api/group-post-approleassignments?view=graph-rest-beta)
- [Delete an appRoleAssignment from a group](https://learn.microsoft.com/en-us/graph/api/group-delete-approleassignments?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `AppRoleAssignment.ReadWrite.All`
- `Directory.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.39.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# To find service principals in your tenant, use:
# Get-MgServicePrincipal -Filter "appId eq '00000003-0000-0000-c000-000000000000'" | Select-Object Id, DisplayName, AppId
# The "Id" property is what you need for resource_id

# Get the Microsoft Graph service principal (resource that defines the permissions)
data "microsoft365_graph_beta_applications_service_principal" "msgraph" {
  filter_type  = "display_name"
  filter_value = "Microsoft Graph"
}

# Example 1: Assign default access to Microsoft Graph
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "graph_default" {
  target_group_id    = "12345678-1234-1234-1234-123456789012"                                          # UUID of your group
  resource_object_id = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id # Object ID of Microsoft Graph service principal
  app_role_id        = "00000000-0000-0000-0000-000000000000"                                          # Default role (basic access)

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Example 2: Assign specific app role to a group
# Common Microsoft Graph App Roles:
# - Directory.Read.All: df021288-bdef-4463-88db-98f22de89214
# - User.Read.All: a154be20-db9c-4678-8ab7-66f6cc099a59
# - Group.Read.All: 5b567255-7703-4780-807c-7be8301ae99b
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "graph_directory_read" {
  target_group_id    = "12345678-1234-1234-1234-123456789012"                                          # UUID of your group
  resource_object_id = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id # Object ID of Microsoft Graph service principal
  app_role_id        = "df021288-bdef-4463-88db-98f22de89214"                                          # Directory.Read.All
}

# Example 3: Assign role to SharePoint Online
# Get the SharePoint Online service principal
# The App ID for SharePoint Online is always: 00000003-0000-0ff1-ce00-000000000000
data "microsoft365_graph_beta_applications_service_principal" "sharepoint" {
  filter_type  = "display_name"
  filter_value = "SharePoint Online"
}

resource "microsoft365_graph_beta_groups_group_app_role_assignment" "sharepoint" {
  target_group_id    = "12345678-1234-1234-1234-123456789012"                                             # UUID of your group
  resource_object_id = data.microsoft365_graph_beta_applications_service_principal.sharepoint.items[0].id # Object ID of SharePoint service principal
  app_role_id        = "678536fe-1083-478a-9c59-b99265e6b0d3"                                             # Example SharePoint app role
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_role_id` (String) The unique identifier (UUID) for the app role defined on the resource service principal to assign to the group. Use '00000000-0000-0000-0000-000000000000' for the default access role.
- `resource_object_id` (String) The unique identifier (UUID) for the resource service principal that has defined the app role. This is the service principal ID of the application.
- `target_group_id` (String) The unique identifier (UUID) for the group to which you are assigning the app role. This is the principal ID.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `creation_timestamp` (String) The date and time the app role assignment was created. Read-only.
- `id` (String) The unique identifier for this app role assignment.
- `principal_display_name` (String) The display name of the group (principal). Read-only.
- `principal_type` (String) The type of the principal. For groups, this will always be 'Group'. Read-only.
- `resource_display_name` (String) The display name of the service principal (resource/application). Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **App Role Assignments**: Used to assign application roles to groups, granting members of the group the permissions defined by that role.
- **Service Principal**: The `resource_object_id` references the service principal (enterprise application) object ID that exposes the app roles. This is the Object ID (not the App ID) of the service principal.
- **Default Role**: Use `00000000-0000-0000-0000-000000000000` as the `app_role_id` to assign the default access role.
- **Target Group**: The `target_group_id` identifies the group receiving the app role assignment.
- **Read-Only Fields**: The `principal_display_name`, `resource_display_name`, and `principal_type` fields are computed and returned by the API.
- **Finding Service Principals**: Use the `microsoft365_graph_beta_applications_service_principal` data source to find service principal Object IDs by App ID or display name.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import scripts for Microsoft 365 Group App Role Assignment
# {group_id}/{assignment_id}

# Import a group app role assignment
terraform import microsoft365_graph_beta_groups_group_app_role_assignment.example {group_id}/{assignment_id}

# Import a read-only app role assignment
terraform import microsoft365_graph_beta_groups_group_app_role_assignment.read_only {group_id}/{assignment_id}

# Import a full-access app role assignment
terraform import microsoft365_graph_beta_groups_group_app_role_assignment.full_access {group_id}/{assignment_id}

# Note: You can import individual app role assignments as needed.
# To find assignment IDs, you can use Microsoft Graph API or PowerShell:
# Get-MgGroupAppRoleAssignment -GroupId {group_id}
```

