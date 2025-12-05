# Unit test: Minimal Agent Identity Blueprint Identifier URI configuration

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri" "test_minimal" {
  blueprint_id   = "11111111-1111-1111-1111-111111111111"
  identifier_uri = "api://22222222-2222-2222-2222-222222222222"

  scope = {
    value = "access_agent"
  }
}
