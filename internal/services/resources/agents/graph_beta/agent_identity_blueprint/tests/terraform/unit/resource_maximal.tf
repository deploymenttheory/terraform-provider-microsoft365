# Maximal Agent Identity Blueprint configuration
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_maximal" {
  display_name = "unit-test-agent-identity-blueprint-maximal"
  description  = "This is a test agent identity blueprint with all optional fields configured"

  sponsor_user_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
  ]
  owner_user_ids = [
    "22222222-2222-2222-2222-222222222222",
    "44444444-4444-4444-4444-444444444444",
  ]
  tags = [
    "terraform",
    "test",
    "agent"
  ]
}
