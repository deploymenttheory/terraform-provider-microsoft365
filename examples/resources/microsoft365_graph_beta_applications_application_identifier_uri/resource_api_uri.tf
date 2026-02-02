resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-api-application"
  description  = "Application exposing an API"
}

# Add an api:// identifier URI
resource "microsoft365_graph_beta_applications_application_identifier_uri" "api_uri" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  identifier_uri = "api://${microsoft365_graph_beta_applications_application.example.app_id}"
}
