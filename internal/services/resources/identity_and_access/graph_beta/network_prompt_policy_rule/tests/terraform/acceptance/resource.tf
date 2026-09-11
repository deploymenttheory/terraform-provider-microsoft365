resource "random_id" "suffix" { byte_length = 4 }

resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy" "test" {
  name           = "tf-api-probe-policy-${random_id.suffix.hex}"
  description    = "Initial description"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy_rule" "test" {
  prompt_policy_id     = microsoft365_graph_beta_identity_and_access_network_prompt_policy.test.id
  name                 = "tf-api-probe-rule-${random_id.suffix.hex}"
  description          = "Initial description"
  action               = "allow"
  priority             = 1000
  enabled              = true
  conversation_schemes = [{ type = "custom", url = "https://example.com/chat", json_path = "$.prompt" }]
}
