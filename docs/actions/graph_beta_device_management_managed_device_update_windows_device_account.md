---
page_title: "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Updates the device account configuration on Windows devices using the /deviceManagement/managedDevices/{managedDeviceId}/updateWindowsDeviceAccount and /deviceManagement/comanagedDevices/{managedDeviceId}/updateWindowsDeviceAccount endpoints. This action is specifically designed for shared Windows devices like Surface Hub and Microsoft Teams Rooms that require device account configuration for Exchange and Skype for Business/Teams integration.
  What This Action Does:
  Updates device account credentialsConfigures Exchange server settingsSets up calendar syncConfigures Teams/SfB settingsManages password rotationUpdates SIP address configuration
  Target Devices:
  Surface Hub: Collaboration devicesMicrosoft Teams Rooms: Meeting room systemsShared Windows devices: Kiosk/common area devices
  Platform Support:
  Windows 10/11: Surface Hub, Teams RoomsWindows 10 IoT: Teams Rooms appliances
  Common Use Cases:
  Update device account passwordReconfigure Exchange serverUpdate calendar sync settingsChange Teams/SfB configurationRotate device credentialsFix authentication issues
  Important Considerations:
  Requires device reboot to applyPassword stored securelyExchange connectivity requiredTeams/SfB license neededAffects device functionality
  Reference: Microsoft Graph API - Update Windows Device Account https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-updatewindowsdeviceaccount?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_update_windows_device_account (Action)

Updates the device account configuration on Windows devices using the `/deviceManagement/managedDevices/{managedDeviceId}/updateWindowsDeviceAccount` and `/deviceManagement/comanagedDevices/{managedDeviceId}/updateWindowsDeviceAccount` endpoints. This action is specifically designed for shared Windows devices like Surface Hub and Microsoft Teams Rooms that require device account configuration for Exchange and Skype for Business/Teams integration.

**What This Action Does:**
- Updates device account credentials
- Configures Exchange server settings
- Sets up calendar sync
- Configures Teams/SfB settings
- Manages password rotation
- Updates SIP address configuration

**Target Devices:**
- **Surface Hub**: Collaboration devices
- **Microsoft Teams Rooms**: Meeting room systems
- **Shared Windows devices**: Kiosk/common area devices

**Platform Support:**
- **Windows 10/11**: Surface Hub, Teams Rooms
- **Windows 10 IoT**: Teams Rooms appliances

**Common Use Cases:**
- Update device account password
- Reconfigure Exchange server
- Update calendar sync settings
- Change Teams/SfB configuration
- Rotate device credentials
- Fix authentication issues

**Important Considerations:**
- Requires device reboot to apply
- Password stored securely
- Exchange connectivity required
- Teams/SfB license needed
- Affects device functionality

