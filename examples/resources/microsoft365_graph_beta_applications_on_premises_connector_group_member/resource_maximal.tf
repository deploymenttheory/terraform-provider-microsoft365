# Create a connector group and move an installed private network connector into it.
resource "microsoft365_graph_beta_applications_on_premises_connector_group" "europe" {
  name   = "Europe connectors"
  region = "eur"
}

resource "microsoft365_graph_beta_applications_on_premises_connector_group_member" "maximal" {
  connector_id       = "00000000-0000-0000-0000-000000000000"
  connector_group_id = microsoft365_graph_beta_applications_on_premises_connector_group.europe.id

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
