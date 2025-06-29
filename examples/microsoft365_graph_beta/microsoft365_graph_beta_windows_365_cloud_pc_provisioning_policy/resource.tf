resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "frontline" {
  display_name              = "test"
  description               = ""
  cloud_pc_naming_template  = "CPC-%USERNAME:5%-%RAND:5%"
  image_id                  = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"
  image_type                = "gallery"
  enable_single_sign_on     = false
  provisioning_type         = "sharedByUser"
  managed_by                = "windows365"

  domain_join_configurations = [
    {
      domain_join_type         = "azureADJoin"
      //type                     = "azureADJoin"
      region_name              = "automatic"
      region_group             = "asia"
    }
  ]

  windows_setting = {
    locale = "en-US"
  }

  # microsoft_managed_desktop = {
  #   managed_type = "starterManaged"
  #   profile      = null
  #   type         = "starterManaged"
  # }

  # autopatch = {
  #   autopatch_group_id = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  # }

  scope_ids = ["9", "8"]
} 