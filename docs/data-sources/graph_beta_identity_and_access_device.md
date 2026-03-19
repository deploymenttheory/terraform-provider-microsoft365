---
page_title: "Data Source: microsoft365_graph_beta_identity_and_access_device"
subcategory: "Identity and Access"
description: |-
  Retrieves Microsoft Entra Devices using the /devices endpoint. Supports flexible lookup by object ID, display name, device ID, or custom OData queries. Can also retrieve device memberships, registered owners, and registered users.
---

# Data Source: microsoft365_graph_beta_identity_and_access_device

Retrieves Microsoft Entra Devices using the `/devices` endpoint. Supports flexible lookup by object ID, display name, device ID, or custom OData queries. Can also retrieve device memberships, registered owners, and registered users.

## API Documentation

- [Device Resource Type](https://learn.microsoft.com/en-us/graph/api/resources/device?view=graph-rest-beta)
- [List Devices](https://learn.microsoft.com/en-us/graph/api/device-list?view=graph-rest-beta)
- [Get Device](https://learn.microsoft.com/en-us/graph/api/device-get?view=graph-rest-beta)
- [List Device MemberOf](https://learn.microsoft.com/en-us/graph/api/device-list-memberof?view=graph-rest-beta)
- [List Registered Owners](https://learn.microsoft.com/en-us/graph/api/device-list-registeredowners?view=graph-rest-beta)
- [List Registered Users](https://learn.microsoft.com/en-us/graph/api/device-list-registeredusers?view=graph-rest-beta)

## Permissions

The following API permissions are required to use this data source:

- `Device.Read.All` (Application or Delegated)
- `Directory.Read.All` (Application or Delegated)

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release with flexible lookup patterns |

## Example Usage

### Example 1: List All Devices

```terraform
# Example: List all devices in the tenant

data "microsoft365_graph_beta_identity_and_access_device" "all" {
  list_all = true
}

# Output the first device's display name
output "first_device_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.all.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.all.items[0].display_name : "No devices found"
}

# Output the total count of devices
output "device_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.all.items)
}
```

### Example 2: Get Device by Object ID

```terraform
# Example: Get a device by its object ID

data "microsoft365_graph_beta_identity_and_access_device" "by_id" {
  object_id = "00000000-0000-0000-0000-000000000000"
}

output "device_display_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_id.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_id.items[0].display_name : null
}

output "device_operating_system" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_id.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_id.items[0].operating_system : null
}
```

### Example 3: Get Devices by Display Name

```terraform
# Example: Get devices by display name

data "microsoft365_graph_beta_identity_and_access_device" "by_name" {
  display_name = "DESKTOP-ABC123"
}

output "device_id" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_name.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_name.items[0].device_id : null
}

output "device_trust_type" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_name.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_name.items[0].trust_type : null
}
```

### Example 4: Get Devices by Device ID

```terraform
# Example: Get devices by device ID (Azure Device Registration Service ID)

data "microsoft365_graph_beta_identity_and_access_device" "by_device_id" {
  device_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}

output "device_display_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_device_id.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_device_id.items[0].display_name : null
}

output "device_is_compliant" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.by_device_id.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.by_device_id.items[0].is_compliant : null
}
```

### Example 5: Get Devices with OData Query

```terraform
# Example: Get devices using an OData query

# Filter for Windows devices that are enabled
data "microsoft365_graph_beta_identity_and_access_device" "windows_enabled" {
  odata_query = "operatingSystem eq 'Windows' and accountEnabled eq true"
}

output "windows_device_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.windows_enabled.items)
}

# Filter for compliant devices
data "microsoft365_graph_beta_identity_and_access_device" "compliant" {
  odata_query = "isCompliant eq true"
}

output "compliant_device_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.compliant.items)
}
```

### Example 6: Get Device with Group Memberships

```terraform
# Example: Get a device and its group memberships

data "microsoft365_graph_beta_identity_and_access_device" "with_groups" {
  object_id      = "00000000-0000-0000-0000-000000000000"
  list_member_of = true
}

output "device_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_groups.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.with_groups.items[0].display_name : null
}

output "group_memberships" {
  value = [for group in data.microsoft365_graph_beta_identity_and_access_device.with_groups.member_of : {
    id           = group.id
    display_name = group.display_name
    type         = group.odata_type
  }]
}

output "group_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_groups.member_of)
}
```

### Example 7: Get Device with Registered Owners

```terraform
# Example: Get a device and its registered owners

data "microsoft365_graph_beta_identity_and_access_device" "with_owners" {
  object_id              = "00000000-0000-0000-0000-000000000000"
  list_registered_owners = true
}

output "device_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_owners.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.with_owners.items[0].display_name : null
}

output "registered_owners" {
  value = [for owner in data.microsoft365_graph_beta_identity_and_access_device.with_owners.registered_owners : {
    id           = owner.id
    display_name = owner.display_name
    type         = owner.odata_type
  }]
}

output "owner_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_owners.registered_owners)
}
```

### Example 8: Get Device with Registered Users

```terraform
# Example: Get a device and its registered users

data "microsoft365_graph_beta_identity_and_access_device" "with_users" {
  object_id             = "00000000-0000-0000-0000-000000000000"
  list_registered_users = true
}

output "device_name" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_users.items) > 0 ? data.microsoft365_graph_beta_identity_and_access_device.with_users.items[0].display_name : null
}

output "registered_users" {
  value = [for user in data.microsoft365_graph_beta_identity_and_access_device.with_users.registered_users : {
    id           = user.id
    display_name = user.display_name
    type         = user.odata_type
  }]
}

output "user_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.with_users.registered_users)
}
```

### Example 9: Get Device with All Related Information

```terraform
# Example: Get a device with all related information

data "microsoft365_graph_beta_identity_and_access_device" "comprehensive" {
  object_id              = "00000000-0000-0000-0000-000000000000"
  list_member_of         = true
  list_registered_owners = true
  list_registered_users  = true
}

# Device information
output "device_info" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items) > 0 ? {
    id                       = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].id
    display_name             = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].display_name
    device_id                = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].device_id
    operating_system         = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].operating_system
    operating_system_version = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].operating_system_version
    is_compliant             = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].is_compliant
    is_managed               = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].is_managed
    trust_type               = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].trust_type
    account_enabled          = data.microsoft365_graph_beta_identity_and_access_device.comprehensive.items[0].account_enabled
  } : null
}

# Group memberships
output "group_memberships" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.comprehensive.member_of)
}

# Registered owners
output "registered_owners" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.comprehensive.registered_owners)
}

# Registered users
output "registered_users" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.comprehensive.registered_users)
}
```

## Lookup Methods

This data source supports multiple lookup methods. **Exactly one** of the following attributes must be specified:

- `list_all` - List all devices in the tenant
- `object_id` - Lookup by specific Microsoft Entra device object ID
- `display_name` - Filter by device display name (exact match via OData filter)
- `device_id` - Filter by device ID (Azure Device Registration Service ID)
- `odata_query` - Custom OData filter expression for advanced queries

### Additional Lookup Options

When using `object_id`, you can optionally enable the following to retrieve related information:

- `list_member_of` - Retrieve groups and administrative units the device is a member of
- `list_registered_owners` - Retrieve the registered owners of the device
- `list_registered_users` - Retrieve the registered users of the device

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `device_id` (String) The unique device identifier set by Azure Device Registration Service. Conflicts with other lookup attributes.
- `display_name` (String) The display name of the device. Conflicts with other lookup attributes.
- `list_all` (Boolean) Retrieve all devices in the tenant. Conflicts with specific lookup attributes.
- `list_member_of` (Boolean) When true and combined with object_id, retrieves the groups and administrative units that the device is a direct member of. Requires object_id to be specified.
- `list_registered_owners` (Boolean) When true and combined with object_id, retrieves the registered owners of the device. Requires object_id to be specified.
- `list_registered_users` (Boolean) When true and combined with object_id, retrieves the registered users of the device. Requires object_id to be specified.
- `object_id` (String) The unique object identifier of the device in Microsoft Entra ID. Conflicts with other lookup attributes.
- `odata_query` (String) Custom OData filter expression for advanced queries (e.g., `operatingSystem eq 'Windows' and accountEnabled eq true`). Conflicts with specific lookup attributes.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for the data source. This is a placeholder attribute required by Terraform.
- `items` (Attributes List) List of devices matching the query criteria. (see [below for nested schema](#nestedatt--items))
- `member_of` (Attributes List) Groups and administrative units that the device is a direct member of. Only populated when list_member_of is true. (see [below for nested schema](#nestedatt--member_of))
- `registered_owners` (Attributes List) The registered owners of the device. Only populated when list_registered_owners is true. (see [below for nested schema](#nestedatt--registered_owners))
- `registered_users` (Attributes List) The registered users of the device. Only populated when list_registered_users is true. (see [below for nested schema](#nestedatt--registered_users))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `account_enabled` (Boolean) true if the account is enabled; otherwise, false.
- `alternative_security_ids` (Attributes List) Alternative security identifiers for the device. (see [below for nested schema](#nestedatt--items--alternative_security_ids))
- `approximate_last_sign_in_date_time` (String) The timestamp of the last sign-in activity.
- `compliance_expiration_date_time` (String) The timestamp when the device is no longer deemed compliant.
- `device_category` (String) User-defined property set by Intune to automatically add devices to groups.
- `device_id` (String) Unique identifier set by Azure Device Registration Service.
- `device_metadata` (String) Metadata for the device.
- `device_ownership` (String) Ownership of the device (unknown, company, personal).
- `device_version` (Number) Version of the device.
- `display_name` (String) The display name for the device.
- `domain_name` (String) The on-premises domain name of the device.
- `enrollment_profile_name` (String) Enrollment profile applied to the device.
- `enrollment_type` (String) Enrollment type of the device.
- `extension_attributes` (Attributes) Extension attributes 1-15 for the device. (see [below for nested schema](#nestedatt--items--extension_attributes))
- `id` (String) The unique identifier for the device object.
- `is_compliant` (Boolean) true if the device complies with Mobile Device Management (MDM) policies.
- `is_managed` (Boolean) true if the device is managed by a Mobile Device Management (MDM) app.
- `is_management_restricted` (Boolean) Indicates whether the device is a member of a restricted management administrative unit.
- `is_rooted` (Boolean) true if the device is rooted or jail-broken.
- `management_type` (String) The management channel of the device.
- `manufacturer` (String) Manufacturer of the device.
- `mdm_app_id` (String) Application identifier used to register device into MDM.
- `model` (String) Model of the device.
- `on_premises_last_sync_date_time` (String) The last time the object was synced with the on-premises directory.
- `on_premises_security_identifier` (String) The on-premises security identifier (SID) for the user who was synchronized from on-premises.
- `on_premises_sync_enabled` (Boolean) true if this object is synced from an on-premises directory.
- `operating_system` (String) The type of operating system on the device.
- `operating_system_version` (String) The version of the operating system on the device.
- `physical_ids` (List of String) Physical identifiers for the device.
- `profile_type` (String) The profile type of the device (RegisteredDevice, SecureVM, Printer, Shared, IoT).
- `registration_date_time` (String) Date and time when the device was registered.
- `system_labels` (List of String) List of labels applied to the device by the system.
- `trust_type` (String) Type of trust for the joined device (Workplace, AzureAd, ServerAd).

<a id="nestedatt--items--alternative_security_ids"></a>
### Nested Schema for `items.alternative_security_ids`

Read-Only:

- `identity_provider` (String) The identity provider for the alternative security identifier.
- `key` (String) The key value of the alternative security identifier.
- `type` (Number) The type of the alternative security identifier.


<a id="nestedatt--items--extension_attributes"></a>
### Nested Schema for `items.extension_attributes`

Read-Only:

- `extension_attribute1` (String) Extension attribute 1.
- `extension_attribute10` (String) Extension attribute 10.
- `extension_attribute11` (String) Extension attribute 11.
- `extension_attribute12` (String) Extension attribute 12.
- `extension_attribute13` (String) Extension attribute 13.
- `extension_attribute14` (String) Extension attribute 14.
- `extension_attribute15` (String) Extension attribute 15.
- `extension_attribute2` (String) Extension attribute 2.
- `extension_attribute3` (String) Extension attribute 3.
- `extension_attribute4` (String) Extension attribute 4.
- `extension_attribute5` (String) Extension attribute 5.
- `extension_attribute6` (String) Extension attribute 6.
- `extension_attribute7` (String) Extension attribute 7.
- `extension_attribute8` (String) Extension attribute 8.
- `extension_attribute9` (String) Extension attribute 9.



<a id="nestedatt--member_of"></a>
### Nested Schema for `member_of`

Read-Only:

- `display_name` (String) The display name of the directory object.
- `id` (String) The unique identifier of the directory object.
- `odata_type` (String) The OData type of the directory object (e.g., #microsoft.graph.group).


<a id="nestedatt--registered_owners"></a>
### Nested Schema for `registered_owners`

Read-Only:

- `display_name` (String) The display name of the directory object.
- `id` (String) The unique identifier of the directory object.
- `odata_type` (String) The OData type of the directory object (e.g., #microsoft.graph.user).


<a id="nestedatt--registered_users"></a>
### Nested Schema for `registered_users`

Read-Only:

- `display_name` (String) The display name of the directory object.
- `id` (String) The unique identifier of the directory object.
- `odata_type` (String) The OData type of the directory object (e.g., #microsoft.graph.user).
