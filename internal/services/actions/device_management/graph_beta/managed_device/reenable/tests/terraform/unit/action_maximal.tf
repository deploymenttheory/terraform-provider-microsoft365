action "microsoft365_graph_beta_device_management_managed_device_reenable" "maximal" {
  config {
    managed_device_ids = [
      "00000000-0000-0000-0000-000000000001",
      "00000000-0000-0000-0000-000000000002"
    ]
    comanaged_device_ids = [
      "00000000-0000-0000-0000-000000000003"
    ]
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

