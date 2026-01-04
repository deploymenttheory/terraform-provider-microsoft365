action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "minimal" {
  config {
    managed_devices = [
      {
        device_id = "00000000-0000-0000-0000-000000000001"
        action    = "refreshMachinePolicy"
      }
    ]
  }
}

