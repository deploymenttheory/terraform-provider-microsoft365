resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_002" {
  display_name     = "unit-test-windows-quality-update-policy-002-maximal"
  description      = "Scenario 2: Maximal configuration without assignments"
  hotpatch_enabled = true

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

