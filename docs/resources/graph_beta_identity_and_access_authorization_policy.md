---
page_title: "microsoft365_graph_beta_identity_and_access_authorization_policy Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages the tenant-wide Microsoft Entra authorization policy using the Microsoft Graph beta /policies/authorizationPolicy endpoint.
  This is a singleton resource — one policy exists per tenant. The create and update operations use PATCH to configure the existing policy. On destroy, Terraform removes the resource from state only, leaving the policy unchanged. Manage only one instance per tenant across all Terraform states. Omitted optional settings retain their existing service values.
---

# microsoft365_graph_beta_identity_and_access_authorization_policy (Resource)

Manages the tenant-wide Microsoft Entra authorization policy using the Microsoft Graph beta `/policies/authorizationPolicy` endpoint.

This is a **singleton resource** — one policy exists per tenant. The `create` and `update` operations use PATCH to configure the existing policy. On `destroy`, Terraform removes the resource from state only, leaving the policy unchanged. Manage only one instance per tenant across all Terraform states. Omitted optional settings retain their existing service values.

## Microsoft Documentation

- [authorizationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/authorizationpolicy?view=graph-rest-beta)
- [Get authorizationPolicy](https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-get?view=graph-rest-beta&tabs=http)
- [Update authorizationPolicy](https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-update?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**

- `Policy.Read.All`
- `Policy.ReadWrite.Authorization`

**Optional:**

- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| Unreleased | Experimental | Initial release |

## Example Usage

### Basic Configuration

```terraform
resource "microsoft365_graph_beta_identity_and_access_authorization_policy" "example" {
  allowed_to_sign_up_email_based_subscriptions    = false
  allowed_to_use_sspr                             = true
  allow_email_verified_users_to_join_organization = false
  allow_user_consent_for_risky_apps               = false
  block_msol_powershell                           = true
  allow_invites_from                              = "adminsAndGuestInviters"
  default_user_role_permissions = {
    allowed_to_create_agent_identity_blueprints     = false
    allowed_to_create_apps                          = false
    allowed_to_create_security_groups               = false
    allowed_to_create_tenants                       = false
    allowed_to_read_bitlocker_keys_for_owned_device = false
    allowed_to_read_other_users                     = true
  }
}
```

### Complete Configuration

```terraform
resource "microsoft365_graph_beta_identity_and_access_authorization_policy" "example" {
  display_name                                              = "Authorization Policy"
  description                                               = "Tenant authorization settings"
  allowed_to_sign_up_email_based_subscriptions              = false
  allowed_to_use_sspr                                       = true
  allow_email_verified_users_to_join_organization           = false
  allow_invites_from                                        = "none"
  allow_user_consent_for_risky_apps                         = false
  block_msol_powershell                                     = true
  enabled_preview_features                                  = []
  guest_user_role_id                                        = "2af84b1e-32c8-42b7-82bc-daa82404023b"
  permission_grant_policy_ids_assigned_to_default_user_role = ["ManagePermissionGrantsForSelf.microsoft-user-default-low"]
  default_user_role_permissions = {
    allowed_to_create_agent_identity_blueprints     = false
    allowed_to_create_apps                          = false
    allowed_to_create_security_groups               = false
    allowed_to_create_tenants                       = false
    allowed_to_read_bitlocker_keys_for_owned_device = false
    allowed_to_read_other_users                     = true
  }
}
```

### Disable User Consent to Apps

```terraform
resource "microsoft365_graph_beta_identity_and_access_authorization_policy" "example" {
  allowed_to_sign_up_email_based_subscriptions              = false
  allowed_to_use_sspr                                       = true
  allow_email_verified_users_to_join_organization           = false
  allow_user_consent_for_risky_apps                         = false
  block_msol_powershell                                     = true
  allow_invites_from                                        = "none"
  enabled_preview_features                                  = []
  permission_grant_policy_ids_assigned_to_default_user_role = []
  default_user_role_permissions = {
    allowed_to_create_agent_identity_blueprints     = false
    allowed_to_create_apps                          = false
    allowed_to_create_security_groups               = false
    allowed_to_create_tenants                       = false
    allowed_to_read_bitlocker_keys_for_owned_device = false
    allowed_to_read_other_users                     = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allow_email_verified_users_to_join_organization` (Boolean) Whether users can join the tenant through email verification.
- `allow_user_consent_for_risky_apps` (Boolean) Whether users can consent to risky applications.
- `allowed_to_sign_up_email_based_subscriptions` (Boolean) Whether users can sign up for email-based subscriptions.
- `allowed_to_use_sspr` (Boolean) Whether tenant administrators can use self-service password reset.
- `block_msol_powershell` (Boolean) Whether user-based access to the legacy MSOnline PowerShell service endpoint is blocked.
- `default_user_role_permissions` (Attributes) Customizable permissions for the default user role. All Boolean permissions must be explicitly configured. (see [below for nested schema](#nestedatt--default_user_role_permissions))

### Optional

- `allow_invites_from` (String) Who can invite guests: `none`, `adminsAndGuestInviters`, `adminsGuestInvitersAndAllMembers`, or `everyone`.
- `description` (String) Description of the authorization policy.
- `display_name` (String) Display name of the authorization policy.
- `enabled_preview_features` (Set of String) Features enabled for private preview on the tenant. Use an empty set to clear the list.
- `guest_user_role_id` (String) Role template ID granted to guests: User (`a0b1b346-4d3e-4e8b-98f8-753987be4970`), Guest User (`10dae51f-b6af-4016-8d66-8c2a99b929b3`), or Restricted Guest User (`2af84b1e-32c8-42b7-82bc-daa82404023b`).
- `permission_grant_policy_ids_assigned_to_default_user_role` (Set of String) Permission grant policies assigned to the default user role. Values use `managePermissionGrantsForSelf.{id}` or `managePermissionGrantsForOwnedResource.{id}`. An empty set disables user consent to apps.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The fixed singleton identifier `authorizationPolicy`.

<a id="nestedatt--default_user_role_permissions"></a>
### Nested Schema for `default_user_role_permissions`

Required:

- `allowed_to_create_agent_identity_blueprints` (Boolean) Whether users can create agent identity blueprints.
- `allowed_to_create_apps` (Boolean) Whether users can register applications.
- `allowed_to_create_security_groups` (Boolean) Whether users can create security groups.
- `allowed_to_create_tenants` (Boolean) Whether users can create Microsoft Entra tenants.
- `allowed_to_read_bitlocker_keys_for_owned_device` (Boolean) Whether users can read BitLocker recovery keys for their owned devices.
- `allowed_to_read_other_users` (Boolean) Whether users can read other users. Microsoft advises keeping this permission enabled.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Singleton Resource**: Manage one instance per tenant across all Terraform states. The ID is always `authorizationPolicy`.
- **Create and Update Behavior**: Both operations use PATCH to configure the existing policy, then GET to refresh state.
- **Destroy Behavior**: Only removes Terraform state. The authorization policy and its settings remain in Microsoft Entra ID.
- **Boolean Settings**: All Boolean attributes are required, including every field in `default_user_role_permissions`.
- **Optional Settings**: Omitted optional attributes retain their existing service values. Set a collection to `[]` to clear it explicitly.

## Import

Import is supported using the following syntax:

```shell
# Import the tenant-wide authorization policy using its fixed singleton ID.
terraform import microsoft365_graph_beta_identity_and_access_authorization_policy.example authorizationPolicy
```
