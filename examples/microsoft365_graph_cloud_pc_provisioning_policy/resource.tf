resource "microsoft365_graph_cloud_pc_provisioning_policy" "example" {
  display_name           = "Example Cloud PC Provisioning Policy"
  description            = "This is an example Cloud PC provisioning policy"
  cloud_pc_naming_template = "CPC-%USERNAME:4%-%RAND:5%"
  
  image_id               = "Microsoftwindowsdesktop_windows-ent-cpc_21h1-ent-cpc-m365"
  image_type             = "gallery"
  
  provisioning_type      = "dedicated"
  
  enable_single_sign_on  = true
  local_admin_enabled    = false

  domain_join_configurations {
    domain_join_type           = "azureADJoin"
    on_premises_connection_id  = "12345678-1234-1234-1234-123456789012"
    region_name                = "eastus"
  }

  microsoft_managed_desktop {
    managed_type = "premiumManaged"
    profile      = "Standard"
  }

  windows_setting {
    locale = "en-US"
  }

  timeouts {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}