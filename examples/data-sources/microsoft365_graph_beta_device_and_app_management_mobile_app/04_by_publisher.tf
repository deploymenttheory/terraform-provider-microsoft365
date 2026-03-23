# Filter apps by publisher name (case-insensitive partial match)
# Uses server-side OData filtering for optimal performance
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_publisher" {
  publisher = "Adobe" # Finds all apps published by Adobe

  timeouts = {
    read = "10s"
  }
}

output "adobe_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_publisher.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      developer    = app.developer
      is_assigned  = app.is_assigned
      categories   = app.categories
    }
  ]
  description = "All apps from Adobe with key details"
}
