---
page_title: "microsoft365_graph_beta_agents_agent_collection Resource - terraform-provider-microsoft365"
subcategory: "Agents"

description: |-
  Manages an Agent Collection in the Microsoft Entra Agent Registry using the /agentRegistry/agentCollections endpoint. An agent collection represents a grouping of agent instances for organizational and access control purposes.
  Reserved Collections: Two system-reserved collections are always available per tenant:
  Global (ID: 00000000-0000-0000-0000-000000000001): Tenant-wide pool of generally available agentsQuarantined (ID: 00000000-0000-0000-0000-000000000002): Holding area for blocked/review-pending agents
  Reserved collections cannot be updated or deleted. Attempting to create a collection with a reserved name returns a 409 Conflict error.
  For more information, see the agentCollection resource type https://learn.microsoft.com/en-us/graph/api/resources/agentcollection?view=graph-rest-beta.
---

# microsoft365_graph_beta_agents_agent_collection (Resource)

Manages an Agent Collection in the Microsoft Entra Agent Registry using the `/agentRegistry/agentCollections` endpoint. An agent collection represents a grouping of agent instances for organizational and access control purposes.

**Reserved Collections**: Two system-reserved collections are always available per tenant:
- **Global** (ID: `00000000-0000-0000-0000-000000000001`): Tenant-wide pool of generally available agents
- **Quarantined** (ID: `00000000-0000-0000-0000-000000000002`): Holding area for blocked/review-pending agents

Reserved collections cannot be updated or deleted. Attempting to create a collection with a reserved name returns a 409 Conflict error.

For more information, see the [agentCollection resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentcollection?view=graph-rest-beta).

## Microsoft Documentation

- [agentCollection resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentcollection?view=graph-rest-beta)
- [Create agentCollection](https://learn.microsoft.com/en-us/graph/api/agentregistry-post-agentcollections?view=graph-rest-beta&tabs=http)
- [Get agentCollection](https://learn.microsoft.com/en-us/graph/api/agentcollection-get?view=graph-rest-beta&tabs=http)
- [Update agentCollection](https://learn.microsoft.com/en-us/graph/api/agentcollection-update?view=graph-rest-beta&tabs=http)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Read**: `AgentCollection.Read.All`
- **Write**: `AgentCollection.ReadWrite.All`, `AgentCollection.ReadWrite.ManagedBy`

Additional lesser-privileged permissions scoped to special collections:
- `AgentCollection.Read.Global` and `AgentCollection.ReadWrite.Global` for the **Global** collection
- `AgentCollection.Read.Quarantined` and `AgentCollection.ReadWrite.Quarantined` for the **Quarantined** collection

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.38.0 | Experimental | Initial release |

## Important Notes

~> **Known Issue: Delete Operation Not Working** The Microsoft Graph API DELETE endpoint for agent collections is currently not functioning correctly. When Terraform attempts to destroy an agent collection, the API returns success but the resource is not actually removed. This is a known issue with the Microsoft Graph API and we are waiting for Microsoft to fix it. In the meantime, you may need to manually clean up agent collections or use `terraform state rm` to remove them from state without attempting deletion.

### Reserved Collections

Two system-reserved collections are always available per tenant and cannot be managed via Terraform:

| Collection  | ID                                   | Purpose                                          |
|-------------|--------------------------------------|--------------------------------------------------|
| Global      | 00000000-0000-0000-0000-000000000001 | Tenant-wide pool of generally available agents   |
| Quarantined | 00000000-0000-0000-0000-000000000002 | Holding area for blocked/review-pending agents   |

- Reserved collections cannot be updated or deleted (returns 403 Forbidden)
- Creating a collection with a reserved displayName returns 409 Conflict

## Example Usage

### Minimal Example

```terraform
# Minimal Agent Collection configuration
# Creates an agent collection with required fields only
resource "microsoft365_graph_beta_agents_agent_collection" "minimal" {
  display_name = "My Agent Collection"
  owner_ids    = ["00000000-0000-0000-0000-000000000000"]
}
```

### Maximal Example

```terraform
# Maximal Agent Collection configuration
# Creates an agent collection with all available fields configured
resource "microsoft365_graph_beta_agents_agent_collection" "maximal" {
  display_name = "IT Automation Agent Collection"
  owner_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]
  description       = "A collection of IT automation agents for managing infrastructure and support workflows"
  managed_by        = "00000000-0000-0000-0000-000000000003"
  originating_store = "Terraform"

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

- `display_name` (String) Friendly name of the collection.
- `owner_ids` (Set of String) List of object IDs for the owners of the agent collection.

### Optional

- `description` (String) Description / purpose of the collection.
- `managed_by` (String) **appId** (referred to as **Application (client) ID** on the Microsoft Entra admin center) of the service principal managing this agent collection.
- `originating_store` (String) Source system/store where the collection originated. For example Copilot Studio. Changing this value will trigger resource recreation.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_by` (String) Object ID of the user or app that created the agent collection. Read-only.
- `created_date_time` (String) Timestamp when agent collection was created. Read-only.
- `id` (String) Unique identifier for the collection. Key. Inherited from entity.
- `last_modified_date_time` (String) Timestamp of last modification. Read-only.

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

# Import an existing Agent Collection using the Object ID (id)
# The ID can be found via the Graph API:
# GET https://graph.microsoft.com/beta/agentRegistry/agentCollections

# Note: Reserved collections (Global and Quarantined) cannot be managed via Terraform
# - Global: 00000000-0000-0000-0000-000000000001
# - Quarantined: 00000000-0000-0000-0000-000000000002

terraform import microsoft365_graph_beta_agents_agent_collection.example 00000000-0000-0000-0000-000000000000
```
