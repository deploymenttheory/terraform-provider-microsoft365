# Look up an application by its exact display name (case-sensitive)
data "microsoft365_graph_beta_applications_application" "by_display_name" {
  display_name = "My Application"
}

# Output the application details
output "app_by_display_name" {
  value = {
    id          = data.microsoft365_graph_beta_applications_application.by_display_name.id
    app_id      = data.microsoft365_graph_beta_applications_application.by_display_name.app_id
    description = data.microsoft365_graph_beta_applications_application.by_display_name.description
    tags        = data.microsoft365_graph_beta_applications_application.by_display_name.tags
  }
}
