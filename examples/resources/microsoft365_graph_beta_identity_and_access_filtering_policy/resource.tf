resource "microsoft365_graph_beta_identity_and_access_filtering_policy" "block_policy" {
  name        = "Block Malicious Traffic"
  description = "Policy to block malicious network traffic"
  action      = "block"
}
