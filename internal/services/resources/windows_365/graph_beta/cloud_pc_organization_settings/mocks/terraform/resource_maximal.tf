resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "maximal" {
  enable_mem_auto_enroll = true
  enable_single_sign_on  = true
  os_version             = "windows11"
  user_account_type      = "administrator"

  windows_settings = {
    language = "en-GB"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}