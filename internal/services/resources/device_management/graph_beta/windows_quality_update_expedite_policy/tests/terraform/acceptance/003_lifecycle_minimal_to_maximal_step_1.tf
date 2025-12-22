
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test_003" {
  display_name = "acc-test-expedite-policy-003-${random_string.test_suffix.result}"

  expedited_update_settings = {
    quality_update_release   = "2025-12-09T00:00:00Z"
    days_until_forced_reboot = 2
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


