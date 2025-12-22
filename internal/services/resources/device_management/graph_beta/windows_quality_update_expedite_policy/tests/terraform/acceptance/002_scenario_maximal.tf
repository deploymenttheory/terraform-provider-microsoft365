
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test_002" {
  display_name       = "acc-test-expedite-policy-002-${random_string.test_suffix.result}"
  description        = "Acceptance test maximal configuration"
  role_scope_tag_ids = ["0", "1"]

  expedited_update_settings = {
    quality_update_release  = "2025-11-20T00:00:00Z"
    days_until_forced_reboot = 1
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


