# Minimal Agent Instance configuration
# Creates an agent instance with required fields only
resource "microsoft365_graph_beta_agents_agent_instance" "minimal" {
  display_name      = "My Agent Instance"
  owner_ids         = ["00000000-0000-0000-0000-000000000000"]
  originating_store = "Terraform"

  agent_card_manifest = {
    display_name                         = "My Agent Card"
    description                          = "A minimal agent card manifest description"
    protocol_version                     = "1.0"
    version                              = "1.0.0"
    supports_authenticated_extended_card = false

    capabilities = {
      streaming                = false
      push_notifications       = false
      state_transition_history = false
    }
  }
}
