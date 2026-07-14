resource "microsoft365_graph_beta_device_and_app_management_windows_managed_app_protection" "test_002" {
  display_name = "acc-test-windows-managed-app-protection-002-maximal"
  description  = "Maximal acceptance test configuration for Windows managed app protection"

  print_blocked                               = true
  allowed_inbound_data_transfer_sources       = "none"
  allowed_outbound_data_transfer_destinations = "none"
  allowed_outbound_clipboard_sharing_level    = "none"
  app_action_if_unable_to_authenticate_user   = "block"
  maximum_allowed_device_threat_level         = "low"
  mobile_threat_defense_remediation_action    = "wipe"

  minimum_required_os_version  = "10.0.19041"
  minimum_warning_os_version   = "10.0.18363"
  minimum_wipe_os_version      = "10.0.17763"
  minimum_required_app_version = "1.0.0"
  minimum_warning_app_version  = "1.1.0"
  minimum_wipe_app_version     = "0.9.0"

  period_offline_before_wipe_is_enforced = "P30D"
  period_offline_before_access_check     = "P7D"
}
