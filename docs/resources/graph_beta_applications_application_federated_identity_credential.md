---
page_title: "microsoft365_graph_beta_applications_application_federated_identity_credential Resource - terraform-provider-microsoft365"
subcategory: "Applications"
description: |-
  Manages a federated identity credential for a Microsoft Entra Application using the /applications/{id}/federatedIdentityCredentials endpoint. Federated identity credentials configure a trust relationship between your application and an external identity provider, enabling token-based authentication with the Microsoft identity platform. Maximum of 20 federated identity credentials can be added to an application.
  For more information, see the Microsoft Graph API documentation https://learn.microsoft.com/en-us/graph/api/application-post-federatedidentitycredentials?view=graph-rest-beta.
---

# microsoft365_graph_beta_applications_application_federated_identity_credential (Resource)

Manages a federated identity credential for a Microsoft Entra Application using the `/applications/{id}/federatedIdentityCredentials` endpoint. Federated identity credentials configure a trust relationship between your application and an external identity provider, enabling token-based authentication with the Microsoft identity platform. Maximum of 20 federated identity credentials can be added to an application.

For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/application-post-federatedidentitycredentials?view=graph-rest-beta).

## Microsoft Documentation

- [federatedIdentityCredential resource type](https://learn.microsoft.com/en-us/graph/api/resources/federatedidentitycredential?view=graph-rest-beta)
- [Create federatedIdentityCredential](https://learn.microsoft.com/en-us/graph/api/application-post-federatedidentitycredentials?view=graph-rest-beta&tabs=http)
- [Get federatedIdentityCredential](https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-get?view=graph-rest-beta&tabs=http)
- [Update federatedIdentityCredential](https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-update?view=graph-rest-beta&tabs=http)
- [Delete federatedIdentityCredential](https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Application.Read.All`
- `Directory.Read.All`
- `Application.ReadWrite.All`
- `Directory.ReadWrite.All`

**Optional:**
- `Application.ReadWrite.OwnedBy` (if managing only applications owned by the service principal)

Find out more about the permissions required for managing applications at Microsoft Learn [here](https://learn.microsoft.com/en-us/graph/api/resources/federatedidentitycredential?view=graph-rest-beta).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.43.0 | Experimental | Initial release |

## Important Notes

- **Workload Identity Federation**: Enables passwordless authentication from external identity providers
- **Common Use Cases**:
  - GitHub Actions deploying to Azure
  - Azure Kubernetes Service (AKS) workload identity
  - Google Cloud Platform (GCP) workload identity
  - AWS workload identity
- **Issuer**: Must be a trusted OIDC issuer URL (e.g., GitHub Actions, Kubernetes clusters)
- **Subject**: Identifies the specific workload (repo, service account, etc.)
- **Audiences**: Typically `["api://AzureADTokenExchange"]` for Azure scenarios

## Example Usage

### GitHub Actions Deployment

```terraform
resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-github-deployment-app"
  description  = "Application for GitHub Actions deployments"
}

# Federated credential for GitHub Actions to deploy to Azure
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "github_actions" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  name           = "github-actions-production"
  description    = "GitHub Actions deploying to Production environment"
  issuer         = "https://token.actions.githubusercontent.com"
  subject        = "repo:myorg/myrepo:environment:Production"
  audiences      = ["api://AzureADTokenExchange"]
}
```

### Kubernetes Workload Identity

```terraform
resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-kubernetes-workload-app"
  description  = "Application for Kubernetes workload identity"
}

# Federated credential for Kubernetes workload identity
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "kubernetes" {
  application_id = microsoft365_graph_beta_applications_application.example.id
  name           = "aks-workload-identity"
  description    = "Azure Kubernetes Service workload identity"
  issuer         = "https://eastus.oic.prod-aks.azure.com/00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111/"
  subject        = "system:serviceaccount:default:my-service-account"
  audiences      = ["api://AzureADTokenExchange"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_id` (String) The Object ID (id) of the Application to which this federated identity credential belongs. This is required and cannot be changed after creation.
- `audiences` (Set of String) The audience that can appear in the external token. This field is mandatory and should be set to `api://AzureADTokenExchange` for Microsoft Entra ID. It says what Microsoft identity platform should accept in the aud claim in the incoming token. This value represents Microsoft Entra ID in your external identity provider and has no fixed value across identity providers - you may need to create a new application registration in your identity provider to serve as the audience of this token. This field can only accept a single value and has a limit of 600 characters.
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

```shell
# Import a federated identity credential by composite ID: application_id/credential_id
terraform import microsoft365_graph_beta_applications_application_federated_identity_credential.example "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"
```
