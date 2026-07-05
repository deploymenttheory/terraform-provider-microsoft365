resource "random_uuid" "connector_group_assignment" {}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name            = "acctest-connector-group-assignment-${random_uuid.connector_group_assignment.result}"
  prevent_duplicate_names = false
  hard_delete             = true
}

resource "microsoft365_graph_beta_applications_on_premises_connector_group" "test" {
  name = "acctest-connector-group-assignment-${random_uuid.connector_group_assignment.result}"
}

resource "microsoft365_graph_beta_applications_on_premises_connector_group_assignment" "minimal" {
  application_id     = microsoft365_graph_beta_applications_application.test.id
  connector_group_id = microsoft365_graph_beta_applications_on_premises_connector_group.test.id
}
