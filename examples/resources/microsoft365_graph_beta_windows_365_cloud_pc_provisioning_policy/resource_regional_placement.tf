resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "example" {
  display_name               = "Regional placement"
  image_id                   = "microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc"
  microsoft_managed_desktop  = { managed_type = "notManaged" }
  image_type                 = "gallery"
  windows_setting            = { locale = "en-US" }
  domain_join_configurations = [{ domain_join_type = "azureADJoin", region_name = "westeurope", region_group = "netherlands" }]
}
