# Agent Collection Acceptance Test - Update Minimal Configuration
# This test creates an agent collection with minimal fields for update testing

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
  display_name        = "acc-test-agent-col-upd-user1-${random_string.test_id_update.result}"
  user_principal_name = "acc-test-agent-col-upd-user1-${random_string.test_id_update.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-col-upd-user1-${random_string.test_id_update.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Test Resource - Agent Collection (Update Minimal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_collection" "test_update" {
  display_name = "acc-test-agent-col-update-${random_string.test_id_update.result}"
  owner_ids    = [microsoft365_graph_beta_users_user.dependency_user_update_1.id]
}
