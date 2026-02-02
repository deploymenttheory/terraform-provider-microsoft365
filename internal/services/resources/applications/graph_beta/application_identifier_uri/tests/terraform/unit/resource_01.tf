# Application Identifier URI configuration for unit testing
resource "microsoft365_graph_beta_applications_application_identifier_uri" "test" {
  application_id = "11111111-1111-1111-1111-111111111111"
  identifier_uri = "api://11111111-1111-1111-1111-111111111111"
}
