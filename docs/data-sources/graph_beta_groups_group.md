---
page_title: "microsoft365_graph_beta_groups_group Data Source - terraform-provider-microsoft365"
subcategory: "Groups"

description: |-
  Retrieves information about a Microsoft Entra ID (Azure AD) group using the /groups endpoint. This data source is used to query group details by ID, display name, mail nickname, or advanced OData filtering.
---

# microsoft365_graph_beta_groups_group (Data Source)

Retrieves information about a Microsoft Entra ID (Azure AD) group using the `/groups` endpoint. This data source is used to query group details by ID, display name, mail nickname, or advanced OData filtering.

## Microsoft Documentation

- [group resource type](https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta)
- [Get group](https://learn.microsoft.com/en-us/graph/api/group-get?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `Group.Read.All`
- `Directory.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.42.0-alpha | Experimental | Initial release of groups_group data source |

## Example Usage

### Example 1: Look up group by object_id (with all possible outputs)

```terraform
# Example 1: Look up group by object_id
# This example shows all possible output attributes

data "microsoft365_graph_beta_groups_group" "by_object_id" {
  object_id = "12345678-1234-1234-1234-123456789012"
}

# All available outputs
output "group_id" {
  description = "The unique identifier for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.id
}

output "object_id" {
  description = "The object ID of the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.object_id
}

output "display_name" {
  description = "The display name for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.display_name
}

output "description" {
  description = "The optional description of the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.description
}

output "classification" {
  description = "A classification for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.classification
}

output "mail_nickname" {
  description = "The mail alias for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.mail_nickname
}

output "mail_enabled" {
  description = "Whether the group is mail-enabled"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.mail_enabled
}

output "security_enabled" {
  description = "Whether the group is a security group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.security_enabled
}

output "group_types" {
  description = "List of group types (e.g., DynamicMembership, Unified)"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.group_types
}

output "visibility" {
  description = "Group join policy and content visibility"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.visibility
}

output "assignable_to_role" {
  description = "Whether group can be assigned to Azure AD role"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.assignable_to_role
}

output "membership_rule" {
  description = "The rule for dynamic membership"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.membership_rule
}

output "membership_rule_processing_state" {
  description = "Dynamic membership processing state (On/Paused)"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.membership_rule_processing_state
}

output "created_date_time" {
  description = "When the group was created"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.created_date_time
}

output "mail" {
  description = "The SMTP address for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.mail
}

output "proxy_addresses" {
  description = "Email addresses that direct to the same mailbox"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.proxy_addresses
}

output "assigned_licenses" {
  description = "Licenses assigned to the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.assigned_licenses
}

output "has_members_with_license_errors" {
  description = "Whether members have license errors"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.has_members_with_license_errors
}

output "hide_from_address_lists" {
  description = "Whether hidden from Outlook address lists"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.hide_from_address_lists
}

output "hide_from_outlook_clients" {
  description = "Whether hidden from Outlook clients"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.hide_from_outlook_clients
}

output "onpremises_sync_enabled" {
  description = "Whether synced from on-premises directory"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_sync_enabled
}

output "onpremises_last_sync_date_time" {
  description = "Last sync time from on-premises"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_last_sync_date_time
}

output "onpremises_sam_account_name" {
  description = "On-premises SAM account name"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_sam_account_name
}

output "onpremises_domain_name" {
  description = "On-premises FQDN"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_domain_name
}

output "onpremises_netbios_name" {
  description = "On-premises NetBIOS name"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_netbios_name
}

output "onpremises_security_identifier" {
  description = "On-premises security identifier (SID)"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_security_identifier
}

output "members" {
  description = "List of member object IDs"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.members
}

output "owners" {
  description = "List of owner object IDs"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.owners
}

output "dynamic_membership_enabled" {
  description = "Whether dynamic membership is enabled"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.dynamic_membership_enabled
}
```

### Example 2: Look up group by display_name

```terraform
# Example 2: Look up group by display_name

data "microsoft365_graph_beta_groups_group" "by_display_name" {
  display_name = "My Group Name"
}

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.by_display_name.id
}

output "object_id" {
  value = data.microsoft365_graph_beta_groups_group.by_display_name.object_id
}
```

### Example 3: Look up group by mail_nickname

```terraform
# Example 3: Look up group by mail_nickname

data "microsoft365_graph_beta_groups_group" "by_mail_nickname" {
  mail_nickname = "mygroup"
}

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.by_mail_nickname.id
}

output "display_name" {
  value = data.microsoft365_graph_beta_groups_group.by_mail_nickname.display_name
}
```

### Example 4: Look up group by display_name with additional filters

```terraform
# Example 4: Look up group by display_name with additional filters
# Use mail_enabled and security_enabled as additional filters to narrow results

data "microsoft365_graph_beta_groups_group" "security_group" {
  display_name     = "My Security Group"
  security_enabled = true
  mail_enabled     = false
}

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.security_group.id
}

output "group_types" {
  value = data.microsoft365_graph_beta_groups_group.security_group.group_types
}

output "members_count" {
  description = "Number of members in the group"
  value       = length(data.microsoft365_graph_beta_groups_group.security_group.members)
}
```

### Example 5: Look up group using custom OData query

```terraform
# Example 5: Look up group using custom OData query
# Use this for advanced filtering when standard attributes don't meet your needs

