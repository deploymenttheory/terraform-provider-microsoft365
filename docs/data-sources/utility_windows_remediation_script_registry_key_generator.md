---
page_title: "microsoft365_utility_windows_remediation_script_registry_key_generator Data Source - terraform-provider-microsoft365"
subcategory: "Utility"

description: |-
  Generates PowerShell detection and remediation scripts for Intune Proactive Remediations to manage Windows registry keys and values. This utility helps create standardized scripts for setting registry values in either the current user's context (HKEY_CURRENT_USER) or for all users on a device (HKEY_USERS). The generated scripts follow Microsoft's recommended patterns for Proactive Remediations, including proper error handling, exit codes, and user context management.
---

# microsoft365_utility_windows_remediation_script_registry_key_generator

Generates PowerShell detection and remediation scripts for managing Windows registry keys and values through Intune Proactive Remediations.

This utility data source creates production-ready PowerShell scripts that integrate directly with the `microsoft365_graph_beta_device_management_windows_remediation_script` resource. The generated scripts follow Microsoft's best practices for Proactive Remediations, including proper exit codes, error handling, and user context management.

## Typical Workflow

1. **Define** your registry configuration using this data source
2. **Generate** PowerShell detection and remediation scripts automatically
3. **Deploy** the scripts to Intune using the `microsoft365_graph_beta_device_management_windows_remediation_script` resource
4. **Monitor** compliance through Intune's Proactive Remediations reporting

## Key Features

- **Automated Script Generation**: No need to write PowerShell - define your registry settings declaratively
- **Direct Intune Integration**: Output attributes map directly to the Windows remediation script resource
- **Dual Context Support**: Handles both current user (HKCU) and all users (HKU) scenarios
- **Comprehensive Type Support**: All 6 common Windows registry value types (String, DWORD, QWORD, Multi-String, ExpandString, Binary)
- **Best Practice Compliance**: Includes proper exit codes (0/1), system account exclusions, and error handling
- **Infrastructure as Code**: Manage registry policies alongside other Intune configurations

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.35.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example usage of the Intune Proactive Remediation Registry Script Generator data source
# This data source generates PowerShell detection and remediation scripts for managing Windows registry keys via Intune

# Example 1: Generate scripts and deploy to Intune as a Proactive Remediation
data "microsoft365_utility_windows_remediation_script_registry_key_generator" "private_store" {
  context           = "current_user"
  registry_key_path = "Software\\Policies\\Microsoft\\WindowsStore\\"
  value_name        = "RequirePrivateStoreOnly"
  value_type        = "REG_DWORD"
  value_data        = "1"
}

# Deploy the generated scripts to Intune
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "private_store_remediation" {
  display_name            = "Enforce Private Store Only"
  description             = "Ensures Windows Store is configured to show only the private store"
  publisher               = "IT Security Team"
  run_as_32_bit           = false
  enforce_signature_check = false
  role_scope_tag_ids      = ["0"]
  run_as_account          = "user" // Use "user" context for current_user registry scripts

  # Use the generated detection script
  detection_script_content = data.microsoft365_utility_windows_remediation_script_registry_key_generator.private_store.detection_script

  # Use the generated remediation script
  remediation_script_content = data.microsoft365_utility_windows_remediation_script_registry_key_generator.private_store.remediation_script

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

      daily_schedule = {
        interval = 1
        time     = "09:00:00"
        use_utc  = true
      }
    }
  ]

  timeouts = {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}

# Example 2: Set registry value for all users with SYSTEM context
data "microsoft365_utility_windows_remediation_script_registry_key_generator" "disable_cortana" {
  context           = "all_users"
  registry_key_path = "Software\\Policies\\Microsoft\\Windows\\Windows Search\\"
  value_name        = "AllowCortana"
  value_type        = "REG_DWORD"
  value_data        = "0"
}

