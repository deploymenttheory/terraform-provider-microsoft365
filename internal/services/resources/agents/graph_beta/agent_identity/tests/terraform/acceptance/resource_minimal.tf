# Agent Identity Acceptance Test - Minimal Configuration
# This test creates an agent identity from an agent identity blueprint with full dependency chain

########################################################################################
# Dependencies - Random string for unique naming
########################################################################################

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

########################################################################################
# Dependencies - Users for sponsors and owners
########################################################################################

resource "microsoft365_graph_beta_users_user" "dependency_user_1" {
  display_name        = "acc-test-agent-identity-user1-${random_string.test_id.result}"
  user_principal_name = "acc-test-agent-identity-user1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-identity-user1-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Dependencies - Agent Identity Blueprint
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test" {
  display_name     = "acc-test-agent-identity-blueprint-${random_string.test_id.result}"
  description      = "Blueprint for agent identity acceptance test"
  sponsor_user_ids = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  owner_user_ids   = [microsoft365_graph_beta_users_user.dependency_user_1.id]
}

########################################################################################
# Dependencies - Agent Identity Blueprint Service Principal
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal" "test" {
  app_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test.app_id
}

########################################################################################
# Test Resource - Agent Identity
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity" "test_minimal" {
  display_name                = "acc-test-agent-identity-${random_string.test_id.result}"
  agent_identity_blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test.app_id
  account_enabled             = true
  sponsor_ids                 = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  owner_ids                   = [microsoft365_graph_beta_users_user.dependency_user_1.id]

  depends_on = [
    microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.test
  ]
}
