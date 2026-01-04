action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "minimal" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      }
    ]
  }
}

