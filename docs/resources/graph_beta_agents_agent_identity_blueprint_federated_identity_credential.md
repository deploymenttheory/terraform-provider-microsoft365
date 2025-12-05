---
page_title: "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential Resource - terraform-provider-microsoft365"
subcategory: "Agents"

description: |-
  Manages a Federated Identity Credential for an Agent Identity Blueprint in Microsoft Entra ID using the /applications endpoint. By configuring a trust relationship between your Microsoft Entra agent identity blueprint registration and the identity provider for your compute platform, you can use tokens issued by that platform to authenticate with Microsoft identity platform and call APIs in the Microsoft ecosystem. Maximum of 20 federated identity credentials can be added to an agentIdentityBlueprint.
---

# microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential (Resource)

Manages a Federated Identity Credential for an Agent Identity Blueprint in Microsoft Entra ID using the `/applications` endpoint. By configuring a trust relationship between your Microsoft Entra agent identity blueprint registration and the identity provider for your compute platform, you can use tokens issued by that platform to authenticate with Microsoft identity platform and call APIs in the Microsoft ecosystem. Maximum of 20 federated identity credentials can be added to an agentIdentityBlueprint.

## Microsoft Documentation

- [federatedIdentityCredential resource type](https://learn.microsoft.com/en-us/graph/api/resources/federatedidentitycredential?view=graph-rest-beta)
- [Create federatedIdentityCredential](https://learn.microsoft.com/en-us/graph/api/application-post-federatedidentitycredentials?view=graph-rest-beta&tabs=http)
- [Update federatedIdentityCredential](https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-update?view=graph-rest-beta&tabs=http)
- [Delete federatedIdentityCredential](https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-delete?view=graph-rest-beta&tabs=http)
- [Workload identity federation](https://learn.microsoft.com/en-us/entra/workload-id/workload-identity-federation)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Read Permissions**: `AgentIdentityBlueprint.AddRemoveCreds.All`, `Directory.Read.All`
- **Write Permissions**: `AgentIdentityBlueprint.AddRemoveCreds.All`, `Directory.ReadWrite.All`

Find out more about the permissions required for managing agent identities at microsoft learn [here](https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.38.0 | Experimental | Initial release |

## Important Notes

- This resource creates a federated identity credential for an existing agent identity blueprint
- By configuring a trust relationship between your blueprint and an external identity provider, you can use tokens from that platform to authenticate with Microsoft identity platform
- Maximum of 20 federated identity credentials can be added to an agent identity blueprint
- The combination of `issuer` and `subject` must be unique within the blueprint
- If using `claims_matching_expression`, the `subject` field must be null and vice versa

## Example Usage

### GitHub Actions Example

```terraform
# Example: Create a federated identity credential for GitHub Actions

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "github-actions-agent"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for GitHub Actions workflows"
}

# Create a federated identity credential for GitHub Actions
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "github_actions" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  name         = "github-actions-production"
  issuer       = "https://token.actions.githubusercontent.com"
  subject      = "repo:my-org/my-repo:environment:Production"
  audiences    = ["api://AzureADTokenExchange"]
  description  = "Federated identity credential for GitHub Actions in production environment"
}

# Output the credential details
output "credential_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential.github_actions.id
  description = "The ID of the federated identity credential"
}
```

### Azure Kubernetes Service (AKS) Example

```terraform
# Example: Create a federated identity credential for Azure Kubernetes Service (AKS)

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "aks-workload-agent"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for AKS workloads"
}

# Create a federated identity credential for AKS workload identity
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "aks_workload" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  name         = "aks-workload-identity"
  issuer       = "https://oidc.prod-aks.azure.com/00000000-0000-0000-0000-000000000000/"
  subject      = "system:serviceaccount:default:workload-identity-sa"
  audiences    = ["api://AzureADTokenExchange"]
  description  = "Federated identity credential for AKS workload identity"
}

# Output the credential details
output "credential_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential.aks_workload.id
  description = "The ID of the federated identity credential"
}
```

### AWS Example

```terraform
# Example: Create a federated identity credential for AWS

# First, create or reference an existing agent identity blueprint
resource "microsoft365_graph_beta_agents_agent_identity_blueprint" "example" {
  display_name     = "aws-workload-agent"
  sponsor_user_ids = ["00000000-0000-0000-0000-000000000000"]
  owner_user_ids   = ["00000000-0000-0000-0000-000000000000"]
  description      = "Agent identity blueprint for AWS workloads"
}

# Create a federated identity credential for AWS IAM
resource "microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential" "aws_iam" {
  blueprint_id = microsoft365_graph_beta_agents_agent_identity_blueprint.example.id
  name         = "aws-iam-role"
  issuer       = "https://token.sts.amazonaws.com"
  subject      = "arn:aws:iam::123456789012:role/my-role"
  audiences    = ["api://AzureADTokenExchange"]
  description  = "Federated identity credential for AWS IAM role"
}

# Output the credential details
output "credential_id" {
  value       = microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential.aws_iam.id
  description = "The ID of the federated identity credential"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `audiences` (Set of String) The audience that can appear in the external token. This field is mandatory and should be set to `api://AzureADTokenExchange` for Microsoft Entra ID. It says what Microsoft identity platform should accept in the aud claim in the incoming token. This value represents Microsoft Entra ID in your external identity provider and has no fixed value across identity providers - you may need to create a new application registration in your identity provider to serve as the audience of this token. This field can only accept a single value and has a limit of 600 characters.
- `blueprint_id` (String) The Object ID (id) of the Agent Identity Blueprint to which this federated identity credential belongs. This is required and cannot be changed after creation.
- `issuer` (String) The URL of the external identity provider and must match the issuer claim of the external token being exchanged. The combination of the values of issuer and subject must be unique on the app. It has a limit of 600 characters.
- `name` (String) The unique identifier for the federated identity credential, which has a limit of 120 characters and must be URL friendly. It is immutable once created.

### Optional

- `claims_matching_expression` (String) Nullable. Defaults to null if not set. Enables the use of claims matching expressions against specified claims. If claims_matching_expression is defined, subject must be null. For the list of supported expression syntax and claims, visit the [Flexible FIC reference](https://aka.ms/flexiblefic).
- `description` (String) An optional description of the federated identity credential.
- `subject` (String) Nullable. Defaults to null if not set. The identifier of the external software workload within the external identity provider. Like the audience value, it has no fixed format, as each identity provider uses their own - sometimes a GUID, sometimes a colon delimited identifier, sometimes arbitrary strings. The value here must match the sub claim within the token presented to Microsoft Entra ID. It has a limit of 600 characters. The combination of issuer and subject must be unique on the app. If subject is defined, claims_matching_expression must be null.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for the federated identity credential. Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax (format: `blueprint_id/credential_id`):

```shell
# Import an existing federated identity credential
# The import ID format is: blueprint_id/credential_id
terraform import microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"
```
