---
page_title: "microsoft365_graph_beta_identity_and_access_mobility_management_policy Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages a Microsoft Entra mobility management policy (the MDM or MAM automatic enrollment user scope) using the Microsoft Graph beta /policies/mobileDeviceManagementPolicies and /policies/mobileAppManagementPolicies endpoints. A mobility management policy is created by the service when a mobility management application such as Microsoft Intune is added to the tenant, so this resource adopts an existing policy and never creates or deletes one.
  create adopts the policy and applies the configured settings; update applies changes. Settings that are not configured keep their current tenant value and are never written. applies_to is written in its own request because Microsoft Graph rejects it together with other properties. selected is applied by adding the included_groups references, and only the group references that differ are added or removed. On destroy, Terraform only removes the resource from state and leaves the tenant policy as it is.
  Permissions: Microsoft Graph does not support application permissions for these endpoints. The provider must authenticate with a delegated (work or school account) identity that holds Policy.ReadWrite.MobilityManagement and a role such as Cloud Device Administrator.
---

# microsoft365_graph_beta_identity_and_access_mobility_management_policy (Resource)

Manages a Microsoft Entra mobility management policy (the MDM or MAM automatic enrollment user scope) using the Microsoft Graph beta `/policies/mobileDeviceManagementPolicies` and `/policies/mobileAppManagementPolicies` endpoints. A mobility management policy is created by the service when a mobility management application such as Microsoft Intune is added to the tenant, so this resource adopts an existing policy and never creates or deletes one.

`create` adopts the policy and applies the configured settings; `update` applies changes. Settings that are not configured keep their current tenant value and are never written. `applies_to` is written in its own request because Microsoft Graph rejects it together with other properties. `selected` is applied by adding the `included_groups` references, and only the group references that differ are added or removed. On `destroy`, Terraform only removes the resource from state and leaves the tenant policy as it is.

**Permissions:** Microsoft Graph does not support application permissions for these endpoints. The provider must authenticate with a delegated (work or school account) identity that holds `Policy.ReadWrite.MobilityManagement` and a role such as Cloud Device Administrator.

## Microsoft Documentation

