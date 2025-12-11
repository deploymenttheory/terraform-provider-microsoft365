# Maximal configuration for update testing
resource "microsoft365_graph_beta_agents_agent_instance" "test_update" {
  display_name        = "Update Test Agent Maximal"
  owner_ids           = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
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
    display_name                         = "Update Test Agent Card Maximal"
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
}
