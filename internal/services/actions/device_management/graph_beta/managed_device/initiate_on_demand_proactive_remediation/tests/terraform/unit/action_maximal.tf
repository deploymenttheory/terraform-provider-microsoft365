action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "maximal" {
  config {
    managed_devices = [
      {
        device_id        = "12345678-1234-1234-1234-123456789abc"
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]

    comanaged_devices = [
      {
        device_id        = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        script_policy_id = "bbbbbbbb-2222-2222-2222-cccccccccccc"
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}
