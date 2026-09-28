---
page_title: "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages the assignment of a Microsoft Entra B2B management policy to an application or service principal. The assignment is read from the policy's appliesTo relationship. To import this resource, use the format: b2b_management_policy_id/directory_object_id.
---

# microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment (Resource)

Manages the assignment of a Microsoft Entra B2B management policy to an application or service principal. The assignment is read from the policy's `appliesTo` relationship. To import this resource, use the format: `b2b_management_policy_id/directory_object_id`.

## Microsoft Graph API

The documented [Add appliesTo](https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-post-appliesto?view=graph-rest-beta) and [Remove appliesTo](https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-delete-appliesto?view=graph-rest-beta) endpoints (`/policies/b2bManagementPolicies/{id}/appliesTo/$ref`) currently reject every request with `400 Property 'appliesTo' is read-only and cannot be set.`. This resource therefore writes the assignment from the application or service principal side:

- Create: `POST /beta/{applications|servicePrincipals}/{directoryObjectId}/policies/$ref` with `@odata.id` set to the policy URL.
- Read: `GET /beta/policies/b2bManagementPolicies/{id}/appliesTo` ([List appliesTo](https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-list-appliesto?view=graph-rest-beta)).
- Delete: `DELETE /beta/{applications|servicePrincipals}/{directoryObjectId}/policies/{policyId}/$ref`.

## Behavior Notes

- B2B management policies can only be applied to applications and service principals. The provider reads `GET /beta/directoryObjects/{id}` to detect which one `directory_object_id` refers to, and exposes the result as `directory_object_type`.
- Before the `$ref` POST, the provider waits until both the policy and the directory object are readable, because objects created moments earlier may not have propagated across Microsoft Entra replicas. The replica serving the POST can still lag and reject it with `404 Directory_ObjectNotFound`. Nothing is written in that case, so the POST is retried.
- If the POST is rejected with `400 One or more added object references already exist`, the assignment already existed outside this resource. Create fails and asks for it to be imported, rather than silently adopting an assignment that a later destroy would remove.
- During refresh, an assignment missing from `appliesTo` (or a 404 for the policy) is only treated as removed once it has persisted for 60 seconds, because stale replicas briefly return an empty `appliesTo`.
- When the caller cannot read an application in the policy's `appliesTo`, Microsoft Graph returns HTTP 200 with a truncated body followed by an `InternalServerError` (`signInAudienceRestrictions ... has a null value`). The provider reports this as an error that names the missing `Application.Read.All` permission, rather than dropping the assignment from state.
- Changing `b2b_management_policy_id` or `directory_object_id` replaces the assignment.

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Policy.Read.B2BManagementPolicy`
- `Policy.ReadWrite.B2BManagementPolicy`
- `Application.Read.All`
- `Application.ReadWrite.OwnedBy`: sufficient when the calling service principal owns the target application or service principal. Otherwise use `Application.ReadWrite.All`, which also covers `Application.Read.All`.

**Optional:**
- `None` `[N/A]`

These permissions were verified against a live tenant with application (app-only) tokens:

| Permissions | `GET /directoryObjects/{id}` | `POST`/`DELETE .../policies/$ref` | `GET .../appliesTo` with an application assigned |
| --- | --- | --- | --- |
| `Policy.ReadWrite.B2BManagementPolicy` | 403 | 403 | truncated response |
| `Policy.ReadWrite.B2BManagementPolicy` + `Policy.ReadWrite.ApplicationConfiguration` | 403 | 403 | truncated response |
| `Policy.Read.B2BManagementPolicy` + `Application.Read.All` | 200 | not tested (read-only) | 200 |
| `Policy.ReadWrite.B2BManagementPolicy` + `Application.ReadWrite.OwnedBy` (caller owns the targets) | 200 | 204 | 200 |
| `Policy.ReadWrite.B2BManagementPolicy` + `Application.ReadWrite.All` | 200 | 204 | 200 |

## Example Usage

```terraform
# Example: Apply a B2B management policy to a service principal and an application
# B2B management policies can only be applied to applications and service principals.
# The provider detects which of the two the directory object is.

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "example" {
  display_name = "example-b2b-management-policy"

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.com"]
        }
      }
    })
  ]
}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "service_principal" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.example.id
  directory_object_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # Object ID of the service principal
}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "application" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.example.id
  directory_object_id      = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy" # Object ID of the application
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `b2b_management_policy_id` (String) The unique identifier of the B2B management policy to assign.
- `directory_object_id` (String) The object ID of the application or service principal to apply the B2B management policy to. B2B management policies can only be applied to applications and service principals.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `directory_object_type` (String) The type of the directory object the policy is applied to, as detected from Microsoft Graph. One of `application` or `servicePrincipal`.
- `id` (String) A locally generated composite identifier for this assignment in the format `b2b_management_policy_id/directory_object_id`. The Microsoft Graph API does not return an assignment-specific ID for this resource; this value is constructed by the provider to uniquely identify the assignment within Terraform state. Use this format when importing the resource.

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
# Import using the composite ID format: {b2b_management_policy_id}/{directory_object_id}
# The Microsoft Graph API does not return an assignment-specific ID.

# {b2b_management_policy_id} - GUID of the B2B management policy
# {directory_object_id} - Object ID of the application or service principal
terraform import microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment.example 00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002
```
