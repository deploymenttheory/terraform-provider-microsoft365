resource "microsoft365_graph_beta_device_and_app_management_windows_managed_app_protection" "test_004" {
  display_name = "unit-test-windows-managed-app-protection-lifecycle"
  description  = "Maximal lifecycle test configuration"

  print_blocked                               = true
  allowed_inbound_data_transfer_sources       = "none"
  allowed_outbound_data_transfer_destinations = "none"
  allowed_outbound_clipboard_sharing_level    = "none"
  app_action_if_unable_to_authenticate_user   = "block"
  maximum_allowed_device_threat_level         = "low"
  mobile_threat_defense_remediation_action    = "wipe"

  minimum_required_os_version = "10.0.19041"
  minimum_warning_os_version  = "10.0.18363"
  minimum_wipe_os_version     = "10.0.17763"

  period_offline_before_wipe_is_enforced = "P30D"
  period_offline_before_access_check     = "P7D"
}
