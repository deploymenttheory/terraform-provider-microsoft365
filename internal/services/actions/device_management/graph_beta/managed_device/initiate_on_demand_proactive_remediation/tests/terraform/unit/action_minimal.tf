action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "minimal" {
  config {
    managed_devices = [
      {
        device_id        = "12345678-1234-1234-1234-123456789abc"
        script_policy_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]
  }
}
