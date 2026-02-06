# ==============================================================================
# Test 04: Maximal to Minimal in Steps - Step 2 (Intermediate)
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_04_regression" {
  display_name                             = "acc-test-04-regression-${random_string.suffix.result}"
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
}
