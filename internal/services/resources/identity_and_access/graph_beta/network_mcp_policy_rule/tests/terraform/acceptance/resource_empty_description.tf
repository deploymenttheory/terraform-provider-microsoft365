resource "random_id" "suffix" { byte_length = 4 }

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "test" {
  description = ""
  name        = "tf-api-probe-policy-${random_id.suffix.hex}-updated"

  default_action = "block"
}

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule" "test" {
  mcp_policy_id = microsoft365_graph_beta_identity_and_access_network_mcp_policy.test.id
  description   = ""
  name          = "tf-api-probe-rule-${random_id.suffix.hex}-updated"

  action              = "block"
  priority            = 65001
  enabled             = false
  matching_conditions = { server_urls = { values = ["https://example.invalid/mcp"], match_type = "contains" } }
}
