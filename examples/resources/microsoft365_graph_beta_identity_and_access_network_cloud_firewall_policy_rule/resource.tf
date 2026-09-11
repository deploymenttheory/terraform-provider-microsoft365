resource "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy" "example" {
  name           = "example-cloud-firewall-policy"
  description    = "Managed by Terraform"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy_rule" "example" {
  policy_id   = microsoft365_graph_beta_identity_and_access_network_cloud_firewall_policy.example.id
  name        = "example-cloud-firewall-rule"
  description = "Match example destination addresses"
  priority    = 100
  action      = "block"
  enabled     = false

  sources = {
    addresses = [{
      type   = "ip"
      values = ["192.0.2.0/24"]
    }]
    ports = ["1024-65535"]
  }

  destinations = {
    addresses = [{
      type   = "ip"
      values = ["198.51.100.1", "198.51.100.2"]
    }]
    ports     = ["443"]
    protocols = ["tcp", "udp"]
  }
}
