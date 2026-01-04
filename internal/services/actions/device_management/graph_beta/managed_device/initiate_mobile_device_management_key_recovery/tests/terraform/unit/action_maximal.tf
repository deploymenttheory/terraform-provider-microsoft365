action "microsoft365_graph_beta_device_management_managed_device_initiate_mobile_device_management_key_recovery" "maximal" {
  config {
    managed_device_ids   = ["12345678-1234-1234-1234-123456789abc"]
    comanaged_device_ids = ["aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}
