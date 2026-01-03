action "microsoft365_graph_beta_device_management_managed_device_clean_windows_device" "maximal" {
  managed_devices = [
    {
      device_id      = "12345678-1234-1234-1234-123456789abc"
      keep_user_data = false
    },
    {
      device_id      = "87654321-4321-4321-4321-987654321cba"
      keep_user_data = true
    }
  ]

  comanaged_devices = [
    {
      device_id      = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      keep_user_data = false
    }
  ]

  ignore_partial_failures = false
  validate_device_exists  = true

  timeouts = {
    invoke = "5m"
  }
}