- [mobilityManagementPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/mobilitymanagementpolicy?view=graph-rest-beta)
- [Get mobileDeviceManagementPolicy](https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-get?view=graph-rest-beta)
- [Update mobileDeviceManagementPolicy](https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-update?view=graph-rest-beta)
- [Add includedGroups](https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-post-includedgroups?view=graph-rest-beta)
- [Remove includedGroups](https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-delete-includedgroups?view=graph-rest-beta)
- [Get mobileAppManagementPolicy](https://learn.microsoft.com/en-us/graph/api/mobileappmanagementpolicies-get?view=graph-rest-beta)
- [Update mobileAppManagementPolicy](https://learn.microsoft.com/en-us/graph/api/mobileappmanagementpolicies-update?view=graph-rest-beta)
- [Set up automatic enrollment for Windows devices](https://learn.microsoft.com/en-us/intune/intune-service/enrollment/windows-enroll)

## Microsoft Graph API Permissions

Microsoft Graph does **not** support application permissions for these endpoints. The following **delegated** (work or school account) permissions are needed in order to use this resource, and the signed-in user needs a Microsoft Entra role such as Cloud Device Administrator:

**Required:**
- `Policy.Read.All`
- `Policy.ReadWrite.MobilityManagement`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v1.7.0 | Experimental | Initial release |

## Example Usage

### Minimal

```terraform
# Adopt the Microsoft Intune MDM user scope without changing any setting.
# Settings that are not configured keep their current tenant value.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "mdm" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"
}
```

### Maximal

```terraform
resource "microsoft365_graph_beta_groups_group" "intune_enrollment" {
  display_name     = "Intune automatic enrollment"
  mail_nickname    = "intune-automatic-enrollment"
  mail_enabled     = false
  security_enabled = true
}

# Limit Microsoft Intune MDM automatic enrollment to one group and set every MDM setting.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "mdm" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  compliance_url   = "https://portal.manage.microsoft.com/?portalAction=Compliance"
  discovery_url    = "https://enrollment.manage.microsoft.com/enrollmentserver/discovery.svc"
  terms_of_use_url = "https://portal.manage.microsoft.com/TermsofUse.aspx"

  is_mdm_enrollment_during_registration_disabled = false

  applies_to      = "selected"
  included_groups = [microsoft365_graph_beta_groups_group.intune_enrollment.id]
}

# Turn off the Microsoft Intune MAM user scope.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "mam" {
  policy_type = "app"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  applies_to = "none"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `policy_id` (String) The ID of the mobility management policy, which is the application ID of the mobility management application. For Microsoft Intune this is `0000000a-0000-0000-c000-000000000000`. Changing this forces a new resource.
- `policy_type` (String) The kind of mobility management policy. `device` manages the MDM user scope (`/policies/mobileDeviceManagementPolicies`) and `app` manages the MAM user scope (`/policies/mobileAppManagementPolicies`). Changing this forces a new resource.

### Optional

- `applies_to` (String) The user scope of the policy. Possible values are: `none`, `all`, `selected`. `selected` requires at least one entry in `included_groups`. `all` and `none` remove every included group. **Warning:** `all` makes every user in the tenant in scope for automatic enrollment. When not configured, the current tenant value is kept.
- `compliance_url` (String) The compliance URL of the mobility management application. When not configured, the current tenant value is kept.
- `discovery_url` (String) The discovery URL of the mobility management application. When not configured, the current tenant value is kept.
- `included_groups` (Set of String) The IDs of the Microsoft Entra groups in scope of the policy. Only valid when `applies_to` is `selected`. Group references are added and removed individually, so only the groups that differ are changed. When not configured, the current tenant groups are kept.
- `is_mdm_enrollment_during_registration_disabled` (Boolean) Whether users in an automatic enrollment configuration on Microsoft Entra registered devices are no longer prompted to enroll their device in MDM during the Entra account registration flow. Only valid when `policy_type` is `device`. When not configured, the current tenant value is kept.
- `terms_of_use_url` (String) The terms of use URL of the mobility management application. When not configured, the current tenant value is kept.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `description` (String) The description of the mobility management application. Read-only.
- `display_name` (String) The display name of the mobility management application. Read-only.
- `id` (String) The resource identifier in the format `{policy_type}/{policy_id}`, for example `device/0000000a-0000-0000-c000-000000000000`.
- `is_valid` (Boolean) Whether the policy is valid. Invalid policies cannot be updated. Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Service-Created Policies

- Mobility management policies are created by the service when a mobility management application such as Microsoft Intune is added to the tenant. This resource adopts an existing policy; it never creates one.
- `policy_type = "device"` manages the MDM user scope and `policy_type = "app"` manages the MAM user scope of the same application.
- Manage each policy from only one Terraform resource across all states.

### Scope

- `applies_to = "all"` puts every user in scope for automatic enrollment. Review the impact before applying it.
- `applies_to` cannot be sent with other properties, so the provider writes it in its own request.
- `selected` is applied by adding group references. Only the group references that differ from the tenant are added or removed.
- `none` and `all` remove every included group.
- When `applies_to` is not configured, the scope and included groups are left as they are.

### Unconfigured Settings

- URL settings and `is_mdm_enrollment_during_registration_disabled` that are not configured keep their current tenant value and are never written.

### Destroy

- Destroy only removes the resource from Terraform state. The policy, its scope and its included groups are left unchanged, because resetting the scope would silently stop automatic enrollment for users.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using {policy_type}/{policy_id}. policy_type is device (MDM user scope) or app (MAM user scope);
# policy_id is the application ID of the mobility management application (Microsoft Intune shown).
terraform import microsoft365_graph_beta_identity_and_access_mobility_management_policy.mdm device/0000000a-0000-0000-c000-000000000000
terraform import microsoft365_graph_beta_identity_and_access_mobility_management_policy.mam app/0000000a-0000-0000-c000-000000000000
```
