resource "microsoft365_graph_m365_admin_m365_apps_installation_options" "maximal" {
  update_channel = "semiAnnual"

  apps_for_windows = {
    is_microsoft_365_apps_enabled = false
    is_skype_for_business_enabled = false
  }

  apps_for_mac = {
    is_microsoft_365_apps_enabled = false
    is_skype_for_business_enabled = false
  }

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
} 