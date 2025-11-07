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
# Example 1: Update device account for a single Microsoft Teams Room with Exchange Online
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_teams_room" {

  managed_devices {
    device_id                           = "12345678-1234-1234-1234-123456789abc"
    device_account_email                = "conference-room-01@company.com"
    password                            = var.teams_room_password # Use sensitive variable
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "outlook.office365.com"
    session_initiation_protocol_address = "sip:conference-room-01@company.com"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Update multiple Teams Rooms in bulk
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_multiple_teams_rooms" {

  managed_devices {
    device_id                           = "11111111-1111-1111-1111-111111111111"
    device_account_email                = "meeting-room-a@company.com"
    password                            = var.room_a_password
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "outlook.office365.com"
    session_initiation_protocol_address = "sip:meeting-room-a@company.com"
  }

  managed_devices {
    device_id                           = "22222222-2222-2222-2222-222222222222"
    device_account_email                = "meeting-room-b@company.com"
    password                            = var.room_b_password
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "outlook.office365.com"
    session_initiation_protocol_address = "sip:meeting-room-b@company.com"
  }

  managed_devices {
    device_id                           = "33333333-3333-3333-3333-333333333333"
    device_account_email                = "meeting-room-c@company.com"
    password                            = var.room_c_password
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "outlook.office365.com"
    session_initiation_protocol_address = "sip:meeting-room-c@company.com"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Update co-managed device (managed by both Intune and SCCM)
action "microsoft365_graph_beta_device_management_managed_device_update_windows_device_account" "update_comanaged_device" {

  comanaged_devices {
    device_id                           = "55555555-5555-5555-5555-555555555555"
    device_account_email                = "hybrid-room@company.com"
    password                            = var.hybrid_room_password
    password_rotation_enabled           = true
    calendar_sync_enabled               = true
    exchange_server                     = "mail.company.local"
    session_initiation_protocol_address = "sip:hybrid-room@company.com"
  }

  timeouts = {
    invoke = "5m"
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed Windows devices to update with individual device account configurations. These are devices managed by both Intune and Configuration Manager (SCCM). Configuration is identical to managed_devices.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `managed_devices` (Attributes List) List of managed Windows devices to update with individual device account configurations. Each entry specifies a device ID and its complete device account settings including credentials, Exchange server, and synchronization options. These are devices fully managed by Intune only.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

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

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

