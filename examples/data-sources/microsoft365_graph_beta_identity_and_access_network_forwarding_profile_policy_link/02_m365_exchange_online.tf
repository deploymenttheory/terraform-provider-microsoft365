data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "exchange_online" {
  traffic_forwarding_type = "m365"
  policy_name             = "Exchange Online"
}
