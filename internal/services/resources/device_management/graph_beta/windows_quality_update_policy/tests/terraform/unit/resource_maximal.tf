resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "maximal" {
  display_name     = "Test Maximal Windows Quality Update Policy - Unique"
  description      = "Maximal Windows Quality Update Policy for testing with all features"
  hotpatch_enabled = true
  role_scope_tag_ids = [
    "0",
    "1",
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}


