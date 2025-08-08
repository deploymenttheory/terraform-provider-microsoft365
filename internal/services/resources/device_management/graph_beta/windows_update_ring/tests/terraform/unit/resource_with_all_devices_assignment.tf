resource "microsoft365_graph_beta_device_management_windows_update_ring" "all_devices_assignment" {
  display_name                            = "Test All Devices Assignment Windows Update Ring - Unique"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 0
  feature_updates_deferral_period_in_days = 0
  allow_windows11_upgrade                 = true
  skip_checks_before_restart              = false
  automatic_update_mode                   = "userDefined"
  feature_updates_rollback_window_in_days = 10

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}