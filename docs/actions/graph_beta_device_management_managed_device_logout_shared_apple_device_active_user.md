---
page_title: "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Logs out the currently active user from Shared iPad devices using the /deviceManagement/managedDevices/{managedDeviceId}/logoutSharedAppleDeviceActiveUser endpoint. This action is specifically designed for iPads configured in Shared iPad mode, where multiple users can use the same device while maintaining separate user environments.
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

Logs out the currently active user from Shared iPad devices using the `/deviceManagement/managedDevices/{managedDeviceId}/logoutSharedAppleDeviceActiveUser` endpoint. This action is specifically designed for iPads configured in Shared iPad mode, where multiple users can use the same device while maintaining separate user environments.

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
- [Windows Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=windows)
- [iOS/iPadOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=ios-ipados)
- [macOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=macos)
- [Android Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=android)
- [ChromeOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=chromeos)

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
| **iPadOS** | ✅ Full Support | Shared iPad mode, supervised, DEP/ABM |
| **iOS** | ❌ Not Supported | iPhones don't support Shared mode |
| **macOS** | ❌ Not Supported | Shared mode is iPad-only |
| **Windows** | ❌ Not Supported | Shared mode is iPad-only |
| **Android** | ❌ Not Supported | Shared mode is iPad-only |

### What is Shared iPad Mode?

Shared iPad is an Apple educational and enterprise feature that enables:

**Multi-User Support**
- Multiple users can use the same physical iPad
- Each user maintains separate data, apps, and settings
- Users log in with Apple ID or Managed Apple ID
- Seamless switching between user accounts

**Local Data Caching**
- User data cached locally on device
- Offline access to user content
- Configurable storage per user
- Automatic cloud sync when online

**Management Features**
- Requires supervised iPads
- Must be enrolled via DEP/ABM
- Centralized configuration
- IT control over user quotas
- Session management capabilities

### How This Action Works

**Immediate Effect**
1. Active user is logged out
2. All user apps are closed
3. Device returns to login screen
4. Next user can immediately log in

**Data Preservation**
- User data remains cached on device
- Documents and files preserved
- App data saved
- User settings retained
- Photos and media kept
- User can re-login anytime

**What Doesn't Change**
- Device configuration
- Installed apps
- MDM enrollment
- Other cached users
- Device settings
- Network configuration

### Common Use Cases

| Scenario | Description | Benefit |
|----------|-------------|---------|
| **Classroom Rotation** | Switch students between periods | Efficient device sharing |
| **Lab Management** | Reset devices between sessions | Clean slate for each user |
| **Cart Devices** | Prepare devices for next day | Ready for morning distribution |
| **Emergency Logout** | Remote user session termination | Security/troubleshooting |
| **Scheduled Sessions** | Enforce time-limited access | Automated session management |
| **Device Rotation** | Prepare for next user in queue | Streamlined user switching |

### Requirements

#### Device Requirements
- iPad with Shared iPad mode enabled
- Supervised via DEP/ABM or Apple Configurator
- iPadOS 9.3 or later
- Sufficient storage for cached users
- Online (for remote command delivery)

#### User Requirements
- User must be actively logged in
- Uses Apple ID or Managed Apple ID
- Has active session on device

#### Infrastructure Requirements
- Apple School Manager or Apple Business Manager
- MDM (Intune) enrollment
- Managed Apple ID infrastructure (for education)
- Federation setup (for enterprise, optional)

### Action Behavior

#### Success Scenarios
- User logged out successfully
- Device returns to login screen
- User data preserved on device
- Ready for next user

#### No Effect Scenarios
- Regular (non-shared) iPad
- No user currently logged in
- Device already at login screen
- Device is iPhone (iOS)
- Unsupervised device

#### Failure Scenarios
- Device offline
- Device not in Shared iPad mode
- API communication error
- Insufficient permissions
- Device not supervised

### User Impact

