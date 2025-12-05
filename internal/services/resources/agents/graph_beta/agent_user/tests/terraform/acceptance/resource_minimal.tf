resource "random_string" "minimal_user_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_agent_user" "minimal" {
  display_name        = "acc-test-agent-user-minimal-${random_string.minimal_user_id.result}"
  user_principal_name = "acc-test-agent-user-minimal-${random_string.minimal_user_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-user-minimal-${random_string.minimal_user_id.result}"
  account_enabled     = true
  identity_parent_id  = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}

