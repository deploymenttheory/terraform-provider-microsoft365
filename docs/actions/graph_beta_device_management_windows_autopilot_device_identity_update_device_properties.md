---
page_title: "<no value> <no value> - <no value>"
subcategory: "Device Management"

description: |-
  Updates properties on an Autopilot device in Microsoft Intune. This action allows updating various properties of Autopilot devices including user assignment, group tag, and display name.
---

# <no value> (<no value>)

Updates properties on an Autopilot device in Microsoft Intune using the `/deviceManagement/windowsAutopilotDeviceIdentities/{windowsAutopilotDeviceIdentityId}/updateDeviceProperties` endpoint. This action allows updating various properties of Autopilot devices including user assignment, group tag, and display name.

## Microsoft Documentation

- [updateDeviceProperties action](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-windowsautopilotdeviceidentity-updatedeviceproperties?view=graph-rest-beta)
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
# Update device properties action
action "microsoft365_graph_beta_device_management_update_device_properties" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789012"
  user_principal_name                  = "user@contoso.com"
  addressable_user_name                = "John Doe"
  group_tag                            = "Finance"
  display_name                         = "John's Laptop"

  timeouts = {
    create = "5m"
  }
}
```