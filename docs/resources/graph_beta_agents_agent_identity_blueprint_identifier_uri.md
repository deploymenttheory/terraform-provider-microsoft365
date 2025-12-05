---
page_title: "microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri Resource - terraform-provider-microsoft365"
subcategory: "Agents"

description: |-
  Manages an identifier URI and optional OAuth2 permission scope for an Agent Identity Blueprint in Microsoft Entra ID using the /applications endpoint. This resource configures the identifier URI and optional scope using a PATCH https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta to the application endpoint.
  The identifier URI is used to uniquely identify the agent identity blueprint and is required for receiving incoming requests from users and other agents.
  Note: This resource manages a single identifier URI. To manage multiple URIs, create multiple resource instances.
---

# microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri (Resource)

Manages an identifier URI and optional OAuth2 permission scope for an Agent Identity Blueprint in Microsoft Entra ID using the `/applications` endpoint. This resource configures the identifier URI and optional scope using a [PATCH](https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta) to the application endpoint.

The identifier URI is used to uniquely identify the agent identity blueprint and is required for receiving incoming requests from users and other agents.

**Note:** This resource manages a single identifier URI. To manage multiple URIs, create multiple resource instances.

## Microsoft Documentation

- [application resource type](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta)
- [Update application](https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta&tabs=http)
- [permissionScope resource type](https://learn.microsoft.com/en-us/graph/api/resources/permissionscope?view=graph-rest-beta)
- [Configure identifier URI and scope](https://learn.microsoft.com/en-us/entra/agent-id/identity-platform/create-blueprint?tabs=microsoft-graph-api#configure-identifier-uri-and-scope)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Read Permissions**: `AgentIdentityBlueprint.Read.All`, `Application.Read.All`, `Directory.Read.All`
- **Write Permissions**: `AgentIdentityBlueprint.ReadWrite.All`, `Directory.ReadWrite.All`

Find out more about the permissions required for managing agent identities at microsoft learn [here](https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.38.0 | Experimental | Initial release |

## Important Notes

- This resource manages a **single** identifier URI. To manage multiple URIs, create multiple resource instances.
- The `identifier_uri` must be unique within the tenant.
- Valid URI formats include:
  - `api://<blueprint-id>` (recommended for agent blueprints)
  - `api://<domain>/<path>`
  - `https://<domain>/<path>`
  - `urn:<namespace>:<identifier>`
- The optional `scope` block defines an OAuth2 permission scope that allows applications to access the agent on behalf of signed-in users.
- Changing the `identifier_uri` requires replacement of the resource.

## Example Usage

### Basic Identifier URI

```terraform
# Example: Configure an identifier URI for an Agent Identity Blueprint

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "my-agent-blueprint"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for automated workflows"
}

# Configure the identifier URI using the blueprint's ID
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri" "example" {
  blueprint_id   = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  identifier_uri = "api://${microsoft365_graph_beta_agents_agent_identity_blueprint.example.id}"
}

output "identifier_uri" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri.example.identifier_uri
  description = "The configured identifier URI"
}
```

### Identifier URI with Custom Scope

```terraform
# Example: Configure an identifier URI with custom OAuth2 permission scope

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "my-agent-blueprint"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for automated workflows"
}

# Configure the identifier URI with a custom permission scope
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri" "example" {
  blueprint_id   = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  identifier_uri = "api://${microsoft365_graph_beta_agents_agent_identity_blueprint.example.id}"

  scope = {
    admin_consent_description  = "Allow the application to access the agent on behalf of the signed-in user."
    admin_consent_display_name = "Access agent"
    is_enabled                 = true
    type                       = "User"
    value                      = "access_agent"
  }
}

output "identifier_uri" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri.example.identifier_uri
  description = "The configured identifier URI"
}

output "scope_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri.example.scope.id
  description = "The ID of the OAuth2 permission scope"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `blueprint_id` (String) The unique identifier (Object ID) of the agent identity blueprint to configure. Required.
- `identifier_uri` (String) The identifier URI for the agent identity blueprint. Valid formats include `api://<guid>`, `api://<domain>/<path>`, `https://<domain>/<path>`, or `urn:<namespace>:<identifier>`. Required.

### Optional

- `scope` (Attributes) Optional OAuth2 permission scope configuration. Defines the scope that allows applications to access the agent on behalf of the signed-in user. (see [below for nested schema](#nestedatt--scope))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--scope"></a>
### Nested Schema for `scope`

Optional:

- `admin_consent_description` (String) A description of the delegated permission, intended to be read by an administrator granting the permission. Default: `Allow the application to access the agent on behalf of the signed-in user.`
- `admin_consent_display_name` (String) The display name for the permission shown in the admin consent experience. Default: `Access agent`
- `id` (String) The unique identifier for the OAuth2 permission scope. If not specified, a UUID will be generated.
- `is_enabled` (Boolean) Whether the permission scope is enabled. Default: `true`
- `type` (String) The type of permission. Valid values are `User` or `Admin`. Default: `User`
- `value` (String) The value of the scope claim that the resource application should expect in the OAuth 2.0 access token. Default: `access_agent`


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax (format: `blueprint_id/identifier_uri`):

```shell
# Import an identifier URI using the format: blueprint_id/identifier_uri
# Note: The identifier_uri should be URL-encoded if it contains special characters
terraform import microsoft365_graph_beta_agents_agent_identity_blueprint_identifier_uri.example "00000000-0000-0000-0000-000000000000/api://00000000-0000-0000-0000-000000000000"
```

