action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "disable_validation" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-987654321cba"
    ]

    validate_device_exists = false

    timeouts = {
      invoke = "5m"
    }
  }
}

