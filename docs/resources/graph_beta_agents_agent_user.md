---
page_title: "microsoft365_graph_beta_agents_agent_user Resource - terraform-provider-microsoft365"
subcategory: "Agents"

description: |-
  Manages Microsoft 365 users using the /users/microsoft.graph.agentUser endpoint. This resource is used to represents a specialized subtype of user identity in Microsoft Entra ID designed for AI-powered applications (agents) that need to function as digital workers. Agent users enable agents to access APIs and services that specifically require user identities, receiving tokens with idtyp=user claims.
---

# microsoft365_graph_beta_agents_agent_user (Resource)

Manages Microsoft 365 users using the `/users/microsoft.graph.agentUser` endpoint. This resource is used to represents a specialized subtype of user identity in Microsoft Entra ID designed for AI-powered applications (agents) that need to function as digital workers. Agent users enable agents to access APIs and services that specifically require user identities, receiving tokens with `idtyp=user` claims.

## Microsoft Documentation

- [agentUser resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentuser?view=graph-rest-beta)
- [Create agentUser](https://learn.microsoft.com/en-us/graph/api/agentuser-post?view=graph-rest-beta&tabs=http)
- [Get agentUser](https://learn.microsoft.com/en-us/graph/api/agentuser-get?view=graph-rest-beta&tabs=http)
- [Update agentUser](https://learn.microsoft.com/en-us/graph/api/agentuser-update?view=graph-rest-beta&tabs=http)
- [Delete agentUser](https://learn.microsoft.com/en-us/graph/api/agentuser-delete?view=graph-rest-beta&tabs=http)
- [List sponsors](https://learn.microsoft.com/en-us/graph/api/agentuser-list-sponsors?view=graph-rest-beta&tabs=http)
- [Permanently delete item](https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `User.Read.All`
- `AgentIdUser.ReadWrite.All`
- `Directory.Read.All`
- `AgentIdUser.ReadWrite.IdentityParentedBy`
- `User.ReadWrite.All`
- `User.DeleteRestore.All`

**Optional:**
- `None` `[N/A]`

Find out more about the permissions required for managing agent identities at Microsoft Learn [here](https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.38.0 | Experimental | Initial release. Added hard delete handling option. |

## Important Notes

- Agent users maintain a one-to-one relationship with a parent agent identity and authenticate through that parent's credentials.
- An agent identity must exist before creating an agent user.
- At least one sponsor must be specified when creating an agent user.
- The `agent_identity_id` is the Object ID of the parent agent identity.
- When `hard_delete` is true, the user is permanently deleted; when false, it is soft-deleted and can be restored within 30 days.

## Example Usage

### Minimal Example

```terraform
# Minimal Agent User Example
# This example shows only the required fields for creating an agent user

resource "microsoft365_graph_beta_agents_agent_user" "example" {
  display_name        = "Example Agent User"
  agent_identity_id   = "00000000-0000-0000-0000-000000000000" # ID of parent agent identity
  account_enabled     = true
  user_principal_name = "agent-user@contoso.com" # Must match your tenant's verified domain
  mail_nickname       = "agent-user"
  sponsor_ids         = ["11111111-1111-1111-1111-111111111111"] # User ID of sponsor
  hard_delete         = true
}
```

### Maximal Example

```terraform
# Maximal Agent User Example
# This example shows all available fields for creating an agent user

resource "microsoft365_graph_beta_agents_agent_user" "example" {
  # Required fields
  display_name        = "Example Agent User"
  agent_identity_id   = "00000000-0000-0000-0000-000000000000" # ID of parent agent identity
  account_enabled     = true
  user_principal_name = "agent-user@contoso.com" # Must match your tenant's verified domain
  mail_nickname       = "agent-user"
  sponsor_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
  hard_delete = true

  # Optional name fields
  given_name = "Agent"
  surname    = "User"

  # Optional organizational fields
  job_title       = "AI Agent"
  department      = "Engineering"
  company_name    = "Contoso"
  office_location = "Building A"

  # Optional address fields
  city           = "Seattle"
  state          = "WA"
  country        = "US"
  postal_code    = "98101"
  street_address = "123 Main Street"

  # Optional locale fields
  usage_location     = "US"
  preferred_language = "en-US"

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
```

### Full Example with Dependency Chain

```terraform
# Agent User with Full Dependency Chain
# This example demonstrates the complete resource hierarchy:
# User (sponsor/owner) -> Agent Identity Blueprint -> Service Principal -> Agent Identity -> Agent User

########################################################################################
# Look up existing user to use as sponsor/owner
########################################################################################

data "microsoft365_graph_beta_users_user_by_filter" "sponsor" {
  filter_type  = "display_name"
  filter_value = "Admin User"
}

########################################################################################
# Agent Identity Blueprint
# The blueprint defines the template for agent identities
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "Example Agent Blueprint"
  description      = "Blueprint for example agent identity"
  sponsor_user_ids = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  owner_user_ids   = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  hard_delete      = true
}

########################################################################################
# Agent Identity Blueprint Service Principal
# Required before creating agent identities from the blueprint
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal" "example" {
  app_id      = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
  hard_delete = true
}

########################################################################################
# Agent Identity
# The parent identity that the agent user will be associated with
########################################################################################

resource "microsoft365_graph_beta_agents_agent_identity" "example" {
  display_name                = "Example Agent Identity"
  agent_identity_blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.app_id
  account_enabled             = true
  sponsor_ids                 = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  owner_ids                   = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  hard_delete                 = true

  depends_on = [
    microsoft365_graph_beta_agents_agent_identity_blueprint_service_principal.example
  ]
}

########################################################################################
# Agent User
# The user identity that authenticates through the parent agent identity
########################################################################################

resource "microsoft365_graph_beta_agents_agent_user" "example" {
  display_name        = "Example Agent User"
  agent_identity_id   = microsoft365_graph_beta_agents_agent_identity.example.id
  account_enabled     = true
  user_principal_name = "agent-user@${var.domain}"
  mail_nickname       = "agent-user"
  sponsor_ids         = [data.microsoft365_graph_beta_users_user_by_filter.sponsor.id]
  hard_delete         = true

  # Optional fields
  given_name     = "Agent"
  surname        = "User"
  job_title      = "AI Agent"
  department     = "Engineering"
  usage_location = "US"

  depends_on = [
    microsoft365_graph_beta_agents_agent_identity.example
  ]
}

########################################################################################
# Variables
########################################################################################

variable "domain" {
  description = "The verified domain for the tenant (e.g., contoso.com)"
  type        = string
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_enabled` (Boolean) Set to `true` if the agent user account is enabled; otherwise, `false`. This property is required when a user is created.
- `agent_identity_id` (String) The object ID of the agent identity that this agent user is associated with. This creates a one-to-one relationship with the agent identity. The agent user authenticates through the parent agent identity's credentials. Required.
- `display_name` (String) The name displayed in the address book for the agent user. This value is usually the combination of the user's first name, middle initial, and last name. This property is required when an agent user is created and it cannot be cleared during updates. Maximum length is 256 characters.
- `mail_nickname` (String) The mail alias for the agent user. This property must be specified when a user is created. Maximum length is 64 characters.
- `sponsor_ids` (Set of String) The users and groups responsible for this agent user's privileges in the tenant and keep the agent user's information and access updated. Required.
- `user_principal_name` (String) The user principal name (UPN) of the agent user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. By convention, this should map to the agent user's email name. The general format is alias@domain, where the domain must be present in the tenant's verified domain collection. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization. NOTE: This property can't contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~.

### Optional

- `city` (String) The city in which the agent user is located. Maximum length is 128 characters.
- `company_name` (String) The company name which the agent user is associated. This property can be useful for describing the company that an external user comes from. Maximum length is 64 characters.
- `country` (String) The country/region in which the agent user is located; for example, US or UK.
- `department` (String) The name of the department in which the agent user works. Maximum length is 64 characters.
- `given_name` (String) The given name (first name) of the agent user. Maximum length is 64 characters.
- `hard_delete` (Boolean) When set to `true`, the resource will be permanently deleted from the Entra ID (hard delete) rather than being moved to deleted items (soft delete). This prevents the resource from being restored and immediately frees up the resource name for reuse. When `false` (default), the resource is soft deleted and can be restored within 30 days. Note: This field defaults to `false` on import since the API does not return this value.
- `job_title` (String) The agent user's job title. Maximum length is 128 characters.
- `mail` (String) The SMTP address for the agent user, for example, jeff@contoso.com. Read-only.
- `office_location` (String) The office location in the agent user's place of business. Maximum length is 128 characters.
- `postal_code` (String) The postal code for the agent user's postal address. The postal code is specific to the user's country/region. Maximum length is 40 characters.
- `preferred_language` (String) The preferred language for the agent user. The preferred language format is based on ISO 639-1 Code; for example en-US.
- `state` (String) The state or province in the agent user's address. Maximum length is 128 characters.
- `street_address` (String) The street address of the agent user's place of business. Maximum length is 1024 characters.
- `surname` (String) The user's surname (family name or last name). Maximum length is 64 characters.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `usage_location` (String) A two-letter country code (ISO standard 3166). Required for users that are assigned licenses due to legal requirements to check for availability of services in countries. Examples include: `US`, `JP`, and `GB`. Not nullable.
- `user_type` (String) A string value that can be used to classify user types in your directory, such as `Member` and `Guest`.

### Read-Only

- `created_date_time` (String) The date and time the agent user was created. The value cannot be modified and is automatically populated when the entity is created. Read-only.
- `creation_type` (String) Indicates whether the agent user account was created through one of the following methods: As a regular school or work account (null), As an external account (Invitation), As a local account for an Azure Active Directory B2C tenant (LocalAccount), Through self-service sign-up by an internal user using email verification (EmailVerified), Through self-service sign-up by an external user signing up through a link that is part of a user flow (SelfServiceSignUp). Read-only.
- `id` (String) The unique identifier for the agent user. Read-only.

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
# Import an existing agent user into Terraform
# The import ID format is: {agent_user_id}[:hard_delete=true|false]
#
# Where:
# - {agent_user_id} is the unique identifier for the agent user
# - hard_delete is optional (defaults to false for soft delete only)

# Basic import (hard_delete defaults to false - soft delete only)
terraform import microsoft365_graph_beta_agents_agent_user.example "00000000-0000-0000-0000-000000000000"

# Import with hard_delete enabled (permanently deletes on terraform destroy)
terraform import microsoft365_graph_beta_agents_agent_user.example "00000000-0000-0000-0000-000000000000:hard_delete=true"
```

