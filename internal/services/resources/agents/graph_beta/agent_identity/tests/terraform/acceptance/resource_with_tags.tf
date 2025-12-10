# Agent Identity Acceptance Test - With Tags Configuration
# This test creates an agent identity with tags from an agent identity blueprint

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
  display_name        = "acc-test-agent-identity-tags-user1-${random_string.test_id.result}"
  user_principal_name = "acc-test-agent-identity-tags-user1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-identity-tags-user1-${random_string.test_id.result}"
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
  display_name     = "acc-test-agent-identity-blueprint-tags-${random_string.test_id.result}"
  description      = "Blueprint for agent identity with tags acceptance test"
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
# Test Resource - Agent Identity with Tags
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity" "test_with_tags" {
  display_name                = "acc-test-agent-identity-tags-${random_string.test_id.result}"
  agent_identity_blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test.app_id
  account_enabled             = true
  sponsor_ids                 = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  owner_ids                   = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  tags                        = ["terraform", "acceptance-test", "with-tags"]
  hard_delete                 = true

  depends_on = [
    microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.test
  ]
}