resource "microsoft365_graph_beta_device_management_windows_remediation_script" "disable_cortana" {
  display_name            = "Disable Cortana for All Users"
  description             = "Disables Cortana via registry for all user profiles"
  publisher               = "IT Security Team"
  run_as_32_bit           = false
  enforce_signature_check = false
  role_scope_tag_ids      = ["0"]
  run_as_account          = "system" // Use "system" context for all_users registry scripts

  detection_script_content   = data.microsoft365_utility_windows_remediation_script_registry_key_generator.disable_cortana.detection_script
  remediation_script_content = data.microsoft365_utility_windows_remediation_script_registry_key_generator.disable_cortana.remediation_script

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

      hourly_schedule = {
        interval = 4
      }
    }
  ]
}

# Example 3: String value with environment variables
data "microsoft365_utility_windows_remediation_script_registry_key_generator" "app_install_path" {
  context           = "current_user"
  registry_key_path = "Software\\MyApp\\"
  value_name        = "InstallPath"
  value_type        = "REG_EXPAND_SZ"
  value_data        = "%ProgramFiles%\\MyApp"
}

resource "microsoft365_graph_beta_device_management_windows_remediation_script" "app_install_path" {
  display_name               = "Configure App Install Path"
  description                = "Sets the application install path with environment variable expansion"
  publisher                  = "IT Application Team"
  run_as_32_bit              = false
  enforce_signature_check    = false
  role_scope_tag_ids         = ["0"]
  run_as_account             = "user"
  detection_script_content   = data.microsoft365_utility_windows_remediation_script_registry_key_generator.app_install_path.detection_script
  remediation_script_content = data.microsoft365_utility_windows_remediation_script_registry_key_generator.app_install_path.remediation_script

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"
      run_once_schedule = {
        date    = "2024-12-31"
        time    = "10:00:00"
        use_utc = true
      }
    }
  ]
}

# Example 4: Deploy multiple policies using for_each
variable "registry_policies" {
  description = "Map of registry policies to configure"
  type = map(object({
    display_name      = string
    description       = string
    context           = string
    registry_key_path = string
    value_name        = string
    value_type        = string
    value_data        = string
    run_as_account    = string
  }))
  default = {
    "disable_consumer_features" = {
      display_name      = "Disable Windows Consumer Features"
      description       = "Prevents Windows from installing consumer apps"
      context           = "all_users"
      registry_key_path = "Software\\Policies\\Microsoft\\Windows\\CloudContent\\"
      value_name        = "DisableWindowsConsumerFeatures"
      value_type        = "REG_DWORD"
      value_data        = "1"
      run_as_account    = "system"
    }
    "disable_lockscreen_tips" = {
      display_name      = "Disable Lock Screen Tips"
      description       = "Disables tips and tricks on the lock screen"
      context           = "all_users"
      registry_key_path = "Software\\Microsoft\\Windows\\CurrentVersion\\ContentDeliveryManager\\"
      value_name        = "RotatingLockScreenEnabled"
      value_type        = "REG_DWORD"
      value_data        = "0"
      run_as_account    = "system"
    }
  }
}

# Generate scripts for each policy
data "microsoft365_utility_windows_remediation_script_registry_key_generator" "policies" {
  for_each          = var.registry_policies
  context           = each.value.context
  registry_key_path = each.value.registry_key_path
  value_name        = each.value.value_name
  value_type        = each.value.value_type
  value_data        = each.value.value_data
}

# Deploy each policy to Intune
resource "microsoft365_graph_beta_device_management_windows_remediation_script" "policies" {
  for_each                   = var.registry_policies
  display_name               = each.value.display_name
  description                = each.value.description
  publisher                  = "IT Security Team"
  run_as_32_bit              = false
  enforce_signature_check    = false
  role_scope_tag_ids         = ["0"]
  run_as_account             = each.value.run_as_account
  detection_script_content   = data.microsoft365_utility_windows_remediation_script_registry_key_generator.policies[each.key].detection_script
  remediation_script_content = data.microsoft365_utility_windows_remediation_script_registry_key_generator.policies[each.key].remediation_script

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

      daily_schedule = {
        interval = 1
        time     = "02:00:00"
        use_utc  = true
      }
    }
  ]
}

# Example 5: Multi-string value (e.g., for allowed sites list)
data "microsoft365_utility_windows_remediation_script_registry_key_generator" "allowed_sites" {
  context           = "all_users"
  registry_key_path = "Software\\MyApp\\Security\\"
  value_name        = "AllowedSites"
  value_type        = "REG_MULTI_SZ"
  value_data        = "https://site1.example.com\nhttps://site2.example.com\nhttps://site3.example.com"
}

