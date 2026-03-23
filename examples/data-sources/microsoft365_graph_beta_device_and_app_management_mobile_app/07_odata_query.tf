# Advanced OData filter query
# Uses custom OData expressions for complex filtering
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "odata_custom" {
  odata_query = "startswith(publisher, 'Microsoft') and isAssigned eq true"

  timeouts = {
    read = "20s"
  }
}

output "microsoft_assigned_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.odata_custom.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      is_assigned  = app.is_assigned
    }
  ]
  description = "Microsoft apps that are assigned"
}

# Example: Filter by creation date
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "recent_apps" {
  odata_query = "createdDateTime gt 2024-01-01T00:00:00Z"

  timeouts = {
    read = "20s"
  }
}

output "recent_apps_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app.recent_apps.items)
  description = "Number of apps created after January 1, 2024"
}
