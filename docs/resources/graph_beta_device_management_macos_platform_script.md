---
page_title: "microsoft365_graph_beta_device_management_macos_platform_script Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages macOS shell scripts using the /deviceManagement/deviceShellScripts endpoint. macOS platform scripts enable automated deployment and execution of shell scripts on managed macOS devices with support for scheduled execution, retry logic, and execution context control for system maintenance and configuration tasks.
---

# microsoft365_graph_beta_device_management_macos_platform_script (Resource)

Manages macOS shell scripts using the `/deviceManagement/deviceShellScripts` endpoint. macOS platform scripts enable automated deployment and execution of shell scripts on managed macOS devices with support for scheduled execution, retry logic, and execution context control for system maintenance and configuration tasks.

## Microsoft Documentation

- [deviceShellScript resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-deviceshellscript?view=graph-rest-beta)
- [Create deviceShellScript](https://learn.microsoft.com/en-us/graph/api/intune-devices-deviceshellscript-create?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |

## Example Usage

```terraform
// Example: Device Shell Script Resource

resource "microsoft365_graph_beta_device_management_macos_platform_script" "example" {
  # Required fields
  display_name = "MacOS Shell Script"
  description  = "Example shell script for MacOS devices"

  script_content = <<EOT
    #!/bin/bash
    echo "Hello World"
  EOT

  run_as_account = "system" # Possible values: "system" or "user"
  file_name      = "example_script.sh"

  # Optional fields
  block_execution_notifications = false
  execution_frequency           = "P1D" # ISO 8601 duration format (e.g., P1D for 1 day, PT1H for 1 hour)
  retry_count                   = 3

  # Role scope tag IDs (optional)
  role_scope_tag_ids = ["0"]

  # Optional: Assignments block
  assignments = [
    # Optional: inclusion group assignments
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
  ]

  # Timeouts configuration (optional)
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Name of the macOS Platform Script.
- `file_name` (String) Script file name.
- `run_as_account` (String) Indicates the type of execution context. Possible values are: `system`, `user`.
- `script_content` (String, Sensitive) The script content.

### Optional

- `assignments` (Attributes Set) Assignments for the Windows remediation script. Each assignment specifies the target group and schedule for script execution. (see [below for nested schema](#nestedatt--assignments))
- `block_execution_notifications` (Boolean) Does not notify the user a script is being executed.
- `description` (String) Optional description for the macOS Platform Script.
- `execution_frequency` (String) The interval for script to run in ISO 8601 duration format (e.g., PT1H for 1 hour, P1D for 1 day). If not defined the script will run once.
- `retry_count` (Number) Number of times for the script to be retried if it fails.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier of the macOS platform script.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **macOS Shell Scripts**: This resource manages shell scripts for macOS devices in Microsoft Intune using the deviceShellScript API.
- **Script Execution**: Scripts are executed on macOS devices using the Intune management agent (Microsoft Intune Agent.app).
- **User vs Root Context**: Scripts can run as the current user or with root privileges depending on configuration.
- **Assignment Required**: Scripts must be assigned to device or user groups to be deployed.
- **Script Validation**: Intune provides execution status reporting and logs for troubleshooting.
- **Return Codes**: Scripts should use appropriate exit codes to indicate success or failure.
- **Security Context**: Scripts running with root privileges should be carefully reviewed for security implications.
- **File Size Limits**: Script files must be less than 1 MB in size.
- **Execution Frequency**: Scripts can be configured to run once or repeatedly based on schedule settings.
- **Platform Support**: Supports macOS 12.0 and later versions with Intune management agent installed.
- **Shebang Requirement**: Scripts must begin with a proper shebang (#!/bin/sh, #!/bin/bash, #!/usr/bin/env zsh).

## Import

Import is supported using the following syntax:

```shell
terraform import microsoft365_graph_beta_device_and_app_management_device_shell_script.example device-shell-script-id
```

