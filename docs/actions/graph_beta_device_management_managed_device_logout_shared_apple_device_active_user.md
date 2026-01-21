---
page_title: "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Logs out the currently active user from Shared iPad devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/logoutSharedAppleDeviceActiveUser endpoint. This action is used to manage user sessions on iPads configured in Shared iPad mode, where multiple users can use the same device while maintaining separate user environments.
  What is Shared iPad Mode?
  Educational/enterprise feature for iPadOSMultiple users share single physical deviceEach user has separate data and settingsUsers log in with their Apple ID or Managed Apple IDLocal caching of user data for offline accessRequires supervised iPads enrolled via DEP/ABM
  What This Action Does:
  Logs out currently active user from Shared iPadReturns device to login screenPreserves user data on device (cached locally)Allows next user to log inDoes not remove user from device rosterDoes not delete user's cached data
  Platform Support:
  iPadOS: Full support (Shared iPad mode only)iOS: Not supported (iPhones don't support Shared mode)macOS: Not supportedWindows: Not supportedAndroid: Not supported
  Common Use Cases:
  Classroom management (switching students)End of class period user logoutPreparing device for next userRemote user session managementEnforcing session time limitsCart/lab device rotationEmergency user logoutTroubleshooting user sessions
  Requirements:
  iPad must be in Shared iPad modeDevice must be supervisedMust be enrolled via DEP/ABMUser must be actively logged inDevice must be online
  Important Notes:
  Only affects Shared iPad devicesRegular (non-shared) iPads: action has no effectUser data remains cached on deviceUser can log back in immediatelyUnsaved work may be lostActive apps will close
  Reference: Microsoft Graph API - Logout Shared Apple Device Active User https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-logoutsharedappledeviceactiveuser?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user (Action)

Logs out the currently active user from Shared iPad devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/logoutSharedAppleDeviceActiveUser` endpoint. This action is used to manage user sessions on iPads configured in Shared iPad mode, where multiple users can use the same device while maintaining separate user environments.

**What is Shared iPad Mode?**
- Educational/enterprise feature for iPadOS
- Multiple users share single physical device
- Each user has separate data and settings
- Users log in with their Apple ID or Managed Apple ID
- Local caching of user data for offline access
- Requires supervised iPads enrolled via DEP/ABM

**What This Action Does:**
- Logs out currently active user from Shared iPad
- Returns device to login screen
- Preserves user data on device (cached locally)
- Allows next user to log in
- Does not remove user from device roster
- Does not delete user's cached data

**Platform Support:**
- **iPadOS**: Full support (Shared iPad mode only)
- **iOS**: Not supported (iPhones don't support Shared mode)
- **macOS**: Not supported
- **Windows**: Not supported
- **Android**: Not supported

**Common Use Cases:**
- Classroom management (switching students)
- End of class period user logout
- Preparing device for next user
- Remote user session management
- Enforcing session time limits
- Cart/lab device rotation
- Emergency user logout
- Troubleshooting user sessions

**Requirements:**
- iPad must be in Shared iPad mode
- Device must be supervised
- Must be enrolled via DEP/ABM
- User must be actively logged in
- Device must be online

**Important Notes:**
- Only affects Shared iPad devices
- Regular (non-shared) iPads: action has no effect
- User data remains cached on device
- User can log back in immediately
- Unsaved work may be lost
- Active apps will close

**Reference:** [Microsoft Graph API - Logout Shared Apple Device Active User](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-logoutsharedappledeviceactiveuser?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [logoutSharedAppleDeviceActiveUser action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-logoutsharedappledeviceactiveuser?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Device logout user](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-logout-user)

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
# Example 1: Logout active user from a single shared Apple device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_single" {
  config {
    device_ids = ["12345678-1234-1234-1234-123456789abc"]
  }
}

# Example 2: Logout active users from multiple shared Apple devices
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_multiple" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Logout with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_maximal" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Logout users from all shared iPads
data "microsoft365_graph_beta_device_management_managed_device" "shared_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (managementMode eq 'shared')"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_all_shared_ipads" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.shared_ipads.items : device.id]

    validate_device_exists  = true
    ignore_partial_failures = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Logout users from classroom iPads at end of day
data "microsoft365_graph_beta_device_management_managed_device" "classroom_ipads" {
  filter_type  = "odata"
  odata_filter = "(deviceCategoryDisplayName eq 'Classroom') and (operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_classroom_ipads" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.classroom_ipads.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of Shared iPad device IDs to log out the active user from. Each ID must be a valid GUID format. Multiple devices can be processed in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** This action only works on iPads configured in Shared iPad mode. The action will fail or have no effect on regular (non-shared) iPads, iPhones, or other device types. Ensure all device IDs refer to Shared iPad devices before executing this action.

### Optional

- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are configured for Shared iPad mode before attempting to log out the active user. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


