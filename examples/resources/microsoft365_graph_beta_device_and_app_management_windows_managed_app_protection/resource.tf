resource "microsoft365_graph_beta_device_and_app_management_windows_managed_app_protection" "example" {
  display_name = "Example Windows MAM Policy"
  description  = "Windows Mobile Application Management policy managed by Terraform."

  # Data transfer controls
  allowed_inbound_data_transfer_sources        = "none"
  allowed_outbound_data_transfer_destinations  = "none"
  allowed_outbound_clipboard_sharing_level     = "none"

  # Printing
  print_blocked = true

  # Threat response
  maximum_allowed_device_threat_level      = "notConfigured"
  mobile_threat_defense_remediation_action = "block"

  # Offline behaviour
  period_offline_before_wipe_is_enforced = "P90D"
  period_offline_before_access_check     = "P30D"

  # Optional — uncomment to configure
  # app_action_if_unable_to_authenticate_user = "block"
  # minimum_required_os_version               = "10.0.19041"
  # minimum_warning_os_version                = "10.0.18363"
  # minimum_wipe_os_version                   = "10.0.17763"
  # minimum_required_app_version              = "1.0.0"
  # minimum_required_sdk_version              = "1.0.0"
  # role_scope_tag_ids                        = ["1", "2"]
}