**Immediate**
- Active user session terminated
- All apps close immediately
- Unsaved work is lost
- Device locks to login screen

**Data Impact**
- **Preserved**: User documents, photos, app data, settings
- **Lost**: Unsaved changes in open apps, clipboard content
- **Unaffected**: Other cached users, device configuration

**User Experience**
- Abrupt logout (no warning)
- Must log in again to access data
- All data available upon re-login
- No data loss (except unsaved work)

### Best Practices

**Scheduling**
- Logout during non-instructional time
- Coordinate with class schedules
- Use during breaks/transitions
- Avoid mid-session logouts

**Communication**
- Notify users before logout when possible
- Establish logout policies
- Train users to save work regularly
- Document standard procedures

**Implementation**
- Test with single device first
- Start with small groups
- Monitor for issues
- Have rollback plan
- Provide help desk support

**Automation**
- Integrate with scheduling systems
- Use time-based triggers
- Combine with other MDM actions
- Log all logout operations
- Monitor success rates

### Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Action fails | Device not Shared iPad | Verify Shared iPad configuration |
| No effect | No active user | Confirm user is logged in |
| Timeout | Device offline | Check network connectivity |
| Error | Not supervised | Enroll via DEP/ABM |
| Wrong device | iPhone (iOS) | Only use with iPadOS devices |
| User data lost | Unsaved work | Train users to save frequently |

### Shared iPad Configuration

**Setup Requirements**
1. Apple School Manager or Business Manager account
2. DEP/ABM enrollment for iPads
3. Shared iPad configuration profile
4. User accounts (Managed Apple IDs)
5. Storage allocation per user
6. Maximum cached users setting

**Configuration Options**
- Maximum resident users
- Storage quota per user
- Temporary session mode
- Guest access settings
- User authentication method
- Data retention policies
- Offline access limits

## Example Usage

```terraform
# ============================================================================
# Example 1: Logout active user from single Shared iPad
# ============================================================================
# Use case: End of class period logout
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "single_device" {

  device_ids = ["12345678-1234-1234-1234-123456789abc"]

  timeouts = {
    invoke = "5m"
  }
}

# ============================================================================
# Example 2: Logout active users from multiple Shared iPads
# ============================================================================
# Use case: End of day logout for classroom cart devices
action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "multiple_devices" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 3: Logout all Shared iPads in specific group
# ============================================================================
# Use case: Classroom management for scheduled logout
data "microsoft365_graph_beta_device_management_managed_device" "classroom_shared_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (isSupervised eq true)"
}

# Filter to only devices with "SharediPad" in the name
locals {
  shared_ipad_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.classroom_shared_ipads.items :
    device.id if can(regex("SharediPad", device.device_name))
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_classroom" {

  device_ids = local.shared_ipad_devices

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 4: Logout Shared iPads by device name pattern
# ============================================================================
# Use case: Lab or cart devices with specific naming convention
data "microsoft365_graph_beta_device_management_managed_device" "lab_ipads" {
  filter_type  = "device_name"
  filter_value = "LAB-IPAD-"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_lab_devices" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.lab_ipads.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 5: Logout supervised iPads (potential Shared iPads)
# ============================================================================
# Use case: End of semester cleanup for all supervised iPads
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_supervised" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ipads.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 6: Logout company-owned supervised iPads
# ============================================================================
# Use case: Institutional device rotation
data "microsoft365_graph_beta_device_management_managed_device" "company_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (managedDeviceOwnerType eq 'company') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_logout_shared_apple_device_active_user" "logout_company_ipads" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.company_ipads.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of Shared iPad device IDs to log out the active user from. Each ID must be a valid GUID format. Multiple devices can be processed in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** This action only works on iPads configured in Shared iPad mode. The action will fail or have no effect on regular (non-shared) iPads, iPhones, or other device types. Ensure all device IDs refer to Shared iPad devices before executing this action.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


