---
page_title: "microsoft365_graph_beta_agents_agent_instance Resource - terraform-provider-microsoft365"
subcategory: "Agents"

description: |-
  Manages an Agent Instance in the Microsoft Entra Agent Registry using the /agentRegistry/agentInstances endpoint. This resource is used to represent a specific deployed instance of an AI agent. Agent instances can be associated with an agentCardManifest that defines its capabilities, skills, and metadata.
  For more information, see the agentInstance resource type https://learn.microsoft.com/en-us/graph/api/resources/agentinstance?view=graph-rest-beta.
---

# microsoft365_graph_beta_agents_agent_instance (Resource)

Manages an Agent Instance in the Microsoft Entra Agent Registry using the `/agentRegistry/agentInstances` endpoint. This resource is used to represent a specific deployed instance of an AI agent. Agent instances can be associated with an agentCardManifest that defines its capabilities, skills, and metadata.

For more information, see the [agentInstance resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentinstance?view=graph-rest-beta).

## Microsoft Documentation

- [agentInstance resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentinstance?view=graph-rest-beta)
- [Create agentInstance](https://learn.microsoft.com/en-us/graph/api/agentregistry-post-agentinstances?view=graph-rest-beta&tabs=http)
- [Get agentInstance](https://learn.microsoft.com/en-us/graph/api/agentinstance-get?view=graph-rest-beta&tabs=http)
- [Update agentInstance](https://learn.microsoft.com/en-us/graph/api/agentinstance-update?view=graph-rest-beta&tabs=http)
- [Delete agentInstance](https://learn.microsoft.com/en-us/graph/api/agentregistry-delete-agentinstances?view=graph-rest-beta&tabs=http)
- [agentCardManifest resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentcardmanifest?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `AgentInstance.Read.All`
- `AgentCardManifest.Read.All`
- `AgentInstance.ReadWrite.All`
- `AgentInstance.ReadWrite.ManagedBy`
- `AgentCardManifest.ReadWrite.All`
- `AgentCardManifest.ReadWrite.ManagedBy`

**Optional:**
- `None` `[N/A]`

Find out more about the permissions required for managing agent instances at microsoft learn [here](https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.38.0 | Experimental | Initial release |

## Important Notes

### Fields That Require Resource Recreation

The following fields in `agent_card_manifest` cannot be removed or cleared once set. Changing or removing these values will trigger resource recreation:

- `icon_url`
- `documentation_url`
- `default_input_modes`
- `default_output_modes`
- `provider`

Additionally, the top-level `originating_store` field also requires resource recreation when changed.

## Example Usage

### Minimal Example

```terraform
# Minimal Agent Instance configuration
# Creates an agent instance with required fields only
resource "microsoft365_graph_beta_agents_agent_instance" "minimal" {
  display_name      = "My Agent Instance"
  owner_ids         = ["00000000-0000-0000-0000-000000000000"]
  originating_store = "Terraform"

  agent_card_manifest = {
    display_name                         = "My Agent Card"
    description                          = "A minimal agent card manifest description"
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
```

### Maximal Example

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `agent_card_manifest` (Attributes) The agent card manifest of the agent instance. (see [below for nested schema](#nestedatt--agent_card_manifest))
- `display_name` (String) Display name for the agent instance.
- `owner_ids` (Set of String) List of object IDs for the owners of the agent instance.

### Optional

- `additional_interfaces` (Attributes List) Additional interfaces/transports supported by the agent. (see [below for nested schema](#nestedatt--additional_interfaces))
- `agent_identity_blueprint_id` (String) Object ID of the agentIdentityBlueprint object.
- `agent_identity_id` (String) Object ID of the agentIdentity object.
- `managed_by` (String) **appId** (referred to as **Application (client) ID** on the Microsoft Entra admin center) of the application managing this agent.
- `originating_store` (String) Name of the store/system where agent originated. For example Copilot Studio, or Microsoft Security Copilot etc. Changing this value will force resource recreation.
- `preferred_transport` (String) Preferred transport protocol. The possible values are `JSONRPC`, `GRPC`, and `HTTP+JSON`.
- `source_agent_id` (String) Identifier of the agent in the original source system.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `url` (String) Endpoint URL for the agent instance.

### Read-Only

- `agent_user_id` (String) Object ID of the agentUser associated with the agent. Read-only.
- `created_by` (String) Object ID of the user or application that created the agent instance. Read-only.
- `created_date_time` (String) Timestamp when agent instance was created. Read-only.
- `id` (String) Unique identifier for the agent instance. Key. Inherited from entity.
- `last_modified_date_time` (String) Timestamp of last modification.

<a id="nestedatt--agent_card_manifest"></a>
### Nested Schema for `agent_card_manifest`

Required:

- `capabilities` (Attributes) Capabilities of the agent. (see [below for nested schema](#nestedatt--agent_card_manifest--capabilities))
- `description` (String) Description of the agent card manifest.
- `display_name` (String) Display name for the agent card manifest.
- `protocol_version` (String) Protocol version for the agent card. Must be in either the Major.Minor versioning format X.Y (e.g., 1.0, 2.1) or the semantic versioning format X.Y.Z (e.g., 1.0.0, 2.1.3)
- `supports_authenticated_extended_card` (Boolean) Whether the agent supports authenticated extended card.
- `version` (String) Version of the agent card manifest. Must be in the semantic versioning format X.Y.Z (e.g., 1.0.0, 2.1.3)

Optional:

- `default_input_modes` (Set of String) Default input modes supported by the agent. Changing or removing this value requires resource recreation.
- `default_output_modes` (Set of String) Default output modes supported by the agent. Changing or removingthis value requires resource recreation.
- `documentation_url` (String) URL to the documentation for the agent. Changing or removing this value requires resource recreation.
- `icon_url` (String) URL to the icon for the agent. Changing or removingthis value requires resource recreation.
- `originating_store` (String) Name of the store/system where the manifest originated.
- `owner_ids` (Set of String) List of owner identifiers for the agent card manifest.
- `provider` (Attributes) Provider information for the agent card. Changing this value requires resource recreation. (see [below for nested schema](#nestedatt--agent_card_manifest--provider))
- `skills` (Attributes List) Skills defined in the agent card manifest. (see [below for nested schema](#nestedatt--agent_card_manifest--skills))

Read-Only:

- `id` (String) Unique identifier for the agent card manifest.

<a id="nestedatt--agent_card_manifest--capabilities"></a>
### Nested Schema for `agent_card_manifest.capabilities`

Required:

- `push_notifications` (Boolean) Whether the agent supports push notifications.
- `state_transition_history` (Boolean) Whether the agent supports state transition history.
- `streaming` (Boolean) Whether the agent supports streaming.

Optional:

- `extensions` (Attributes List) Capability extensions for the agent. (see [below for nested schema](#nestedatt--agent_card_manifest--capabilities--extensions))

<a id="nestedatt--agent_card_manifest--capabilities--extensions"></a>
### Nested Schema for `agent_card_manifest.capabilities.extensions`

Required:

- `uri` (String) URI of the extension.

Optional:

- `description` (String) Description of the extension.
- `params` (Map of String) Parameters for the extension.
- `required` (Boolean) Whether the extension is required.



<a id="nestedatt--agent_card_manifest--provider"></a>
### Nested Schema for `agent_card_manifest.provider`

Optional:

- `organization` (String) Organization name of the provider.
- `url` (String) URL of the provider.


<a id="nestedatt--agent_card_manifest--skills"></a>
### Nested Schema for `agent_card_manifest.skills`

Required:

- `display_name` (String) Display name for the skill.
- `id` (String) Unique identifier for the skill.

Optional:

- `description` (String) Description of the skill.
- `examples` (Set of String) Example prompts for the skill.
- `input_modes` (Set of String) Input modes supported by the skill.
- `output_modes` (Set of String) Output modes supported by the skill.
- `tags` (Set of String) Tags associated with the skill.



<a id="nestedatt--additional_interfaces"></a>
### Nested Schema for `additional_interfaces`

Required:

- `transport` (String) Transport protocol. The possible values are `JSONRPC`, `GRPC`, and `HTTP+JSON`.
- `url` (String) URL for the interface.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import an existing Agent Instance using the Object ID (id)
# The ID can be found in the Microsoft Entra admin center or via the Graph API:
# GET https://graph.microsoft.com/beta/agentRegistry/agentInstances

terraform import microsoft365_graph_beta_agents_agent_instance.example 00000000-0000-0000-0000-000000000000
```
