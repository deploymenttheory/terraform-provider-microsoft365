action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "maximal" {
  config {
    provisioning_policy_id = "00000000-0000-0000-0000-000000000001"
    policy_settings        = "singleSignOn"
    reserve_percentage     = 50
    validate_policy_exists = true

    timeouts = {
      invoke = "5m"
    }
  }
}

