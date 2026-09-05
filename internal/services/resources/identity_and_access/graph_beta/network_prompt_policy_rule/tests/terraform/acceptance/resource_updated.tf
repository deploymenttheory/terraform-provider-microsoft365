resource "random_id" "suffix" { byte_length = 4 }

resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy" "test" {
  name = "tf-api-probe-policy-${random_id.suffix.hex}-updated"

  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy_rule" "test" {
  prompt_policy_id = microsoft365_graph_beta_identity_and_access_network_prompt_policy.test.id
  name             = "tf-api-probe-rule-${random_id.suffix.hex}-updated"

  action               = "block"
  priority             = 65001
  enabled              = false
  prompt_logging       = "onBlock"
  conversation_schemes = [{ type = "predefined", scheme_name = "chatGpt" }]
}
