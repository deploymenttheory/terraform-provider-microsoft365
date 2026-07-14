data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "default_acquire" {
  traffic_forwarding_type = "internet"
  policy_name             = "Default Acquire"
}

resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "default_acquire" {
  forwarding_profile_id = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.default_acquire.forwarding_profile_id
  policy_link_id        = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.default_acquire.policy_link_id
  state                 = "enabled"
}
