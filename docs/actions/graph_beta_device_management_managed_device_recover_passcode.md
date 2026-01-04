---
page_title: "microsoft365_graph_beta_device_management_managed_device_recover_passcode Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Recovers passcodes for managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/recoverPasscode endpoint. This action retrieves existing passcodes that are escrowed in Intune, which is different from reset passcode that generates new temporary passcodes. Recover passcode is primarily used for iOS/iPadOS devices where passcodes may be escrowed during enrollment or management.
  Important Notes:
  Retrieves existing escrowed passcode from IntuneDifferent from reset passcode (which creates new passcode)Passcode must have been previously escrowedPrimarily for iOS/iPadOS supervised devicesRetrieved passcode displayed in Intune portalMay not be available for all device types
  Use Cases:
  User forgot their device passcode (iOS/iPadOS)Supervised device lockout recoveryAdministrative access to escrowed passcodesDevice recovery without factory resetEmergency access to locked devicesHelp desk support for locked devices
  Platform Support:
  iOS/iPadOS: Supported (supervised devices with passcode escrow)macOS: Limited (may work with specific configurations)Windows: Not typically supported for passcode recoveryAndroid: Not typically supported for passcode recovery
  Passcode Escrow:
  Passcodes must be escrowed during device enrollmentNot all devices escrow passcodes automaticallySupervised iOS/iPadOS devices typically escrow passcodesCheck device enrollment configuration for escrow settingsRecovery may fail if passcode not escrowed
  Recover vs Reset Passcode:
  Recover: Retrieves existing escrowed passcode (no change to device)Reset: Generates new temporary passcode (device must be unlocked and reset)Use recover first if passcode is escrowedUse reset if recover fails or passcode not escrowed
  Reference: Microsoft Graph API - Recover Passcode https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-recoverpasscode?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_recover_passcode (Action)