resource "microsoft365_graph_beta_device_management_windows_remediation_script" "allowed_sites" {
  display_name               = "Configure Allowed Sites List"
  description                = "Sets the list of allowed sites in the application"
  publisher                  = "IT Security Team"
  run_as_32_bit              = false
  enforce_signature_check    = false
  role_scope_tag_ids         = ["0"]
  run_as_account             = "system"
  detection_script_content   = data.microsoft365_utility_windows_remediation_script_registry_key_generator.allowed_sites.detection_script
  remediation_script_content = data.microsoft365_utility_windows_remediation_script_registry_key_generator.allowed_sites.remediation_script

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

      daily_schedule = {
        interval = 1
        time     = "03:00:00"
        use_utc  = true
      }
    }
  ]
}

# Example 6: Binary data
data "microsoft365_utility_windows_remediation_script_registry_key_generator" "binary_config" {
  context           = "current_user"
  registry_key_path = "Software\\MyApp\\Config\\"
  value_name        = "BinarySettings"
  value_type        = "REG_BINARY"
  value_data        = "01AF3C4D"
}

resource "microsoft365_graph_beta_device_management_windows_remediation_script" "binary_config" {
  display_name               = "Configure Binary Settings"
  description                = "Sets binary configuration data"
  publisher                  = "IT Application Team"
  run_as_32_bit              = false
  enforce_signature_check    = false
  role_scope_tag_ids         = ["0"]
  run_as_account             = "user"
  detection_script_content   = data.microsoft365_utility_windows_remediation_script_registry_key_generator.binary_config.detection_script
  remediation_script_content = data.microsoft365_utility_windows_remediation_script_registry_key_generator.binary_config.remediation_script

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "00000000-0000-0000-0000-000000000000"
      filter_type = "include"

      daily_schedule = {
        interval = 1
        time     = "04:00:00"
        use_utc  = true
      }
    }
  ]
}

# Output examples for manual script review
output "private_store_detection_script" {
  description = "Detection script for private store policy"
  value       = data.microsoft365_utility_windows_remediation_script_registry_key_generator.private_store.detection_script
}

