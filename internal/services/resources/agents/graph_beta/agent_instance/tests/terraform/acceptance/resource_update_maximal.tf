# Agent Instance Acceptance Test - Update Maximal Configuration
# This configuration is used for update testing (maximal state)

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

resource "microsoft365_graph_beta_users_user" "dependency_user_update_2" {
  display_name        = "acc-test-agent-update-user2-${random_string.test_id_update.result}"
  user_principal_name = "acc-test-agent-update-user2-${random_string.test_id_update.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-update-user2-${random_string.test_id_update.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Test Resource - Agent Instance (Update Maximal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_instance" "test_update" {
  display_name = "acc-test-agent-update-maximal-${random_string.test_id_update.result}"
  owner_ids = [
    microsoft365_graph_beta_users_user.dependency_user_update_1.id,
    microsoft365_graph_beta_users_user.dependency_user_update_2.id
  ]
  originating_store   = "Terraform"
  url                 = "https://updated-agent.example.com/api"
  preferred_transport = "HTTP+JSON"

  additional_interfaces = [
    {
      url       = "https://updated-agent.example.com/grpc"
      transport = "GRPC"
    }
  ]

  agent_card_manifest = {
    display_name                         = "acc-test-update-agent-card-maximal-${random_string.test_id_update.result}"
    description                          = "Maximal configuration for update testing with all fields"
    protocol_version                     = "1.0"
    version                              = "2.0.0"
    icon_url                             = "https://updated-agent.example.com/icon.png"
    documentation_url                    = "https://docs.example.com/updated-agent"
    supports_authenticated_extended_card = false

    default_input_modes = [
      "application/json",
      "text/plain"
    ]

    default_output_modes = [
      "application/json",
      "text/html"
    ]

    provider = {
      organization = "Test Organization"
      url          = "https://www.example.com"
    }

    capabilities = {
      streaming                = true
      push_notifications       = true
      state_transition_history = false

      extensions = [
        {
          uri         = "https://example.com/extension"
          description = "Test extension"
          required    = false
        }
      ]
    }

    skills = [
      {
        id           = "test-skill"
        display_name = "Test Skill"
        description  = "A test skill for update testing"

        tags = [
          "test",
          "update"
        ]

        examples = [
          "Run the test skill"
        ]

        input_modes = [
          "application/json"
        ]

        output_modes = [
          "application/json"
        ]
      }
    ]
  }

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
