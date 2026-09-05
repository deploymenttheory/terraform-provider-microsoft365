resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy" "test" {
  name = "tf-api-probe-policy-updated"

  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy_rule" "test" {
  prompt_policy_id = microsoft365_graph_beta_identity_and_access_network_prompt_policy.test.id
  name             = "tf-api-probe-rule-updated"

  action               = "block"
  priority             = 65001
  enabled              = false
  prompt_logging       = "onBlock"
  conversation_schemes = [{ type = "predefined", scheme_name = "chatGpt" }]
}
