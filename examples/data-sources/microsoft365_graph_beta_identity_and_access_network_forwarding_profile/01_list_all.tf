data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "all" {
  filter_type = "all"

  timeouts = {
    read = "3m"
  }
}

output "internet_forwarding_profile" {
  value = [
    for profile in data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile.all.items : {
      id                      = profile.id
      name                    = profile.name
      traffic_forwarding_type = profile.traffic_forwarding_type
      policies = [
        for link in profile.policies : {
          policy_link_id = link.policy_link_id
          policy_id      = link.policy_id
          policy_name    = link.policy_name
          state          = link.state
        }
      ]
    }
    if profile.traffic_forwarding_type == "internet"
  ]
}
