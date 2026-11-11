resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "example" {
  display_name                            = "User settings persistence"
  image_id                                = "microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc"
  microsoft_managed_desktop               = { managed_type = "notManaged" }
  image_type                              = "gallery"
  windows_setting                         = { locale = "en-US" }
  provisioning_type                       = "sharedByEntraGroup"
  user_settings_persistence_configuration = { enabled = true, storage_size_category = "eightGB" }
  domain_join_configurations              = [{ domain_join_type = "azureADJoin", region_name = "automatic", region_group = "usWest" }]
}
