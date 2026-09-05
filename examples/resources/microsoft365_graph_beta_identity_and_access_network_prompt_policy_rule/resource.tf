# Creating policies and rules does not assign them to traffic.
resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy" "example" {
  name        = "Example prompt policy"
  description = "Prompt protection for AI applications"
}

resource "microsoft365_graph_beta_identity_and_access_network_prompt_policy_rule" "example" {
  prompt_policy_id = microsoft365_graph_beta_identity_and_access_network_prompt_policy.example.id
  name             = "Block malicious prompts"
  priority         = 1000
  enabled          = true
  action           = "block"
  prompt_logging   = "onBlock"
  conversation_schemes = [
    {
      type        = "predefined"
      scheme_name = "chatGpt"
    },
    {
      type      = "custom"
      url       = "https://ai.example.com/chat"
      json_path = "$.prompt"
    }
  ]
}
