resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test" {
  display_name     = "Acceptance - Windows Quality Update Policy - Updated"
  description      = "Updated description for acceptance testing"
  hotpatch_enabled = true
  role_scope_tag_ids = [
    "0",
    "1",
  ]
}


