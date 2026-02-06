# ==============================================================================
# Test 05: Minimal Resource with Minimal Assignments (Unit Test)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_05_min_assignments" {
  display_name                             = "Test 05: Minimal Assignments macOS Software Update Configuration"
  update_schedule_type                     = "alwaysUpdate"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"

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
