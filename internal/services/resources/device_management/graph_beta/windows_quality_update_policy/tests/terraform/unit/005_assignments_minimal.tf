resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_005" {
  display_name     = "unit-test-windows-quality-update-policy-005-assignments-minimal"
  hotpatch_enabled = false

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

