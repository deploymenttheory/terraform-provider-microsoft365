---
page_title: "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Bypasses Activation Lock on iOS, iPadOS, and macOS devices using the /deviceManagement/managedDevices/{managedDeviceId}/bypassActivationLock endpoint. Activation Lock is an Apple security feature that prevents unauthorized use of a device after it has been erased. When Find My iPhone/iPad/Mac is enabled and a device is erased, Activation Lock requires the original Apple ID and password before the device can be reactivated. This action generates a bypass code that allows IT administrators to reactivate managed devices without the user's Apple ID credentials.
  What is Activation Lock?
  Security feature built into iOS, iPadOS, and macOSAutomatically enabled when Find My iPhone/iPad/Mac is turned onPrevents device reactivation after factory reset without Apple ID credentialsHelps prevent theft and unauthorized device reuseLinks device to specific Apple ID
  Important Notes:
  Device must be supervised (iOS/iPadOS) or enrolled via DEP/ABM (macOS)Activation Lock must currently be enabled on the deviceGenerates a bypass code stored in Intune for future useBypass code can be retrieved from device properties in Intune portalCode can be used during device setup to bypass Activation Lock screenDoes not disable Find My iPhone/iPad/Mac, only provides bypass capabilityBypass code remains valid until Activation Lock is disabled by user
  Use Cases:
  Wiping and reassigning corporate devices to new employeesRecovering devices from departing employees who forgot to disable Find MyPreparing devices for return to vendor or recyclingEnabling IT to factory reset and redeploy devicesHandling devices with lost or forgotten Apple ID credentialsBulk device preparation and provisioning
  Platform Support:
  iOS: Supported (iOS 7.1+, supervised devices only)iPadOS: Supported (supervised devices only)macOS: Supported (macOS 10.11+, DEP/ABM enrolled devices)Other Platforms: Not supported (Activation Lock is Apple-only feature)
  How to Use Bypass Code:
  Issue bypass command via this actionRetrieve bypass code from Intune portal (device properties)Erase/wipe the device (using wipe action or manually)When device shows Activation Lock screen during setupEnter bypass code in password fieldDevice will bypass Activation Lock and complete setup
  Workflow Example:
  
  1. Employee leaves organization with device in lost mode
  2. IT issues bypass activation lock command (this action)
  3. IT retrieves bypass code from Intune portal
  4. IT wipes device (removes all data)
  5. During setup, device shows Activation Lock screen
  6. IT enters bypass code to unlock device
  7. Device can now be re-enrolled and assigned to new user
  
  Security Considerations:
  Bypass code should be treated as sensitive credentialOnly authorized IT staff should have access to bypass codesDocument usage for compliance and audit purposesConsider implementing approval workflow for bypass requestsVerify device ownership before issuing bypassBypass does not affect device security after reactivation
  Limitations:
  Cannot bypass Activation Lock on personal (non-supervised) iOS/iPadOS devicesCannot bypass Activation Lock on macOS devices not enrolled via DEP/ABMBypass code only works for device it was generated forMust have Activation Lock bypass code retrieved before device is erasedSome older device models may not support this feature
  Best Practices:
  Issue bypass command before wiping device when possibleStore bypass codes securely in password manager or secure vaultDocument which devices have bypass codes generatedInclude activation lock bypass in device offboarding proceduresTest bypass process in controlled environment firstVerify device supervision status before attempting bypassConsider enabling automatic bypass code escrow during enrollment
  Reference: Microsoft Graph API - Bypass Activation Lock https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-bypassactivationlock?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock (Action)

Bypasses Activation Lock on iOS, iPadOS, and macOS devices using the `/deviceManagement/managedDevices/{managedDeviceId}/bypassActivationLock` endpoint. Activation Lock is an Apple security feature that prevents unauthorized use of a device after it has been erased. When Find My iPhone/iPad/Mac is enabled and a device is erased, Activation Lock requires the original Apple ID and password before the device can be reactivated. This action generates a bypass code that allows IT administrators to reactivate managed devices without the user's Apple ID credentials.

**What is Activation Lock?**
- Security feature built into iOS, iPadOS, and macOS
- Automatically enabled when Find My iPhone/iPad/Mac is turned on
- Prevents device reactivation after factory reset without Apple ID credentials
- Helps prevent theft and unauthorized device reuse
- Links device to specific Apple ID

**Important Notes:**
- Device must be supervised (iOS/iPadOS) or enrolled via DEP/ABM (macOS)
- Activation Lock must currently be enabled on the device
- Generates a bypass code stored in Intune for future use
- Bypass code can be retrieved from device properties in Intune portal
- Code can be used during device setup to bypass Activation Lock screen
- Does not disable Find My iPhone/iPad/Mac, only provides bypass capability
- Bypass code remains valid until Activation Lock is disabled by user

**Use Cases:**
- Wiping and reassigning corporate devices to new employees
- Recovering devices from departing employees who forgot to disable Find My
- Preparing devices for return to vendor or recycling
- Enabling IT to factory reset and redeploy devices
- Handling devices with lost or forgotten Apple ID credentials
- Bulk device preparation and provisioning

**Platform Support:**
- **iOS**: Supported (iOS 7.1+, supervised devices only)
- **iPadOS**: Supported (supervised devices only)
- **macOS**: Supported (macOS 10.11+, DEP/ABM enrolled devices)
- **Other Platforms**: Not supported (Activation Lock is Apple-only feature)

