---
page_title: "microsoft365_graph_beta_identity_and_access_administrative_unit Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages an administrative unit in Microsoft Entra ID using the /directory/administrativeUnits endpoint. Administrative units provide a conceptual container for user, group, and device directory objects, allowing delegation of administrative responsibilities.
---

# microsoft365_graph_beta_identity_and_access_administrative_unit (Resource)

Manages an administrative unit in Microsoft Entra ID using the `/directory/administrativeUnits` endpoint. Administrative units provide a conceptual container for user, group, and device directory objects, allowing delegation of administrative responsibilities.

## Microsoft Documentation

- [Administrative units documentation](https://learn.microsoft.com/en-us/entra/identity/role-based-access-control/administrative-units)
- [administrativeUnit resource type](https://learn.microsoft.com/en-us/graph/api/resources/administrativeunit?view=graph-rest-beta)
- [Create administrativeUnit](https://learn.microsoft.com/en-us/graph/api/directory-post-administrativeunits?view=graph-rest-beta)
- [Get administrativeUnit](https://learn.microsoft.com/en-us/graph/api/administrativeunit-get?view=graph-rest-beta)
- [Update administrativeUnit](https://learn.microsoft.com/en-us/graph/api/administrativeunit-update?view=graph-rest-beta)
- [Delete administrativeUnit](https://learn.microsoft.com/en-us/graph/api/administrativeunit-delete?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `AdministrativeUnit.ReadWrite.All`
- `Directory.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.1.0-alpha | Experimental | Initial release |

## Example Usage

### AU001: Basic Administrative Unit

```terraform
# AU001: Basic Administrative Unit
# Creates a simple administrative unit with assigned membership
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au001_basic" {
  display_name = "Finance Department"
  description  = "Administrative unit for Finance department users and resources"
}
```

### AU002: Hidden Membership Administrative Unit

```terraform
# AU002: Hidden Membership Administrative Unit
# Creates an administrative unit with hidden membership where only members
# can see other members of the unit
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au002_hidden" {
  display_name = "Executive Team"
  description  = "Administrative unit for executive team with hidden membership"
  visibility   = "HiddenMembership"
}
```

### AU003: Dynamic Membership Administrative Unit

```terraform
# AU003: Dynamic Membership Administrative Unit
# Creates an administrative unit with dynamic membership based on user attributes
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au003_dynamic" {
  display_name                     = "US-Based Users"
  description                      = "Administrative unit for all users located in the United States"
  membership_type                  = "Dynamic"
  membership_rule                  = "(user.country -eq \"United States\")"
  membership_rule_processing_state = "On"
}
```

### AU004: Restricted Management Administrative Unit

```terraform
# AU004: Restricted Management Administrative Unit
# Creates an administrative unit with restricted member management
# Only administrators with specific permissions can manage members
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au004_restricted" {
  display_name                    = "Managed Devices"
  description                     = "Administrative unit for managed devices with restricted management"
  is_member_management_restricted = true
}
```

### AU005: Administrative Unit with Hard Delete

```terraform
# AU005: Administrative Unit with Hard Delete
# Creates an administrative unit that will be permanently deleted when destroyed
# instead of being moved to the deleted items container
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au005_hard_delete" {
  display_name = "Temporary Project Team"
  description  = "Administrative unit for temporary project that will be permanently deleted"
  hard_delete  = true
}
```

### AU006: Dynamic Department-Based Administrative Unit

```terraform
# AU006: Dynamic Department-Based Administrative Unit
# Creates an administrative unit that automatically includes users from a specific department
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au006_sales" {
  display_name                     = "Sales Department"
  description                      = "Administrative unit for all Sales department users"
  membership_type                  = "Dynamic"
  membership_rule                  = "(user.department -eq \"Sales\")"
  membership_rule_processing_state = "On"
  visibility                       = "Public"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Display name for the administrative unit. Maximum length is 256 characters.

### Optional

- `description` (String) An optional description for the administrative unit.
- `hard_delete` (Boolean) When `true`, the administrative unit will be permanently deleted (hard delete) during destroy. When `false` (default), the administrative unit will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. Note: This field defaults to `false` on import since the API does not return this value.
- `is_member_management_restricted` (Boolean) `true` if members of this administrative unit should be treated as sensitive, which requires specific permissions to manage. If not set, the default value is `false`. Use this property to define administrative units with roles that don't inherit from tenant-level administrators, and where the management of individual member objects is limited to administrators scoped to a restricted management administrative unit. This property is immutable and can't be changed later.
- `membership_rule` (String) The dynamic membership rule for the administrative unit. For more information about the rules you can use for dynamic administrative units and dynamic groups, see [Manage rules for dynamic membership groups in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity/users/groups-dynamic-membership).
- `membership_rule_processing_state` (String) Controls whether the dynamic membership rule is actively processed. Set to `On` to activate the dynamic membership rule, or `Paused` to stop updating membership dynamically.
- `membership_type` (String) Indicates the membership type for the administrative unit. The possible values are: `Dynamic`, `Assigned`. If not set, the default behavior is assigned.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `visibility` (String) Controls whether the administrative unit and its members are hidden or public. Can be set to `HiddenMembership` or `Public`. If not set, the default behavior is public. When set to `HiddenMembership`, only members of the administrative unit can list other members of the administrative unit. This property is immutable and can't be changed later.

### Read-Only

- `id` (String) Unique identifier for the administrative unit. Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Membership Types
- **Assigned** (default): Members are manually added to the administrative unit
- **Dynamic**: Members are automatically added based on a membership rule

### Dynamic Membership Rules
- Dynamic membership rules use the same syntax as dynamic groups
- Rules are based on user or device attributes (e.g., department, country, device type)
- See [Manage rules for dynamic membership groups](https://learn.microsoft.com/en-us/entra/identity/users/groups-dynamic-membership) for rule syntax
- The `membership_rule_processing_state` must be set to `On` to activate the rule
- Set to `Paused` to temporarily stop processing the membership rule

### Visibility
- **Public** (default): The administrative unit and its members are visible to all users
- **HiddenMembership**: Only members of the administrative unit can see other members
- This property is immutable and cannot be changed after creation

### Member Management Restriction
- When `is_member_management_restricted` is `true`, only administrators with specific permissions can manage members
- This provides additional security for sensitive administrative units
- Does not affect the ability to manage the administrative unit itself

### Deletion Behavior
- **Soft Delete** (default, `hard_delete = false`): The administrative unit is moved to the deleted items container and can be restored within 30 days
- **Hard Delete** (`hard_delete = true`): The administrative unit is permanently deleted and cannot be restored
- On import, `hard_delete` defaults to `false` since the API does not return this value

### Best Practices
- Use descriptive names that clearly identify the purpose of the administrative unit
- Consider using hidden membership for sensitive groups like executives or security teams
- Use dynamic membership to automatically maintain membership based on user attributes
- Test dynamic membership rules before setting `membership_rule_processing_state` to `On`
- Document the purpose and scope of each administrative unit
- Use `hard_delete = false` (default) for production resources to enable recovery

### Common Use Cases
- **Departmental Organization**: Create administrative units for each department to delegate user management
- **Geographic Delegation**: Use dynamic rules to group users by country or region
- **Project Teams**: Create temporary administrative units for project-based teams
- **Device Management**: Group devices for delegated device administration
- **Compliance Scoping**: Isolate users or devices that require specific compliance policies

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import with soft delete (default)
terraform import microsoft365_graph_beta_identity_and_access_administrative_unit.example 00000000-0000-0000-0000-000000000000

# Import with hard delete enabled
terraform import microsoft365_graph_beta_identity_and_access_administrative_unit.example 00000000-0000-0000-0000-000000000000:hard_delete=true
```
