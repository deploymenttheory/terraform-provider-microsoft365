# Example: Using the application datasource to reference an existing application
# and create a service principal for it

# Look up an existing application
data "microsoft365_graph_beta_applications_application" "existing" {
  display_name = "My Existing Application"
}

# Create a service principal for the application
resource "microsoft365_graph_beta_applications_service_principal" "sp" {
  app_id                       = data.microsoft365_graph_beta_applications_application.existing.app_id
  account_enabled              = true
  app_role_assignment_required = true

  tags = [
    "WindowsAzureActiveDirectoryIntegratedApp",
    "Production"
  ]
}

# Output the relationship
output "application_and_sp" {
  value = {
    application_id           = data.microsoft365_graph_beta_applications_application.existing.id
    application_app_id       = data.microsoft365_graph_beta_applications_application.existing.app_id
    application_name         = data.microsoft365_graph_beta_applications_application.existing.display_name
    service_principal_id     = microsoft365_graph_beta_applications_service_principal.sp.id
    service_principal_app_id = microsoft365_graph_beta_applications_service_principal.sp.app_id
  }
}
