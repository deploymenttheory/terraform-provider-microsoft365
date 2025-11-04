---
page_title: "microsoft365_graph_beta_groups_group Resource - terraform-provider-microsoft365"
subcategory: "Groups"
description: |-
  Manages Azure AD/Entra groups using the /groups endpoint. This resource enables creation and management of security groups, Microsoft 365 groups, and distribution groups with support for dynamic membership, role assignment capabilities, and comprehensive group configuration options for organizational identity and access management.
---

# microsoft365_graph_beta_groups_group (Resource)

Manages Azure AD/Entra groups using the `/groups` endpoint. This resource enables creation and management of security groups, Microsoft 365 groups, and distribution groups with support for dynamic membership, role assignment capabilities, and comprehensive group configuration options for organizational identity and access management.

## Microsoft Documentation

- [group resource type](https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta)
- [Create group](https://learn.microsoft.com/en-us/graph/api/group-post-groups?view=graph-rest-beta)
- [Update group](https://learn.microsoft.com/en-us/graph/api/group-update?view=graph-rest-beta)
- [Delete group](https://learn.microsoft.com/en-us/graph/api/group-delete?view=graph-rest-beta)

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
# Example 1: Basic Security Group with Assigned Membership
# Creates a standard security group where members are manually assigned.
# This is the most common type of security group used for access control.
resource "microsoft365_graph_beta_groups_group" "security_basic" {
  display_name     = "Engineering Team"
  mail_nickname    = "engineering-team"
  mail_enabled     = false
  security_enabled = true
  description      = "Security group for engineering team members"
}

# Example 2: Security Group with Dynamic User Membership
# Creates a security group that automatically adds/removes users based on a membership rule.
# Useful for automatically managing group membership based on user attributes.
resource "microsoft365_graph_beta_groups_group" "security_dynamic_users" {
  display_name                     = "Active Employees"
  mail_nickname                    = "active-employees"
  mail_enabled                     = false
  security_enabled                 = true
  description                      = "Security group containing all active employees"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}

# Example 3: Security Group with Dynamic Device Membership
# Creates a security group that automatically includes devices based on a membership rule.
# Ideal for device management scenarios like Conditional Access or Intune policies.
resource "microsoft365_graph_beta_groups_group" "security_dynamic_devices" {
  display_name                     = "Corporate Managed Devices"
  mail_nickname                    = "corporate-devices"
  mail_enabled                     = false
  security_enabled                 = true
  description                      = "Security group containing all corporate managed devices"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(device.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}

# Example 4: Role-Assignable Security Group
# Creates a security group that can be assigned to Entra ID roles.
# Note: Requires elevated permissions and visibility must be "Private".
# Once created, is_assignable_to_role cannot be changed.
resource "microsoft365_graph_beta_groups_group" "security_role_assignable" {
  display_name          = "Privileged Access Administrators"
  mail_nickname         = "privileged-admins"
  mail_enabled          = false
  security_enabled      = true
  description           = "Security group for privileged access administration"
  is_assignable_to_role = true
  visibility            = "Private"
}

# Example 5: Microsoft 365 Group with Dynamic User Membership
# Creates a Microsoft 365 group (formerly Office 365 group) with automatic membership.
# Includes Teams, SharePoint, Outlook, and other Microsoft 365 services.
resource "microsoft365_graph_beta_groups_group" "m365_dynamic_users" {
  display_name                     = "Marketing Department"
  mail_nickname                    = "marketing-dept"
  mail_enabled                     = true
  security_enabled                 = true
  group_types                      = ["Unified", "DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
  visibility                       = "Private"
}

# Example 6: Microsoft 365 Group with Role Assignment
# Creates a Microsoft 365 group that can be assigned to Entra ID roles.
# Combines collaboration features with privileged access management.
# Note: Requires elevated permissions and visibility must be "Private".
resource "microsoft365_graph_beta_groups_group" "m365_role_assignable" {
  display_name          = "Executive Leadership Team"
  mail_nickname         = "executive-team"
  mail_enabled          = true
  security_enabled      = true
  group_types           = ["Unified"]
  description           = "Microsoft 365 group for executive leadership"
  is_assignable_to_role = true
  visibility            = "Private"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name for the group. This property is required when a group is created and can't be cleared during updates. Maximum length is 256 characters.
- `mail_enabled` (Boolean) Specifies whether the group is mail-enabled. Required.
- `mail_nickname` (String) The mail alias for the group, unique for Microsoft 365 groups in the organization. Maximum length is 64 characters. This property can contain only characters in the ASCII character set 0 - 127 except the following: @ () \ [] " ; : <> , SPACE.
- `security_enabled` (Boolean) Specifies whether the group is a security group. Required.

### Optional

- `description` (String) An optional description for the group.
- `group_members` (Set of String) The members of the group at creation time. A maximum of 20 relationships, such as owners and members, can be added as part of group creation. Specify the user IDs (GUIDs) of the users who should be members of the group. Additional members can be added after creation using the `/groups/{id}/members/$ref` endpoint or JSON batching.
- `group_owners` (Set of String) The owners of the group at creation time. A maximum of 20 relationships, such as owners and members, can be added as part of group creation. Specify the user IDs (GUIDs) of the users who should be owners of the group. Note: A non-admin user cannot add themselves to the group owners collection. Owners can be added after creation using the `/groups/{id}/owners/$ref` endpoint.
- `group_types` (Set of String) Specifies the group type and its membership. If the collection contains 'Unified', the group is a Microsoft 365 group; otherwise, it's either a security group or a distribution group. If the collection includes 'DynamicMembership', the group has dynamic membership; otherwise, membership is static.
- `is_assignable_to_role` (Boolean) Indicates whether this group can be assigned to a Microsoft Entra role. This property can only be set while creating the group and is immutable. If set to true, the securityEnabled property must also be set to true, visibility must be Hidden, and the group can't be a dynamic group. Default is false.
- `membership_rule` (String) The rule that determines members for this group if the group is a dynamic group (groupTypes contains DynamicMembership). For more information about the syntax of the membership rule, see Membership Rules syntax.
- `membership_rule_processing_state` (String) Indicates whether the dynamic membership processing is on or paused. Possible values are 'On' or 'Paused'. Only applicable for dynamic groups (when groupTypes contains DynamicMembership).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `visibility` (String) Specifies the group join policy and group content visibility for groups. Possible values are: `Private`, `Public`, or `HiddenMembership`. `HiddenMembership` can be set only for Microsoft 365 groups when the groups are created and cannot be updated later. Other values of visibility can be updated after group creation. If visibility value is not specified during group creation, a security group is created as `Private` by default, and a Microsoft 365 group is `Public`. Groups assignable to roles are always `Private`. Returned by default. Nullable.

### Read-Only

- `created_date_time` (String) Timestamp of when the group was created. The value can't be modified and is automatically populated when the group is created. Read-only.
- `id` (String) The unique identifier for the group. Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Group Types**: This resource supports security groups, Microsoft 365 groups, and distribution groups.
- **Dynamic Membership**: Groups can have dynamic membership based on user or device attributes when `group_types` includes "DynamicMembership".
- **Role Assignment**: Groups can be made assignable to Azure AD roles by setting `is_assignable_to_role` to true (only during creation).
- **Mail Features**: Microsoft 365 groups automatically get mail functionality when `mail_enabled` is true and `group_types` includes "Unified".
- **Visibility**: Controls who can see and join the group - Private, Public, or HiddenMembership.
- **Character Restrictions**: The `mail_nickname` field has strict character restrictions (ASCII only, excluding special characters).
- **Length Limits**: Display names are limited to 256 characters, mail nicknames to 64 characters.
- **Language Codes**: Preferred language should follow ISO 639-1 format (e.g., "en-US").
- **Theme Colors**: Available themes are Teal, Purple, Green, Blue, Pink, Orange, or Red.
- **Immutable Properties**: Some properties like `is_assignable_to_role` can only be set during creation.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# {group_id}
terraform import microsoft365_graph_beta_group.example 00000000-0000-0000-0000-000000000000
``` 