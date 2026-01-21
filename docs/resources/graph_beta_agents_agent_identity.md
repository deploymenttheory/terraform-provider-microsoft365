---
page_title: "microsoft365_graph_beta_agents_agent_identity Resource - terraform-provider-microsoft365"
subcategory: "Agents"

description: |-
  Manages an Agent Identity in Microsoft Entra ID using the /servicePrincipals/microsoft.graph.agentIdentity endpoint. This resource is used to represent a service principal for an AI agent instance, created from an agent identity blueprint. Agent identities inherit settings from their blueprint and can be assigned permissions and credentials.
  For more information, see the Agent Identity documentation https://learn.microsoft.com/en-us/graph/api/resources/agentidentity?view=graph-rest-beta.
---

# microsoft365_graph_beta_agents_agent_identity (Resource)

Manages an Agent Identity in Microsoft Entra ID using the `/servicePrincipals/microsoft.graph.agentIdentity` endpoint. This resource is used to represent a service principal for an AI agent instance, created from an agent identity blueprint. Agent identities inherit settings from their blueprint and can be assigned permissions and credentials.

For more information, see the [Agent Identity documentation](https://learn.microsoft.com/en-us/graph/api/resources/agentidentity?view=graph-rest-beta).

## Microsoft Documentation

- [agentIdentity resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentidentity?view=graph-rest-beta)
- [Create agentIdentity](https://learn.microsoft.com/en-us/graph/api/agentidentity-post?view=graph-rest-beta&tabs=http)
- [Get agentIdentity](https://learn.microsoft.com/en-us/graph/api/agentidentity-get?view=graph-rest-beta&tabs=http)
- [Update agentIdentity](https://learn.microsoft.com/en-us/graph/api/agentidentity-update?view=graph-rest-beta&tabs=http)
- [List owners](https://learn.microsoft.com/en-us/graph/api/agentidentity-list-owners?view=graph-rest-beta&tabs=http)
- [List sponsors](https://learn.microsoft.com/en-us/graph/api/agentidentity-list-sponsors?view=graph-rest-beta&tabs=http)
- [Permanently delete item](https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `AgentInstance.Read.All`
- `Directory.Read.All`
- `AgentInstance.ReadWrite.All`
- `Directory.ReadWrite.All`
- `AgentIdentity.DeleteRestore.All`

**Optional:**
- `None` `[N/A]`

Find out more about the permissions required for managing agent identities at Microsoft Learn [here](https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.38.0 | Experimental | Initial release |

## Important Notes

- Agent identities are deleted in two steps: first soft-deleted, then permanently removed from deleted items via the `/directory/deleteditems/{id}` endpoint.
- An agent identity blueprint and its service principal must exist before creating an agent identity.
- At least one sponsor and one owner must be specified when creating an agent identity.
- The `agent_identity_blueprint_id` is the `app_id` (Application/Client ID) of the blueprint, not its Object ID.

## Example Usage

### Basic Example with Tags

```terraform
# Example: Basic Agent Identity with Tags
#
# This example shows the minimal configuration for an agent identity
# with optional tags for categorization.
#
# Prerequisites:
# - An existing Agent Identity Blueprint with app_id
# - The Agent Identity Blueprint must have a service principal created
# - At least one user to assign as sponsor and owner

resource "microsoft365_graph_beta_agents_agent_identity" "basic" {
  display_name                = "My Agent Identity"
  agent_identity_blueprint_id = "00000000-0000-0000-0000-000000000000" # Replace with blueprint app_id
  account_enabled             = true
  sponsor_ids                 = ["00000000-0000-0000-0000-000000000001"] # Replace with user IDs
  owner_ids                   = ["00000000-0000-0000-0000-000000000001"] # Replace with user IDs
  tags                        = ["production", "customer-service", "ai-agent"]

  # When true, permanently deletes from Entra ID on destroy (cannot be restored)
  # When false, moves to deleted items (can be restored within 30 days)
  hard_delete = true
}
```

### Full Example with Dependency Chain

```terraform
# Example: Agent Identity with Full Dependency Chain
#
# This example shows the complete setup including:
# - A user to act as sponsor and owner
# - An Agent Identity Blueprint
# - The Blueprint's Service Principal
# - The Agent Identity itself

# Create a user to be the sponsor
resource "microsoft365_graph_beta_users_user" "agent_sponsor" {
  display_name        = "Agent Sponsor User"
  user_principal_name = "agent-sponsor@yourdomain.com"
  mail_nickname       = "agent-sponsor"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Create a user to be the owner
resource "microsoft365_graph_beta_users_user" "agent_owner" {
  display_name        = "Agent Owner User"
  user_principal_name = "agent-owner@yourdomain.com"
  mail_nickname       = "agent-owner"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Create an agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "Customer Service Agent Blueprint"
  description      = "Blueprint for customer service AI agents"
  sponsor_user_ids = [microsoft365_graph_beta_users_user.agent_sponsor.id]
  owner_user_ids   = [microsoft365_graph_beta_users_user.agent_owner.id]
  tags             = ["customer-service", "production"]
  hard_delete      = true
}

# Create the service principal for the blueprint (required before creating agent identities)
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal" "example" {
  app_id      = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
  hard_delete = true
}

# Create an agent identity from the blueprint
resource "microsoft365_graph_beta_agents_agent_identity" "example" {
  display_name                = "Customer Service Agent 01"
  agent_identity_blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
  account_enabled             = true
  sponsor_ids                 = [microsoft365_graph_beta_users_user.agent_sponsor.id]
  owner_ids                   = [microsoft365_graph_beta_users_user.agent_owner.id]
  tags                        = ["customer-service", "agent-instance"]
  hard_delete                 = true

  depends_on = [
    microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.example
  ]
}

# Outputs
output "agent_identity_id" {
  description = "The ID of the created agent identity"
  value       = microsoft365_graph_beta_agents_agent_identity.example.id
}

output "agent_identity_display_name" {
  description = "The display name of the agent identity"
  value       = microsoft365_graph_beta_agents_agent_identity.example.display_name
}

output "agent_identity_service_principal_type" {
  description = "The service principal type of the agent identity"
  value       = microsoft365_graph_beta_agents_agent_identity.example.service_principal_type
}

output "blueprint_app_id" {
  description = "The app ID of the agent identity blueprint"
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_enabled` (Boolean) Set whether the agent identity is enabled. If `false`, the agent identity cannot authenticate or access resources.
- `agent_identity_blueprint_id` (String) The application (client) ID of the agent identity blueprint from which this agent identity is created. Required. This is the `app_id` of the `microsoft365_graph_beta_agents_agent_identity_blueprint` resource.
- `display_name` (String) The display name for the agent identity. Maximum length is 256 characters. Required.
- `owner_ids` (Set of String) The user IDs of the owners for the agent identity. At least one owner is required when creating an agent identity. Owners are users who have full control over the agent identity.
- `sponsor_ids` (Set of String) The user IDs of the sponsors for the agent identity. At least one sponsor is required when creating an agent identity. Sponsors are users who can approve or oversee the agent identity.

### Optional

- `hard_delete` (Boolean) When set to `true`, the resource will be permanently deleted from the Entra ID (hard delete) rather than being moved to deleted items (soft delete). This prevents the resource from being restored and immediately frees up the resource name for reuse. When `false` (default), the resource is soft deleted and can be restored within 30 days. Note: This field defaults to `false` on import since the API does not return this value.
- `tags` (Set of String) Custom strings that can be used to categorize and identify the agent identity.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_by_app_id` (String) The application ID of the application that created this agent identity. Read-only.
- `created_date_time` (String) The date and time when the agent identity was created. Read-only.
- `disabled_by_microsoft_status` (String) Indicates whether Microsoft has disabled the agent identity. Possible values are: `null`, `NotDisabled`, `DisabledDueToViolationOfServicesAgreement`. Read-only.
- `id` (String) The unique identifier for the agent identity service principal. Read-only.
- `service_principal_type` (String) The type of the service principal. For agent identities, this is always `ServiceIdentity`. Read-only.

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
# Import an existing agent identity into Terraform
# The import ID format is: {agent_identity_id}/{agent_identity_blueprint_id}[:hard_delete=true|false]
#
# Where:
# - {agent_identity_id} is the Object ID of the agent identity service principal
# - {agent_identity_blueprint_id} is the Application (client) ID of the blueprint
# - hard_delete is optional (defaults to false for soft delete only)

# Basic import (hard_delete defaults to false - soft delete only)
terraform import microsoft365_graph_beta_agents_agent_identity.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"

# Import with hard_delete enabled (permanently deletes on terraform destroy)
terraform import microsoft365_graph_beta_agents_agent_identity.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111:hard_delete=true"
```

