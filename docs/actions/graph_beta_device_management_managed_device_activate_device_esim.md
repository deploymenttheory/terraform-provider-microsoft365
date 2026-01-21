---
page_title: "microsoft365_graph_beta_device_management_managed_device_activate_device_esim Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Activates eSIM cellular data plans on iOS and iPadOS devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/activateDeviceEsim and /deviceManagement/comanagedDevices/{managedDeviceId}/activateDeviceEsim endpoints. This action is used to remotely activate eSIM cellular plans without physical SIM cards, making it easier to manage connectivity for users.
---

# microsoft365_graph_beta_device_management_managed_device_activate_device_esim (Action)

Activates eSIM cellular data plans on iOS and iPadOS devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/activateDeviceEsim` and `/deviceManagement/comanagedDevices/{managedDeviceId}/activateDeviceEsim` endpoints. This action is used to remotely activate eSIM cellular plans without physical SIM cards, making it easier to manage connectivity for users.

## Microsoft Documentation

### Graph API References
- [activateDeviceEsim action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-activatedeviceesim?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions
- [Update cellular data plan (eSIM activation)](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/update-cellular-data-plan)

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
# Example 1: Activate eSIM on a single device
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_single" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        carrier_url = "https://carrier.example.com/esim/activate?token=abc123xyz"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 2: Activate eSIM on multiple devices with different carriers
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_multiple" {
  config {
    managed_devices = [
      {
        device_id   = "12345678-1234-1234-1234-123456789abc"
        carrier_url = "https://carrier-a.example.com/esim/activate?code=device1"
      },
      {
        device_id   = "87654321-4321-4321-4321-ba9876543210"
        carrier_url = "https://carrier-b.example.com/esim/activate?code=device2"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
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
  config {
    managed_devices = [
      for device_id, carrier_url in var.new_devices_with_esim : {
        device_id   = device_id
        carrier_url = carrier_url
      }
    ]

    timeouts = {
      invoke = "15m"
    }
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
  config {
    managed_devices = [
      for device_id, carrier_url in local.ipad_carrier_urls : {
        device_id   = device_id
        carrier_url = carrier_url
      }
    ]

    timeouts = {
      invoke = "20m"
    }
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
  config {
    managed_devices = [
      for device_id, carrier_url in local.international_devices : {
        device_id   = device_id
        carrier_url = carrier_url
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 6: Activate eSIM on co-managed device
action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_comanaged" {
  config {
    comanaged_devices = [
      {
        device_id   = "abcdef12-3456-7890-abcd-ef1234567890"
        carrier_url = "https://carrier.example.com/esim/activate?device=comanaged001"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 7: Activate eSIM on Windows devices with cellular modems
data "microsoft365_graph_beta_device_management_managed_device" "windows_cellular" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and model contains 'LTE'"
}

action "microsoft365_graph_beta_device_management_managed_device_activate_device_esim" "activate_windows_cellular" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.windows_cellular.items : {
        device_id   = device.id
        carrier_url = format("https://carrier.example.com/esim/activate?imei=%s", device.imei)
      }
    ]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Output examples - demonstrating how to reference action configuration
output "devices_activated_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_activate_device_esim.activate_multiple.config.managed_devices)
  description = "Number of devices that had eSIM activation initiated"
}

output "activation_summary" {
  value = {
    managed   = length(action.microsoft365_graph_beta_device_management_managed_device_activate_device_esim.bulk_activation.config.managed_devices)
    comanaged = length(action.microsoft365_graph_beta_device_management_managed_device_activate_device_esim.activate_comanaged.config.comanaged_devices)
  }
  description = "Count of eSIM activations by device type"
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of iOS/iPadOS co-managed devices to activate eSIM on. These are devices managed by both Intune and Configuration Manager (SCCM). Devices must have eSIM hardware capability (iPhone XS and later, cellular iPads with eSIM support). (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some devices fail eSIM activation. Failed devices will be reported as warnings instead of errors. Default: `false` (action fails if any device fails).
- `managed_devices` (Attributes List) List of iOS/iPadOS managed devices to activate eSIM on. These are devices fully managed by Intune only. Devices must have eSIM hardware capability (iPhone XS and later, cellular iPads with eSIM support). (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist before attempting activation. Disabling this can speed up planning but may result in runtime errors for non-existent devices. Default: `true`.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `carrier_url` (String) The activation server URL provided by your mobile carrier for eSIM activation. This URL is carrier-specific and contains the activation profile. Example: `"https://carrier.example.com/esim/activate?code=xyz789"`
- `device_id` (String) The unique identifier (GUID) of the iOS/iPadOS co-managed device to activate eSIM on. Device must have eSIM hardware capability (iPhone XS+, cellular iPad with eSIM). Example: `"12345678-1234-1234-1234-123456789abc"`


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `carrier_url` (String) The activation server URL provided by your mobile carrier for eSIM activation. This URL is carrier-specific and contains the activation profile. Example: `"https://carrier.example.com/esim/activate?token=abc123"`
- `device_id` (String) The unique identifier (GUID) of the iOS/iPadOS device to activate eSIM on. Device must have eSIM hardware capability (iPhone XS+, cellular iPad with eSIM). Example: `"12345678-1234-1234-1234-123456789abc"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
