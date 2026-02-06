# ==============================================================================
# Test 08: Maximal Assignments to Minimal Assignments - Step 2 (Minimal) (Unit Test)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_08_assignments_regression" {
  display_name                             = "Test 08: Assignments Regression macOS Software Update Configuration"
  update_schedule_type                     = "alwaysUpdate"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "66666666-6666-6666-6666-666666666666"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
