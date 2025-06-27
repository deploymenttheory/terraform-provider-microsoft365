---
page_title: "microsoft365_graph_beta_groups_group_owner_assignment Resource - terraform-provider-microsoft365"
subcategory: "Groups"
description: |-
  Manages Azure AD/Entra group owner assignments using the /groups/{group-id}/owners endpoint. This resource enables adding and removing users or service principals as owners of security groups and Microsoft 365 groups.
  Owner Type Support by Group Type:
  Security Groups: Users and Service principalsMicrosoft 365 Groups: Users and Service principalsMail-enabled Security Groups: Read-only, cannot add ownersDistribution Groups: Read-only, cannot add owners
  Important Notes:
  Owners are allowed to modify the group objectThe last owner (user object) of a group cannot be removedIf you update group owners and created a team for the group, it can take up to 2 hours for owners to sync with Microsoft TeamsIf you want the owner to make changes in a team (e.g., creating a Planner plan), the owner also needs to be added as a group/team member
  Required Permissions by Owner Type:
  Users: Group.ReadWrite.All or Directory.ReadWrite.AllService Principals: Group.ReadWrite.All or Directory.ReadWrite.AllRole-assignable Groups: Additional RoleManagement.ReadWrite.Directory permission required
  Supported Microsoft Entra Roles:
  Group owners (can modify all types of group owners)Groups Administrator (can modify all types of group owners)User Administrator (can modify user owners only)Directory Writers (can modify user owners only)Exchange Administrator (Microsoft 365 groups only)SharePoint Administrator (Microsoft 365 groups only)Teams Administrator (Microsoft 365 groups only)Yammer Administrator (Microsoft 365 groups only)Intune Administrator (security groups only)
---

# microsoft365_graph_beta_groups_group_owner_assignment (Resource)

Manages Azure AD/Entra group owner assignments using the `/groups/{group-id}/owners` endpoint. This resource enables adding and removing users or service principals as owners of security groups and Microsoft 365 groups.

**Owner Type Support by Group Type:**
- **Security Groups**: Users and Service principals
- **Microsoft 365 Groups**: Users and Service principals
- **Mail-enabled Security Groups**: Read-only, cannot add owners
- **Distribution Groups**: Read-only, cannot add owners

**Important Notes:**
- Owners are allowed to modify the group object
- The last owner (user object) of a group cannot be removed
- If you update group owners and created a team for the group, it can take up to 2 hours for owners to sync with Microsoft Teams
- If you want the owner to make changes in a team (e.g., creating a Planner plan), the owner also needs to be added as a group/team member

**Required Permissions by Owner Type:**
- **Users**: `Group.ReadWrite.All` or `Directory.ReadWrite.All`
- **Service Principals**: `Group.ReadWrite.All` or `Directory.ReadWrite.All`
- **Role-assignable Groups**: Additional `RoleManagement.ReadWrite.Directory` permission required

**Supported Microsoft Entra Roles:**
- Group owners (can modify all types of group owners)
- Groups Administrator (can modify all types of group owners)
- User Administrator (can modify user owners only)
- Directory Writers (can modify user owners only)
- Exchange Administrator (Microsoft 365 groups only)
- SharePoint Administrator (Microsoft 365 groups only)
- Teams Administrator (Microsoft 365 groups only)
- Yammer Administrator (Microsoft 365 groups only)
- Intune Administrator (security groups only)

## Microsoft Documentation

