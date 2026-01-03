action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "maximal" {
  device_ids = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-987654321cba",
    "11111111-2222-3333-4444-555555555555",
  ]

  ignore_partial_failures = false
  validate_device_exists  = true

  timeouts = {
    invoke = "5m"
  }
}
