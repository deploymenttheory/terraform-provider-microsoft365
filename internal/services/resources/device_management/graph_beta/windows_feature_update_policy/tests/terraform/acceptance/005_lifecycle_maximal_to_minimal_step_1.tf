resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Test Case: Lifecycle Maximal To Minimal Step 1
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "test_004" {
  display_name                                            = "acc-test-005-lifecycle-maximal-to-minimal-${random_string.test_suffix.result}"
  description                                             = "Maximal lifecycle test configuration"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = true
  install_latest_windows10_on_windows11_ineligible_device = true

  rollout_settings = {
    offer_start_date_time_in_utc = "2030-01-13T00:00:00Z"
    offer_end_date_time_in_utc   = "2030-01-14T00:00:00Z"
    offer_interval_in_days       = 1
  }

  role_scope_tag_ids = ["0", "1"]
}
