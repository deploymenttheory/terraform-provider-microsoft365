resource "microsoft365_graph_beta_identity_and_access_network_forwarding_options" "test" {
  skip_dns_lookup_state = var.existing_dns_state
}