- [Add owners](https://learn.microsoft.com/en-us/graph/api/group-post-owners?view=graph-rest-beta)
- [List owners](https://learn.microsoft.com/en-us/graph/api/group-list-owners?view=graph-rest-beta)
- [Remove owner](https://learn.microsoft.com/en-us/graph/api/group-delete-owners?view=graph-rest-beta)
- [Groups overview](https://learn.microsoft.com/en-us/graph/api/resources/groups-overview?view=graph-rest-beta)
- [group resource type](https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `Group.ReadWrite.All`, `Directory.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.15.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example 1: Add a user as an owner to a security group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "user_to_security_group" {
  group_id          = "1132b215-826f-42a9-8cfe-1643d19d17fd" # Security group UUID
  owner_id          = "2243c326-937g-53f0-c9df-2e68f106b901" # User UUID
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 2: Add a user as an owner to a Microsoft 365 group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "user_to_m365_group" {
  group_id          = "3354d437-048h-64g1-d0ef-3f79g217c012" # Microsoft 365 group UUID
  owner_id          = "4465e548-159i-75h2-e1fg-4g80h328d123" # User UUID
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 3: Add a service principal as an owner to a security group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "service_principal_to_security_group" {
  group_id          = "5576f659-260j-86i3-f2gh-5i91i439e234" # Security group UUID
  owner_id          = "6687g760-371k-97j4-g3hi-6j02j540f345" # Service principal UUID
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 4: Add a service principal as an owner to a Microsoft 365 group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "service_principal_to_m365_group" {
  group_id          = "7798h871-482l-08k5-h4ij-7k13k651g456" # Microsoft 365 group UUID
  owner_id          = "8809i982-593m-19l6-i5jk-8l24l762h567" # Service principal UUID
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 5: Using data sources to get IDs dynamically
data "microsoft365_graph_beta_groups_group" "target_group" {
  display_name = "Sales Team Security Group"
}

data "microsoft365_graph_beta_user" "target_user" {
  user_principal_name = "john.doe@contoso.com"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "dynamic_assignment" {
  group_id          = data.microsoft365_graph_beta_groups_group.target_group.id
  owner_id          = data.microsoft365_graph_beta_user.target_user.id
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 6: Multiple owner assignments to the same group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "multiple_user_owners" {
  for_each = toset([
    "3354d437-048h-64g1-d0ef-3f79g217c012", # User 1
    "4465e548-159i-75h2-e1fg-4g80h328d123", # User 2
    "5576f659-260j-86i3-f2gh-5i91i439e234"  # User 3
  ])

  group_id          = "7798h871-482l-08k5-h4ij-7k13k651g456" # Target security group
  owner_id          = each.value
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 7: Mixed owner types (users and service principals) for the same group
locals {
  owners = [
    {
      id   = "9910j093-604n-20m7-j6kl-9m35m873i678"
      type = "User"
    },
    {
      id   = "0021k104-715o-31n8-k7lm-0n46n984j789"
      type = "User"
    },
    {
      id   = "1132l215-826p-42o9-l8mn-1d57o095k890"
      type = "ServicePrincipal"
    }
  ]
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "mixed_owner_types" {
  for_each = { for idx, owner in local.owners : "${owner.type}_${idx}" => owner }

  group_id          = "2243m326-937q-53p0-m9no-2e68p106l901" # Target group
  owner_id          = each.value.id
  owner_object_type = each.value.type

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 8: Creating a group and immediately adding owners
resource "microsoft365_graph_beta_groups_group" "example_group" {
  display_name     = "Example Project Team"
  mail_nickname    = "example-project-team"
  description      = "Security group for Example Project Team"
  security_enabled = true
  mail_enabled     = false
  group_types      = []
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "project_team_owners" {
  depends_on = [microsoft365_graph_beta_groups_group.example_group]

  for_each = toset([
    "1132b215-826f-42a9-8cfe-1643d19d17fd", # Project Lead
    "2243c326-937g-53f0-c9df-2e68f106b901", # Team Manager
  ])

  group_id          = microsoft365_graph_beta_groups_group.example_group.id
  owner_id          = each.value
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 9: Service principal from managed identity as group owner
data "azuread_service_principal" "managed_identity" {
  display_name = "my-app-managed-identity"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "managed_identity_owner" {
  group_id          = "3354d437-048h-64g1-d0ef-3f79g217c012" # Target group
  owner_id          = data.azuread_service_principal.managed_identity.object_id
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 10: Conditional owner assignment based on group type
data "microsoft365_graph_beta_groups_group" "conditional_group" {
  display_name = "Conditional Target Group"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "conditional_owner" {
  # Only add owner if the group is a security group
  count = contains(data.microsoft365_graph_beta_groups_group.conditional_group.group_types, "Unified") ? 0 : 1

  group_id          = data.microsoft365_graph_beta_groups_group.conditional_group.id
  owner_id          = "4465e548-159i-75h2-e1fg-4g80h328d123" # User UUID
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (String) The unique identifier (UUID) for the group.
- `owner_id` (String) The unique identifier (UUID) for the owner to be added to the group. This can be a user or service principal.
- `owner_object_type` (String) The type of object being added as an owner. This determines the correct Microsoft Graph API endpoint to use. Valid values: 'User', 'ServicePrincipal'. Both security groups and Microsoft 365 groups support both types.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for this group owner assignment. This is a composite ID formed by combining the group ID and owner ID.
- `owner_display_name` (String) The display name of the owner. Read-only.
- `owner_type` (String) The type of the owner object as returned by Microsoft Graph (e.g., 'User', 'ServicePrincipal'). Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Owner Type Support**: The types of owners that can be added depend on the target group type:
  - **Security Groups**: Users and Service principals
  - **Microsoft 365 Groups**: Users and Service principals
  - **Mail-enabled Security Groups**: Read-only, cannot add owners
  - **Distribution Groups**: Read-only, cannot add owners

- **Owner Permissions**: Group owners are allowed to modify the group object and can perform administrative tasks on the group.

- **Last Owner Protection**: The last owner (user object) of a group cannot be removed. Microsoft Graph API will prevent this operation.

- **Teams Integration**: If you update group owners for a group that has an associated Microsoft Teams team:
  - It can take up to 2 hours for the owners to be synchronized with Microsoft Teams
  - If you want the owner to be able to make changes in the team (e.g., creating a Planner plan), the owner also needs to be added as a group/team member

- **Object Type Validation**: The `owner_object_type` field determines the correct Microsoft Graph API endpoint to use:
  - `User`: Uses `/users/{id}` endpoint
  - `ServicePrincipal`: Uses `/servicePrincipals/{id}` endpoint

- **Permissions by Owner Type**: The same permissions are required for both owner types:
  - **Users**: `Group.ReadWrite.All` or `Directory.ReadWrite.All`
  - **Service Principals**: `Group.ReadWrite.All` or `Directory.ReadWrite.All`
  - **Role-assignable Groups**: Additional `RoleManagement.ReadWrite.Directory` permission required

- **Microsoft Entra Roles**: The following least privileged roles are supported for this operation:
  - **Group owners** (can modify all types of group owners)
  - **Groups Administrator** (can modify all types of group owners)
  - **User Administrator** (can modify user owners only)
  - **Directory Writers** (can modify user owners only)
  - **Exchange Administrator** (Microsoft 365 groups only)
  - **SharePoint Administrator** (Microsoft 365 groups only)
  - **Teams Administrator** (Microsoft 365 groups only)
  - **Yammer Administrator** (Microsoft 365 groups only)
  - **Intune Administrator** (security groups only)

- **API Behavior**: The resource uses the `/groups/{group-id}/owners/$ref` endpoint with `@odata.id` references
- **Validation**: Client-side validation prevents common API errors and provides clear error messages
- **Composite ID**: The resource ID is automatically generated as `{group_id}/{owner_id}`
- **Idempotent**: The Microsoft Graph API handles duplicate owner additions appropriately

## Owner Type Compatibility Matrix

| Object Type | Security Group | Microsoft 365 Group | Mail-enabled Security | Distribution |
|-------------|----------------|---------------------|----------------------|-------------|
| User | ✅ Supported | ✅ Supported | ❌ Read-only | ❌ Read-only |
| Service Principal | ✅ Supported | ✅ Supported | ❌ Read-only | ❌ Read-only |

## Import

Group owner assignments can be imported using the composite ID format: `{group_id}/{owner_id}`

```bash
terraform import microsoft365_graph_beta_groups_group_owner_assignment.example "12345678-1234-1234-1234-123456789012/87654321-4321-4321-4321-210987654321"
```

Where:
- `12345678-1234-1234-1234-123456789012` is the group ID
- `87654321-4321-4321-4321-210987654321` is the owner ID 