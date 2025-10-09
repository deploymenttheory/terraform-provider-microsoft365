resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "minimal" {
  display_name             = "unit-test-minimal"
  description              = "test"
  cloud_pc_naming_template = "CPC-%USERNAME:5%-%RAND:5%"
  provisioning_type        = "dedicated"
  image_id                 = "microsoftwindowsdesktop_windows-ent-cpc_win11-25h2-ent-cpc"
  image_type               = "gallery"
  enable_single_sign_on    = true
  managed_by               = "windows365"

  windows_setting = {
    locale = "en-US"
  }

  microsoft_managed_desktop = {
    managed_type = "notManaged"
  }

  domain_join_configurations = [
    {
      domain_join_type = "azureADJoin"
      region_group     = "japan"
      region_name      = "automatic"
    }
  ]

  scope_ids = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}