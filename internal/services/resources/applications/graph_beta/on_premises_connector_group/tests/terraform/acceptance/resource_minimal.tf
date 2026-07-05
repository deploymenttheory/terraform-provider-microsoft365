resource "random_uuid" "connector_group" {}

resource "microsoft365_graph_beta_applications_on_premises_connector_group" "minimal" {
  name = "acctest-connector-group-${random_uuid.connector_group.result}"
}
