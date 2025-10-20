---
page_title: "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Creates a device log collection request for Windows managed devices using the /deviceManagement/managedDevices/{managedDeviceId}/createDeviceLogCollectionRequest and /deviceManagement/comanagedDevices/{managedDeviceId}/createDeviceLogCollectionRequest endpoints. This action initiates the collection of diagnostic logs from Windows devices, which are essential for troubleshooting device issues, analyzing compliance problems, and supporting technical investigations. The collected logs are uploaded to Intune and can be downloaded for analysis. This action is critical for IT support teams when diagnosing device-specific problems or investigating security incidents.
  Important Notes:
  Only applicable to Windows devices (Windows 10/11)Device must be online to receive collection requestLog collection runs on the device and uploads resultsLogs are available in Intune portal after collection completesCollection includes system logs, event logs, and diagnostic dataLog files have expiration dates for security
  Use Cases:
  Troubleshooting device configuration issuesInvestigating compliance failures or policy problemsSupporting help desk tickets requiring detailed diagnosticsAnalyzing application deployment failuresSecurity incident investigation and forensicsProactive monitoring and preventive maintenance
  Platform Support:
  Windows: Fully supported (Windows 10 version 1709 or later, Windows 11)Other Platforms: Not supported (macOS, iOS/iPadOS, Android use different logging mechanisms)
  Reference: Microsoft Graph API - Create Device Log Collection Request https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-createdevicelogcollectionrequest?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request (Action)

Creates a device log collection request for Windows managed devices using the `/deviceManagement/managedDevices/{managedDeviceId}/createDeviceLogCollectionRequest` and `/deviceManagement/comanagedDevices/{managedDeviceId}/createDeviceLogCollectionRequest` endpoints. This action initiates the collection of diagnostic logs from Windows devices, which are essential for troubleshooting device issues, analyzing compliance problems, and supporting technical investigations. The collected logs are uploaded to Intune and can be downloaded for analysis. This action is critical for IT support teams when diagnosing device-specific problems or investigating security incidents.

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

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`, `DeviceManagementManagedDevices.ReadWrite.All`
- **Delegated**: `DeviceManagementConfiguration.ReadWrite.All`, `DeviceManagementManagedDevices.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **Windows** | ✅ Full Support | Windows 10 version 1709 or later, Windows 11 |
| **macOS** | ❌ Not Supported | Different logging mechanism (use MDM logs) |
| **iOS/iPadOS** | ❌ Not Supported | Different logging mechanism (use iOS logs) |
| **Android** | ❌ Not Supported | Different logging mechanism (use Android logs) |

### What is Device Log Collection?

Device Log Collection is an action that:
- Initiates diagnostic log gathering from Windows devices
- Collects system logs, event logs, and diagnostic data
- Uploads collected logs to Intune for analysis
- Provides a response with collection status and details
- Makes logs available for download in the Intune portal
- Supports troubleshooting and investigation workflows

### When to Collect Device Logs

- Troubleshooting device configuration or policy issues
- Investigating compliance violations or policy failures
- Diagnosing application deployment problems
- Security incident investigation and forensic analysis
- Proactive monitoring of device health
- Support ticket escalation requiring detailed diagnostics
- Analyzing Windows update or patch deployment failures

### What Happens When Log Collection is Requested

- Intune sends log collection request to the Windows device
- Device gathers specified diagnostic logs based on template
- Logs are compressed into a ZIP file
- Device uploads logs to Intune storage
- Collection status is tracked and made available
- Logs become available for download in Intune portal
- Collection request has an expiration date
- Response includes collection ID, status, timestamps, and other details

### Log Collection Response

The action returns a response for each device with:
- **Collection ID**: Unique identifier for tracking the collection
- **Status**: Current state (pending, completed, failed)
- **Requested/Received Times**: Timestamps for tracking progress
- **Expiration Date**: When the collected logs will be deleted
- **Initiator**: User who initiated the collection
- **Size Information**: Size of collected logs (when available)
- **Error Code**: If collection fails, error code for troubleshooting

### Accessing Collected Logs

1. Navigate to the device in Microsoft Intune admin center
2. Select "Device diagnostics" or "Collect logs" from device actions
3. View collection status and download links
4. Download collected log files (typically ZIP format)
5. Extract and analyze using appropriate tools (Event Viewer, text editors, etc.)

## Example Usage

