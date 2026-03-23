# Get a specific mobile app by its ID
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_id" {
  app_id = "b395af0b-910f-40f9-ad74-1cb84406a20f" # Replace with actual app ID

  timeouts = {
    read = "10s"
  }
}

output "app_by_id" {
  value       = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_id.items[0], null)
  description = "Complete details of the specific app"
}

output "app_by_id_name" {
  value       = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_id.items[0].display_name, null)
  description = "Display name of the app"
}
