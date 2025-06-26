---
page_title: "microsoft365_graph_beta_groups_group_member_assignment Resource - terraform-provider-microsoft365"
subcategory: "Groups"
description: |-
  Manages Azure AD/Entra group member assignments using the /groups/{group-id}/members endpoint. This resource enables adding and removing users, groups, service principals, devices, and organizational contacts as members of security groups and Microsoft 365 groups.
  Member Type Support by Group Type:
  Security Groups: Users, other Security groups, Devices, Service principals, and Organizational contactsMicrosoft 365 Groups: Only Users are supportedMail-enabled Security Groups: Read-only, cannot add membersDistribution Groups: Read-only, cannot add members
  Important Notes:
  The resource automatically validates member compatibility with the target group typeWhen adding a Group as a member, both the target and member groups must be Security groupsMicrosoft 365 groups cannot be members of any group type
  Required Permissions by Member Type:
  Users: GroupMember.ReadWrite.AllGroups: GroupMember.ReadWrite.AllDevices: GroupMember.ReadWrite.All + Device.ReadWrite.AllService Principals: GroupMember.ReadWrite.All + Application.ReadWrite.AllOrganizational Contacts: GroupMember.ReadWrite.All + OrgContact.Read.AllRole-assignable Groups: Additional RoleManagement.ReadWrite.Directory permission required
---

# microsoft365_graph_beta_groups_group_member_assignment (Resource)

Manages Azure AD/Entra group member assignments using the `/groups/{group-id}/members` endpoint. This resource enables adding and removing users, groups, service principals, devices, and organizational contacts as members of security groups and Microsoft 365 groups.

**Member Type Support by Group Type:**
- **Security Groups**: Users, other Security groups, Devices, Service principals, and Organizational contacts
- **Microsoft 365 Groups**: Only Users are supported
- **Mail-enabled Security Groups**: Read-only, cannot add members
- **Distribution Groups**: Read-only, cannot add members

**Important Notes:**
- The resource automatically validates member compatibility with the target group type
- When adding a Group as a member, both the target and member groups must be Security groups
- Microsoft 365 groups cannot be members of any group type

**Required Permissions by Member Type:**
- **Users**: `GroupMember.ReadWrite.All`
- **Groups**: `GroupMember.ReadWrite.All`
- **Devices**: `GroupMember.ReadWrite.All` + `Device.ReadWrite.All`
- **Service Principals**: `GroupMember.ReadWrite.All` + `Application.ReadWrite.All`
- **Organizational Contacts**: `GroupMember.ReadWrite.All` + `OrgContact.Read.All`
- **Role-assignable Groups**: Additional `RoleManagement.ReadWrite.Directory` permission required

## Microsoft Documentation

