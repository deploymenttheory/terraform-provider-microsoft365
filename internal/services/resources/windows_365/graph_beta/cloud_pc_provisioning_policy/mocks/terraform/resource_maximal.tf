resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "maximal" {
  display_name             = "Test Maximal Provisioning Policy - Unique"
  description              = "Maximal policy for testing with all features"
  cloud_pc_naming_template = "CPC-MAX-%USERNAME:5%-%RAND:5%"
  provisioning_type        = "dedicated"
  image_id                 = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"
  image_type               = "gallery"
  enable_single_sign_on    = true
  local_admin_enabled      = true
  managed_by               = "windows365"

  domain_join_configurations = [
    {
      domain_join_type          = "hybridAzureADJoin"
      on_premises_connection_id = "33333333-3333-3333-3333-333333333333"
      region_name               = "automatic"
      region_group              = "usWest"
    }
  ]

  windows_setting = {
    locale = "en-US"
  }

  microsoft_managed_desktop = {
    managed_type = "notManaged"
    profile      = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  }

  autopatch = {
    autopatch_group_id = "4aa9b805-9494-4eed-a04b-ed51ec9e631e"
  }

  autopilot_configuration = {
    device_preparation_profile_id   = "12345678-1234-1234-1234-123456789012"
    application_timeout_in_minutes  = 60
    on_failure_device_access_denied = true
  }

  apply_to_existing_cloud_pcs = {
    microsoft_entra_single_sign_on_for_all_devices        = false
    region_or_azure_network_connection_for_all_devices    = true
    region_or_azure_network_connection_for_select_devices = false
  }

  scope_ids = ["9", "8"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}