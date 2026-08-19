resource "microsoft365_graph_beta_device_and_app_management_android_managed_app_protection" "test_004" {
  display_name = "acc-test-android-managed-app-protection-004-lifecycle"
  description  = "Maximal lifecycle acceptance test configuration"

  allowed_inbound_data_transfer_sources       = "none"
  allowed_outbound_data_transfer_destinations = "none"
  allowed_outbound_clipboard_sharing_level    = "blocked"
  data_backup_blocked                         = true
  screen_capture_blocked                      = true
  print_blocked                               = true
  encrypt_app_data                            = true

  pin_required       = true
  minimum_pin_length = 6
  maximum_pin_retries = 10
  pin_character_set  = "alphanumericAndSymbol"

  minimum_required_os_version  = "9.0"
  minimum_required_app_version = "2.0.0"

  period_offline_before_wipe_is_enforced = "P30D"
  period_offline_before_access_check     = "P30D"
}
