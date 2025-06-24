resource "microsoft365_graph_m365_admin_m365_apps_installation_options" "minimal" {
  update_channel = "current"
  
  apps_for_windows = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
  
  apps_for_mac = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }
} 