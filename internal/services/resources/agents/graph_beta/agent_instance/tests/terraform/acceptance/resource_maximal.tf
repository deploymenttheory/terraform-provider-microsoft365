
########################################################################################
# Agent Instance Acceptance Test - Maximal Configuration
# This test creates an agent instance with all available fields populated
########################################################################################

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
  display_name        = "acc-test-agent-max-user1-${random_string.test_id_maximal.result}"
  user_principal_name = "acc-test-agent-max-user1-${random_string.test_id_maximal.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-max-user1-${random_string.test_id_maximal.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_users_user" "dependency_user_maximal_2" {
  display_name        = "acc-test-agent-max-user2-${random_string.test_id_maximal.result}"
  user_principal_name = "acc-test-agent-max-user2-${random_string.test_id_maximal.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-max-user2-${random_string.test_id_maximal.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

########################################################################################
# Test Resource - Agent Instance (Maximal)
########################################################################################

resource "microsoft365_graph_beta_agents_agent_instance" "test_maximal" {
  display_name = "IT Service Desk Agent - ${random_string.test_id_maximal.result}"
  owner_ids = [
    microsoft365_graph_beta_users_user.dependency_user_maximal_1.id,
    microsoft365_graph_beta_users_user.dependency_user_maximal_2.id
  ]
  originating_store   = "Deployment Theory"
  url                 = "https://servicedesk.deploymenttheory.com/api"
  preferred_transport = "HTTP+JSON"

  additional_interfaces = [
    {
      url       = "https://servicedesk.deploymenttheory.com/grpc"
      transport = "GRPC"
    },
    {
      url       = "https://servicedesk.deploymenttheory.com/jsonrpc"
      transport = "JSONRPC"
    }
  ]

  agent_card_manifest = {
    display_name                         = "IT Service Desk Agent"
    description                          = "An intelligent IT service desk agent that helps users troubleshoot common IT issues, submit support tickets, check ticket status, and find solutions from the knowledge base."
    protocol_version                     = "1.0"
    version                              = "2.0.0"
    icon_url                             = "https://servicedesk.example.com/assets/agent-icon.png"
    documentation_url                    = "https://docs.example.com/servicedesk-agent"
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
      organization = "Deployment Theory"
      url          = "https://www.deploymenttheory.com"
    }

    capabilities = {
      streaming                = true
      push_notifications       = true
      state_transition_history = false

      extensions = [
        {
          uri         = "https://servicedesk.example.com/a2a/capabilities/ticketing"
          description = "Integration with IT ticketing system for creating and managing support requests"
          required    = false
        },
      ]
    }

    skills = [
      {
        id           = "troubleshoot-issues"
        display_name = "IT Troubleshooter"
        description  = "Diagnose and provide solutions for common IT issues including password resets, VPN connectivity, printer problems, and software installation."

        tags = [
          "support",
          "troubleshooting",
          "it-help"
        ]

        examples = [
          "My VPN is not connecting",
          "How do I reset my password?",
          "My printer is not working"
        ]

        input_modes = [
          "application/json",
          "text/plain"
        ]

        output_modes = [
          "application/json",
          "application/vnd.geo+json",
          "text/html"
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
