---
page_title: "microsoft365_graph_beta_device_management_managed_device_remote_lock Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Remotely locks managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/remoteLock endpoint. This action immediately locks the device screen, requiring the user to enter their passcode to unlock it. This is useful for securing lost or stolen devices, or for security compliance scenarios. This action supports remotely locking multiple devices in a single operation.
  Important Notes:
  The device must be online and able to receive the commandThe device will lock immediately when it receives the commandThe user's existing passcode remains unchangedThe user will need to enter their passcode to unlock the deviceFor lost/stolen devices, consider using remote lock before more drastic measuresThis action does not remove any data from the device
  Use Cases:
  Lost or stolen device - immediate security actionSecurity incident - prevent unauthorized accessCompliance enforcement - ensure device is securedUnattended device in public locationUser reported potential device compromise
  Platform Support:
  iOS/iPadOS: Fully supported (iOS 9.0+)Android: Supported on Android Enterprise devices (work profile and fully managed)Windows: Supported on Windows 10/11 devicesmacOS: Supported on managed Mac computers
  Reference: Microsoft Graph API - Remote Lock https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-remotelock?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_remote_lock (Action)

Remotely locks managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/remoteLock` endpoint. This action immediately locks the device screen, requiring the user to enter their passcode to unlock it. This is useful for securing lost or stolen devices, or for security compliance scenarios. This action supports remotely locking multiple devices in a single operation.

**Important Notes:**
- The device must be online and able to receive the command
- The device will lock immediately when it receives the command
- The user's existing passcode remains unchanged
- The user will need to enter their passcode to unlock the device
- For lost/stolen devices, consider using remote lock before more drastic measures
- This action does not remove any data from the device

**Use Cases:**
- Lost or stolen device - immediate security action
- Security incident - prevent unauthorized access
- Compliance enforcement - ensure device is secured
- Unattended device in public location
- User reported potential device compromise

**Platform Support:**
- **iOS/iPadOS**: Fully supported (iOS 9.0+)
- **Android**: Supported on Android Enterprise devices (work profile and fully managed)
- **Windows**: Supported on Windows 10/11 devices
- **macOS**: Supported on managed Mac computers

**Reference:** [Microsoft Graph API - Remote Lock](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-remotelock?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [remoteLock action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-remotelock?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Device remote lock - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-remote-lock?pivots=ios)
- [Device remote lock - macOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-remote-lock?pivots=macos)
- [Device remote lock - Android](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-remote-lock?pivots=android)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.PrivilegedOperations.All`
- **Delegated**: `DeviceManagementManagedDevices.PrivilegedOperations.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Full Support | Windows 10/11 |
| **macOS** | ✅ Full Support | All supported versions |
| **iOS** | ✅ Full Support | All supported versions |
| **iPadOS** | ✅ Full Support | All supported versions |
| **Android** | ✅ Full Support | Fully managed and work profile |

### How Remote Lock Works

- Command sent to device immediately if online
- Device locks screen instantly upon receiving command
- User must enter existing passcode to unlock
- Lock is enforced by operating system
- Command queued if device offline
- Device locks when it next checks in
- No data is erased, device remains functional once unlocked

### Common Use Cases

- Lost device scenarios
- Stolen device immediate response
- Security incident response
- Unauthorized device access prevention
- Compliance enforcement
- Compromised credential response
- Emergency device isolation
- User-requested device locking

### User Impact

- Screen locks immediately
- Existing passcode required to unlock
- No data loss
- Notifications may still appear on lock screen
- Device remains connected to network
- Background processes continue
- No warning provided to user

### Best Practices

- Document reason for lock in incident log
- Combine with locate device for lost devices
- Follow up with user communication
- Consider lost mode for iOS/iPadOS devices
- Use as first step in security incident response
- Monitor device status after lock command
- Have unlock procedure documented
- Consider follow-up actions (wipe, retire) if needed

## Example Usage

```terraform
# Example 1: Remote lock a single device (lost device scenario)
action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_lost_device" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Remote lock multiple devices
action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Lock all devices for a specific user (security incident)
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "user_id"
  filter_value = "compromised.user@example.com"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_compromised_user" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Lock non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_non_compliant_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 5: Lock iOS devices reported as lost
data "microsoft365_graph_beta_device_management_managed_device" "ios_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_ios" {

  # In production, you would have additional filtering for "lost" status
  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.ios_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Emergency lock all corporate Windows devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and managedDeviceOwnerType eq 'company'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_corporate_windows" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_windows.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 7: Lock Android Enterprise devices
data "microsoft365_graph_beta_device_management_managed_device" "android_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Android'"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_android" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.android_devices.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 8: Lock devices by device name pattern (department-specific)
data "microsoft365_graph_beta_device_management_managed_device" "department_devices" {
  filter_type  = "device_name"
  filter_value = "SALES-"
}

action "microsoft365_graph_beta_device_management_managed_device_remote_lock" "lock_sales_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.department_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "locked_device_count" {
  value       = length(action.lock_batch.device_ids)
  description = "Number of devices that received remote lock command"
}

output "emergency_locked_count" {
  value       = length(action.lock_compromised_user.device_ids)
  description = "Number of devices locked in emergency scenario"
}

# Important Note: 
# - Devices lock IMMEDIATELY when they receive the command
# - Users must enter their existing passcode to unlock
# - For lost devices, follow up with locate/wipe if needed
# - Document the reason for locking devices for compliance/audit purposes
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to remotely lock. Each ID must be a valid GUID format. Multiple devices can be locked in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** Devices will lock immediately when they receive the command. Ensure you have authorization to lock these devices. For lost or stolen devices, this provides an immediate security measure while you determine next steps (locate, wipe, etc.).

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

