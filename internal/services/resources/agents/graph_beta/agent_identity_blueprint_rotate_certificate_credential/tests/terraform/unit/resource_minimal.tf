provider "microsoft365" {
  use_mock = true
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_blueprint" {
  display_name                  = "test-blueprint-unit"
  description                   = "Test blueprint for unit testing"
  unique_name                   = "test-blueprint-unit"
  unique_name_include_app_id    = false
  visibility                    = "managedTenantsOnly"
  owner_user_ids                = []
  sponsor_user_ids              = []
  allow_external_token_issuance = false
  application_audience          = "AzureADMyOrg"
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_key_credential" "test_minimal" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_blueprint.id
  display_name = "unit-test-key-credential"
  type         = "AsymmetricX509Cert"
  usage        = "Verify"
  key          = "MIIDXTCCAkWgAwIBAgIJAJC1HiI..."                     # Base64 encoded public key
  proof        = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6IjJ..." # JWT proof token
}
