---
page_title: "microsoft365_graph_beta_device_management_managed_device_activate_device_esim Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Activates eSIM on managed cellular devices using the /deviceManagement/managedDevices/{managedDeviceId}/activateDeviceEsim and /deviceManagement/comanagedDevices/{managedDeviceId}/activateDeviceEsim endpoints. This action enables eSIM functionality on compatible devices by providing a carrier activation URL. eSIM (embedded SIM) technology allows devices to connect to cellular networks without a physical SIM card, providing greater flexibility for device deployment and carrier management. This action supports activating eSIM on multiple devices in a single operation with per-device carrier URL configuration.
  Important Notes:
  Only applicable to devices with eSIM hardware capabilityRequires carrier-specific activation URLDevice must support eSIM technologyCarrier must support eSIM activationDevice must be online to receive activationEach device requires its own carrier activation URL
  Use Cases:
  Initial eSIM activation on new devicesSwitching carriers on eSIM-capable devicesBulk eSIM deployment for corporate devicesRemote eSIM provisioning for field devicesInternational device deployment with local carriers
  Platform Support:
  iOS/iPadOS: Supported on eSIM-capable devices (iPhone XS and later, cellular iPads)Windows: Supported on eSIM-capable Windows devices with cellular modemsAndroid: Support varies by device manufacturer and Android versionOther Platforms: Not applicable
  Reference: Microsoft Graph API - Activate Device eSIM https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-activatedeviceesim?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_activate_device_esim (Action)

Activates eSIM on managed cellular devices using the `/deviceManagement/managedDevices/{managedDeviceId}/activateDeviceEsim` and `/deviceManagement/comanagedDevices/{managedDeviceId}/activateDeviceEsim` endpoints. This action enables eSIM functionality on compatible devices by providing a carrier activation URL. eSIM (embedded SIM) technology allows devices to connect to cellular networks without a physical SIM card, providing greater flexibility for device deployment and carrier management. This action supports activating eSIM on multiple devices in a single operation with per-device carrier URL configuration.

**Important Notes:**
- Only applicable to devices with eSIM hardware capability
- Requires carrier-specific activation URL
- Device must support eSIM technology
- Carrier must support eSIM activation
- Device must be online to receive activation
- Each device requires its own carrier activation URL

**Use Cases:**
- Initial eSIM activation on new devices
- Switching carriers on eSIM-capable devices
- Bulk eSIM deployment for corporate devices
- Remote eSIM provisioning for field devices
- International device deployment with local carriers

**Platform Support:**
- **iOS/iPadOS**: Supported on eSIM-capable devices (iPhone XS and later, cellular iPads)
- **Windows**: Supported on eSIM-capable Windows devices with cellular modems
- **Android**: Support varies by device manufacturer and Android version
- **Other Platforms**: Not applicable

