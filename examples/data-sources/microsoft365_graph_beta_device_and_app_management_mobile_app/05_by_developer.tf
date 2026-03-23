# Filter apps by developer name (case-insensitive partial match)
# Uses server-side OData filtering for optimal performance
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_developer" {
  developer = "Microsoft" # Finds all apps developed by Microsoft

  timeouts = {
    read = "10s"
  }
}

output "microsoft_developed_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_developer.items : {
      id           = app.id
      display_name = app.display_name
      developer    = app.developer
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ]
  description = "All apps developed by Microsoft"
}
