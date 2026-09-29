---
page_title: "microsoft365_graph_beta_identity_and_access_b2b_management_policy Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft Entra B2B management policies using the /policies/b2bManagementPolicies endpoint. A B2B management policy controls B2B collaboration settings such as the domains that can be invited.
---

# microsoft365_graph_beta_identity_and_access_b2b_management_policy (Resource)

Manages Microsoft Entra B2B management policies using the `/policies/b2bManagementPolicies` endpoint. A B2B management policy controls B2B collaboration settings such as the domains that can be invited.

## Organization default

Only the policy with `is_organization_default = true` takes effect tenant-wide (for example, invitations to a blocked domain fail). A policy that is not the default has no tenant-wide effect.

- Only one policy can be the default. Creating a second one fails with `400 Another object with the same value for property isOrganizationDefault already exists`.
- Microsoft Graph cannot change `isOrganizationDefault` on an existing policy (it returns `500`), so changing `is_organization_default` replaces the policy. The policy gets a new ID, and assignments that reference it are replaced as well.
- Moving the default to another policy in one apply works. Enforcement lags the API by roughly 10–20 seconds: in a live run, the old restriction was still enforced about 20 seconds after it was deleted, and the new one took effect after about 10 seconds.
- If the new default happens to be created before the old one is deleted, the apply fails with the `400` above. Apply again to finish.
- Do not use `create_before_destroy` on these policies. When the default moves, the new default is created while the old one still exists, so its creation fails. The old default's own replacement can still complete, which leaves the tenant without a default policy until a later apply succeeds (observed live).
- Deleting the default policy removes the tenant's B2B domain restrictions.

## Behavior Notes

- The documented `description` property is not supported: Microsoft Graph rejects any write that contains it with `404 Request_ResourceNotFound`.
- A refresh treats a `404` as deletion only after it persists for 60 seconds, because Microsoft Entra replicas briefly return `404` after writes.

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Policy.Read.B2BManagementPolicy`
- `Policy.ReadWrite.B2BManagementPolicy`

**Optional:**
- `None` `[N/A]`

## Example Usage

Each example is a separate scenario. Non-default policies do not change tenant-wide invitation restrictions; use the organization-default example to enforce a policy.

### Allow Invitations from Selected Domains

```terraform
# Stage an allow-list policy for selected partner domains.
# Set is_organization_default to true to enforce it tenant-wide.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "allow_domains" {
  display_name            = "B2B invitations - allowed partner domains"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          AllowedDomains = ["example.com", "example.net"]
        }
      }
    })
  ]
}
```

### Block Invitations from Selected Domains

```terraform
# Stage a block-list policy while allowing invitations to other domains.
# AllowedDomains and BlockedDomains are alternative modes; configure only one.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "block_domains" {
  display_name            = "B2B invitations - blocked domains"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.net", "example.org"]
        }
      }
    })
  ]
}
```

### Allow Invitations from Any Domain

```terraform
# An empty domain policy specifies no domain restrictions.
# This non-default example does not change the active organization policy.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "any_domain" {
  display_name            = "B2B invitations - any domain"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {}
      }
    })
  ]
}
```

### Enforce a Policy Organization-Wide

```terraform
# Only one organization-default B2B management policy can exist.
# If one already exists, import it instead of creating a second default.
# Changing is_organization_default replaces the policy; do not use create_before_destroy.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "organization_default" {
  display_name            = "B2B invitations - organization default"
  is_organization_default = true

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          AllowedDomains = ["example.com", "example.net"]
        }
      }
    })
  ]
}
```

### Supply Empty Additional Settings

```terraform
# Keep an allow list and explicitly supply empty preview and auto-redemption lists.
# The resource accepts these additional settings in the same JSON definition.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "other_settings" {
  display_name            = "B2B invitations - explicit empty additional settings"
  is_organization_default = false

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          AllowedDomains = ["example.com", "example.net"]
        }
        PreviewPolicy = {
          Features = []
        }
        AutoRedeemPolicy = {
          AdminConsentedForUsersIntoTenantIds = []
          NoAADConsentForUsersFromTenantsIds  = []
        }
      }
    })
  ]
}
```

### Supply a Definition as JSON Text

```terraform
# Use a heredoc when maintaining the definition as JSON rather than HCL.
# JSON formatting and whitespace do not cause a recurring Terraform diff.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "heredoc" {
  display_name            = "B2B invitations - JSON definition"
  is_organization_default = false

  definition = [<<-JSON
    {
      "B2BManagementPolicy": {
        "InvitationsAllowedAndBlockedDomainsPolicy": {
          "AllowedDomains": [
            "example.com",
            "example.net"
          ]
        }
      }
    }
  JSON
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `definition` (List of String) A collection containing the policy definition as a JSON string, for example `{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["example.com"]}}}`.
- `display_name` (String) The display name of the B2B management policy.

### Optional

- `is_organization_default` (Boolean) If `true`, the policy applies tenant-wide as the organization default. Only one policy can be the organization default. Microsoft Graph cannot change this on an existing policy, so changing it replaces the policy. Defaults to `false`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier of the B2B management policy.

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
# Import using the B2B management policy ID (GUID) from Microsoft Graph

# {id} - The unique identifier (GUID) of the B2B management policy
terraform import microsoft365_graph_beta_identity_and_access_b2b_management_policy.organization_default 00000000-0000-0000-0000-000000000000
```
