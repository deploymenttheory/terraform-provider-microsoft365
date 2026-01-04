action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "maximal" {
  config {
    managed_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "refreshMachinePolicy"
      },
      {
        device_id = "00000000-0000-0000-0000-000000000002"
        action    = "appEvaluation"
      }
    ]
    comanaged_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000003"
        action    = "quickScan"
      }
    ]
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

