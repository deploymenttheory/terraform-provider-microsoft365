---
page_title: "microsoft365_graph_beta_device_management_windows_remediation_script Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Windows remediation scripts using the /deviceManagement/deviceHealthScripts endpoint. Remediation scripts enable proactive detection and automatic remediation of common issues on Windows devices through PowerShell detection scripts paired with remediation scripts that execute when problems are identified.
---

# microsoft365_graph_beta_device_management_windows_remediation_script (Resource)

Manages Windows remediation scripts using the `/deviceManagement/deviceHealthScripts` endpoint. Remediation scripts enable proactive detection and automatic remediation of common issues on Windows devices through PowerShell detection scripts paired with remediation scripts that execute when problems are identified.

## Microsoft Documentation

- [deviceHealthScript resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicehealthscript?view=graph-rest-beta)
- [Create deviceHealthScript](https://learn.microsoft.com/en-us/graph/api/intune-devices-devicehealthscript-create?view=graph-rest-beta)

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
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "basic_example" {
  display_name            = "windows remediation script with no assignments"
  description             = "Simple script applied to no devices"
  publisher               = "Contoso IT"
  run_as_account          = "system"
  run_as_32_bit           = false
  enforce_signature_check = true
  role_scope_tag_ids      = [8, 9] # Optional

  detection_script_content = <<-EOT
    # Detection script logic
    if (Test-Path "C:\Temp\issues.txt") {
      Write-Host "Issue detected"
      Exit 1
    } else {
      Write-Host "No issues found"
      Exit 0
    }
  EOT

  remediation_script_content = <<-EOT
    # Remediation script logic
    Remove-Item "C:\Temp\issues.txt" -Force
    Write-Host "Issue remediated"
    Exit 0
  EOT

  assignment {
    all_devices = false
    all_users   = false
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "microsoft365_graph_beta_device_management_windows_remediation_script" "example_with_filters" {
  display_name            = "windows remediation script with no assignments"
  description             = "Simple script applied to no devices"
  publisher               = "Contoso IT"
  run_as_account          = "system"
  run_as_32_bit           = false
  enforce_signature_check = true
  role_scope_tag_ids      = [8, 9] # Optional

  detection_script_content = <<-EOT
    # Detection script logic
    if (Test-Path "C:\Temp\issues.txt") {
      Write-Host "Issue detected"
      Exit 1
    } else {
      Write-Host "No issues found"
      Exit 0
    }
  EOT

  remediation_script_content = <<-EOT
    # Remediation script logic
    Remove-Item "C:\Temp\issues.txt" -Force
    Write-Host "Issue remediated"
    Exit 0
  EOT

  assignment {
    all_devices             = true
    all_devices_filter_type = "include"
    all_devices_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"

    all_users             = true
    all_users_filter_type = "exclude"
    all_users_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}


resource "microsoft365_graph_beta_device_management_windows_remediation_script" "windows_remediation_script_with_scoping" {
  display_name            = "windows remediation script with assignment options"
  description             = "Simple script applied to scoped devices"
  publisher               = "Contoso IT"
  run_as_account          = "system"
  run_as_32_bit           = false
  enforce_signature_check = true
  role_scope_tag_ids      = [8, 9] # Optional

  detection_script_content = <<-EOT
    # Detection script logic
    if (Test-Path "C:\Temp\issues.txt") {
      Write-Host "Issue detected"
      Exit 1
    } else {
      Write-Host "No issues found"
      Exit 0
    }
  EOT

  remediation_script_content = <<-EOT
    # Remediation script logic
    Remove-Item "C:\Temp\issues.txt" -Force
    Write-Host "Issue remediated"
    Exit 0
  EOT

  assignment {
    all_devices = false
    all_users   = false

    include_groups = [
      {
        group_id                   = "bae7a85a-8284-4f58-9873-a84bd4d22585"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "once"
          date          = "2025-05-01"
          time          = "14:30"
          use_utc       = true
        }
      },
      {
        group_id                   = "6117fcd2-2812-44b2-a0d7-3c57ca81c015"
        include_groups_filter_type = "include"
        include_groups_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "daily"
          interval      = "1"
          time          = "14:30"
          use_utc       = true
        }
      },
      {
        group_id                   = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
        include_groups_filter_type = "exclude"
        include_groups_filter_id   = "43cb3789-2d36-4fb6-aa4d-0c3678b064e7"
        run_remediation_script     = true
        run_schedule = {
          schedule_type = "hourly"
          interval      = "1"
        }
      }
    ]

    exclude_group_ids = [
      "b8c661c2-fa9a-4351-af86-adc1729c343f",
      "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f"
    ]
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `detection_script_content` (String, Sensitive) The entire content of the detection PowerShell script.
- `display_name` (String) Name of the device health script.
- `publisher` (String) Name of the device health script publisher.
- `remediation_script_content` (String, Sensitive) The entire content of the remediation PowerShell script.
- `run_as_account` (String) Indicates the type of execution context. Possible values are: system, user.

### Optional

- `assignment` (Block List) List of assignment configurations for the device health script (see [below for nested schema](#nestedblock--assignment))
- `description` (String) Description of the device health script.
- `detection_script_parameters` (Attributes List) List of ComplexType DetectionScriptParameters objects. (see [below for nested schema](#nestedatt--detection_script_parameters))
- `enforce_signature_check` (Boolean) Indicate whether the script signature needs be checked.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `run_as_32_bit` (Boolean) Indicate whether PowerShell script(s) should run as 32-bit.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The timestamp of when the device health script was created. This property is read-only.
- `device_health_script_type` (String) DeviceHealthScriptType for the script policy. Possible values are: deviceHealthScript, managedInstallerScript.
- `highest_available_version` (String) Highest available version for a Microsoft Proprietary script.
- `id` (String) Unique identifier for the device health script.
- `is_global_script` (Boolean) Determines if this is Microsoft Proprietary Script. Proprietary scripts are read-only.
- `last_modified_date_time` (String) The timestamp of when the device health script was modified. This property is read-only.
- `version` (String) Version of the device health script.

<a id="nestedblock--assignment"></a>
### Nested Schema for `assignment`

Optional:

- `all_devices` (Boolean) Assign to all devices. Cannot be used with all_users or include_groups.
- `all_devices_filter_id` (String) Filter ID for all devices assignment.
- `all_devices_filter_type` (String) Filter type for all devices assignment. Can be 'include' or 'exclude'.
- `all_users` (Boolean) Assign to all users. Cannot be used with all_devices or include_groups.
- `all_users_filter_id` (String) Filter ID for all users assignment.
- `all_users_filter_type` (String) Filter type for all users assignment. Can be 'include' or 'exclude'.
- `exclude_group_ids` (Set of String) Group IDs to exclude from the assignment.
- `include_groups` (Attributes Set) Groups to include in the assignment. Cannot be used with all_devices or all_users. (see [below for nested schema](#nestedatt--assignment--include_groups))

<a id="nestedatt--assignment--include_groups"></a>
### Nested Schema for `assignment.include_groups`

Required:

- `group_id` (String) Group ID to include.

Optional:

- `include_groups_filter_id` (String) Filter ID for include group assignment.
- `include_groups_filter_type` (String) Filter type for include group assignment. Can be 'include' or 'exclude'.
- `run_remediation_script` (Boolean) Whether to run the remediation script for this group assignment.
- `run_schedule` (Attributes) Run schedule for this group assignment. (see [below for nested schema](#nestedatt--assignment--include_groups--run_schedule))

<a id="nestedatt--assignment--include_groups--run_schedule"></a>
### Nested Schema for `assignment.include_groups.run_schedule`

Required:

- `schedule_type` (String) Type of schedule. Can be 'daily', 'hourly', or 'once'.

Optional:

- `date` (String) Date for once schedule (e.g., '2025-05-01').
- `interval` (Number) Repeat interval for the schedule.For 'daily' the interal represents days, for 'hourly' the interval represents hours.
- `time` (String) Time of day for daily and once schedules (e.g., '14:30').
- `use_utc` (Boolean) Whether to use UTC time.




<a id="nestedatt--detection_script_parameters"></a>
### Nested Schema for `detection_script_parameters`

Required:

- `name` (String) The name of the param

Optional:

- `apply_default_value_when_not_assigned` (Boolean) Whether Apply DefaultValue When Not Assigned
- `description` (String) The description of the param
- `is_required` (Boolean) Whether the param is required


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_management_windows_remediation_script.example windows-remediation-script-id
```