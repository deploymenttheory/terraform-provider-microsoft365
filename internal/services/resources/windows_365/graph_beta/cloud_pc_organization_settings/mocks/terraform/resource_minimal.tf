resource "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings" "minimal" {
  enable_mem_auto_enroll = false
  enable_single_sign_on  = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}