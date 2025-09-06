resource "microsoft365_graph_beta_device_management_windows_enrollment_status_page" "with_assignments" {
  display_name                                                        = "acc-test-windows-enrollment-status-page-with-assignments-${random_string.test_suffix.result}"
  description                                                         = "Test description for enrollment status page with assignments"
  show_installation_progress                                          = true
  custom_error_message                                                = "Contact IT support for assistance"
  install_quality_updates                                             = true
  install_progress_timeout_in_minutes                                 = 120
  allow_log_collection_on_install_failure                             = true
  only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe = true

  block_device_use_until_all_apps_and_profiles_are_installed = false // this set to false enables the fields below to work
  allow_device_reset_on_install_failure                      = true  // can only be set to true if block_device_use_until_all_apps_and_profiles_are_installed is false
  allow_device_use_on_install_failure                        = true  // can only be set to true if block_device_use_until_all_apps_and_profiles_are_installed is false


  selected_mobile_app_ids = [ // if not set, this sets the field to 'all' in the gui. // can only be set to true if block_device_use_until_all_apps_and_profiles_are_installed is false
    "e4938228-aab3-493b-a9d5-8250aa8e9d55",
    "e83d36e1-3ff2-4567-90d9-940919184ad5",
    "cd4486df-05cc-42bd-8c34-67ac20e10166",
  ]

  only_fail_selected_blocking_apps_in_technician_phase = true // can only be set to true if block_device_use_until_all_apps_and_profiles_are_installed is false and selected_mobile_app_ids is set

  role_scope_tag_ids = [
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id,
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id
  ]

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    }
  ]
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}