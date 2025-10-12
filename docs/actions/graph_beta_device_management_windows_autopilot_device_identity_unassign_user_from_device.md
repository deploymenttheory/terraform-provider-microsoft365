---
page_title: "<no value> <no value> - <no value>"
subcategory: "Device Management"

description: |-
  Retrieves audit events from Microsoft 365 managed tenants as an ephemeral resource.
---

# <no value> (<no value>)

Retrieves audit events from Microsoft 365 managed tenants as an ephemeral resource. This does not persist in state and fetches fresh data on each execution.

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
action "microsoft365_graph_beta_device_management_windows_autopilot_device_identity_unassign_user_from_device" "example" {
  windows_autopilot_device_identity_id = "12345678-1234-1234-1234-123456789012"

  timeouts = {
    create = "5m"
  }
}
```

Schema documentation not available.