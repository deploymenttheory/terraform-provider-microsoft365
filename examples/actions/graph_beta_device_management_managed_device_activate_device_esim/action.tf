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
  type = map(string)
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

