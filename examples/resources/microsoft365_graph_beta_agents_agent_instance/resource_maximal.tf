# Maximal Agent Instance configuration
# Creates an agent instance with all available fields configured
resource "microsoft365_graph_beta_agents_agent_instance" "maximal" {
  display_name = "IT Service Desk Agent"
  owner_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]
  originating_store   = "Deployment Theory"
  url                 = "https://servicedesk.example.com/api"
  preferred_transport = "HTTP+JSON"

  # Optional: Link to agent identity resources
  # source_agent_id             = "00000000-0000-0000-0000-000000000000"
  # agent_identity_blueprint_id = "00000000-0000-0000-0000-000000000000"
  # agent_identity_id           = "00000000-0000-0000-0000-000000000000"
  # managed_by                  = "00000000-0000-0000-0000-000000000000"

  additional_interfaces = [
    {
      url       = "https://servicedesk.example.com/grpc"
      transport = "GRPC"
    },
    {
      url       = "https://servicedesk.example.com/jsonrpc"
      transport = "JSONRPC"
    }
  ]

  agent_card_manifest = {
    display_name                         = "IT Service Desk Agent"
    description                          = "An intelligent IT service desk agent that helps users troubleshoot common IT issues, submit support tickets, check ticket status, and find solutions from the knowledge base."
    protocol_version                     = "1.0"
    version                              = "2.0.0"
    supports_authenticated_extended_card = false

    # Note: Once set, these fields cannot be removed without recreating the resource
    icon_url          = "https://servicedesk.example.com/assets/agent-icon.png"
    documentation_url = "https://docs.example.com/servicedesk-agent"

    # Note: Once set, these fields cannot be removed without recreating the resource
    default_input_modes = [
      "application/json",
      "text/plain"
    ]

    default_output_modes = [
      "application/json",
      "text/html"
    ]

    # Note: Once set, this block cannot be removed without recreating the resource
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
        }
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
