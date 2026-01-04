# Example 1: Apply Cloud PC provisioning policy - Minimal
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_policy" {
  config {
    provisioning_policy_id = "12345678-1234-1234-1234-123456789abc"
  }
}

# Example 2: Apply Cloud PC provisioning policy with timeout
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_with_timeout" {
  config {
    provisioning_policy_id = "12345678-1234-1234-1234-123456789abc"

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Apply Cloud PC provisioning policy with all options - Maximal
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_maximal" {
  config {
    provisioning_policy_id = "12345678-1234-1234-1234-123456789abc"
    policy_settings        = "singleSignOn"
    reserve_percentage     = 50
    validate_policy_exists = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Apply provisioning policy from data source
data "microsoft365_graph_beta_device_management_cloud_pc_provisioning_policy" "standard_policy" {
  filter_type  = "display_name"
  filter_value = "Standard User Policy"
}

action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_from_data" {
  config {
    provisioning_policy_id = data.microsoft365_graph_beta_device_management_cloud_pc_provisioning_policy.standard_policy.id

    validate_policy_exists = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Apply policy with custom reserve percentage
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_with_reserve" {
  config {
    provisioning_policy_id = "12345678-1234-1234-1234-123456789abc"
    reserve_percentage     = 25

    timeouts = {
      invoke = "5m"
    }
  }
}

# Output examples
output "applied_policy_id" {
  value       = action.microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy.apply_policy.config.provisioning_policy_id
  description = "The provisioning policy ID that was applied"
}
