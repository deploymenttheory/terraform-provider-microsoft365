data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "internet" {
  filter_type  = "traffic_forwarding_type"
  filter_value = "internet"
}
