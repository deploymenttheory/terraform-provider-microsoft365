# Confirm availability in this tenant before selecting the developer-ready image.
data "microsoft365_graph_beta_windows_365_cloud_pc_gallery_image" "developer_ready" {
  filter_type  = "id"
  filter_value = "microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc-devready"
}
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "developer_ready" {
  display_name               = "Developer-ready Cloud PCs"
  image_id                   = data.microsoft365_graph_beta_windows_365_cloud_pc_gallery_image.developer_ready.items[0].id
  microsoft_managed_desktop  = { managed_type = "notManaged" }
  image_type                 = "gallery"
  domain_join_configurations = [{ domain_join_type = "azureADJoin", region_name = "automatic", region_group = "usWest" }]
  windows_setting            = { locale = "en-US" }
}
