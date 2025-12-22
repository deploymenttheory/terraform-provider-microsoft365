resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test_003" {
  display_name       = "unit-test-windows-quality-update-expedite-policy-003-lifecycle"
  description        = "Lifecycle Step 2: Updated to maximal configuration"
  role_scope_tag_ids = ["0", "1"]

  expedited_update_settings = {
    quality_update_release   = "2025-11-20T00:00:00Z"
    days_until_forced_reboot = 1
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