**Reference:** [Microsoft Graph API - Update Windows Device Account](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-updatewindowsdeviceaccount?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [updateWindowsDeviceAccount action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-updatewindowsdeviceaccount?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.PrivilegedOperations.All`
- **Delegated**: `DeviceManagementManagedDevices.PrivilegedOperations.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Full Support | Windows devices |
| **macOS** | ❌ Not Supported | Not supported |
| **iOS/iPadOS** | ❌ Not Supported | Not supported |
| **Android** | ❌ Not Supported | Not supported |

## Example Usage

```terraform
# Example 1: Update Windows device account on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_single" {
  config {
    managed_devices = [
      {
        device_id                 = "12345678-1234-1234-1234-123456789abc"
        device_account_email      = "conference-room-01@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]
  }
}

# Example 2: Update multiple Windows device accounts
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_multiple" {
  config {
    managed_devices = [
      {
        device_id                 = "12345678-1234-1234-1234-123456789abc"
        device_account_email      = "conference-room-01@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      },
      {
        device_id                 = "87654321-4321-4321-4321-ba9876543210"
        device_account_email      = "conference-room-02@company.com"
        password                  = "SecurePassword456!"
        password_rotation_enabled = false
        calendar_sync_enabled     = false
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Update with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_maximal" {
  config {
    managed_devices = [
      {
        device_id                 = "12345678-1234-1234-1234-123456789abc"
        device_account_email      = "conference-room-01@company.com"
        password                  = "SecurePassword123!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]

    comanaged_devices = [
      {
        device_id                 = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        device_account_email      = "meeting-room-03@company.com"
        password                  = "SecurePassword789!"
        password_rotation_enabled = true
        calendar_sync_enabled     = true
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Update Surface Hub devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "surface_hubs" {
  filter_type  = "odata"
  odata_filter = "model eq 'Surface Hub'"
}

action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_surface_hubs" {
  config {
    managed_devices = [
      for idx, device in data.microsoft365_graph_beta_device_management_managed_device.surface_hubs.items : {
        device_id                 = device.id
        device_account_email      = format("hub-%02d@company.com", idx + 1)
        password                  = format("SecurePass%03d!", idx + 1)
        password_rotation_enabled = true
        calendar_sync_enabled     = true
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

- `comanaged_devices` (Attributes List) List of co-managed Windows devices to update with individual device account configurations. These are devices managed by both Intune and Configuration Manager (SCCM). Configuration is identical to managed_devices.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to update. When `false` (default), the action will fail if any device update fails. Use this flag when updating multiple devices and you want the action to succeed even if some updates fail.
- `managed_devices` (Attributes List) List of managed Windows devices to update with individual device account configurations. Each entry specifies a device ID and its complete device account settings including credentials, Exchange server, and synchronization options. These are devices fully managed by Intune only.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist and are Windows devices before attempting to update them. When `false`, device validation is skipped and the action will attempt to update devices directly. Disabling validation can improve performance but may result in errors if devices don't exist or are not Windows devices.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `calendar_sync_enabled` (Boolean) Whether calendar synchronization is enabled. See managed_devices.calendar_sync_enabled for details.
- `device_account_email` (String) The email address of the device account. See managed_devices.device_account_email for details.
- `device_id` (String) The co-managed device ID (GUID). Example: `12345678-1234-1234-1234-123456789abc`
- `password` (String) The password for the device account. See managed_devices.password for details.
- `password_rotation_enabled` (Boolean) Whether automatic password rotation is enabled. See managed_devices.password_rotation_enabled for details.

Optional:

- `exchange_server` (String) The Exchange server address. See managed_devices.exchange_server for details.
- `session_initiation_protocol_address` (String) The SIP address for Teams/SfB. See managed_devices.session_initiation_protocol_address for details.


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `calendar_sync_enabled` (Boolean) Whether calendar synchronization is enabled for the device. This determines if the device will sync its calendar from Exchange.

- **`true`**: Enable calendar sync (shows meetings, availability)
- **`false`**: Disable calendar sync (no meeting information displayed)

**Use Cases:**
- Teams Rooms: Typically enabled (display meeting schedule)
- Surface Hub: Typically enabled (meeting coordination)
- Kiosk devices: May be disabled (no calendar needed)
- `device_account_email` (String) The email address of the device account (resource mailbox). This is typically a room mailbox in Exchange for Teams Rooms or Surface Hub. Example: `conference-room-01@company.com` or `surfacehub-lobby@company.com`

**Requirements:**
- Must be a valid email address
- Must exist in Exchange/Microsoft 365
- Should be a room or equipment mailbox
- Requires appropriate licenses
- `device_id` (String) The managed device ID (GUID) of the Windows device to update. Example: `12345678-1234-1234-1234-123456789abc`
- `password` (String) The password for the device account. This password is used to authenticate the device with Exchange and Teams/Skype for Business services. The password is transmitted securely and stored encrypted.

**Best Practices:**
- Use a strong, complex password
- Consider enabling password rotation
- Store securely (use Terraform sensitive values)
- Rotate regularly for security
- Follow organizational password policies
- `password_rotation_enabled` (Boolean) Whether automatic password rotation is enabled for the device account. When enabled, the device will automatically rotate its password periodically.

- **`true`**: Enable automatic password rotation (recommended for security)
- **`false`**: Disable automatic password rotation (manual management required)

**Note:** When enabled, ensure the device account has appropriate permissions in Active Directory to change its own password.

Optional:

- `exchange_server` (String) The Exchange server address for mailbox connectivity. This can be an on-premises Exchange server or Exchange Online (Microsoft 365).

**Examples:**
- Exchange Online: `outlook.office365.com`
- On-premises: `mail.company.com` or `exchange.company.local`
- Autodiscover: Leave blank to use autodiscover

**Note:** If not specified, the device will attempt to use Exchange autodiscover to locate the appropriate Exchange server automatically.
- `session_initiation_protocol_address` (String) The Session Initiation Protocol (SIP) address for Teams/Skype for Business connectivity. This is the SIP URI for the device account, typically matching the email address but with 'sip:' prefix.

**Format:** `sip:username@domain.com`

**Examples:**
- `sip:conference-room-01@company.com`
- `sip:surfacehub-lobby@company.com`

**Requirements:**
- Required for Teams/Skype for Business functionality
- Must match the device account UPN or email
- Account must be enabled for Teams/SfB
- Requires appropriate licensing


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

