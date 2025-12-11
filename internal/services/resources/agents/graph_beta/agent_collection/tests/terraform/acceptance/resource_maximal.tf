# Agent Collection Acceptance Test - Maximal Configuration
# This test creates an agent collection with all available fields populated

########################################################################################
# Dependencies - Random string for unique naming
########################################################################################

resource "random_string" "test_id_maximal" {
  length  = 8
  special = false
  upper   = false
}

########################################################################################
# Dependencies - Users for owners
########################################################################################

resource "microsoft365_graph_beta_users_user" "dependency_user_maximal_1" {
  display_name        = "acc-test-agent-col-max-user1-${random_string.test_id_maximal.result}"
  user_principal_name = "acc-test-agent-col-max-user1-${random_string.test_id_maximal.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-col-max-user1-${random_string.test_id_maximal.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_users_user" "dependency_user_maximal_2" {
  display_name        = "acc-test-agent-col-max-user2-${random_string.test_id_maximal.result}"
  user_principal_name = "acc-test-agent-col-max-user2-${random_string.test_id_maximal.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-col-max-user2-${random_string.test_id_maximal.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Test Resource - Agent Collection (Maximal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_collection" "test_maximal" {
  display_name = "IT Automation Agent Collection - ${random_string.test_id_maximal.result}"
  owner_ids = [
    microsoft365_graph_beta_users_user.dependency_user_maximal_1.id,
    microsoft365_graph_beta_users_user.dependency_user_maximal_2.id
  ]
  description       = "A collection of IT automation agents for managing infrastructure and support workflows"
  originating_store = "Deployment Theory"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
