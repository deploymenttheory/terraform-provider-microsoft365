data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "internet" {
  filter_type  = "name"
  filter_value = "Internet traffic forwarding profile"
}
