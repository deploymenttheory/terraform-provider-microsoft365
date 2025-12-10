# Acceptance test: Minimal Agent Identity Blueprint Identifier URI configuration
# Full dependency chain: random_string -> users -> agent_identity_blueprint -> identifier_uri

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "dependency_user_1" {
  display_name        = "acc-test-id-uri-user1-${random_string.test_id.result}"
  user_principal_name = "acc-test-id-uri-user1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-id-uri-user1-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_blueprint" {
  display_name = "acc-test-blueprint-id-uri-${random_string.test_id.result}"
  description  = "Agent identity blueprint for identifier URI acceptance test"
  sponsor_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_1.id,
  ]
  owner_user_ids = [
    microsoft365_graph_beta_users_user.dependency_user_1.id,
  ]
  hard_delete = true
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri" "test_minimal" {
  blueprint_id   = microsoft365_graph_beta_agents_agent_identity_blueprint.test_blueprint.id
  identifier_uri = "api://${microsoft365_graph_beta_agents_agent_identity_blueprint.test_blueprint.id}"

  scope = {
    admin_consent_description  = "Allow the application to access the agent on behalf of the signed-in user."
    admin_consent_display_name = "Access agent"
    is_enabled                 = true
    type                       = "User"
    value                      = "access_agent"
  }
}

