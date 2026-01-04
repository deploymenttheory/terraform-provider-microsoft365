action "microsoft365_graph_beta_device_management_managed_device_deprovision" "maximal" {
  config {
    managed_devices = [
      {
        device_id          = "12345678-1234-1234-1234-123456789abc"
        deprovision_reason = "Transitioning to new management"
      },
      {
        device_id          = "87654321-4321-4321-4321-987654321cba"
        deprovision_reason = "Device repurposing"
      }
    ]

    comanaged_devices = [
      {
        device_id          = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        deprovision_reason = "Removing co-management"
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

