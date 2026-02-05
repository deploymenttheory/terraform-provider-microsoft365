# ==============================================================================
# Make Update Available As Soon As Possible
# ==============================================================================
# This example demonstrates how to deploy a Windows feature update immediately
# without any rollout schedule. The update will be made available to devices
# as soon as possible.

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "immediate" {
  display_name                                            = "Windows 11 25H2 - Immediate Deployment"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false
}
