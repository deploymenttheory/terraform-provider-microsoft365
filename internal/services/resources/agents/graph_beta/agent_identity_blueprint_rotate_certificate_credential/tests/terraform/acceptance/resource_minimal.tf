# Note: This acceptance test requires:
# 1. An existing agent identity blueprint with at least one valid certificate
# 2. A proof of possession JWT token signed with the private key of an existing certificate
#
# Agent identity blueprints that don't have any existing valid certificates cannot use
# the addKey/removeKey service actions. Use the Update agent identity blueprint operation instead.
#
# For more information on generating proof of possession tokens, see:
# https://learn.microsoft.com/en-us/graph/application-rollkey-prooftoken

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

data "microsoft365_graph_beta_identity_governance_user" "owner" {
  user_principal_name = "admin@yourtenant.onmicrosoft.com"
}

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "test_blueprint" {
  display_name                  = "acc-test-key-cred-${random_string.suffix.result}"
  description                   = "Test blueprint for key credential acceptance testing"
  unique_name                   = "acc-test-key-cred-${random_string.suffix.result}"
  unique_name_include_app_id    = false
  visibility                    = "managedTenantsOnly"
  owner_user_ids                = [data.microsoft365_graph_beta_identity_governance_user.owner.id]
  sponsor_user_ids              = [data.microsoft365_graph_beta_identity_governance_user.owner.id]
  allow_external_token_issuance = false
  application_audience          = "AzureADMyOrg"
}

# Note: This test requires manually providing a valid key and proof token
# as they must be generated from an existing certificate on the blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_key_credential" "test_minimal" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_blueprint.id
  display_name = "acc-test-key-credential-${random_string.suffix.result}"
  type         = "AsymmetricX509Cert"
  usage        = "Verify"

  # These values must be replaced with actual Base64-encoded certificate and JWT proof token
  # The key should be the public key of the certificate being added
  # The proof should be a JWT signed with an existing valid certificate's private key
  key   = var.certificate_public_key
  proof = var.proof_of_possession_token
}

variable "certificate_public_key" {
  description = "Base64-encoded public key of the certificate to add"
  type        = string
  sensitive   = true
}

variable "proof_of_possession_token" {
  description = "JWT proof of possession token signed with an existing valid certificate"
  type        = string
  sensitive   = true
}
