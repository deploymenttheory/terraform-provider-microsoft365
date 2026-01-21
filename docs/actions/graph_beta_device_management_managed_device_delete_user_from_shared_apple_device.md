---
page_title: "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Deletes a user and their cached data from Shared iPad devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/deleteUserFromSharedAppleDevice and /deviceManagement/comanagedDevices/{managedDeviceId}/deleteUserFromSharedAppleDevice endpoints. This action is used to permanently remove specified user accounts and associated cached data from Shared iPads, freeing up storage space. The action permanently removes the specified user's account and all associated cached data from the Shared iPad.
  What This Action Does:
  Permanently deletes user from Shared iPad device rosterRemoves all cached user data (documents, photos, app data)Frees up device storage spacePrevents user from logging back into that specific deviceDoes not affect user's account or data in the cloudCannot be undone (user must be re-added if needed)
  Platform Support:
  iPadOS: Full support (Shared iPad mode only)iOS: Not supported (iPhones don't support Shared mode)Other platforms: Not supported
  Reference: Microsoft Graph API - Delete User From Shared Apple Device https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device (Action)

Deletes a user and their cached data from Shared iPad devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/deleteUserFromSharedAppleDevice` and `/deviceManagement/comanagedDevices/{managedDeviceId}/deleteUserFromSharedAppleDevice` endpoints. This action is used to permanently remove specified user accounts and associated cached data from Shared iPads, freeing up storage space. The action permanently removes the specified user's account and all associated cached data from the Shared iPad.

**What This Action Does:**
- Permanently deletes user from Shared iPad device roster
- Removes all cached user data (documents, photos, app data)
- Frees up device storage space
- Prevents user from logging back into that specific device
- Does not affect user's account or data in the cloud
- Cannot be undone (user must be re-added if needed)

**Platform Support:**
- **iPadOS**: Full support (Shared iPad mode only)
- **iOS**: Not supported (iPhones don't support Shared mode)
- **Other platforms**: Not supported

**Reference:** [Microsoft Graph API - Delete User From Shared Apple Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [deleteUserFromSharedAppleDevice action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Device remove user](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-remove-user)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this action:

**Required:**
- `DeviceManagementManagedDevices.PrivilegedOperations.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


## Example Usage

```terraform
# Example 1: Delete user from a single shared Apple device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "delete_single" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user@example.com"
      }
    ]
  }
}

# Example 2: Delete users from multiple shared Apple devices
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "delete_multiple" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user1@example.com"
      },
      {
        device_id           = "87654321-4321-4321-4321-ba9876543210"
        user_principal_name = "user2@example.com"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Delete with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "delete_maximal" {
  config {
    managed_devices = [
      {
        device_id           = "12345678-1234-1234-1234-123456789abc"
        user_principal_name = "user1@example.com"
      }
    ]

    comanaged_devices = [
      {
        device_id           = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        user_principal_name = "user2@example.com"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Delete departing user from all shared iPads
data "microsoft365_graph_beta_device_management_managed_device" "shared_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (managementMode eq 'shared')"
}

action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "delete_departing_user" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.shared_ipads.items : {
        device_id           = device.id
        user_principal_name = "departing.user@example.com"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed device-user pairs. Co-managed devices are managed by both Intune and Configuration Manager (SCCM). (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_devices` (Attributes List) List of managed device-user pairs. Managed devices are fully managed by Intune only. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Shared iPad devices before attempting user deletion. Disabling this can speed up planning but may result in runtime errors for non-existent or non-Shared iPad devices. Default: `true`.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The co-managed device ID (GUID) of the Shared iPad.
- `user_principal_name` (String) The user principal name (UPN) to delete from the device. Example: user@domain.com


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The managed device ID (GUID) of the Shared iPad.
- `user_principal_name` (String) The user principal name (UPN) to delete from the device. Example: user@domain.com


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


