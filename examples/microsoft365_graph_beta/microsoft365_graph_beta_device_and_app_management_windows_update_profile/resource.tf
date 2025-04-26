resource "microsoft365_graph_beta_device_and_app_management_windows_driver_update_profile" "manual_example" {
  display_name              = "Windows Driver Updates - Production"
  description               = "Driver update profile for production machines"
  approval_type             = "manual"
  role_scope_tag_ids      = [8, 9]
}

resource "microsoft365_graph_beta_device_and_app_management_windows_driver_update_profile" "automatic_example" {
  display_name              = "Windows Driver Updates - Production"
  description               = "Driver update profile for production machines"
  approval_type             = "automatic"
  deployment_deferral_in_days = 14
  role_scope_tag_ids      = [8, 9]

  # Optional - Timeouts
  timeouts = {
    create = "1m"
    read   = "1m"
    update = "30s"
    delete = "1m"
  }
}