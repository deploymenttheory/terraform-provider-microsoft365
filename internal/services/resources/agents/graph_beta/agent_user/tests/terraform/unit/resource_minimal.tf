# Minimal Agent User configuration - only required fields
resource "microsoft365_graph_beta_agents_agent_user" "test_minimal" {
  display_name        = "Unit Test Agent User"
  agent_identity_id   = "11111111-1111-1111-1111-111111111111"
  account_enabled     = true
  user_principal_name = "testagentuser@contoso.onmicrosoft.com"
  mail_nickname       = "testagentuser"
  sponsor_ids         = ["22222222-2222-2222-2222-222222222222"]
  hard_delete         = true
}