- [Add member](https://learn.microsoft.com/en-us/graph/api/group-post-members?view=graph-rest-beta)
- [List members](https://learn.microsoft.com/en-us/graph/api/group-list-members?view=graph-rest-beta)
- [Remove member](https://learn.microsoft.com/en-us/graph/api/group-delete-members?view=graph-rest-beta)
- [Group membership overview](https://learn.microsoft.com/en-us/graph/api/resources/groups-overview?view=graph-rest-beta#group-membership)
- [group resource type](https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `GroupMember.ReadWrite.All`, `Directory.ReadWrite.All`, `Device.ReadWrite.All`, `Application.ReadWrite.All`, `OrgContact.Read.All`, `RoleManagement.ReadWrite.Directory`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.15.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example 1: Add a user to a security group
resource "microsoft365_graph_beta_groups_group_member_assignment" "user_to_security_group" {
  group_id           = "1132b215-826f-42a9-8cfe-1643d19d17fd" # Security group UUID
  member_id          = "2243c326-937g-53f0-c9df-2e68f106b901" # User UUID
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 2: Add a user to a Microsoft 365 group
resource "microsoft365_graph_beta_groups_group_member_assignment" "user_to_m365_group" {
  group_id           = "3354d437-048h-64g1-d0ef-3f79g217c012" # Microsoft 365 group UUID
  member_id          = "4465e548-159i-75h2-e1fg-4g80h328d123" # User UUID
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 3: Add a security group to another security group (nested groups)
resource "microsoft365_graph_beta_groups_group_member_assignment" "group_to_security_group" {
  group_id           = "5576f659-260j-86i3-f2gh-5i91i439e234" # Parent security group UUID
  member_id          = "6687g760-371k-97j4-g3hi-6j02j540f345" # Member security group UUID
  member_object_type = "Group"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 4: Add a device to a security group
resource "microsoft365_graph_beta_groups_group_member_assignment" "device_to_security_group" {
  group_id           = "7798h871-482l-08k5-h4ij-7k13k651g456" # Security group UUID
  member_id          = "8809i982-593m-19l6-i5jk-8l24l762h567" # Device UUID
  member_object_type = "Device"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 5: Add a service principal to a security group
resource "microsoft365_graph_beta_groups_group_member_assignment" "service_principal_to_security_group" {
  group_id           = "9910j093-604n-20m7-j6kl-9m35m873i678" # Security group UUID
  member_id          = "0021k104-715o-31n8-k7lm-0n46n984j789" # Service principal UUID
  member_object_type = "ServicePrincipal"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 6: Add an organizational contact to a security group
resource "microsoft365_graph_beta_groups_group_member_assignment" "contact_to_security_group" {
  group_id           = "1132l215-826p-42o9-l8mn-1d57o095k890" # Security group UUID
  member_id          = "2243m326-937q-53p0-m9no-2e68p106l901" # Organizational contact UUID
  member_object_type = "OrganizationalContact"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 7: Using data sources to get IDs dynamically
data "microsoft365_graph_beta_groups_group" "target_group" {
  display_name = "Sales Team Security Group"
}

data "microsoft365_graph_beta_user" "target_user" {
  user_principal_name = "john.doe@contoso.com"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "dynamic_assignment" {
  group_id           = data.microsoft365_graph_beta_groups_group.target_group.id
  member_id          = data.microsoft365_graph_beta_user.target_user.id
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 8: Multiple member assignments to the same group
resource "microsoft365_graph_beta_groups_group_member_assignment" "multiple_users" {
  for_each = toset([
    "3354d437-048h-64g1-d0ef-3f79g217c012", # User 1
    "4465e548-159i-75h2-e1fg-4g80h328d123", # User 2
    "5576f659-260j-86i3-f2gh-5i91i439e234"  # User 3
  ])

  group_id           = "7798h871-482l-08k5-h4ij-7k13k651g456" # Target security group
  member_id          = each.value
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 9: Creating a group and immediately adding members
resource "microsoft365_graph_beta_groups_group" "example_group" {
  display_name     = "Example Project Team"
  mail_nickname    = "example-project-team"
  description      = "Security group for Example Project Team members"
  security_enabled = true
  mail_enabled     = false
  group_types      = []
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "project_team_members" {
  depends_on = [microsoft365_graph_beta_groups_group.example_group]

  for_each = toset([
    "1132b215-826f-42a9-8cfe-1643d19d17fd", # Project Manager
    "2243c326-937g-53f0-c9df-2e68f106b901", # Developer 1
    "3354d437-048h-64g1-d0ef-3f79g217c012"  # Developer 2
  ])

  group_id           = microsoft365_graph_beta_groups_group.example_group.id
  member_id          = each.value
  member_object_type = "User"

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
- `member_id` (String) The unique identifier (UUID) for the member to be added to the group. This can be a user, group, device, service principal, or organizational contact.
- `member_object_type` (String) The type of object being added as a member. This determines the correct Microsoft Graph API endpoint to use. Valid values: 'User', 'Group', 'Device', 'ServicePrincipal', 'OrganizationalContact'. Note: Microsoft 365 groups only support 'User' and 'Group' (where Group must be a security group), while security groups support all types.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for this group member assignment. This is a composite ID formed by combining the group ID and member ID.
- `member_display_name` (String) The display name of the member. Read-only.
- `member_type` (String) The type of the member object as returned by Microsoft Graph (e.g., 'User', 'Group', 'Device', 'ServicePrincipal', 'OrganizationalContact'). Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Member Type Support**: The types of members that can be added depend on the target group type:
  - **Security Groups**: Users, other Security groups, Devices, Service principals, and Organizational contacts
  - **Microsoft 365 Groups**: Only Users are supported
  - **Mail-enabled Security Groups**: Read-only, cannot add members
  - **Distribution Groups**: Read-only, cannot add members

- **Group-to-Group Membership**: When adding a Group as a member:
  - Only Security groups can be added to Security groups
  - Microsoft 365 groups cannot be members of any group type
  - The resource automatically validates group type compatibility

- **Object Type Validation**: The `member_object_type` field determines the correct Microsoft Graph API endpoint to use:
  - `User`: Uses `/directoryObjects/{id}` endpoint
  - `Group`: Uses `/groups/{id}` endpoint  
  - `Device`: Uses `/devices/{id}` endpoint
  - `ServicePrincipal`: Uses `/servicePrincipals/{id}` endpoint
  - `OrganizationalContact`: Uses `/contacts/{id}` endpoint

- **Permissions by Member Type**: Different permissions are required depending on the member type:
  - **Users**: `GroupMember.ReadWrite.All`
  - **Groups**: `GroupMember.ReadWrite.All`
  - **Devices**: `GroupMember.ReadWrite.All` + `Device.ReadWrite.All`
  - **Service Principals**: `GroupMember.ReadWrite.All` + `Application.ReadWrite.All`
  - **Organizational Contacts**: `GroupMember.ReadWrite.All` + `OrgContact.Read.All`
  - **Role-assignable Groups**: Additional `RoleManagement.ReadWrite.Directory` permission required

- **API Behavior**: The resource uses the `/groups/{group-id}/members/$ref` endpoint with `@odata.id` references
- **Validation**: Client-side validation prevents common API errors and provides clear error messages
- **Composite ID**: The resource ID is automatically generated as `{group_id}/{member_id}`
- **Idempotent**: The Microsoft Graph API handles duplicate member additions appropriately

## Member Type Compatibility Matrix

| Object Type | Security Group | Microsoft 365 Group | Mail-enabled Security | Distribution |
|-------------|----------------|---------------------|----------------------|-------------|
| User | ✅ Supported | ✅ Supported | ❌ Read-only | ❌ Read-only |
| Security Group | ✅ Supported | ❌ Not allowed | ❌ Read-only | ❌ Read-only |
| Microsoft 365 Group | ❌ Not allowed | ❌ Not allowed | ❌ Read-only | ❌ Read-only |
| Device | ✅ Supported | ❌ Not allowed | ❌ Read-only | ❌ Read-only |
| Service Principal | ✅ Supported | ❌ Not allowed | ❌ Read-only | ❌ Read-only |
| Organizational Contact | ✅ Supported | ❌ Not allowed | ❌ Read-only | ❌ Read-only |

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using composite ID format: {group_id}/{member_id}
terraform import microsoft365_graph_beta_groups_group_member_assignment.example "1132b215-826f-42a9-8cfe-1643d19d17fd/2243c326-937g-53f0-c9df-2e68f106b901"
``` 