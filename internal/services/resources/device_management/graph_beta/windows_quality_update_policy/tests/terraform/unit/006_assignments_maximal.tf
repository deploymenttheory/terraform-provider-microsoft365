resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_006" {
  display_name     = "unit-test-windows-quality-update-policy-006-assignments-maximal"
  hotpatch_enabled = false

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "77777777-7777-7777-7777-777777777777"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

