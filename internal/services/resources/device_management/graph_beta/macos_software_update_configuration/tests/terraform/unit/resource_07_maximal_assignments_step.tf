# ==============================================================================
# Test 07: Minimal Assignments to Maximal Assignments - Step 2 (Maximal) (Unit Test)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_07_assignments_progression" {
  display_name                             = "Test 07: Assignments Progression macOS Software Update Configuration"
  update_schedule_type                     = "alwaysUpdate"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "55555555-5555-5555-5555-555555555555"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
