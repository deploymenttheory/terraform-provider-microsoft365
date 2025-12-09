# Minimal Agent Identity configuration
resource "microsoft365_graph_beta_agents_agent_identity" "test_minimal" {
  display_name                = "Unit Test Agent Identity"
  agent_identity_blueprint_id = "11111111-1111-1111-1111-111111111111"
  account_enabled             = true
  sponsor_ids                 = ["22222222-2222-2222-2222-222222222222"]
  owner_ids                   = ["33333333-3333-3333-3333-333333333333"]
}
