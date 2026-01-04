action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "maximal" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      },
      {
        device_id     = "87654321-4321-4321-4321-987654321cba"
        template_type = "unknownFutureValue"
      }
    ]

    comanaged_devices = [
      {
        device_id     = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        template_type = "predefined"
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

