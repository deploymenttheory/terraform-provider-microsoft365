---
page_title: "<no value> <no value> - <no value>"
subcategory: "Device Management"

description: |-
  Allows the next enrollment for an Autopilot device in Microsoft Intune. This action enables the device to be enrolled again.
---

# <no value> (<no value>)

Allows the next enrollment for an Autopilot device in Microsoft Intune using the `/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/allowNextEnrollment` endpoint. This action enables the device to be enrolled again.

## Microsoft Documentation

- [allowNextEnrollment action](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-windowsautopilotdeviceidentity-allownextenrollment?view=graph-rest-beta)
- [windowsAutopilotDeviceIdentity resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-windowsautopilotdeviceidentity?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementServiceConfig.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.31.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Allow next enrollment for an Autopilot device
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_allow_next_enrollment" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789abc"

  timeouts = {
    create = "5m"
  }
}
```