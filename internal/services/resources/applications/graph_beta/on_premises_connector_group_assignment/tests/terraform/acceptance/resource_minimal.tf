resource "random_uuid" "connector_group_assignment" {}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name            = "acctest-connector-group-assignment-${random_uuid.connector_group_assignment.result}"
  prevent_duplicate_names = false
  hard_delete             = true
}

resource "microsoft365_graph_beta_applications_on_premises_connector_group" "test" {
  name = "acctest-connector-group-assignment-${random_uuid.connector_group_assignment.result}"
}

resource "microsoft365_graph_beta_applications_application_on_premises_publishing" "test" {
  application_id                     = microsoft365_graph_beta_applications_application.test.id
  application_type                   = "quickaccessapp"
  external_authentication_type       = "passthru"
  internal_url                       = "https://acctest-connector-group-assignment-${random_uuid.connector_group_assignment.result}.example.com"
  is_on_prem_publishing_enabled      = true
  is_translate_host_header_enabled   = true
  is_translate_links_in_body_enabled = true
}

resource "microsoft365_graph_beta_applications_on_premises_connector_group_assignment" "minimal" {
  application_id     = microsoft365_graph_beta_applications_application.test.id
  connector_group_id = microsoft365_graph_beta_applications_on_premises_connector_group.test.id

  depends_on = [
    microsoft365_graph_beta_applications_application_on_premises_publishing.test
  ]
}
