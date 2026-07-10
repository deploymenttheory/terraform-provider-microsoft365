data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "internet" {
  traffic_forwarding_type = "internet"
}

locals {
  internet_profile = one(data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile.internet.items)
  default_acquire_policy_link = one([
    for link in local.internet_profile.policies : link
    if link.policy_name == "Default Acquire"
  ])
}

resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "default_acquire" {
  forwarding_profile_id = local.internet_profile.id
  policy_link_id        = local.default_acquire_policy_link.policy_link_id
  state                 = "enabled"
}
