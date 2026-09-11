resource "random_id" "suffix" { byte_length = 4 }

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "test" {
  name           = "tf-api-probe-policy-${random_id.suffix.hex}"
  description    = "Initial description"
  default_action = "allow"
}
