resource "microsoft365_graph_beta_device_management_windows_enrollment_status_page" "minimal" {
  display_name                                                        = "unit-test-windows-enrollment-status-page-minimal"
  description                                                         = "Test description for minimal enrollment status page"
  show_installation_progress                                          = true
  custom_error_message                                                = "Contact IT support for assistance"
  install_quality_updates                                             = true
  install_progress_timeout_in_minutes                                 = 120
  allow_log_collection_on_install_failure                             = true
  only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe = true

  block_device_use_until_all_apps_and_profiles_are_installed = true  // this set to false enables the fields below to work
  allow_device_reset_on_install_failure                      = false // can only be set to true if block_device_use_until_all_apps_and_profiles_are_installed is false
  allow_device_use_on_install_failure                        = false // can only be set to true if block_device_use_until_all_apps_and_profiles_are_installed is false

  only_fail_selected_blocking_apps_in_technician_phase = false // can only be set to true if block_device_use_until_all_apps_and_profiles_are_installed is false and selected_mobile_app_ids is set

  role_scope_tag_ids = ["0"]
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}