# Minimal configuration for update testing
resource "microsoft365_graph_beta_agents_agent_instance" "test_update" {
  display_name      = "Update Test Agent Minimal"
  owner_ids         = ["11111111-1111-1111-1111-111111111111"]
  originating_store = "Terraform"

  agent_card_manifest = {
    display_name                         = "Update Test Agent Card Minimal"
    description                          = "Minimal configuration for update testing"
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
