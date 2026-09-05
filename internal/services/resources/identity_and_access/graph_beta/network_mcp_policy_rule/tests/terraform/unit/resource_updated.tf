resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "test" {
  name = "tf-api-probe-policy-updated"

  default_action = "block"
}

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule" "test" {
  mcp_policy_id = microsoft365_graph_beta_identity_and_access_network_mcp_policy.test.id
  name          = "tf-api-probe-rule-updated"

  action              = "block"
  priority            = 65001
  enabled             = false
  matching_conditions = { server_urls = { values = ["https://example.invalid/mcp"], match_type = "contains" } }
}