data "microsoft365_graph_beta_groups_group" "by_odata_query" {
  odata_query = "displayName eq 'My Group' and securityEnabled eq true"
}

output "group_id" {
  value = data.microsoft365_graph_beta_groups_group.by_odata_query.id
}

output "display_name" {
  value = data.microsoft365_graph_beta_groups_group.by_odata_query.display_name
}

# Example: More complex OData query
data "microsoft365_graph_beta_groups_group" "dynamic_group" {
  odata_query = "startswith(displayName, 'DYN-') and securityEnabled eq true and mailEnabled eq false"
}

output "dynamic_group_id" {
  value = data.microsoft365_graph_beta_groups_group.dynamic_group.id
}
```

## Notes

- One of `object_id`, `display_name`, `mail_nickname`, or `odata_query` must be specified.
- `object_id`, `display_name`, `mail_nickname`, and `odata_query` are mutually exclusive.
- `mail_enabled` and `security_enabled` can be used as additional filters when combined with `display_name` or `mail_nickname`.
- For advanced filtering scenarios, use the `odata_query` attribute with custom OData filter expressions.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) The display name for the group.
- `mail_enabled` (Boolean) Whether the group is mail-enabled. Can be used as an additional filter when combined with other lookup attributes.
- `mail_nickname` (String) The mail alias for the group, unique in the organisation. One of `object_id`, `display_name`, or `mail_nickname` must be specified.
- `object_id` (String) The object ID of the group. One of `object_id`, `display_name`, `mail_nickname`, or `odata_query` must be specified.
- `odata_query` (String) Custom OData filter query. Use this for advanced filtering when the standard lookup attributes don't meet your needs. Cannot be combined with `object_id`, `display_name`, or `mail_nickname`. Example: `displayName eq 'My Group' and securityEnabled eq true`
- `security_enabled` (Boolean) Whether the group is a security group. Can be used as an additional filter when combined with other lookup attributes.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `assignable_to_role` (Boolean) Indicates whether this group can be assigned to an Azure AD role. Can only be set during group creation and cannot be changed afterwards.
- `assigned_licenses` (Attributes List) The licenses that are assigned to the group for group-based licensing. (see [below for nested schema](#nestedatt--assigned_licenses))
- `classification` (String) A classification for the group (such as low, medium or high business impact). Valid values are defined by creating a ClassificationList setting value in the directory.
- `created_date_time` (String) The date and time when the group was created in RFC3339 format.
- `description` (String) The optional description of the group.
- `dynamic_membership_enabled` (Boolean) Whether the group has dynamic membership enabled.
- `group_types` (Set of String) A list of group types configured for the group. Possible values include:
  - `DynamicMembership`: Denotes a group with dynamic membership
  - `Unified`: Specifies a Microsoft 365 group
- `has_members_with_license_errors` (Boolean) Indicates whether there are members in this group that have license errors from group-based license assignment. This property is never returned on a GET operation unless explicitly requested via $select.
- `hide_from_address_lists` (Boolean) True if the group is not displayed in certain parts of the Outlook UI: the Address Book, address lists for selecting message recipients, and the Browse Groups dialog for searching groups; otherwise false. Default value is false.
- `hide_from_outlook_clients` (Boolean) True if the group is not displayed in Outlook clients, such as Outlook for Windows and Outlook on the web; otherwise false. Default value is false.
- `id` (String) The unique identifier for the group.
- `mail` (String) The SMTP address for the group.
- `members` (Set of String) List of object IDs of the group members.
- `membership_rule` (String) The rule that determines members for a dynamic membership group. Only populated when `dynamic_membership_enabled` is `true`.
- `membership_rule_processing_state` (String) Indicates whether the dynamic membership is processing. Possible values are:
  - `On`: Dynamic membership is active
  - `Paused`: Dynamic membership is paused
- `onpremises_domain_name` (String) The on-premises FQDN (dnsDomainName), synchronised from the on-premises directory when Azure AD Connect is used.
- `onpremises_last_sync_date_time` (String) The last time the group was synced from the on-premises directory in RFC3339 format.
- `onpremises_netbios_name` (String) The on-premises NetBIOS name, synchronised from the on-premises directory when Azure AD Connect is used.
- `onpremises_sam_account_name` (String) The on-premises SAM account name, synchronised from the on-premises directory when Azure AD Connect is used.
- `onpremises_security_identifier` (String) The on-premises security identifier (SID), synchronised from the on-premises directory when Azure AD Connect is used.
- `onpremises_sync_enabled` (Boolean) Whether this group is synchronised from an on-premises directory. Possible values are `true` (synced), `false` (no longer synced), or null (never synced).
- `owners` (Set of String) List of object IDs of the group owners.
- `proxy_addresses` (Set of String) Email addresses for the group that direct to the same group mailbox.
- `visibility` (String) The group join policy and group content visibility. Possible values are:
  - `Private`: Only members can view content
  - `Public`: Anyone can view content
  - `Hiddenmembership`: Only members can see membership (Microsoft 365 groups only)

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--assigned_licenses"></a>
### Nested Schema for `assigned_licenses`

Read-Only:

- `disabled_plans` (Set of String) A collection of the unique identifiers for plans that have been disabled.
- `sku_id` (String) The unique identifier (GUID) for the service SKU.