Recovers passcodes for managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/recoverPasscode` endpoint. This action retrieves existing passcodes that are escrowed in Intune, which is different from reset passcode that generates new temporary passcodes. Recover passcode is primarily used for iOS/iPadOS devices where passcodes may be escrowed during enrollment or management.

**Important Notes:**
- Retrieves existing escrowed passcode from Intune
- Different from reset passcode (which creates new passcode)
- Passcode must have been previously escrowed
- Primarily for iOS/iPadOS supervised devices
- Retrieved passcode displayed in Intune portal
- May not be available for all device types

**Use Cases:**
- User forgot their device passcode (iOS/iPadOS)
- Supervised device lockout recovery
- Administrative access to escrowed passcodes
- Device recovery without factory reset
- Emergency access to locked devices
- Help desk support for locked devices

**Platform Support:**
- **iOS/iPadOS**: Supported (supervised devices with passcode escrow)
- **macOS**: Limited (may work with specific configurations)
- **Windows**: Not typically supported for passcode recovery
- **Android**: Not typically supported for passcode recovery

**Passcode Escrow:**
- Passcodes must be escrowed during device enrollment
- Not all devices escrow passcodes automatically
- Supervised iOS/iPadOS devices typically escrow passcodes
- Check device enrollment configuration for escrow settings
- Recovery may fail if passcode not escrowed

**Recover vs Reset Passcode:**
- **Recover**: Retrieves existing escrowed passcode (no change to device)
- **Reset**: Generates new temporary passcode (device must be unlocked and reset)
- Use recover first if passcode is escrowed
- Use reset if recover fails or passcode not escrowed

**Reference:** [Microsoft Graph API - Recover Passcode](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-recoverpasscode?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [recoverPasscode action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-recoverpasscode?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Device remove passcode](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-remove-passcode)

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

| Platform | Support | Escrow Requirements |
|----------|---------|---------------------|
| **iOS** | ✅ Full Support | Supervised, DEP/ABM enrolled, escrow enabled |
| **iPadOS** | ✅ Full Support | Supervised, DEP/ABM enrolled, escrow enabled |
| **macOS** | ⚠️ Limited | Specific configurations, rarely used |
| **Windows** | ❌ Not Supported | Windows passcode recovery not available |
| **Android** | ❌ Not Supported | Android passcode recovery not available |

### Recover vs Reset Passcode

| Action | What It Does | Best For | Requirements |
|--------|--------------|----------|--------------|
| **Recover** | Retrieves existing escrowed passcode | Supervised iOS/iPadOS with escrow | Passcode must be escrowed |
| **Reset** | Generates new temporary passcode | Any supported device | Device must be online |

### What is Passcode Escrow?

Passcode escrow is a security feature that:
- Stores encrypted copy of device passcode in Intune
- Configured during device enrollment setup
- Primarily available for supervised iOS/iPadOS devices
- Requires specific enrollment profile settings
- Enables IT to recover (not reset) user passcodes
- Useful for emergency device access

### Escrow Configuration Requirements

For passcode recovery to work, devices must:
1. Be supervised (iOS/iPadOS)
2. Enrolled via DEP/ABM or Apple Configurator
3. Have passcode escrow enabled in enrollment profile
4. Have user-set passcode after escrow enabled
5. Passcode must be actively escrowed (not expired)

### When to Use Recover vs Reset

**Use Recover Passcode When:**
- Device is supervised iOS/iPadOS
- Passcode escrow is confirmed enabled
- You want the user's original passcode
- No device configuration change desired
- User forgot passcode but device enrolled correctly

**Use Reset Passcode When:**
- Passcode recovery fails (not escrowed)
- Device is unsupervised
- Platform doesn't support escrow (Windows, Android)
- You need to force a passcode change
- Passcode escrow wasn't configured

### How to Verify Passcode Escrow

Before attempting recovery:
1. Check device enrollment profile settings
2. Verify "Escrow passcode" is enabled
3. Confirm device is supervised
4. Check device was enrolled with correct profile
5. Verify enrollment wasn't bypassed

### Retrieving Recovered Passcode

After successful recovery:
1. Navigate to Microsoft Intune admin center
2. Select Devices > All devices
3. Choose the device
4. View device properties or hardware information
5. Look for "Passcode" or "Recovery" section
6. Passcode displayed as plain text or retrievable
7. Securely communicate to authorized user

### Common Failure Reasons

| Error | Cause | Solution |
|-------|-------|----------|
| Not escrowed | Passcode never saved to Intune | Use reset passcode instead |
| Unsupervised | Device not in supervised mode | Re-enroll via DEP/ABM |
| Wrong profile | Enrolled without escrow enabled | Check enrollment profile settings |
| Expired escrow | Passcode changed after enrollment | May need to reset instead |
| Wrong platform | Windows/Android attempted | Use reset for these platforms |

### Security Considerations

- **Sensitive Data**: Recovered passcodes are actual user passcodes
- **Access Control**: Strictly limit who can recover passcodes
- **Verification**: Verify user identity before providing passcode
- **Communication**: Never send passcodes via email or unsecured chat
- **Documentation**: Log all passcode recovery requests
- **Compliance**: Ensure recovery aligns with privacy policies
- **Audit Trail**: Maintain records of who recovered what/when

### Best Practices

- Try recover before reset (preserves user experience)
- Verify escrow status before attempting recovery
- Have reset passcode as fallback option
- Implement approval workflow for recovery requests
- Train help desk on when to use recover vs reset
- Document business justification for recovery
- Secure passcode communication channels
- Monitor for patterns of repeated recovery
- Review enrollment profiles regularly
- Test escrow functionality with test devices

## Example Usage

```terraform
# Example 1: Recover passcode for a single iOS device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "single_device" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Recover passcodes for multiple supervised iOS devices
action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "multiple_devices" {
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

# Example 3: Recover passcodes with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_with_validation" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Recover passcodes for supervised iOS devices using data source
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_supervised_ios" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Recover passcodes for supervised iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_classroom_ipads" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ipads.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 6: Recover passcode for specific user's supervised iOS device
data "microsoft365_graph_beta_device_management_managed_device" "user_ios_device" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'user@example.com') and (operatingSystem eq 'iOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_user_device" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_ios_device.items : device.id]

    timeouts = {
      invoke = "5m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to recover passcodes for. Each ID must be a valid GUID format. Multiple devices can have passcodes recovered in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** This action retrieves existing escrowed passcodes. If a passcode was not escrowed during device enrollment, the recovery will fail. Check device properties in Intune to verify passcode escrow status before attempting recovery.

### Optional

- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist, are iOS/iPadOS devices, and are supervised before attempting to recover passcodes. Disabling this can speed up planning but may result in runtime errors for non-existent, unsupported, or unsupervised devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


