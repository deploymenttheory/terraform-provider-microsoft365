# Look up an application using OData query filtering by tags
# This example finds an application with a specific tag
data "microsoft365_graph_beta_applications_application" "by_tags" {
  odata_query = "tags/any(t:t eq 'Production')"
}

# Output the application details with tags
output "app_by_tags" {
  value = {
    id           = data.microsoft365_graph_beta_applications_application.by_tags.id
    app_id       = data.microsoft365_graph_beta_applications_application.by_tags.app_id
    display_name = data.microsoft365_graph_beta_applications_application.by_tags.display_name
    tags         = data.microsoft365_graph_beta_applications_application.by_tags.tags
  }
}
