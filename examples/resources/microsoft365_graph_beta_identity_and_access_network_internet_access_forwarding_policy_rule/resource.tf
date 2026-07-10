data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "custom_acquire" {
  traffic_forwarding_type = "internet"
  policy_name             = "Custom Acquire"
}

data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "custom_bypass" {
  traffic_forwarding_type = "internet"
  policy_name             = "Custom bypass"
}

resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "fqdn" {
  forwarding_policy_id = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_acquire.policy_id

  name      = "Example Internet Access FQDN rule"
  action    = "forward"
  rule_type = "fqdn"
  ports     = ["80", "443"]
  protocol  = "tcp"

  destinations = [
    {
      type  = "fqdn"
      value = "example.com"
    }
  ]
}

resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "cidr" {
  forwarding_policy_id = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_bypass.policy_id

  name      = "Example Internet Access CIDR bypass rule"
  action    = "bypass"
  rule_type = "ip_subnet"
  ports     = ["443"]
  protocol  = "tcp"

  destinations = [
    {
      type  = "ip_subnet"
      value = "192.0.2.0/24"
    }
  ]
}

resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "ip_range" {
  forwarding_policy_id = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_bypass.policy_id

  name      = "Example Internet Access IP range bypass rule"
  action    = "bypass"
  rule_type = "ip_range"
  ports     = ["443"]
  protocol  = "udp"

  destinations = [
    {
      type          = "ip_range"
      begin_address = "192.0.2.10"
      end_address   = "192.0.2.20"
    }
  ]
}
