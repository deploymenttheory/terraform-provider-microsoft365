resource "microsoft365_graph_beta_users_agent_user" "minimal" {
  display_name        = "unit-test-agent-user-minimal"
  user_principal_name = "unit-test-agent-user-minimal@deploymenttheory.com"
  mail_nickname       = "unit-test-agent-user-minimal"
  account_enabled     = true
  identity_parent_id  = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}

