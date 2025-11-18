resource "microsoft365_graph_beta_identity_and_access_network_filtering_policy" "block_policy" {
  name        = "Block Malicious Traffic"
  description = "Policy to block malicious network traffic"
  action      = "block"
}

resource "microsoft365_graph_beta_identity_and_access_network_filtering_policy" "allow_policy" {
  name        = "Allow Traffic from Trusted Sources"
  description = "Policy to allow traffic from trusted sources"
  action      = "allow"
}