output "private_store_remediation_script" {
  description = "Remediation script for private store policy"
  value       = data.microsoft365_utility_windows_remediation_script_registry_key_generator.private_store.remediation_script
}
```

## Argument Reference

* `context` - (Required) The execution context for the registry operation. Valid values:
  - `current_user`: Targets HKEY_CURRENT_USER for the logged-on user. Maps to `run_as_account = "user"` in the remediation script resource.
  - `all_users`: Targets HKEY_USERS for all user profiles on the device. Maps to `run_as_account = "system"` in the remediation script resource.

* `registry_key_path` - (Required) The registry key path relative to the user hive (without HKCU or HKU prefix).
  Example: `Software\\Policies\\Microsoft\\WindowsStore\\` or `Software\\MyApp\\Settings`.
  **Note**: Use double backslashes (`\\`) to escape path separators in HCL.

* `value_name` - (Required) The name of the registry value to manage. Use `(Default)` to target the default value of a key.

* `value_type` - (Required) The registry value type. Valid values:
  - `REG_SZ`: String value
  - `REG_DWORD`: 32-bit integer value (0 to 4,294,967,295)
  - `REG_QWORD`: 64-bit integer value
  - `REG_MULTI_SZ`: Multi-string value (separate strings with `\n` in `value_data`)
  - `REG_EXPAND_SZ`: Expandable string value (supports environment variables like `%ProgramFiles%`)
  - `REG_BINARY`: Binary data (provide as hex string in `value_data`, e.g., '01AF3C')

* `value_data` - (Required) The desired value data. Format depends on `value_type`:
  - `REG_SZ`, `REG_EXPAND_SZ`: String (e.g., `"Enabled"`)
  - `REG_DWORD`, `REG_QWORD`: Decimal integer as string (e.g., `"1"`, `"255"`)
  - `REG_MULTI_SZ`: Newline-separated strings (e.g., `"Line1\nLine2\nLine3"`)
  - `REG_BINARY`: Hexadecimal string (e.g., `"01AF3C"`)

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of this data source (same as `context` value).

* `detection_script` - Generated PowerShell detection script. This script:
  - Checks if the registry key and value exist and match the desired state
  - Returns exit code `0` if compliant (no action needed)
  - Returns exit code `1` if non-compliant (triggers remediation)
  - Use as the value for `detection_script_content` in `microsoft365_graph_beta_device_management_windows_remediation_script`

* `remediation_script` - Generated PowerShell remediation script. This script:
  - Creates the registry key if it doesn't exist
  - Creates or updates the registry value to the desired state
  - Returns exit code `0` on successful remediation
  - Use as the value for `remediation_script_content` in `microsoft365_graph_beta_device_management_windows_remediation_script`


### Context and run_as_account Mapping

| Data Source `context` | Resource `run_as_account` | Registry Hive | Use Case |
|-----------------------|---------------------------|---------------|----------|
| `current_user` | `user` | HKEY_CURRENT_USER | Per-user settings, user preferences |
| `all_users` | `system` | HKEY_USERS (all profiles) | Machine-wide settings, apply to all users |

## Important Considerations

### System Account Exclusions

When using `all_users` context, the generated scripts automatically exclude system accounts:
- `S-1-5-18` (LOCAL SYSTEM)
- `S-1-5-19` (LOCAL SERVICE)
- `S-1-5-20` (NETWORK SERVICE)
- `.DEFAULT` (Default user profile)
- `*_Classes` (Registry classes hives)

### Registry Key Paths

- Paths are relative to the user hive - do not include `HKCU:\` or `HKU:\` prefixes
- Always use double backslashes in HCL: `Software\\Path\\To\\Key\\`
- Trailing backslash is optional but recommended for clarity

### Permissions and Elevation

- **User context** (`current_user`): Scripts run as the logged-on user, suitable for HKCU settings
- **System context** (`all_users`): Scripts run as SYSTEM, suitable for applying settings to all user profiles
- Some registry keys may require additional permissions even in SYSTEM context

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `context` (String) The execution context for the registry operation. Valid values:
  - `current_user`: Applies to the currently logged-on user's registry hive (HKEY_CURRENT_USER)
  - `all_users`: Applies to all user profiles on the device (HKEY_USERS), excluding system accounts
- `registry_key_path` (String) The registry key path relative to the user hive (without HKCU or HKU prefix). Example: `Software\Policies\Microsoft\WindowsStore\` or `Software\MyApp\Settings`. Use double backslashes (`\\`) to escape path separators.
- `value_data` (String) The desired value data. Format depends on value_type:
  - `REG_SZ`, `REG_EXPAND_SZ`: String (e.g., 'Enabled')
  - `REG_DWORD`: Decimal integer (e.g., '1', '0', '255')
  - `REG_QWORD`: Decimal integer (e.g., '1234567890')
  - `REG_MULTI_SZ`: Multiple strings separated by newlines
  - `REG_BINARY`: Hexadecimal string (e.g., '01AF3C')
- `value_name` (String) The name of the registry value to manage. Use `(Default)` for the default value of a key.
- `value_type` (String) The registry value type. Valid values:
  - `REG_SZ`: String value
  - `REG_DWORD`: 32-bit integer value
  - `REG_QWORD`: 64-bit integer value
  - `REG_MULTI_SZ`: Multi-string value (separate strings with newlines)
  - `REG_EXPAND_SZ`: Expandable string value (can contain environment variables)
  - `REG_BINARY`: Binary data (provide as hex string, e.g., '01AF3C')

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `detection_script` (String) Generated PowerShell detection script that checks if the registry value exists and matches the desired state. Returns exit code 0 if compliant, 1 if remediation is needed. Use this as the detection script in Intune Proactive Remediations.
- `id` (String) Unique identifier for the data source (computed)
- `remediation_script` (String) Generated PowerShell remediation script that creates or updates the registry key and value to the desired state. Use this as the remediation script in Intune Proactive Remediations.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
