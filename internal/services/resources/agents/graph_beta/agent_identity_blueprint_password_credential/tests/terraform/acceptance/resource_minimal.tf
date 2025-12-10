# Acceptance test: Minimal Agent Identity Blueprint Password Credential configuration
# Full dependency chain: random_string -> users -> agent_identity_blueprint -> password_credential

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "dependency_user_1" {
  display_name        = "acc-test-pwd-cred-user1-${random_string.test_id.result}"
  user_principal_name = "acc-test-pwd-cred-user1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-pwd-cred-user1-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_users_user" "dependency_user_2" {
  display_name        = "acc-test-pwd-cred-user2-${random_string.test_id.result}"
  user_principal_name = "acc-test-pwd-cred-user2-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-pwd-cred-user2-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_blueprint" {
  display_name = "acc-test-blueprint-pwd-cred-${random_string.test_id.result}"
  description  = "Agent identity blueprint for password credential acceptance test"
  sponsor_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_1.id,
    microsoft365_graph_beta_users_user.dependency_user_2.id,
  ]
  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_1.id,
    microsoft365_graph_beta_users_user.dependency_user_2.id,
  ]
  hard_delete = true
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential" "test_minimal" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_blueprint.id
  display_name = "acc-test-password-credential-${random_string.test_id.result}"
}


