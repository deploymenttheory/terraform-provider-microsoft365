resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "win_365_with_hybrid_ad_join" {
  display_name             = "test"
  description              = ""
  cloud_pc_naming_template = "CPC-%USERNAME:5%-%RAND:5%"
  image_id                 = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"
  image_type               = "gallery"
  enable_single_sign_on    = true
  provisioning_type        = "dedicated"
  managed_by               = "windows365"

  domain_join_configurations = [
    {
      domain_join_type = "azureADJoin"
      region_name      = "automatic"
      region_group     = "asia"
    }
  ]

  windows_setting = {
    locale = "en-US"
  }

  microsoft_managed_desktop = {
    managed_type = "starterManaged"
    profile      = null
    type         = "starterManaged"
  }

  autopatch = {
    autopatch_group_id = "00000000-0000-0000-0000-000000000000"
  }

  apply_to_existing_cloud_pcs = {
    microsoft_entra_single_sign_on_for_all_devices        = false
    region_or_azure_network_connection_for_all_devices    = true
    region_or_azure_network_connection_for_select_devices = false
  }

  scope_ids = ["9", "8"]

  assignments = [
    {
      group_id = "af5dbc68-0ee3-485c-b85c-e27bbfff44c2"
      # Only group_id for dedicated
    }
  ]
}

resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "frontline_with_entra_id" {
  display_name             = "test"
  description              = ""
  cloud_pc_naming_template = "CPC-%USERNAME:5%-%RAND:5%"
  image_id                 = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"
  image_type               = "gallery"
  enable_single_sign_on    = false
  provisioning_type        = "sharedByUser"
  managed_by               = "windows365"

  domain_join_configurations = [
    {
      domain_join_type = "azureADJoin"
      region_name      = "automatic"
      region_group     = "asia"
    }
  ]

  windows_setting = {
    locale = "en-US"
  }

  microsoft_managed_desktop = {
    managed_type = "notManaged"
  }

  scope_ids = ["9", "8"]

  apply_to_existing_cloud_pcs = {
    microsoft_entra_single_sign_on_for_all_devices        = false
    region_or_azure_network_connection_for_all_devices    = true
    region_or_azure_network_connection_for_select_devices = false
  }

  assignments = [
    {
      group_id                = "00000000-0000-0000-0000-000000000000"
      service_plan_id         = data.microsoft365_graph_beta_windows_365_cloud_pc_frontline_service_plan.all.items[0].id
      allotment_license_count = 1
      allotment_display_name  = "Frontline Allotment"
    }
  ]
}