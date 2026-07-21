resource "microsoft365_graph_beta_device_and_app_management_android_managed_app_protection" "test_002" {
  display_name = "acc-test-android-managed-app-protection-002-maximal"
  description  = "Maximal acceptance test configuration for Android managed app protection"

  allowed_inbound_data_transfer_sources       = "none"
  allowed_outbound_data_transfer_destinations = "none"
  allowed_outbound_clipboard_sharing_level    = "blocked"
  data_backup_blocked                         = true
  screen_capture_blocked                      = true
  print_blocked                               = true
  save_as_blocked                             = true
  contact_sync_blocked                        = true
  fingerprint_blocked                         = true
  encrypt_app_data                            = true

  pin_required        = true
  minimum_pin_length  = 6
  maximum_pin_retries = 10
  simple_pin_blocked  = true
  pin_character_set   = "alphanumericAndSymbol"
  period_before_pin_reset = "P30D"

  minimum_required_os_version  = "9.0"
  minimum_warning_os_version   = "8.0"
  minimum_required_app_version = "2.0.0"
  minimum_warning_app_version  = "1.9.0"
  minimum_required_patch_version = "2024-01-01"
  minimum_warning_patch_version  = "2023-12-01"

  period_offline_before_wipe_is_enforced = "P30D"
  period_offline_before_access_check     = "P30D"
  period_online_before_access_check      = "PT30M"

  allowed_data_storage_locations = ["oneDriveForBusiness", "sharePoint"]
}
