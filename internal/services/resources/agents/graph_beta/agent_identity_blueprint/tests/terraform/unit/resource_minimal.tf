# Minimal Agent Identity Blueprint configuration
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_minimal" {
  display_name     = "unit-test-agent-identity-blueprint-minimal"
  sponsor_user_ids = ["11111111-1111-1111-1111-111111111111"]
  owner_user_ids   = ["11111111-1111-1111-1111-111111111111"]
}