**How to Use Bypass Code:**
1. Issue bypass command via this action
2. Retrieve bypass code from Intune portal (device properties)
3. Erase/wipe the device (using wipe action or manually)
4. When device shows Activation Lock screen during setup
5. Enter bypass code in password field
6. Device will bypass Activation Lock and complete setup

**Workflow Example:**
```
1. Employee leaves organization with device in lost mode
2. IT issues bypass activation lock command (this action)
3. IT retrieves bypass code from Intune portal
4. IT wipes device (removes all data)
5. During setup, device shows Activation Lock screen
6. IT enters bypass code to unlock device
7. Device can now be re-enrolled and assigned to new user
```

**Security Considerations:**
- Bypass code should be treated as sensitive credential
- Only authorized IT staff should have access to bypass codes
- Document usage for compliance and audit purposes
- Consider implementing approval workflow for bypass requests
- Verify device ownership before issuing bypass
- Bypass does not affect device security after reactivation

**Limitations:**
- Cannot bypass Activation Lock on personal (non-supervised) iOS/iPadOS devices
- Cannot bypass Activation Lock on macOS devices not enrolled via DEP/ABM
- Bypass code only works for device it was generated for
- Must have Activation Lock bypass code retrieved before device is erased
- Some older device models may not support this feature

**Best Practices:**
- Issue bypass command before wiping device when possible
- Store bypass codes securely in password manager or secure vault
- Document which devices have bypass codes generated
- Include activation lock bypass in device offboarding procedures
- Test bypass process in controlled environment first
- Verify device supervision status before attempting bypass
- Consider enabling automatic bypass code escrow during enrollment

**Reference:** [Microsoft Graph API - Bypass Activation Lock](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-bypassactivationlock?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [bypassActivationLock action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-bypassactivationlock?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Device activation lock disable - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-activation-lock-disable?pivots=ios)
- [Device activation lock disable - macOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/device-activation-lock-disable?pivots=macos)

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
| **iOS** | ✅ Full Support | iOS 7.1+, supervised devices only |
| **iPadOS** | ✅ Full Support | Supervised devices only |
| **macOS** | ✅ Full Support | macOS 10.11+, DEP/ABM enrolled |
| **Windows** | ❌ Not Supported | Activation Lock is Apple-only |
| **Android** | ❌ Not Supported | Activation Lock is Apple-only |

### What is Activation Lock?

Activation Lock is an Apple security feature that:
- Prevents unauthorized device reactivation after factory reset
- Automatically enabled when Find My iPhone/iPad/Mac is turned on
- Requires original Apple ID and password to reactivate
- Links device to specific Apple ID

### How Bypass Works

1. Issue bypass command via this action
2. Intune generates unique bypass code for device
3. Bypass code stored in Intune device properties
4. Retrieve code from Intune admin portal
5. Factory reset/wipe the device
6. During setup, device shows Activation Lock screen
7. Enter bypass code to unlock device
8. Device completes setup without user's Apple ID

### Important Considerations

- **Supervised iOS/iPadOS**: Devices must be supervised (DEP/ABM or Apple Configurator)
- **macOS DEP**: Best results with DEP/ABM enrolled macOS devices
- **Code Security**: Bypass codes are sensitive credentials - treat like passwords
- **Single Use**: Bypass code only works for specific device it was generated for
- **Physical Access**: Code must be entered during device setup at Activation Lock screen

## Example Usage

```terraform
# Example 1: Bypass Activation Lock for a single device
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_single_device" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 2: Bypass Activation Lock for multiple devices (batch processing)
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_batch" {
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

# Example 3: Bypass Activation Lock for supervised iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS') and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_supervised_ios" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 4: Bypass Activation Lock for DEP-enrolled macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "dep_macos" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'macOS') and (deviceEnrollmentType eq 'deviceEnrollmentProgram')"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_dep_macos" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.dep_macos.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Bypass Activation Lock for departing employee's Apple devices
data "microsoft365_graph_beta_device_management_managed_device" "departing_user_apple_devices" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS' or operatingSystem eq 'macOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_departing_user" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_user_apple_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Bypass Activation Lock for corporate-owned Apple devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_apple" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS' or operatingSystem eq 'macOS') and managedDeviceOwnerType eq 'company'"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_corporate_apple" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_apple.items : device.id]

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 7: Bypass Activation Lock for devices with specific model (e.g., iPhone 13)
data "microsoft365_graph_beta_device_management_managed_device" "iphone_13_devices" {
  filter_type  = "odata"
  odata_filter = "model eq 'iPhone 13' and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_iphone_13" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.iphone_13_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "bypassed_device_count" {
  value       = length(action.bypass_batch.config.device_ids)
  description = "Number of devices that received Activation Lock bypass command"
}

output "bypassed_corporate_count" {
  value       = length(action.bypass_corporate_apple.config.device_ids)
  description = "Number of corporate Apple devices with bypass codes generated"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to generate Activation Lock bypass codes for. Each ID must be a valid GUID format. Multiple devices can have bypass codes generated in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** Devices must be supervised iOS/iPadOS devices or DEP/ABM enrolled macOS devices with Activation Lock enabled. The bypass code will be stored in Intune and can be retrieved from device properties in the admin portal. This code is required to reactivate the device after a factory reset if Activation Lock is enabled.

### Optional

- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some devices fail Activation Lock bypass. Failed devices will be reported as warnings instead of errors. Default: `false` (action fails if any device fails).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and support Activation Lock before attempting bypass. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

