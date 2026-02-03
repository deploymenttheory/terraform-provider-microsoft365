# Look up an application and output its API permissions configuration
data "microsoft365_graph_beta_applications_application" "with_api_permissions" {
  display_name = "My API Application"
}

# Output API configuration details
output "api_configuration" {
  value = {
    id                   = data.microsoft365_graph_beta_applications_application.with_api_permissions.id
    display_name         = data.microsoft365_graph_beta_applications_application.with_api_permissions.display_name
    identifier_uris      = data.microsoft365_graph_beta_applications_application.with_api_permissions.identifier_uris
    api                  = data.microsoft365_graph_beta_applications_application.with_api_permissions.api
    app_roles            = data.microsoft365_graph_beta_applications_application.with_api_permissions.app_roles
    required_resource_access = data.microsoft365_graph_beta_applications_application.with_api_permissions.required_resource_access
  }
  sensitive = true # API configuration may contain sensitive scope information
}
