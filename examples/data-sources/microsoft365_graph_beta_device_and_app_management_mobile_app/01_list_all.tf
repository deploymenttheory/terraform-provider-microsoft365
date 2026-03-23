# Get all mobile apps from Intune
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "all_apps" {
  list_all = true
  timeouts = {
    read = "10s"
  }
}

output "all_intune_apps_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items)
  description = "Total number of mobile apps in Intune"
}

output "all_intune_apps_summary" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.all_apps.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ]
  description = "Summary of all apps with key fields"
}
