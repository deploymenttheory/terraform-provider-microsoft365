# Agent Instance Acceptance Test - Update Minimal Configuration
# This configuration is used for update testing (minimal state)

########################################################################################
# Dependencies - Random string for unique naming
########################################################################################

resource "random_string" "test_id_update" {
  length  = 8
  special = false
  upper   = false
}

########################################################################################
# Dependencies - Users for owners
########################################################################################

resource "microsoft365_graph_beta_users_user" "dependency_user_update_1" {
  display_name        = "acc-test-agent-update-user1-${random_string.test_id_update.result}"
  user_principal_name = "acc-test-agent-update-user1-${random_string.test_id_update.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-update-user1-${random_string.test_id_update.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Test Resource - Agent Instance (Update Minimal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_instance" "test_update" {
  display_name      = "acc-test-agent-update-${random_string.test_id_update.result}"
  owner_ids         = [microsoft365_graph_beta_users_user.dependency_user_update_1.id]
  originating_store = "Terraform"

  agent_card_manifest = {
    display_name                         = "acc-test-update-agent-card-${random_string.test_id_update.result}"
    description                          = "Minimal configuration for update testing"
    protocol_version                     = "1.0"
    version                              = "1.0.0"
    supports_authenticated_extended_card = false

    capabilities = {
      streaming                = false
      push_notifications       = false
      state_transition_history = false
    }
  }
}
