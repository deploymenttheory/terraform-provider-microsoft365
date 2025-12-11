# Agent Instance Acceptance Test - Minimal Configuration
# This test creates an agent instance with minimal required fields

########################################################################################
# Dependencies - Random string for unique naming
########################################################################################

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

########################################################################################
# Dependencies - Users for owners
########################################################################################

resource "microsoft365_graph_beta_users_user" "dependency_user_1" {
  display_name        = "acc-test-agent-instance-user1-${random_string.test_id.result}"
  user_principal_name = "acc-test-agent-instance-user1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-instance-user1-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Test Resource - Agent Instance (Minimal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_instance" "test_minimal" {
  display_name      = "acc-test-agent-instance-1-${random_string.test_id.result}"
  owner_ids         = [microsoft365_graph_beta_users_user.dependency_user_1.id]
  originating_store = "Terraform"

  agent_card_manifest = {
    display_name                         = "acc-test-agent-card-${random_string.test_id.result}"
    description                          = "Acceptance test agent card manifest description"
    protocol_version                     = "1.0"
    version                              = "1.0.1"
    supports_authenticated_extended_card = false

    capabilities = {
      streaming                = true
      push_notifications       = false
      state_transition_history = false
    }
  }
}
