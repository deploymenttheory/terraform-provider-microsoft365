data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "custom_acquire" {
  traffic_forwarding_type = "internet"
  policy_name             = "Custom Acquire"
}

output "custom_acquire_policy_link" {
  value = {
    forwarding_profile_id = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_acquire.forwarding_profile_id
    policy_link_id        = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_acquire.policy_link_id
    policy_id             = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_acquire.policy_id
    state                 = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_acquire.state
  }
}
