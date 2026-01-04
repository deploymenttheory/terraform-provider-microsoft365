action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "ignore_failures" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-987654321cba"
    ]

    ignore_partial_failures = true

    timeouts = {
      invoke = "5m"
    }
  }
}

