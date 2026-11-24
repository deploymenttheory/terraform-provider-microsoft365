# Whole Middle East geography: every region group and region in the geography.
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "middle_east" {
  display_name              = "Middle East geography"
  image_id                  = "microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc-m365"
  image_type                = "gallery"
  microsoft_managed_desktop = { managed_type = "notManaged" }
  windows_setting           = { locale = "en-US" }

  domain_join_configurations = [
    {
      domain_join_type         = "azureADJoin"
      geographic_location_type = "middleEast"
      region_group             = "automatic"
      region_name              = "automatic"
    }
  ]
}

# UAE region group within the Middle East geography.
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "uae" {
  display_name              = "Middle East - UAE"
  image_id                  = "microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc-m365"
  image_type                = "gallery"
  microsoft_managed_desktop = { managed_type = "notManaged" }
  windows_setting           = { locale = "en-US" }

  domain_join_configurations = [
    {
      domain_join_type         = "azureADJoin"
      geographic_location_type = "middleEast"
      region_group             = "uae"
      region_name              = "automatic"
    }
  ]
}

# Whole Africa geography (this is how the Intune portal selects South Africa).
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "africa" {
  display_name              = "Africa geography"
  image_id                  = "microsoftwindowsdesktop_windows-ent-cpc_win11-26h2-ent-cpc-m365"
  image_type                = "gallery"
  microsoft_managed_desktop = { managed_type = "notManaged" }
  windows_setting           = { locale = "en-US" }

  domain_join_configurations = [
    {
      domain_join_type         = "azureADJoin"
      geographic_location_type = "africa"
      region_group             = "automatic"
      region_name              = "automatic"
    }
  ]
}
