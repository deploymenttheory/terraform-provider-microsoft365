resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test_006" {
  display_name       = "unit-test-windows-quality-update-expedite-policy-006-assignments-maximal"
  description        = "Maximal configuration with multiple assignments"
  role_scope_tag_ids = ["0", "1"]

  expedited_update_settings = {
    quality_update_release   = "2025-11-20T00:00:00Z"
    days_until_forced_reboot = 1
  }

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

