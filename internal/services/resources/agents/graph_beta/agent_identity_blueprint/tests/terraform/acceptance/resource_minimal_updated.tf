# Updated Minimal Agent Identity Blueprint configuration for acceptance testing
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_minimal" {
  display_name = "acc-test-agent-identity-blueprint-minimal-updated-${random_string.test_id.result}"
  description  = "Updated description for acceptance test"

  sponsor_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_1.id,
  ]
  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_1.id,
  ]
  hard_delete = true
}

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "dependency_user_1" {
  display_name        = "acc-test-blueprint-user1-${random_string.test_id.result}"
  user_principal_name = "acc-test-blueprint-user1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-blueprint-user1-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

resource "microsoft365_graph_beta_users_user" "dependency_user_2" {
  display_name        = "acc-test-blueprint-user2-${random_string.test_id.result}"
  user_principal_name = "acc-test-blueprint-user2-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-blueprint-user2-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
  hard_delete = true
}

