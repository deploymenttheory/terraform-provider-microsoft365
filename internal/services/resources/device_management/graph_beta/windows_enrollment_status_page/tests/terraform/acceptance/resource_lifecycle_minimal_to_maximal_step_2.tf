resource "random_string" "test_suffix" {
  length  = 8
  upper   = false
  special = false
}

resource "microsoft365_graph_beta_device_management_windows_enrollment_status_page" "lifecycle" {
  display_name                                                        = "acc-test-windows-enrollment-status-page-lifecycle-${random_string.test_suffix.result}"
  description                                                         = "Test description for lifecycle enrollment status page"
  show_installation_progress                                          = true
  custom_error_message                                                = "Contact IT support for assistance"
  install_quality_updates                                             = true
  install_progress_timeout_in_minutes                                 = 120
  allow_log_collection_on_install_failure                             = true
  only_show_page_to_devices_provisioned_by_out_of_box_experience_oobe = true

  block_device_use_until_all_apps_and_profiles_are_installed = false
  allow_device_reset_on_install_failure                      = true
  allow_device_use_on_install_failure                        = true

  selected_mobile_app_ids = [
    "e4938228-aab3-493b-a9d5-8250aa8e9d55",
    "e83d36e1-3ff2-4567-90d9-940919184ad5",
    "cd4486df-05cc-42bd-8c34-67ac20e10166",
  ]

  only_fail_selected_blocking_apps_in_technician_phase = true

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