**Reference:** [Microsoft Graph API - Activate Device eSIM](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-activatedeviceesim?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [activateDeviceEsim action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-activatedeviceesim?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune eSIM and Cellular Management
- [Activate device eSim action - iOS](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/update-cellular-data-plan

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.PrivilegedOperations.All`
- **Delegated**: `DeviceManagementManagedDevices.PrivilegedOperations.All`

## Version History

|| Version | Status | Notes |
||---------|--------|-------|
|| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

|| Platform | Support | Requirements |
||----------|---------|--------------|
|| **iPad** | ✅ Full Support | iPadOS 12.1+ |
|| **iPhone** | ✅ Full Support | iOS 12.1+ |
|| **Android** | ❌ Not Supported | eSIM activation not supported |
|| **Windows** | ❌ Not Supported | eSIM activation not supported |
|| **macOS** | ❌ Not Supported | eSIM activation not supported |

### What is eSIM?

eSIM is an embedded SIM that:
- Enables cellular connectivity without physical SIM cards
- Allows remote activation of carrier profiles
- Supports profile switching without device restart
- Reduces physical SIM inventory management
- Provides more secure device management options
- Works through Intune management policies

### How eSIM Activation Works

1. IT configures eSIM cellular plans in carrier systems
2. Intune receives eSIM activation server URLs from carriers
3. Intune pushes eSIM activation details to iOS/iPadOS devices
4. This action triggers device to activate eSIM profile
5. Device contacts activation server and downloads eSIM profile
6. eSIM profile is installed on device's embedded SIM
7. Cellular connectivity becomes active automatically
8. Device appears in Settings app's cellular section

### Prerequisites

- Device must be iOS/iPadOS device with eSIM hardware capability
- Device must be enrolled in Intune
- Device must have WiFi or existing cellular connectivity
- eSIM activation server URL must be configured in Intune
- eSIM plan must be provisioned by carrier
- Device must be in line of sight with internet connectivity

### Important Considerations

- **iOS/iPadOS Only**: This action only works on iOS 12.1+ and iPadOS 12.1+ devices
- **eSIM Hardware Required**: Device must have embedded SIM capability (not all iOS devices do)
- **Connectivity Required**: Device needs WiFi or existing cellular to download eSIM profile
- **Activation Server**: Requires activation server URL from carrier
- **One Profile at a Time**: Action activates a single eSIM profile per invocation
- **Dual SIM Support**: Some devices support multiple eSIM profiles or eSIM + physical SIM
- **Permanent Activation**: Once activated, eSIM profile remains on device until manually removed
- **No User Interaction**: Activation happens automatically without requiring user action


## Example Usage

```terraform
# Example 1: Activate eSIM on a single device
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_single" {

  managed_devices {
    device_id   = "12345678-1234-1234-1234-123456789abc"
    carrier_url = "https://carrier.example.com/esim/activate?token=abc123xyz"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Activate eSIM on multiple devices with different carriers
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_multiple" {

  managed_devices {
    device_id   = "12345678-1234-1234-1234-123456789abc"
    carrier_url = "https://carrier-a.example.com/esim/activate?code=device1"
  }

  managed_devices {
    device_id   = "87654321-4321-4321-4321-ba9876543210"
    carrier_url = "https://carrier-b.example.com/esim/activate?code=device2"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Bulk eSIM activation for new device deployment
variable "new_devices_with_esim" {
  description = "Map of device IDs to carrier activation URLs"
  type        = map(string)
  default = {
    "11111111-1111-1111-1111-111111111111" = "https://carrier.com/activate?id=1"
    "22222222-2222-2222-2222-222222222222" = "https://carrier.com/activate?id=2"
    "33333333-3333-3333-3333-333333333333" = "https://carrier.com/activate?id=3"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "bulk_activation" {

  dynamic "managed_devices" {
    for_each = var.new_devices_with_esim
    content {
      device_id   = managed_devices.key
      carrier_url = managed_devices.value
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Activate eSIM on cellular iPads
data "microsoft365_graph_beta_device_management_managed_device" "cellular_ipads" {
  filter_type  = "odata"
  odata_filter = "deviceType eq 'iPad' and model contains 'Cellular'"
}

locals {
  # Carrier URLs would typically come from carrier API or provisioning system
  ipad_carrier_urls = {
    for device in data.microsoft365_graph_beta_device_management_managed_device.cellular_ipads.items :
    device.id => "https://carrier.example.com/esim/activate?sn=${device.serial_number}"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_ipads" {

  dynamic "managed_devices" {
    for_each = local.ipad_carrier_urls
    content {
      device_id   = managed_devices.key
      carrier_url = managed_devices.value
    }
  }

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Activate eSIM for international deployment
locals {
  international_devices = {
    # Europe region devices - Local carrier
    "aaaa1111-1111-1111-1111-111111111111" = "https://eu-carrier.example.com/esim/activate?region=eu&token=abc"
    # Asia region devices - Local carrier
    "bbbb2222-2222-2222-2222-222222222222" = "https://asia-carrier.example.com/esim/activate?region=asia&token=def"
    # Americas region devices - Local carrier
    "cccc3333-3333-3333-3333-333333333333" = "https://us-carrier.example.com/esim/activate?region=us&token=ghi"
  }
}

action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "international_activation" {

  dynamic "managed_devices" {
    for_each = local.international_devices
    content {
      device_id   = managed_devices.key
      carrier_url = managed_devices.value
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Activate eSIM on co-managed device
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_comanaged" {

  comanaged_devices {
    device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
    carrier_url = "https://carrier.example.com/esim/activate?device=comanaged001"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Activate eSIM on Windows devices with cellular modems
data "microsoft365_graph_beta_device_management_managed_device" "windows_cellular" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and model contains 'LTE'"
}

action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_windows_cellular" {

  dynamic "managed_devices" {
    for_each = { for device in data.microsoft365_graph_beta_device_management_managed_device.windows_cellular.items : device.id => device }
    content {
      device_id   = managed_devices.key
      carrier_url = format("https://carrier.example.com/esim/activate?imei=%s", managed_devices.value.imei)
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Output examples
output "devices_activated_count" {
  value       = length(action.activate_multiple.managed_devices)
  description = "Number of devices that had eSIM activation initiated"
}

output "activation_summary" {
  value = {
    managed   = length(action.bulk_activation.managed_devices)
    comanaged = length(action.activate_comanaged.comanaged_devices)
  }
  description = "Count of eSIM activations by device type"
}

# Important Notes:
# eSIM Activation Features:
# - Enables cellular connectivity without physical SIM cards
# - Supports multiple carrier profiles on compatible devices
# - Allows remote provisioning and carrier switching
# - Each device requires carrier-specific activation URL
# - Device must have eSIM hardware capability
# - Carrier must support eSIM technology
#
# When to Activate eSIM:
# - Initial deployment of eSIM-capable devices
# - Switching carriers on existing devices
# - International deployments with local carriers
# - Field devices requiring remote cellular setup
# - Corporate devices needing managed connectivity
# - Devices supporting dual SIM (physical + eSIM)
#
# What Happens When eSIM is Activated:
# - Device receives carrier activation profile
# - eSIM downloads and installs carrier settings
# - Device activates cellular service
# - eSIM appears in device settings
# - Device can connect to carrier network
# - May require device restart on some platforms
# - Physical SIM (if present) remains functional
#
# Platform-Specific Support:
# - iOS/iPadOS: iPhone XS and later, cellular iPad models
# - Windows: Surface Pro X, Surface Pro 9 5G, other eSIM-capable PCs
# - Android: Varies by manufacturer and model
# - Must verify device eSIM capability before activation
# - Check carrier eSIM support for target region
#
# Carrier URL Requirements:
# - Provided by mobile carrier or MVNO
# - Format varies by carrier
# - May include activation tokens or codes
# - Often time-limited or single-use
# - Secure URL for profile download
# - Contains encrypted activation profile
#
# Best Practices:
# - Verify device eSIM capability before activation
# - Obtain valid carrier URLs before deployment
# - Test activation on pilot devices first
# - Coordinate with carrier for bulk activations
# - Document carrier URLs and activation dates
# - Plan for activation failures and retries
# - Consider time zones for international deployments
#
# eSIM vs Physical SIM:
# - No physical card required
# - Remote provisioning and management
# - Faster deployment at scale
# - Easier carrier switching
# - Supports multiple profiles (eSIM + physical)
# - Reduced logistics and shipping costs
# - Better for international deployments
#
# Device Requirements:
# - eSIM-capable hardware
# - Supported cellular modem
# - Compatible with carrier's network
# - Proper Intune enrollment
# - Online connectivity for activation
# - Sufficient storage for eSIM profile
#
# Carrier Considerations:
# - Must support eSIM technology
# - Provides activation URLs or QR codes
# - May have regional restrictions
# - Different pricing than physical SIM
# - Support for device management platforms
# - Activation process varies by carrier
#
# Troubleshooting:
# - Verify device eSIM support
# - Check carrier URL validity
# - Ensure device is online
# - Verify carrier network availability
# - Check for device restrictions
# - Review carrier activation logs
# - Contact carrier support if needed
#
# Security Considerations:
# - Carrier URLs may contain sensitive tokens
# - Store URLs securely
# - Use HTTPS for all carrier communications
# - Limit URL exposure and sharing
# - Monitor for unauthorized activations
# - Requires appropriate Intune permissions
# - Audit eSIM activation events
#
# Related Actions:
# - Device enrollment: Initial device setup
# - Network configuration: VPN and WiFi settings
# - Compliance policies: Ensure device requirements
# - Inventory management: Track eSIM-capable devices
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-activatedeviceesim?view=graph-rest-beta
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Block List) List of co-managed devices to activate eSIM on. These are devices managed by both Intune and Configuration Manager (SCCM). Each entry specifies a device ID and the carrier activation URL.

**Examples:**
```hcl
comanaged_devices = [
  {
    device_id   = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    carrier_url = "https://carrier.example.com/esim/activate?code=xyz789"
  }
]
```

**Platform Support:** Windows 10/11 with cellular modems (primary), limited iOS/Android support

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. Device must be online and support eSIM technology. (see [below for nested schema](#nestedblock--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some devices fail eSIM activation. Failed devices will be reported as warnings instead of errors. Default: `false` (action fails if any device fails).
- `managed_devices` (Block List) List of managed devices to activate eSIM on. These are devices fully managed by Intune only. Each entry specifies a device ID and the carrier-specific activation URL.

**Examples:**
```hcl
managed_devices = [
  {
    device_id   = "12345678-1234-1234-1234-123456789abc"
    carrier_url = "https://carrier.example.com/esim/activate?token=abc123"
  },
  {
    device_id   = "87654321-4321-4321-4321-987654321cba"
    carrier_url = "https://carrier.example.com/esim/activate?token=def456"
  }
]
```

**Platform Support:** iOS (iPhone XS+), Windows 10/11 with cellular, Android (varies by manufacturer)

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. Device must be online and support eSIM technology. (see [below for nested schema](#nestedblock--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist before attempting activation. Disabling this can speed up planning but may result in runtime errors for non-existent devices. Default: `true`.

<a id="nestedblock--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `carrier_url` (String) The carrier activation URL for this co-managed device. Example: `"https://carrier.example.com/esim/activate?code=xyz789"`
- `device_id` (String) The unique identifier (GUID) of the co-managed device to activate eSIM on. Example: `"12345678-1234-1234-1234-123456789abc"`


<a id="nestedblock--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `carrier_url` (String) The carrier-specific activation URL for this device's eSIM. This URL is provided by the mobile carrier and contains the activation profile. Format varies by carrier. Example: `"https://carrier.example.com/esim/activate?token=abc123"`
- `device_id` (String) The unique identifier (GUID) of the managed device to activate eSIM on. Device must have eSIM hardware capability. Example: `"12345678-1234-1234-1234-123456789abc"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
