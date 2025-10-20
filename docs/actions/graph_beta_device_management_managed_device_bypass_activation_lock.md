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

  device_ids = [
    "12345678-1234-1234-1234-123456789abc"
  ]

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Bypass Activation Lock for multiple devices (batch processing)
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_batch" {

  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210",
    "abcdef12-3456-7890-abcd-ef1234567890"
  ]

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Bypass Activation Lock for supervised iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS') and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_supervised_ios" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items : device.id]

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Bypass Activation Lock for DEP-enrolled macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "dep_macos" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'macOS') and (deviceEnrollmentType eq 'deviceEnrollmentProgram')"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_dep_macos" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.dep_macos.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Bypass Activation Lock for departing employee's Apple devices
data "microsoft365_graph_beta_device_management_managed_device" "departing_user_apple_devices" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS' or operatingSystem eq 'macOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_departing_user" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_user_apple_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Bypass Activation Lock for corporate-owned Apple devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_apple" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS' or operatingSystem eq 'macOS') and managedDeviceOwnerType eq 'company'"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_corporate_apple" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_apple.items : device.id]

  timeouts = {
    invoke = "20m"
  }
}

# Example 7: Bypass Activation Lock for devices with specific model (e.g., iPhone 13)
data "microsoft365_graph_beta_device_management_managed_device" "iphone_13_devices" {
  filter_type  = "odata"
  odata_filter = "model eq 'iPhone 13' and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_iphone_13" {

  device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.iphone_13_devices.items : device.id]

  timeouts = {
    invoke = "10m"
  }
}

# Output examples
output "bypassed_device_count" {
  value       = length(action.bypass_batch.device_ids)
  description = "Number of devices that received Activation Lock bypass command"
}

output "bypassed_corporate_count" {
  value       = length(action.bypass_corporate_apple.device_ids)
  description = "Number of corporate Apple devices with bypass codes generated"
}

# Important Notes:
# 
# What is Activation Lock?
# - Apple security feature that prevents device reactivation after factory reset
# - Automatically enabled when Find My iPhone/iPad/Mac is turned on
# - Requires original Apple ID and password to reactivate device
# - Protects against unauthorized device use after theft or loss
#
# When to Use Bypass Activation Lock:
# - Preparing corporate devices for reassignment to new employees
# - Recovering devices from departing employees who forgot to disable Find My
# - Preparing devices for return to vendor or recycling
# - Handling devices with lost or forgotten Apple ID credentials
# - Bulk device refresh/replacement projects
# - Emergency device recovery scenarios
#
# Platform Requirements:
# - iOS/iPadOS: Devices MUST be supervised (enrolled via DEP/ABM or Apple Configurator)
# - macOS: Devices should be enrolled via DEP/ABM (Automated Device Enrollment)
# - Windows/Android: Not supported (Activation Lock is Apple-only feature)
#
# How Activation Lock Bypass Works:
# 1. Issue bypass command via this action
# 2. Intune generates a unique bypass code for each device
# 3. Bypass code is stored in Intune device properties
# 4. Retrieve bypass code from Intune admin portal
# 5. Factory reset/wipe the device
# 6. During setup, device shows Activation Lock screen
# 7. Enter bypass code in password field to bypass lock
# 8. Device completes setup without requiring user's Apple ID
#
# Workflow Example:
# Step 1: Employee leaves organization, device in their possession
# Step 2: IT issues bypass activation lock command (this action)
# Step 3: IT retrieves bypass code from Intune portal
# Step 4: IT wipes device remotely or obtains physical possession
# Step 5: During device setup, Activation Lock screen appears
# Step 6: IT enters bypass code to unlock device
# Step 7: Device can now be re-enrolled and assigned to new user
#
# Retrieving Bypass Codes:
# - Navigate to Intune admin center (https://intune.microsoft.com)
# - Go to Devices > All devices
# - Select the device
# - Under "Hardware" section, find "Activation Lock bypass code"
# - Copy code and securely store or use immediately
# - Code format: Usually 6-8 alphanumeric characters
#
# Security Considerations:
# - Bypass codes are sensitive credentials - treat like passwords
# - Only authorized IT staff should have access to codes
# - Document bypass code usage for compliance/audit purposes
# - Consider implementing approval workflow for bypass requests
# - Verify device ownership before issuing bypass
# - Bypass codes should be stored in secure password vault
# - Rotate or regenerate codes periodically if stored long-term
#
# Limitations:
# - Cannot bypass Activation Lock on personal (non-supervised) iOS/iPadOS devices
# - Cannot bypass on macOS devices not enrolled via DEP/ABM (limited support)
# - Bypass code only works for specific device it was generated for
# - Must retrieve bypass code before device is erased (code stored in Intune)
# - Some older device models may not support Activation Lock bypass
# - Device must have had Activation Lock enabled for bypass to be relevant
#
# Best Practices:
# - Issue bypass command BEFORE wiping device when possible
# - Store bypass codes securely in enterprise password manager
# - Document which devices have bypass codes generated
# - Include Activation Lock bypass in offboarding procedures
# - Test bypass process in controlled environment first
# - Verify device supervision status before attempting bypass
# - Consider enabling automatic bypass code escrow during enrollment
# - Train help desk staff on bypass code retrieval and usage
# - Maintain audit log of bypass code usage
#
# Common Issues and Solutions:
# 
# Issue: Bypass command fails with "Device not supervised"
# Solution: iOS/iPadOS devices must be supervised. Re-enroll via DEP/ABM or Apple Configurator
#
# Issue: Bypass command fails with "Activation Lock not enabled"
# Solution: Device doesn't have Find My enabled. No bypass needed for this device
#
# Issue: Bypass code doesn't work during device setup
# Solution: Verify you copied the code correctly. Try entering code in all caps or lowercase
#
# Issue: Can't find bypass code in Intune portal
# Solution: Code may take a few minutes to appear. Refresh device properties page
#
# Issue: macOS device doesn't accept bypass code
# Solution: Verify device was enrolled via DEP/ABM. Manually enrolled Macs have limited support
#
# Integration with Other Actions:
# - Often used before or after wipe action
# - Can be combined with retire for less aggressive device cleanup
# - May be used with disable lost mode if device is in lost mode
# - Should be part of comprehensive device lifecycle management
#
# Compliance and Legal Considerations:
# - Ensure you have legal right to bypass device (corporate ownership)
# - Document business justification for bypass in audit logs
# - Consider privacy implications in different jurisdictions
# - Review employment agreements regarding device management
# - Maintain records of device ownership and bypass authorization
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-bypassactivationlock?view=graph-rest-beta
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `device_ids` (List of String) List of managed device IDs to generate Activation Lock bypass codes for. Each ID must be a valid GUID format. Multiple devices can have bypass codes generated in a single action. Example: `["12345678-1234-1234-1234-123456789abc", "87654321-4321-4321-4321-ba9876543210"]`

**Important:** Devices must be supervised iOS/iPadOS devices or DEP/ABM enrolled macOS devices with Activation Lock enabled. The bypass code will be stored in Intune and can be retrieved from device properties in the admin portal. This code is required to reactivate the device after a factory reset if Activation Lock is enabled.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

