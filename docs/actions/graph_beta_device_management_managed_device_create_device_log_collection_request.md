---
page_title: "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Creates a device log collection request for Windows managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/createDeviceLogCollectionRequest and /deviceManagement/comanagedDevices/{managedDeviceId}/createDeviceLogCollectionRequest endpoints. This action is used to initiate diagnostic log collection from Windows devices for troubleshooting device issues, analyzing compliance problems, and supporting technical investigations. The collected logs are uploaded to Intune and can be downloaded for analysis. This action is critical for IT support teams when diagnosing device-specific problems or investigating security incidents.
  Important Notes:
  Only applicable to Windows devices (Windows 10/11)Device must be online to receive collection requestLog collection runs on the device and uploads resultsLogs are available in Intune portal after collection completesCollection includes system logs, event logs, and diagnostic dataLog files have expiration dates for security
  Use Cases:
  Troubleshooting device configuration issuesInvestigating compliance failures or policy problemsSupporting help desk tickets requiring detailed diagnosticsAnalyzing application deployment failuresSecurity incident investigation and forensicsProactive monitoring and preventive maintenance
  Platform Support:
  Windows: Fully supported (Windows 10 version 1709 or later, Windows 11)Other Platforms: Not supported (macOS, iOS/iPadOS, Android use different logging mechanisms)
  Reference: Microsoft Graph API - Create Device Log Collection Request https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-createdevicelogcollectionrequest?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request (Action)

Creates a device log collection request for Windows managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/createDeviceLogCollectionRequest` and `/deviceManagement/comanagedDevices/{managedDeviceId}/createDeviceLogCollectionRequest` endpoints. This action is used to initiate diagnostic log collection from Windows devices for troubleshooting device issues, analyzing compliance problems, and supporting technical investigations. The collected logs are uploaded to Intune and can be downloaded for analysis. This action is critical for IT support teams when diagnosing device-specific problems or investigating security incidents.

**Important Notes:**
- Only applicable to Windows devices (Windows 10/11)
- Device must be online to receive collection request
- Log collection runs on the device and uploads results
- Logs are available in Intune portal after collection completes
- Collection includes system logs, event logs, and diagnostic data
- Log files have expiration dates for security

**Use Cases:**
- Troubleshooting device configuration issues
- Investigating compliance failures or policy problems
- Supporting help desk tickets requiring detailed diagnostics
- Analyzing application deployment failures
- Security incident investigation and forensics
- Proactive monitoring and preventive maintenance

**Platform Support:**
- **Windows**: Fully supported (Windows 10 version 1709 or later, Windows 11)
- **Other Platforms**: Not supported (macOS, iOS/iPadOS, Android use different logging mechanisms)

**Reference:** [Microsoft Graph API - Create Device Log Collection Request](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-createdevicelogcollectionrequest?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [createDeviceLogCollectionRequest action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-createdevicelogcollectionrequest?view=graph-rest-beta)
- [deviceLogCollectionRequest resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicelogcollectionrequest?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Log Collection Guides
- [Collect diagnostics from a Windows device](https://learn.microsoft.com/en-us/mem/intune/remote-actions/collect-diagnostics)
- [Windows device diagnostics](https://learn.microsoft.com/en-us/mem/intune/fundamentals/collect-diagnostics)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this action:

**Required:**
- `DeviceManagementConfiguration.ReadWrite.All`
- `DeviceManagementManagedDevices.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |

## Example Usage

```terraform
# Example 1: Create device log collection request for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "create_single" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      }
    ]
  }
}

# Example 2: Create log collection requests for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "create_multiple" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      },
      {
        device_id     = "87654321-4321-4321-4321-ba9876543210"
        template_type = "custom"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Create with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "create_maximal" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      }
    ]

    comanaged_devices = [
      {
        device_id     = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        template_type = "predefined"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Create log collection requests for non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "noncompliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "create_noncompliant" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.noncompliant_devices.items : {
        device_id     = device.id
        template_type = "predefined"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "20m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed devices to collect logs from. These are Windows devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some devices fail log collection request. Failed devices will be reported as warnings instead of errors. Default: `false` (action fails if any device fails).
- `managed_devices` (Attributes List) List of managed devices to collect logs from. These are Windows devices fully managed by Intune only. Each device can have its own template type configuration.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices before attempting log collection. Disabling this can speed up planning but may result in runtime errors for non-existent or non-Windows devices. Default: `true`.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to collect logs from. This must be a Windows device running Windows 10 version 1709 or later, or Windows 11.

Example: `"abcdef12-3456-7890-abcd-ef1234567890"`

Optional:

- `template_type` (String) The template type for the log collection. Determines the scope and type of logs collected.

Valid values:
- `"predefined"` (default): Uses the standard predefined log collection template
- `"unknownFutureValue"`: Reserved for future expansion

If not specified, defaults to `"predefined"`.


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to collect logs from. This must be a Windows device running Windows 10 version 1709 or later, or Windows 11.

Example: `"12345678-1234-1234-1234-123456789abc"`

Optional:

- `template_type` (String) The template type for the log collection. Determines the scope and type of logs collected.

Valid values:
- `"predefined"` (default): Uses the standard predefined log collection template that includes common system and diagnostic logs
- `"unknownFutureValue"`: Reserved for future expansion

If not specified, defaults to `"predefined"`.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

