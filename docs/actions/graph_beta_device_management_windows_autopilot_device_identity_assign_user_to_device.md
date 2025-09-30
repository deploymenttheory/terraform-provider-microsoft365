---
page_title: "<no value> <no value> - <no value>"
subcategory: "Device Management"

description: |-
  Assigns a user to an Autopilot device in Microsoft Intune. This action assigns user to Autopilot devices for streamlined device setup and management.
---

# <no value> (<no value>)

Assigns a user to an Autopilot device in Microsoft Intune using the `/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/assignUserToDevice` endpoint. This action assigns user to Autopilot devices for streamlined device setup and management.

## Microsoft Documentation

- [assignUserToDevice action](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-windowsautopilotdeviceidentity-assignusertodevice?view=graph-rest-beta)
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
# Assign user to device action
action "microsoft365_graph_beta_device_management_assign_user_to_device" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789012"
  user_principal_name                  = "user@contoso.com"
  addressable_user_name                = "John Doe"

  timeouts = {
    create = "5m"
  }
}
```