resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "test_005" {
  display_name = "unit-test-windows-quality-update-expedite-policy-005-assignments-minimal"

  expedited_update_settings = {
    quality_update_release   = "2025-11-20T00:00:00Z"
    days_until_forced_reboot = 1
  }

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

