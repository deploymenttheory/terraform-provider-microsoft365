resource "microsoft365_graph_device_and_app_management_cloud_pc_user_setting" "example" {
  display_name      = "Windows 365 User Setting"
  local_admin_enabled = true
  reset_enabled     = false

  restore_point_setting {
    frequency_type      = "sixHours"
    user_restore_enabled = true
  }

  timeouts {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}