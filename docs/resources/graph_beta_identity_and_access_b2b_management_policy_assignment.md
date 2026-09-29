---
page_title: "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Applies a Microsoft Entra B2B management policy to an application or service principal. Import ID format: b2b_management_policy_id/directory_object_id.
---

# microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment (Resource)

Applies a Microsoft Entra B2B management policy to an application or service principal. Import ID format: `b2b_management_policy_id/directory_object_id`.

## Microsoft Graph API

The documented `appliesTo/$ref` endpoints ([Add](https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-post-appliesto?view=graph-rest-beta), [Remove](https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-delete-appliesto?view=graph-rest-beta)) reject writes with `400 Property 'appliesTo' is read-only`. This resource therefore uses the application or service principal side:

- Create: `POST /beta/{applications|servicePrincipals}/{id}/policies/$ref`
- Read: `GET /beta/policies/b2bManagementPolicies/{id}/appliesTo`
- Delete: `DELETE /beta/{applications|servicePrincipals}/{id}/policies/{policyId}/$ref`

## Behavior Notes

- The target must be an application or a service principal; the type is detected and exposed as `directory_object_type`.
- An assignment that already exists outside Terraform is not adopted: create fails with the ID to import.
- Replacing the policy (for example, changing its `is_organization_default`) replaces the assignment too, because the policy ID changes.
- A refresh treats a missing assignment as removed only after it stays missing for 60 seconds, because Microsoft Entra replicas briefly return stale results after writes.

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Policy.Read.B2BManagementPolicy`
- `Policy.ReadWrite.B2BManagementPolicy`
- `Application.Read.All`
- `Application.ReadWrite.OwnedBy` if the caller owns the target, otherwise `Application.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

Without `Application.Read.All`, Microsoft Graph returns a truncated `appliesTo` response once an application is assigned; the provider reports this as a permission error.

## Example Usage

Each example includes its policy and target resources. Use the target object ID for assignments, rather than the application (client) ID.

The service-principal examples also use the HashiCorp `time` provider to allow the new application registration to replicate.

### Assign to an Application

```terraform
# Create the application registration used by this example.
resource "microsoft365_graph_beta_applications_application" "application" {
  display_name = "B2B policy target - application"
}

# This non-default policy does not change tenant-wide invitation restrictions.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "application" {
  display_name            = "B2B policy - application"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.net"]
        }
      }
    })
  ]
}

# Use the application object ID, not its app_id (client ID).
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "application" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.application.id
  directory_object_id      = microsoft365_graph_beta_applications_application.application.id
}
```

### Assign to a Service Principal

```terraform
# Create the application registration used by this example.
resource "microsoft365_graph_beta_applications_application" "service_principal" {
  display_name = "B2B policy target - service principal"
}

# Allow the application registration to replicate before creating its service principal.
resource "time_sleep" "service_principal_application_replication" {
  depends_on      = [microsoft365_graph_beta_applications_application.service_principal]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal" "service_principal" {
  app_id = microsoft365_graph_beta_applications_application.service_principal.app_id

  depends_on = [time_sleep.service_principal_application_replication]
}

# This non-default policy does not change tenant-wide invitation restrictions.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "service_principal" {
  display_name            = "B2B policy - service principal"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.net"]
        }
      }
    })
  ]
}

# Use the service principal object ID, not its app_id (client ID).
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "service_principal" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.service_principal.id
  directory_object_id      = microsoft365_graph_beta_applications_service_principal.service_principal.id
}
```

### Assign to an Application and Its Service Principal

```terraform
# Create the application registration used by this example.
resource "microsoft365_graph_beta_applications_application" "both_targets" {
  display_name = "B2B policy target - both targets"
}

# Allow the application registration to replicate before creating its service principal.
resource "time_sleep" "both_application_replication" {
  depends_on      = [microsoft365_graph_beta_applications_application.both_targets]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal" "both_targets" {
  app_id = microsoft365_graph_beta_applications_application.both_targets.app_id

  depends_on = [time_sleep.both_application_replication]
}

# This non-default policy does not change tenant-wide invitation restrictions.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "both_targets" {
  display_name            = "B2B policy - both targets"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.net"]
        }
      }
    })
  ]
}

# Use the application object ID, not its app_id (client ID).
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "both_application" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.both_targets.id
  directory_object_id      = microsoft365_graph_beta_applications_application.both_targets.id
}

# Use the service principal object ID, not its app_id (client ID).
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "both_service_principal" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.both_targets.id
  directory_object_id      = microsoft365_graph_beta_applications_service_principal.both_targets.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `b2b_management_policy_id` (String) The unique identifier of the B2B management policy to assign.
- `directory_object_id` (String) The object ID of the application or service principal to apply the policy to.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `directory_object_type` (String) The detected type of `directory_object_id`: `application` or `servicePrincipal`.
- `id` (String) Composite ID `b2b_management_policy_id/directory_object_id`; Microsoft Graph has no assignment ID.

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
terraform import microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment.application 00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002
```
