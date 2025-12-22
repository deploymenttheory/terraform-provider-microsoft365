resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test_003" {
  display_name = "unit-test-windows-quality-update-expedite-policy-003-lifecycle"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

