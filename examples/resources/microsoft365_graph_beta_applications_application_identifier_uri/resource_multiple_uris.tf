resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-multi-uri-app"
  description  = "Application with multiple identifier URIs"
}

# First identifier URI - api:// format
resource "microsoft365_graph_beta_applications_application_identifier_uri" "api_uri" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  identifier_uri = "api://${microsoft365_graph_beta_applications_application.example.app_id}"
}

# Second identifier URI - https:// format (must use verified domain)
resource "microsoft365_graph_beta_applications_application_identifier_uri" "https_uri" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  identifier_uri = "https://myverifieddomain.com/my-app"

  depends_on = [microsoft365_graph_beta_applications_application_identifier_uri.api_uri]
}
