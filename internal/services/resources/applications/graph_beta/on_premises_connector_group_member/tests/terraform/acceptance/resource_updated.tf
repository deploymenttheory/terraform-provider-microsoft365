variable "connector_id" {
  description = "ID of an installed private network connector. Set with TF_VAR_connector_id."
  type        = string
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_on_premises_connector_group" "first" {
  name = "acc-test-connector-group-member-a-${random_string.suffix.result}"
}

resource "microsoft365_graph_beta_applications_on_premises_connector_group" "second" {
  name = "acc-test-connector-group-member-b-${random_string.suffix.result}"
}

resource "microsoft365_graph_beta_applications_on_premises_connector_group_member" "test" {
  connector_id       = var.connector_id
  connector_group_id = microsoft365_graph_beta_applications_on_premises_connector_group.second.id
}
