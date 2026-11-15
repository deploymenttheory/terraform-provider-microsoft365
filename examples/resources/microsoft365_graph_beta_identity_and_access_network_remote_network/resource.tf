data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "m365" {
  traffic_forwarding_type = "m365"
}

data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "internet" {
  traffic_forwarding_type = "internet"
}

resource "microsoft365_graph_beta_identity_and_access_network_remote_network" "example" {
  name   = "example-branch-office"
  region = "eastUS"

  # Only the default Microsoft 365 and Internet Access profiles can be assigned, and they must be enabled.
  forwarding_profile_ids = [
    data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile.m365.items[0].id,
    data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile.internet.items[0].id,
  ]

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
