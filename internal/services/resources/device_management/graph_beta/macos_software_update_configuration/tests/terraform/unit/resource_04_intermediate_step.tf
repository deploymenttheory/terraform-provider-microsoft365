# ==============================================================================
# Test 04: Maximal to Minimal in Steps - Step 2 (Intermediate) (Unit Test)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_04_regression" {
  display_name                             = "Test 04: Regression macOS Software Update Configuration"
  description                              = "Intermediate configuration removing some features"
  update_schedule_type                     = "updateDuringTimeWindows"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"

  custom_update_time_windows = [
    {
      start_day  = "monday"
      end_day    = "friday"
      start_time = "02:00:00"
      end_time   = "06:00:00"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
