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
