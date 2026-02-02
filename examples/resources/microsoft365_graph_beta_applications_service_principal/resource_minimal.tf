resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-application"
  description  = "Application for service principal"
}

# Create service principal for the application
resource "microsoft365_graph_beta_applications_service_principal" "example" {
  app_id = microsoft365_graph_beta_applications_application.example.app_id
}
