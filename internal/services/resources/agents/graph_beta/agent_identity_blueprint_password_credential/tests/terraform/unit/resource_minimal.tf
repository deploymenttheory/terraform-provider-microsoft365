# Unit test: Minimal Agent Identity Blueprint Password Credential configuration
# Dependencies: agent_identity_blueprint

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_blueprint" {
  display_name     = "unit-test-agent-identity-blueprint-minimal"
  sponsor_user_ids = ["11111111-1111-1111-1111-111111111111"]
  owner_user_ids   = ["11111111-1111-1111-1111-111111111111"]
  hard_delete      = true
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_password_credential" "test_minimal" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_blueprint.id
  display_name = "unit-test-password-credential"
}

