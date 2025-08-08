resource "microsoft365_graph_beta_device_management_windows_update_ring" "all_assignment_types" {
  display_name         = "Test All Assignment Types Windows Update Ring - Unique"
  description          = "Windows update ring with all assignment types for unit testing"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 7
  feature_updates_deferral_period_in_days = 14
  allow_windows11_upgrade                 = true
  skip_checks_before_restart              = false
  automatic_update_mode                   = "autoInstallAtMaintenanceTime"
  active_hours_start                      = "09:00:00"
  active_hours_end                        = "17:00:00"
  feature_updates_rollback_window_in_days = 30

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}