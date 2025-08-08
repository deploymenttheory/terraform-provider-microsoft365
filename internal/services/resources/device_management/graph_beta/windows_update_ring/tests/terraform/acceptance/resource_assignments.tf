resource "microsoft365_graph_beta_device_management_windows_update_ring" "assignments" {
  display_name         = "Test All Assignment Types Windows Update Ring"
  description          = "Windows update ring with comprehensive assignments for acceptance testing"
  microsoft_update_service_allowed        = true
  drivers_excluded                        = false
  quality_updates_deferral_period_in_days = 7
  feature_updates_deferral_period_in_days = 14
  allow_windows11_upgrade                 = true
  skip_checks_before_restart              = false
  automatic_update_mode                   = "userDefined"
  feature_updates_rollback_window_in_days = 30

  role_scope_tag_ids = [microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id, microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    }
  ]

  
  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}