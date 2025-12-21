resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_001" {
  display_name     = "unit-test-windows-quality-update-policy-001-minimal"
  hotpatch_enabled = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

