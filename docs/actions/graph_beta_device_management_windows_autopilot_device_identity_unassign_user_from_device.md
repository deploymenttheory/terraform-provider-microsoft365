---
page_title: "<no value> <no value> - <no value>"
subcategory: "Device Management"

description: |-
  Unassigns a user from an Autopilot device in Microsoft Intune. This action removes the user assignment from Autopilot devices.
---

# <no value> (<no value>)

Unassigns a user from an Autopilot device in Microsoft Intune using the `/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/unassignUserFromDevice` endpoint. This action removes the user assignment from Autopilot devices.

## Microsoft Documentation

- [unassignUserFromDevice action](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-windowsautopilotdeviceidentity-unassignuserfromdevice?view=graph-rest-beta)
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
# Unassign user from device action
action "microsoft365_graph_beta_device_management_unassign_user_from_device" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789012"

  timeouts = {
    create = "5m"
  }
}
```