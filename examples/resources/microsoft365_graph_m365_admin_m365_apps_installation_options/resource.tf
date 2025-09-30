resource "microsoft365_graph_device_and_app_management_m365_apps_installation_options" "example" {
  update_channel = "monthlyEnterprise"

  apps_for_windows = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }

  apps_for_mac = {
    is_microsoft_365_apps_enabled = true
    is_skype_for_business_enabled = true
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}