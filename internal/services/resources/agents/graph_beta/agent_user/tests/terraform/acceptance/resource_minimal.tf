# Agent User Acceptance Test - Minimal Configuration
# This test creates an agent user from an agent identity with full dependency chain

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
  display_name        = "acc-test-agent-user-sponsor1-${random_string.test_id.result}"
  user_principal_name = "acc-test-agent-user-sponsor1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-user-sponsor1-${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Dependencies - Agent Identity Blueprint
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test" {
  display_name     = "acc-test-agent-user-blueprint-${random_string.test_id.result}"
  description      = "Blueprint for agent user acceptance test"
  sponsor_user_ids = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  owner_user_ids   = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  hard_delete      = true
}

########################################################################################
# Dependencies - Agent Identity Blueprint Service Principal
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal" "test" {
  app_id      = microsoft365_graph_beta_agents_agent_identity_blueprint.test.app_id
  hard_delete = true
}

########################################################################################
# Pause - Wait for blueprint service principal to propagate
########################################################################################

resource "time_sleep" "wait_for_blueprint_service_principal" {
  depends_on      = [microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.test]
  create_duration = "15s"
}

########################################################################################
# Dependencies - Agent Identity (required parent for agent user)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity" "test" {
  display_name                = "acc-test-agent-identity-${random_string.test_id.result}"
  agent_identity_blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test.app_id
  account_enabled             = true
  sponsor_ids                 = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  owner_ids                   = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  hard_delete                 = true

  depends_on = [
    time_sleep.wait_for_blueprint_service_principal
  ]
}

########################################################################################
# Pause - Wait for eventual consistency before creating agent user
########################################################################################

resource "time_sleep" "wait_for_agent_identity" {
  depends_on      = [microsoft365_graph_beta_agents_agent_identity.test]
  create_duration = "10s"
}

########################################################################################
# Test Resource - Agent User (Minimal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_user" "test_minimal" {
  display_name        = "acc-test-agent-user-${random_string.test_id.result}"
  agent_identity_id   = microsoft365_graph_beta_agents_agent_identity.test.id
  account_enabled     = true
  user_principal_name = "acc-test-agent-user-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-user-${random_string.test_id.result}"
  sponsor_ids         = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  hard_delete         = true

  depends_on = [
    time_sleep.wait_for_agent_identity
  ]
}
