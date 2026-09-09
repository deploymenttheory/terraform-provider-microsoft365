resource "microsoft365_graph_beta_device_and_app_management_ios_managed_app_protection" "test_002" {
  display_name = "acc-test-ios-managed-app-protection-002-maximal"
  description  = "Maximal acceptance test configuration for iOS managed app protection"

  allowed_inbound_data_transfer_sources       = "none"
  allowed_outbound_data_transfer_destinations = "none"
  allowed_outbound_clipboard_sharing_level    = "blocked"
  data_backup_blocked                         = true
  print_blocked                               = true
  save_as_blocked                             = true
  contact_sync_blocked                        = true
  fingerprint_blocked                         = true
  face_id_blocked                             = true
  managed_browser_to_open_links_required      = true
  managed_browser                             = "microsoftEdge"
  custom_browser_protocol                     = "msedge"
  third_party_keyboards_blocked               = true
  filter_open_in_to_only_managed_apps         = true
  app_data_encryption_type                    = "afterDeviceRestart"

  pin_required            = true
  minimum_pin_length      = 6
  maximum_pin_retries     = 10
  simple_pin_blocked      = true
  pin_character_set       = "alphanumericAndSymbol"
  period_before_pin_reset = "P30D"

  minimum_required_os_version  = "15.0"
  minimum_warning_os_version   = "14.0"
  minimum_required_app_version = "2.0.0"
  minimum_warning_app_version  = "1.9.0"

  period_offline_before_wipe_is_enforced = "P30D"
  period_offline_before_access_check     = "P30D"
  period_online_before_access_check      = "PT30M"

  allowed_data_storage_locations = ["oneDriveForBusiness", "sharePoint"]
}