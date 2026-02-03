# Look up an application using a simple OData query filter
# This example finds an application that starts with a specific prefix
data "microsoft365_graph_beta_applications_application" "by_odata_simple" {
  odata_query = "startswith(displayName, 'Contoso')"
}

# Output the application details
output "app_by_odata_simple" {
  value = {
    id           = data.microsoft365_graph_beta_applications_application.by_odata_simple.id
    app_id       = data.microsoft365_graph_beta_applications_application.by_odata_simple.app_id
    display_name = data.microsoft365_graph_beta_applications_application.by_odata_simple.display_name
  }
}
