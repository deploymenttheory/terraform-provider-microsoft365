terraform {
  required_providers {
    microsoft365 = {
      source = "deploymenttheory/microsoft365"
    }
  }
}

provider "microsoft365" {
  # Provider configuration
}

# Example 1: Apply single sign-on settings to existing Cloud PCs
action "apply_sso_settings" {
  action = microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy

  provisioning_policy_id = "12345678-1234-1234-1234-123456789abc"
  policy_settings        = "singleSignOn"

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Apply region settings to existing Cloud PCs
action "apply_region_settings" {
  action = microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy

  provisioning_policy_id = "12345678-1234-1234-1234-123456789abc"
  policy_settings        = "region"

  timeouts = {
    invoke = "5m"
  }
}

