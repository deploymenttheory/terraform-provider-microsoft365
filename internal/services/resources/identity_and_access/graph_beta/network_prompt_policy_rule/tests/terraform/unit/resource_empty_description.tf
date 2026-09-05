resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy" "test" {
  description = ""
  name        = "tf-api-probe-policy-updated"

  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy_rule" "test" {
  prompt_policy_id = microsoft365_graph_beta_identity_and_access_network_prompt_policy.test.id
  description      = ""
  name             = "tf-api-probe-rule-updated"

  action               = "block"
  priority             = 65001
  enabled              = false
  conversation_schemes = [{ type = "custom", url = "https://example.com/chat", json_path = "$.prompt" }]
}