```terraform
# Example 1: Create log collection request for a single Windows device
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "collect_single" {

  managed_devices {
    device_id = "12345678-1234-1234-1234-123456789abc"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Collect logs from multiple Windows devices
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "collect_multiple" {

  managed_devices {
    device_id = "12345678-1234-1234-1234-123456789abc"
  }

  managed_devices {
    device_id     = "87654321-4321-4321-4321-ba9876543210"
    template_type = "predefined"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Troubleshooting help desk tickets
variable "troubleshooting_device_ids" {
  description = "Device IDs requiring log collection for support tickets"
  type        = list(string)
  default = [
    "aaaa1111-1111-1111-1111-111111111111",
    "bbbb2222-2222-2222-2222-222222222222"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "helpdesk_diagnostics" {

  dynamic "managed_devices" {
    for_each = var.troubleshooting_device_ids
    content {
      device_id = managed_devices.value
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Collect logs from all non-compliant Windows devices
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "compliance_investigation" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.non_compliant_windows.items
    content {
      device_id     = managed_devices.value.id
      template_type = "predefined"
    }
  }

  timeouts = {
    invoke = "20m"
  }
}

# Example 5: Collect logs for security incident investigation
locals {
  suspected_compromised_devices = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "security_investigation" {

  dynamic "managed_devices" {
    for_each = local.suspected_compromised_devices
    content {
      device_id = managed_devices.value
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Collect logs from co-managed device
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "comanaged_diagnostics" {

  comanaged_devices {
    device_id     = "abcdef12-3456-7890-abcd-ef1234567890"
    template_type = "predefined"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Collect logs from devices in specific Azure AD group
data "microsoft365_graph_beta_device_management_managed_device" "finance_windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and azureADDeviceId ne null"
}

locals {
  finance_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.finance_windows_devices.items : device.id]
}

action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "department_diagnostics" {

  dynamic "managed_devices" {
    for_each = local.finance_device_ids
    content {
      device_id = managed_devices.value
    }
  }

  timeouts = {
    invoke = "20m"
  }
}

# Example 8: Proactive monitoring - collect logs from devices with recent errors
data "microsoft365_graph_beta_device_management_managed_device" "devices_with_errors" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and lastSyncDateTime gt 2025-01-01T00:00:00Z"
}

action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "proactive_monitoring" {

  dynamic "managed_devices" {
    for_each = { for device in data.microsoft365_graph_beta_device_management_managed_device.devices_with_errors.items : device.id => device }
    content {
      device_id     = managed_devices.key
      template_type = "predefined"
    }
  }

  timeouts = {
    invoke = "30m"
  }
}

# Output examples
output "log_collection_summary" {
  value = {
    managed_devices   = length(action.collect_multiple.managed_devices)
    comanaged_devices = length(action.comanaged_diagnostics.comanaged_devices)
  }
  description = "Count of devices with log collection requests initiated"
}

# Important Notes:
# Device Log Collection Features:
# - Collects comprehensive diagnostic logs from Windows devices
# - Logs uploaded to Intune and available in portal
# - Includes system logs, event logs, and diagnostic data
# - Log files have expiration dates
# - Essential for troubleshooting and investigations
# - Supports predefined log collection templates
#
# When to Collect Device Logs:
# - Troubleshooting device configuration issues
# - Investigating policy application failures
# - Diagnosing app deployment problems
# - Security incident investigation
# - Compliance violation analysis
# - Proactive device health monitoring
# - Support ticket escalation requiring detailed diagnostics
#
# What Happens When Log Collection is Requested:
# - Device receives collection request from Intune
# - Device gathers specified diagnostic logs
# - Logs are compressed and uploaded to Intune
# - Collection status tracked in response
# - Logs available for download in Intune portal
# - Collection request has expiration date
# - User is typically not notified of collection
#
# Platform Requirements:
# - Windows: Fully supported (Windows 10 version 1709+, Windows 11)
# - Device must be enrolled in Intune
# - Device must be online to receive request
# - Other platforms: Not supported for this API
#
# Template Types:
# - predefined (default): Standard log collection template
#   - Includes common system and diagnostic logs
#   - Event logs from key sources
#   - Configuration data
#   - Policy application logs
# - unknownFutureValue: Reserved for future expansion
#
# Best Practices:
# - Only collect logs when necessary
# - Document business justification
# - Consider data privacy implications
# - Review collected logs promptly
# - Delete logs after analysis
# - Monitor collection success rates
# - Set appropriate timeouts for large deployments
#
# Log Collection Response:
# - Collection ID for tracking
# - Status (pending, completed, failed)
# - Requested and received timestamps
# - Expiration date/time
# - Initiator information
# - File size information
# - Error codes if collection fails
#
# Accessing Collected Logs:
# - Navigate to device in Intune portal
# - View "Device diagnostics" or "Collect logs"
# - Download collected log files
# - Logs typically in ZIP format
# - Extract and analyze using appropriate tools
# - Event logs viewable in Event Viewer
#
# Log Contents:
# - Windows Event Logs (System, Application, Security)
# - MDM diagnostic logs
# - Group Policy logs
# - Certificate information
# - Network configuration
# - Installed applications list
# - Device hardware info
# - Policy application status
#
# Privacy and Security:
# - Logs may contain sensitive information
# - User activity may be logged
# - Follow data protection regulations
# - Secure access to downloaded logs
# - Document log access and usage
# - Delete logs when no longer needed
# - Audit log collection activities
#
# Troubleshooting:
# - Verify device is Windows 10/11
# - Check device is online and syncing
# - Ensure device has sufficient storage
# - Verify Intune connectivity
# - Review collection status in portal
# - Check for error codes in response
# - Allow time for collection to complete
#
# Collection Status Values:
# - pending: Request sent, collection not started
# - completed: Logs collected and uploaded
# - failed: Collection encountered error
# - Check managedDeviceId in response
# - Monitor expirationDateTimeUTC
# - Review errorCode if present
#
# Common Use Cases:
# - App installation failures
# - Policy not applying correctly
# - Device compliance issues
# - Network connectivity problems
# - Certificate errors
# - Windows update failures
# - Security baseline violations
# - Performance degradation
#
# Limitations:
# - Windows devices only
# - Requires online connectivity
# - Storage space on device needed
# - Collection may take time
# - Log files have size limits
# - Logs expire after period
# - Some logs require admin rights
#
# Related Actions:
# - Device sync: Ensure device is current
# - Remote lock: Secure device during investigation
# - Wipe: If security breach confirmed
# - Compliance checks: Review policy status
#
# Integration Points:
# - Help desk ticketing systems
# - SIEM for security analysis
# - Monitoring and alerting platforms
# - Compliance reporting tools
# - Automation workflows
#
# Reference:
# https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-createdevicelogcollectionrequest?view=graph-rest-beta
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Block List) List of co-managed devices to collect logs from. These are Windows devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--comanaged_devices))
- `managed_devices` (Block List) List of managed devices to collect logs from. These are Windows devices fully managed by Intune only. Each device can have its own template type configuration.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedblock--comanaged_devices"></a>
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


<a id="nestedblock--managed_devices"></a>
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

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

