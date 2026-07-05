resource "random_uuid" "connector_group" {}

resource "microsoft365_graph_beta_applications_on_premises_connector_group" "with_region" {
  name   = "acctest-connector-group-region-${random_uuid.connector_group.result}"
  region = "nam"
}
