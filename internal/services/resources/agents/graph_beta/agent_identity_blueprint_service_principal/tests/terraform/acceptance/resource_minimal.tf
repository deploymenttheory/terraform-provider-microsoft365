# Minimal Agent Identity Blueprint Service Principal configuration for acceptance testing
# Note: This requires an existing agent identity blueprint to be created first

########################################################################################
# Dependencies
########################################################################################
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


# First create an agent identity blueprint (dependency)
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_dependency" {
  display_name = "acc-test-agent-identity-blueprint-sp-dependency-${random_string.test_id.result}"
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

########################################################################################
# Test Resource
########################################################################################
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal" "test_minimal" {
  app_id      = microsoft365_graph_beta_agents_agent_identity_blueprint.test_dependency.app_id
  hard_delete = true
}

