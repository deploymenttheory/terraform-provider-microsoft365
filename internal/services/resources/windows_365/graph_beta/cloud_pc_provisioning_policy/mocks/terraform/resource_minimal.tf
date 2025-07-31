resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "minimal" {
  display_name = "Test Minimal Provisioning Policy - Unique"
  image_id     = "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc"

  microsoft_managed_desktop = {
    # Uses default values: managed_type = "notManaged", profile = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  }

  windows_setting = {
    locale = "en-US"
  }

  domain_join_configurations = []

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}