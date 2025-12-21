resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_004" {
  display_name     = "unit-test-windows-quality-update-policy-004-lifecycle"
  description      = "Lifecycle Step 1: Starting with maximal configuration"
  hotpatch_enabled = true

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

