---
page_title: "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Triggers Configuration Manager client actions on Windows managed and co-managed devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/triggerConfigurationManagerAction and /deviceManagement/comanagedDevices/{managedDeviceId}/triggerConfigurationManagerAction endpoints. This action is used to remotely invoke specific Configuration Manager (SCCM) operations on devices that have the Configuration Manager client installed. This is particularly useful for co-managed devices where Intune and Configuration Manager work together to manage devices. Actions include policy refresh, application evaluation, antivirus scans, and more.
  Important Notes:
  Requires Configuration Manager client installed on devicePrimarily used for co-managed devices (Intune + Configuration Manager)Device must be online to receive the action triggerDifferent actions available for different management scenariosActions execute on the Configuration Manager client side
  Use Cases:
  Force policy refresh after configuration changesTrigger application deployment evaluationInitiate antivirus scans remotelyWake up clients for scheduled operationsUpdate Windows Defender signaturesSynchronize device state with Configuration Manager
  Platform Support:
  Windows: Fully supported (devices with Configuration Manager client)Other Platforms: Not supported (Configuration Manager is Windows-only)
  Reference: Microsoft Graph API - Trigger Configuration Manager Action https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-triggerconfigurationmanageraction?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action (Action)

Triggers Configuration Manager client actions on Windows managed and co-managed devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/triggerConfigurationManagerAction` and `/deviceManagement/comanagedDevices/{managedDeviceId}/triggerConfigurationManagerAction` endpoints. This action is used to remotely invoke specific Configuration Manager (SCCM) operations on devices that have the Configuration Manager client installed. This is particularly useful for co-managed devices where Intune and Configuration Manager work together to manage devices. Actions include policy refresh, application evaluation, antivirus scans, and more.

**Important Notes:**
- Requires Configuration Manager client installed on device
- Primarily used for co-managed devices (Intune + Configuration Manager)
- Device must be online to receive the action trigger
- Different actions available for different management scenarios
- Actions execute on the Configuration Manager client side

**Use Cases:**
- Force policy refresh after configuration changes
- Trigger application deployment evaluation
- Initiate antivirus scans remotely
- Wake up clients for scheduled operations
- Update Windows Defender signatures
- Synchronize device state with Configuration Manager

**Platform Support:**
- **Windows**: Fully supported (devices with Configuration Manager client)
- **Other Platforms**: Not supported (Configuration Manager is Windows-only)

**Reference:** [Microsoft Graph API - Trigger Configuration Manager Action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-triggerconfigurationmanageraction?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [triggerConfigurationManagerAction action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-triggerconfigurationmanageraction?view=graph-rest-beta)
- [configurationManagerAction resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-configurationmanageraction?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Configuration Manager and Co-Management Guides
- [Co-management for Windows devices](https://learn.microsoft.com/en-us/mem/configmgr/comanage/overview)
- [How to enable co-management](https://learn.microsoft.com/en-us/mem/configmgr/comanage/how-to-enable)
- [Co-management workloads](https://learn.microsoft.com/en-us/mem/configmgr/comanage/workloads)

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
# Example 1: Trigger Configuration Manager action on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "trigger_single" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
        action    = "refreshMachinePolicy"
      }
    ]
  }
}

# Example 2: Trigger multiple Configuration Manager actions
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "trigger_multiple" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
        action    = "refreshMachinePolicy"
      },
      {
        device_id = "87654321-4321-4321-4321-ba9876543210"
        action    = "refreshUserPolicy"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Trigger with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "trigger_maximal" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
        action    = "appEvaluation"
      }
    ]

    comanaged_devices = [
      {
        device_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        action    = "refreshMachinePolicy"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Trigger policy refresh on all co-managed devices
data "microsoft365_graph_beta_device_management_managed_device" "comanaged_devices" {
  filter_type  = "odata"
  odata_filter = "managementAgent eq 'configurationManagerClientMdm'"
}

action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "refresh_all_comanaged" {
  config {
    comanaged_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.comanaged_devices.items : {
        device_id = device.id
        action    = "refreshMachinePolicy"
      }
    ]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Attributes List) List of co-managed devices to trigger Configuration Manager actions on. These are Windows devices managed by both Intune and Configuration Manager (SCCM). This is the most common scenario for this action.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) When set to `true`, the action will complete successfully even if some devices fail to trigger Configuration Manager actions. When `false` (default), the action will fail if any device action trigger fails. Use this flag when triggering actions on multiple devices and you want the action to succeed even if some triggers fail.
- `managed_devices` (Attributes List) List of managed devices to trigger Configuration Manager actions on. These are Windows devices fully managed by Intune that also have the Configuration Manager client installed.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) When set to `true` (default), the action will validate that all specified devices exist and are Windows devices before attempting to trigger actions. When `false`, device validation is skipped and the action will attempt to trigger actions directly. Disabling validation can improve performance but may result in errors if devices don't exist or are not Windows devices.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `action` (String) The Configuration Manager action to trigger on this device.

Valid values:
- `"refreshMachinePolicy"`: Refresh the device's machine-level policies
- `"refreshUserPolicy"`: Refresh the current user's policies
- `"wakeUpClient"`: Wake up the Configuration Manager client
- `"appEvaluation"`: Trigger application deployment evaluation
- `"quickScan"`: Initiate a quick antivirus scan
- `"fullScan"`: Initiate a full antivirus scan
- `"windowsDefenderUpdateSignatures"`: Update Windows Defender signatures

Example: `"appEvaluation"`
- `device_id` (String) The unique identifier (GUID) of the co-managed device to trigger the action on. This must be a Windows device with Configuration Manager client installed.

Example: `"abcdef12-3456-7890-abcd-ef1234567890"`


<a id="nestedatt--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `action` (String) The Configuration Manager action to trigger on this device.

Valid values:
- `"refreshMachinePolicy"`: Refresh the device's machine-level policies from Configuration Manager
- `"refreshUserPolicy"`: Refresh the current user's policies from Configuration Manager
- `"wakeUpClient"`: Wake up the Configuration Manager client for immediate activity
- `"appEvaluation"`: Trigger application deployment evaluation cycle
- `"quickScan"`: Initiate a quick antivirus scan using Windows Defender
- `"fullScan"`: Initiate a full antivirus scan using Windows Defender
- `"windowsDefenderUpdateSignatures"`: Update Windows Defender antivirus signatures

Example: `"refreshMachinePolicy"`
- `device_id` (String) The unique identifier (GUID) of the managed device to trigger the action on. This must be a Windows device with Configuration Manager client installed.

Example: `"12345678-1234-1234-1234-123456789abc"`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

