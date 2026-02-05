resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Test Case: Make Available On Specific Date
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_001" {
  display_name                                            = "acc-test-002-make-available-on-specific-date-${random_string.test_suffix.result}"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = false
  install_latest_windows10_on_windows11_ineligible_device = false

  rollout_settings = {
    offer_start_date_time_in_utc = "2030-01-13T00:00:00Z"
  }
}
