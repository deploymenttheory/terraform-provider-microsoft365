resource "microsoft365_graph_beta_windows_365_user_setting" "minimal" {
  display_name         = "Test Minimal User Setting"
  local_admin_enabled  = false
  reset_enabled        = false
  self_service_enabled = false

  restore_point_setting = {
    frequency_in_hours   = 12
    frequency_type       = "default"
    user_restore_enabled = false
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}