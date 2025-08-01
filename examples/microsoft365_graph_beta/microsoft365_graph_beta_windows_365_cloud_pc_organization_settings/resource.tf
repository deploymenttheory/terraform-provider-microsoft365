resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "example" {
  enable_mem_auto_enroll = true
  enable_single_sign_on  = true
  os_version             = "windows11"
  user_account_type      = "standardUser"
  windows_settings = {
    language = "en-US"
  }
} 