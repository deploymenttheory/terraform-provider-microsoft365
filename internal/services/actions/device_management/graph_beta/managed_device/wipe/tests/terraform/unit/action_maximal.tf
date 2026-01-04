action "microsoft365_graph_beta_device_management_managed_device_wipe" "maximal" {
  config {
    device_ids = [
      "00000000-0000-0000-0000-000000000001",
      "00000000-0000-0000-0000-000000000002"
    ]
    keep_enrollment_data    = true
    keep_user_data          = false
    macos_unlock_code       = "123456"
    obliteration_behavior   = "default"
    persist_esim_data_plan  = true
    use_protected_wipe      = false
    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

