# ==============================================================================
# Test 01: Deploy Minimal Configuration (Unit Test)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_macos_software_update_configuration" "test_01_minimal" {
  display_name                             = "Test 01: Minimal macOS Software Update Configuration"
  update_schedule_type                     = "alwaysUpdate"
  critical_update_behavior                 = "installASAP"
  config_data_update_behavior              = "installASAP"
  firmware_update_behavior                 = "installASAP"
  all_other_update_behavior                = "installASAP"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
