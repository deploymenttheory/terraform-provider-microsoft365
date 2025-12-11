# Agent Collection Acceptance Test - Minimal Configuration
# This test creates an agent collection with minimal required fields

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
  display_name        = "acc-test-agent-collection-user1-${random_string.test_id.result}"
  user_principal_name = "acc-test-agent-collection-user1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-collection-user1-${random_string.test_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Dependencies - Agent Instance (Minimal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_instance" "test_minimal" {
  display_name = "acc-test-agent-instance-${random_string.test_id.result}"
  owner_ids    = [microsoft365_graph_beta_users_user.dependency_user_1.id]

  lifecycle {
    ignore_changes = [agent_card_manifest]
  }
}

########################################################################################
# Dependencies - Agent Collection (Minimal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_collection" "test_minimal" {
  display_name = "acc-test-agent-collection-${random_string.test_id.result}"
  owner_ids    = [microsoft365_graph_beta_users_user.dependency_user_1.id]
}

########################################################################################
# Test Resource - Agent Collection Assignment (Minimal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_collection_assignment" "test_minimal" {
  agent_instance_id   = microsoft365_graph_beta_agents_agent_instance.test_minimal.id
  agent_collection_id = microsoft365_graph_beta_agents_agent_collection.test_minimal.id
}