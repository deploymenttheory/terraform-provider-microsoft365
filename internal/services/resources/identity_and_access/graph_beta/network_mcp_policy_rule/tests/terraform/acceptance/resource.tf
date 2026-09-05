resource "random_id" "suffix" { byte_length = 4 }

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "test" {
  name           = "tf-api-probe-policy-${random_id.suffix.hex}"
  description    = "Initial description"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule" "test" {
  mcp_policy_id       = microsoft365_graph_beta_identity_and_access_network_mcp_policy.test.id
  name                = "tf-api-probe-rule-${random_id.suffix.hex}"
  description         = "Initial description"
  action              = "allow"
  priority            = 1000
  enabled             = true
  matching_conditions = { tool_matching = { names = { values = ["tf_probe_tool"], match_type = "exactMatch" }, methods = "call" } }
}
