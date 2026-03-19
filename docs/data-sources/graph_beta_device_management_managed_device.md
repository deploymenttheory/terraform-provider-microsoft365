---
page_title: "microsoft365_graph_beta_device_management_managed_device Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"
description: |-
  Retrieves Managed Devices from Microsoft Intune using the /deviceManagement/managedDevices endpoint. Supports flexible lookup by device ID, device name, operating system, OS version, Azure AD device ID, serial number, user principal name, or custom OData queries.
---

# microsoft365_graph_beta_device_management_managed_device (Data Source)

Retrieves managed devices from Microsoft Intune using the `/deviceManagement/managedDevices` endpoint. This data source provides flexible lookup capabilities including list all, device ID, device name, operating system, OS version, Azure AD device ID, serial number, user display name, and custom OData queries.

This data source allows you to list and filter managed devices in your tenant, providing comprehensive details such as device name, operating system, compliance state, user information, hardware details, and more.

## Microsoft Documentation

- [List managedDevices](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-list?view=graph-rest-beta)
- [Get managedDevice](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-get?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementManagedDevices.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.18.0-alpha | Experimental | Initial release |
| v0.19.0-alpha | Experimental | Added flexible lookup patterns (device_id, device_name, operating_system, os_version, azure_ad_device_id, serial_number, user_display_name, odata_query, list_all) |
| v0.50.0-alpha | Preview | **Breaking Change**: Refactored to flexible lookup pattern following groups data source design. Removed `filter_type` and `filter_value` attributes. Added mutually exclusive convenience lookup attributes: `device_id`, `device_name`, `operating_system`, `os_version`, `azure_ad_device_id`, `serial_number`, `user_principal_name`, and `list_all`. Consolidated all individual OData parameters (`odata_filter`, `odata_top`, `odata_orderby`, `odata_select`, `odata_count`, `odata_search`, `odata_expand`) into single `odata_query` attribute for custom filter expressions. Added operating system validator to ensure only valid device types are used. |

## Example Usage

### Example 1: List All Managed Devices

```terraform
# Example: List all managed devices
data "microsoft365_graph_beta_device_management_managed_device" "all" {
  list_all = true
}

# Output: Total number of managed devices
output "all_managed_devices_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.all.items)
  description = "Total number of managed devices in the tenant"
}

# Output: Device names
output "all_device_names" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.all.items :
    device.device_name
  ]
  description = "List of all device names"
}
```

### Example 2: Get Device by ID

```terraform
# Example: Get a specific managed device by ID
data "microsoft365_graph_beta_device_management_managed_device" "by_id" {
  device_id = "00000000-0000-0000-0000-000000000001"
}

# Output: Device information by ID
output "device_by_id_info" {
  value = length(data.microsoft365_graph_beta_device_management_managed_device.by_id.items) > 0 ? {
    id            = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].id
    name          = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].device_name
    os            = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].operating_system
    os_version    = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].os_version
    enrolled      = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].enrolled_date_time
    compliance    = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].compliance_state
    serial_number = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].serial_number
    manufacturer  = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].manufacturer
    model         = data.microsoft365_graph_beta_device_management_managed_device.by_id.items[0].model
  } : null
  description = "Detailed device information for the specified device ID"
}
```

### Example 3: Get Devices by Device Name

```terraform
# Example: Get managed devices by device name
data "microsoft365_graph_beta_device_management_managed_device" "by_name" {
  device_name = "DESKTOP-WIN-001"
}

# Output: Devices matching the name filter
output "devices_by_name" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.by_name.items :
    {
      id          = device.id
      device_name = device.device_name
      os          = device.operating_system
      user        = device.user_principal_name
    }
  ]
  description = "List of devices matching the specified device name"
}
```

### Example 4: Get Devices by Operating System

```terraform
# Example: Get managed devices by operating system
data "microsoft365_graph_beta_device_management_managed_device" "windows" {
  operating_system = "Windows"
}

# Output: Windows devices
output "windows_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.windows.items :
    {
      id               = device.id
      device_name      = device.device_name
      operating_system = device.operating_system
      os_version       = device.os_version
      compliance_state = device.compliance_state
      last_sync        = device.last_sync_date_time
    }
  ]
  description = "List of Windows devices with key information"
}

# Output: Count of Windows devices
output "windows_devices_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.windows.items)
  description = "Total number of Windows devices"
}
```

### Example 5: Get Devices by OS and Version

```terraform
# Example: Get managed devices by operating system and version
data "microsoft365_graph_beta_device_management_managed_device" "windows_10" {
  operating_system = "Windows"
  os_version       = "10.0.19045"
}

# Output: Windows 10 devices with specific version
output "windows_10_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.windows_10.items :
    {
      id          = device.id
      device_name = device.device_name
      os_version  = device.os_version
      compliance  = device.compliance_state
      encrypted   = device.is_encrypted
    }
  ]
  description = "List of Windows 10 devices with the specified version"
}
```

### Example 6: Get Device by Azure AD Device ID

```terraform
# Example: Get managed device by Azure AD Device ID
data "microsoft365_graph_beta_device_management_managed_device" "by_aad_id" {
  azure_ad_device_id = "aaaaaaaa-0000-0000-0000-000000000001"
}

# Output: Device information by Azure AD Device ID
output "device_by_aad_id" {
  value = length(data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items) > 0 ? {
    intune_id          = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].id
    azure_ad_device_id = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].azure_ad_device_id
    device_name        = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].device_name
    registration_state = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].device_registration_state
    user               = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].user_principal_name
  } : null
  description = "Device information using Azure AD Device ID lookup"
}
```

### Example 7: Get Devices by Serial Number

```terraform
# Example: Get managed devices by serial number
data "microsoft365_graph_beta_device_management_managed_device" "by_serial" {
  serial_number = "SN-WIN-001"
}

# Output: Devices with matching serial number
output "devices_by_serial" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.by_serial.items :
    {
      id            = device.id
      device_name   = device.device_name
      serial_number = device.serial_number
      manufacturer  = device.manufacturer
      model         = device.model
      user          = device.user_display_name
    }
  ]
  description = "List of devices matching the specified serial number"
}

# Note: Multiple devices may share the same serial number in some scenarios
output "serial_number_device_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.by_serial.items)
  description = "Number of devices with this serial number"
}
```

### Example 8: Get Devices by User Principal Name

```terraform
# Example: Get managed devices by user principal name (email)
data "microsoft365_graph_beta_device_management_managed_device" "by_user" {
  user_principal_name = "user@contoso.com"
}

# Output: Devices assigned to the specified user
output "user_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.by_user.items :
    {
      id                = device.id
      device_name       = device.device_name
      user_display_name = device.user_display_name
      user_email        = device.user_principal_name
      operating_system  = device.operating_system
      compliance_state  = device.compliance_state
      last_sync         = device.last_sync_date_time
    }
  ]
  description = "List of devices assigned to the specified user"
}
```

### Example 9: OData Query - Simple Filter

```terraform
# Example: Get managed devices using a simple OData filter
data "microsoft365_graph_beta_device_management_managed_device" "compliant_devices" {
  odata_query = "complianceState eq 'compliant'"
}

# Output: Compliant devices
output "compliant_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.compliant_devices.items :
    {
      id               = device.id
      device_name      = device.device_name
      compliance_state = device.compliance_state
      operating_system = device.operating_system
    }
  ]
  description = "List of compliant devices"
}

# Output: Compliant device count
output "compliant_devices_count" {
  value       = length(data.microsoft365_graph_beta_device_management_managed_device.compliant_devices.items)
  description = "Total number of compliant devices"
}
```

### Example 10: OData Query - AND Filter

```terraform
# Example: Get managed devices using OData filter with AND logic
data "microsoft365_graph_beta_device_management_managed_device" "compliant_windows" {
  odata_query = "operatingSystem eq 'Windows' and complianceState eq 'compliant'"
}

# Output: Compliant Windows devices
output "compliant_windows_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.compliant_windows.items :
    {
      id               = device.id
      device_name      = device.device_name
      operating_system = device.operating_system
      os_version       = device.os_version
      compliance_state = device.compliance_state
      is_encrypted     = device.is_encrypted
    }
  ]
  description = "List of compliant Windows devices"
}
```

### Example 11: OData Query - StartsWith Filter

```terraform
# Example: Get managed devices using OData startswith filter
data "microsoft365_graph_beta_device_management_managed_device" "desktop_devices" {
  odata_query = "startswith(deviceName,'DESKTOP')"
}

# Output: Desktop devices
output "desktop_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.desktop_devices.items :
    {
      id          = device.id
      device_name = device.device_name
      os          = device.operating_system
      user        = device.user_principal_name
    }
  ]
  description = "List of devices with names starting with 'DESKTOP'"
}
```

### Example 12: OData Query - OR Filter

```terraform
# Example: Get managed devices using OData filter with OR logic
data "microsoft365_graph_beta_device_management_managed_device" "ios_or_android" {
  odata_query = "operatingSystem eq 'iOS' or operatingSystem eq 'Android'"
}

# Output: iOS and Android devices
output "mobile_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.ios_or_android.items :
    {
      id               = device.id
      device_name      = device.device_name
      operating_system = device.operating_system
      os_version       = device.os_version
      model            = device.model
      manufacturer     = device.manufacturer
    }
  ]
  description = "List of iOS and Android mobile devices"
}

# Output: Device count by OS
output "mobile_device_summary" {
  value = {
    total = length(data.microsoft365_graph_beta_device_management_managed_device.ios_or_android.items)
    ios = length([
      for device in data.microsoft365_graph_beta_device_management_managed_device.ios_or_android.items :
      device if device.operating_system == "iOS"
    ])
    android = length([
      for device in data.microsoft365_graph_beta_device_management_managed_device.ios_or_android.items :
      device if device.operating_system == "Android"
    ])
  }
  description = "Summary of mobile devices by operating system"
}
```

## Lookup Methods

This data source supports multiple lookup methods. **Exactly one** of the following attributes must be specified:

- `list_all` - List all managed devices in the tenant
- `device_id` - Lookup by specific Intune device ID
- `device_name` - Filter by device name (exact match via OData filter)
- `operating_system` - Filter by operating system (e.g., "Windows", "iOS", "Android")
- `os_version` - Filter by OS version (must be used with `operating_system`)
- `azure_ad_device_id` - Lookup by Azure AD device ID
- `serial_number` - Filter by device serial number
- `user_principal_name` - Filter by user principal name (email address)
- `odata_query` - Custom OData filter expression for advanced queries

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `azure_ad_device_id` (String) The Azure AD device ID (Entra ID device object ID). Conflicts with other lookup attributes.
- `device_id` (String) The unique identifier of the managed device (Intune device ID). Conflicts with other lookup attributes.
- `device_name` (String) The name of the managed device. Conflicts with other lookup attributes.
- `list_all` (Boolean) Retrieve all managed devices. Conflicts with specific lookup attributes.
- `odata_query` (String) Custom OData filter expression for advanced filtering. Example: `operatingSystem eq 'Windows' and osVersion startswith '10'`. Conflicts with other lookup attributes.
- `operating_system` (String) Filter devices by operating system. Valid values: 'Windows', 'iOS', 'Android', 'macOS', 'androidForWork', 'androidEnterprise', 'chromeOS', 'linux', 'visionOS', 'tvOS', 'unknown', 'cloudPC'. Can be combined with os_version. Conflicts with device_id, list_all, and odata_query.
- `os_version` (String) Filter devices by OS version. Can be combined with operating_system. Conflicts with device_id, list_all, and odata_query.
- `serial_number` (String) The serial number of the managed device. Conflicts with other lookup attributes.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `user_principal_name` (String) The user principal name (email) of the user associated with the managed device. Filters devices by userPrincipalName. Conflicts with other lookup attributes.

### Read-Only

- `items` (Attributes List) The list of managed devices that match the filter criteria. (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `aad_registered` (Boolean) Whether the device is Azure AD registered.
- `activation_lock_bypass_code` (String) Activation lock bypass code for the device.
- `android_security_patch_level` (String) Android security patch level on the device.
- `autopilot_enrolled` (Boolean) Whether the device is enrolled in Autopilot.
- `azure_active_directory_device_id` (String) Azure Active Directory device ID.
- `azure_ad_device_id` (String) Azure AD device ID (legacy field).
- `azure_ad_registered` (Boolean) Whether the device is Azure AD registered (legacy field).
- `bootstrap_token_escrowed` (Boolean) Whether the bootstrap token is escrowed for the device.
- `chassis_type` (String) Chassis type of the device (e.g., desktop, laptop).
- `chrome_os_device_info` (Attributes List) List of Chrome OS device information properties. (see [below for nested schema](#nestedatt--items--chrome_os_device_info))
- `compliance_grace_period_expiration_date_time` (String) Expiration date and time for the compliance grace period.
- `compliance_state` (String) Compliance state of the device (e.g., compliant, noncompliant).
- `configuration_manager_client_enabled_features` (Attributes) Configuration Manager client enabled features. (see [below for nested schema](#nestedatt--items--configuration_manager_client_enabled_features))
- `configuration_manager_client_health_state` (Attributes) Configuration Manager client health state. (see [below for nested schema](#nestedatt--items--configuration_manager_client_health_state))
- `configuration_manager_client_information` (Attributes) Configuration Manager client information. (see [below for nested schema](#nestedatt--items--configuration_manager_client_information))
- `device_action_results` (Attributes List) List of device action results for the device. (see [below for nested schema](#nestedatt--items--device_action_results))
- `device_category_display_name` (String) Display name of the device category.
- `device_enrollment_type` (String) Type of device enrollment (e.g., userEnrollment, deviceEnrollmentManager).
- `device_firmware_configuration_interface_managed` (Boolean) Whether the device firmware configuration interface is managed.
- `device_health_attestation_state` (Attributes) Device health attestation state. (see [below for nested schema](#nestedatt--items--device_health_attestation_state))
- `device_identity_attestation_detail` (Attributes) Device identity attestation detail. (see [below for nested schema](#nestedatt--items--device_identity_attestation_detail))
- `device_name` (String) The name of the device as displayed in Intune.
- `device_registration_state` (String) Registration state of the device.
- `device_type` (String) Type of the device (e.g., windowsRT, windows).
- `eas_activated` (Boolean) Whether Exchange ActiveSync is activated on the device.
- `eas_activation_date_time` (String) Date and time when Exchange ActiveSync was activated.
- `eas_device_id` (String) Exchange ActiveSync device ID.
- `email_address` (String) Email address associated with the device.
- `enrolled_date_time` (String) Date and time when the device was enrolled.
- `enrollment_profile_name` (String) Enrollment profile name for the device.
- `ethernet_mac_address` (String) Ethernet MAC address of the device.
- `exchange_access_state` (String) Exchange access state for the device.
- `exchange_access_state_reason` (String) Reason for the Exchange access state.
- `exchange_last_successful_sync_date_time` (String) Last successful Exchange sync date and time.
- `free_storage_space_in_bytes` (Number) Free storage space in bytes.
- `hardware_information` (Attributes) Hardware information for the device. (see [below for nested schema](#nestedatt--items--hardware_information))
- `iccid` (String) Integrated Circuit Card Identifier (ICCID) for the device's SIM card.
- `id` (String) The unique identifier for the managed device.
- `imei` (String) International Mobile Equipment Identity (IMEI) of the device.
- `is_encrypted` (Boolean) Whether the device storage is encrypted.
- `is_supervised` (Boolean) Whether the device is supervised (Apple devices only).
- `jail_broken` (String) Indicates if the device is jailbroken (for iOS devices).
- `join_type` (String) Join type of the device (e.g., azureADJoined).
- `last_sync_date_time` (String) Last time the device synced with Intune.
- `lost_mode_state` (String) State of lost mode on the device (e.g., enabled, disabled).
- `managed_device_name` (String) Managed device name.
- `managed_device_owner_type` (String) Managed device owner type (e.g., company, personal).
- `management_agent` (String) Management agent used for the device (e.g., mdm, eas).
- `management_certificate_expiration_date` (String) Expiration date of the management certificate.
- `management_features` (String) Management features enabled on the device.
- `management_state` (String) Management state of the device (e.g., retirePending, managed).
- `manufacturer` (String) Device manufacturer.
- `meid` (String) Mobile Equipment Identifier (MEID) of the device.
- `model` (String) Device model.
- `notes` (String) Notes associated with the device.
- `operating_system` (String) Operating system of the device.
- `os_version` (String) Operating system version of the device.
- `owner_type` (String) Owner type of the device (e.g., company, personal).
- `partner_reported_threat_state` (String) Partner reported threat state.
- `phone_number` (String) Phone number associated with the device.
- `physical_memory_in_bytes` (Number) Physical memory in bytes on the device.
- `prefer_mdm_over_group_policy_applied_date_time` (String) Date and time when MDM was preferred over group policy.
- `processor_architecture` (String) Processor architecture of the device (e.g., x86, x64).
- `remote_assistance_session_error_details` (String) Error details for the remote assistance session.
- `remote_assistance_session_url` (String) URL for the remote assistance session.
- `require_user_enrollment_approval` (Boolean) Whether user enrollment approval is required.
- `retire_after_date_time` (String) Date and time after which the device will be retired.
- `role_scope_tag_ids` (List of String) List of role scope tag IDs assigned to the device.
- `security_patch_level` (String) Security patch level of the device.
- `serial_number` (String) Device serial number.
- `sku_family` (String) SKU family of the device.
- `sku_number` (Number) SKU number of the device.
- `specification_version` (String) Specification version of the device.
- `subscriber_carrier` (String) Mobile carrier for the device's SIM card.
- `total_storage_space_in_bytes` (Number) Total storage space in bytes.
- `udid` (String) Unique Device Identifier (UDID) for the device.
- `user_display_name` (String) Display name of the user associated with the device.
- `user_id` (String) The unique identifier for the user associated with the device.
- `user_principal_name` (String) User principal name associated with the device.
- `users_logged_on` (Attributes List) List of users currently logged on to the device. (see [below for nested schema](#nestedatt--items--users_logged_on))
- `wi_fi_mac_address` (String) Wi-Fi MAC address of the device.
- `windows_active_malware_count` (Number) Count of active malware instances on the device.
- `windows_remediated_malware_count` (Number) Count of remediated malware instances on the device.

<a id="nestedatt--items--chrome_os_device_info"></a>
### Nested Schema for `items.chrome_os_device_info`

Read-Only:

- `name` (String) Name of the Chrome OS device property.
- `updatable` (Boolean) Whether the Chrome OS device property is updatable.
- `value` (String) Value of the Chrome OS device property.
- `value_type` (String) Type of the value for the Chrome OS device property.


<a id="nestedatt--items--configuration_manager_client_enabled_features"></a>
### Nested Schema for `items.configuration_manager_client_enabled_features`

Read-Only:

- `compliance_policy` (Boolean) Whether compliance policy is enabled.
- `device_configuration` (Boolean) Whether device configuration is enabled.
- `endpoint_protection` (Boolean) Whether endpoint protection is enabled.
- `inventory` (Boolean) Whether inventory is enabled.
- `modern_apps` (Boolean) Whether modern apps are enabled.
- `office_apps` (Boolean) Whether Office apps are enabled.
- `resource_access` (Boolean) Whether resource access is enabled.
- `windows_update_for_business` (Boolean) Whether Windows Update for Business is enabled.


<a id="nestedatt--items--configuration_manager_client_health_state"></a>
### Nested Schema for `items.configuration_manager_client_health_state`

Read-Only:

- `error_code` (Number) Error code for the Configuration Manager client health state.
- `last_sync_date_time` (String) Last sync date and time for the Configuration Manager client.
- `state` (String) Health state of the Configuration Manager client.


<a id="nestedatt--items--configuration_manager_client_information"></a>
### Nested Schema for `items.configuration_manager_client_information`

Read-Only:

- `client_identifier` (String) Client identifier for the Configuration Manager client.
- `client_version` (String) Version of the Configuration Manager client.
- `is_blocked` (Boolean) Whether the Configuration Manager client is blocked.


<a id="nestedatt--items--device_action_results"></a>
### Nested Schema for `items.device_action_results`

Read-Only:

- `action_name` (String) Name of the action performed on the device.
- `action_state` (String) State of the action (e.g., pending, completed).
- `last_updated_date_time` (String) Last update time of the action.
- `start_date_time` (String) Start time of the action.


<a id="nestedatt--items--device_health_attestation_state"></a>
### Nested Schema for `items.device_health_attestation_state`

Read-Only:

- `attestation_identity_key` (String) Attestation identity key.
- `bit_locker_status` (String) BitLocker status for device health attestation.
- `boot_app_security_version` (String) Boot app security version for device health attestation.
- `boot_debugging` (String) Boot debugging status for device health attestation.
- `boot_manager_security_version` (String) Boot manager security version for device health attestation.
- `boot_manager_version` (String) Boot manager version for device health attestation.
- `boot_revision_list_info` (String) Boot revision list info for device health attestation.
- `code_integrity` (String) Code integrity status for device health attestation.
- `code_integrity_check_version` (String) Code integrity check version for device health attestation.
- `code_integrity_policy` (String) Code integrity policy for device health attestation.
- `content_namespace_url` (String) Content namespace URL for device health attestation.
- `content_version` (String) Content version for device health attestation.
- `data_excution_policy` (String) Data execution policy for device health attestation.
- `device_health_attestation_status` (String) Device health attestation status.
- `early_launch_anti_malware_driver_protection` (String) Early launch anti-malware driver protection status.
- `firmware_protection` (String) Firmware protection status.
- `health_attestation_supported_status` (String) Health attestation supported status.
- `health_status_mismatch_info` (String) Health status mismatch info for device health attestation.
- `issued_date_time` (String) Issued date and time for device health attestation.
- `last_update_date_time` (String) Last update date and time for device health attestation.
- `memory_access_protection` (String) Memory access protection status.
- `memory_integrity_protection` (String) Memory integrity protection status.
- `operating_system_kernel_debugging` (String) Operating system kernel debugging status for device health attestation.
- `operating_system_rev_list_info` (String) Operating system revision list info for device health attestation.
- `pcr0` (String) PCR0 value for device health attestation.
- `pcr_hash_algorithm` (String) PCR hash algorithm for device health attestation.
- `reset_count` (Number) Reset count for device health attestation.
- `restart_count` (Number) Restart count for device health attestation.
- `safe_mode` (String) Safe mode status for device health attestation.
- `secure_boot` (String) Secure boot status for device health attestation.
- `secure_boot_configuration_policy_finger_print` (String) Secure boot configuration policy fingerprint.
- `secured_core_pc` (String) Secured core PC status.
- `system_management_mode` (String) System management mode status.
- `test_signing` (String) Test signing status for device health attestation.
- `tpm_version` (String) TPM version for device health attestation.
- `virtual_secure_mode` (String) Virtual secure mode status for device health attestation.
- `virtualization_based_security` (String) Virtualization based security status.
- `windows_pe` (String) Windows PE status for device health attestation.


<a id="nestedatt--items--device_identity_attestation_detail"></a>
### Nested Schema for `items.device_identity_attestation_detail`

Read-Only:

- `device_identity_attestation_status` (String) Device identity attestation status.


<a id="nestedatt--items--hardware_information"></a>
### Nested Schema for `items.hardware_information`

Read-Only:

- `battery_charge_cycles` (Number) Number of battery charge cycles.
- `battery_health_percentage` (Number) Battery health as a percentage.
- `battery_level_percentage` (Number) Battery level as a percentage.
- `battery_serial_number` (String) Serial number of the device's battery.
- `cellular_technology` (String) Cellular technology used by the device (e.g., LTE, 5G).
- `device_full_qualified_domain_name` (String) Fully qualified domain name of the device.
- `device_guard_local_system_authority_credential_guard_state` (String) Device Guard LSA Credential Guard state.
- `device_guard_virtualization_based_security_hardware_requirement_state` (String) Device Guard VBS hardware requirement state.
- `device_guard_virtualization_based_security_state` (String) Device Guard VBS state.
- `device_licensing_last_error_code` (Number) Last error code for device licensing.
- `device_licensing_last_error_description` (String) Last error description for device licensing.
- `device_licensing_status` (String) Device licensing status.
- `esim_identifier` (String) eSIM identifier for the device.
- `free_storage_space` (Number) Free storage space on the device in bytes.
- `imei` (String) International Mobile Equipment Identity (IMEI) of the device.
- `ip_address_v4` (String) IPv4 address of the device.
- `is_encrypted` (Boolean) Whether the device storage is encrypted.
- `is_shared_device` (Boolean) Whether the device is a shared device.
- `is_supervised` (Boolean) Whether the device is supervised (Apple devices only).
- `manufacturer` (String) Device manufacturer.
- `meid` (String) Mobile Equipment Identifier (MEID) of the device.
- `model` (String) Device model.
- `operating_system_edition` (String) Edition of the device's operating system.
- `operating_system_language` (String) Language of the device's operating system.
- `operating_system_product_type` (Number) Product type of the operating system.
- `os_build_number` (String) Operating system build number.
- `phone_number` (String) Phone number associated with the device.
- `product_name` (String) Product name of the device.
- `resident_users_count` (Number) Number of resident users on the device.
- `serial_number` (String) Device serial number.
- `shared_device_cached_users` (Attributes List) List of users cached on a shared device. (see [below for nested schema](#nestedatt--items--hardware_information--shared_device_cached_users))
- `subnet_address` (String) Subnet address of the device.
- `subscriber_carrier` (String) Mobile carrier for the device's SIM card.
- `system_management_bios_version` (String) System Management BIOS version.
- `total_storage_space` (Number) Total storage space on the device in bytes.
- `tpm_manufacturer` (String) TPM manufacturer.
- `tpm_specification_version` (String) TPM specification version.
- `tpm_version` (String) TPM version.
- `wifi_mac` (String) Wi-Fi MAC address of the device.
- `wired_ipv4_addresses` (List of String) List of wired IPv4 addresses for the device.

<a id="nestedatt--items--hardware_information--shared_device_cached_users"></a>
### Nested Schema for `items.hardware_information.shared_device_cached_users`

Read-Only:

- `data_quota` (Number) Data quota for the user in MB.
- `data_to_sync` (Boolean) Whether there is data to sync for the user.
- `data_used` (Number) Data used by the user in MB.
- `user_principal_name` (String) User principal name of the cached user.



<a id="nestedatt--items--users_logged_on"></a>
### Nested Schema for `items.users_logged_on`

Read-Only:

- `last_log_on_date_time` (String) Last logon date and time for the user.
- `user_id` (String) User ID of the logged on user. 