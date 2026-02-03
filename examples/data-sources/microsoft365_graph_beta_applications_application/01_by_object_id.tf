# Look up an application by its object ID (most direct and efficient method)
data "microsoft365_graph_beta_applications_application" "by_object_id" {
  object_id = "00000000-0000-0000-0000-000000000000" # Replace with actual object ID
}

# Output the application details
output "app_by_object_id" {
  value = {
    id           = data.microsoft365_graph_beta_applications_application.by_object_id.id
    app_id       = data.microsoft365_graph_beta_applications_application.by_object_id.app_id
    display_name = data.microsoft365_graph_beta_applications_application.by_object_id.display_name
  }
}
