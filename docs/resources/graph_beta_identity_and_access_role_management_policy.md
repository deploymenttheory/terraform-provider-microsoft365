---
page_title: "microsoft365_graph_beta_identity_and_access_role_management_policy Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages the rules of an existing Privileged Identity Management (PIM) role management policy using the /policies/roleManagementPolicies endpoint. Microsoft Entra creates one policy for every role and scope, so this resource adopts the existing policy, identified either by role_definition_id (with scope_id and scope_type) or by policy_id, and updates only the rules listed in the configuration. Destroying the resource removes it from Terraform state and leaves the policy rules as they are.
---

# microsoft365_graph_beta_identity_and_access_role_management_policy (Resource)

Manages the rules of an existing Privileged Identity Management (PIM) role management policy using the `/policies/roleManagementPolicies` endpoint. Microsoft Entra creates one policy for every role and scope, so this resource adopts the existing policy, identified either by `role_definition_id` (with `scope_id` and `scope_type`) or by `policy_id`, and updates only the rules listed in the configuration. Destroying the resource removes it from Terraform state and leaves the policy rules as they are.

## Microsoft Documentation

- [Configure Microsoft Entra role settings in PIM](https://learn.microsoft.com/en-us/entra/id-governance/privileged-identity-management/pim-how-to-change-default-settings)
- [Overview of role management policies in PIM APIs](https://learn.microsoft.com/en-us/graph/identity-governance-pim-rules-overview)
- [unifiedRoleManagementPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicy?view=graph-rest-beta)
- [Get unifiedRoleManagementPolicy](https://learn.microsoft.com/en-us/graph/api/unifiedrolemanagementpolicy-get?view=graph-rest-beta)
- [Update unifiedRoleManagementPolicy](https://learn.microsoft.com/en-us/graph/api/unifiedrolemanagementpolicy-update?view=graph-rest-beta)
- [List roleManagementPolicyAssignments](https://learn.microsoft.com/en-us/graph/api/policyroot-list-rolemanagementpolicyassignments?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `RoleManagementPolicy.Read.Directory`
- `RoleManagementPolicy.ReadWrite.Directory`

**Optional:**
- `RoleManagementPolicy.Read.AzureADGroup` (PIM for Groups policies, `scope_type = "Group"`)
- `RoleManagementPolicy.ReadWrite.AzureADGroup` (PIM for Groups policies, `scope_type = "Group"`)

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v1.7.0 | Experimental | Initial release |

## Example Usage

### Limit the duration of active assignments

```terraform
# Limit permanent active assignments of the Usage Summary Reports Reader role to 90 days.
# The policy is found through the role definition at the directory scope ("/").
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "usage_summary_reports_reader" {
  role_definition_id = "75934031-6c7e-415a-99d7-48dbd49e875e"

  expiration_rules = {
    Expiration_Admin_Assignment = {
      is_expiration_required = true
      maximum_duration       = "P90D"
    }
  }
}
```

### Harden role activation

```terraform
# Harden activation of the Security Administrator role: activations last at most 4 hours, require MFA,
# a justification and a ticket, and require the "c1" Conditional Access authentication context.
# Approvers receive critical notifications only, sent to an additional mailbox.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "security_administrator" {
  role_definition_id = "194ae4cb-b126-40b2-bd5b-6091b380977d"
  scope_id           = "/"
  scope_type         = "DirectoryRole"

  expiration_rules = {
    Expiration_EndUser_Assignment = {
      is_expiration_required = true
      maximum_duration       = "PT4H"
    }
  }

  enablement_rules = {
    Enablement_EndUser_Assignment = {
      enabled_rules = ["MultiFactorAuthentication", "Justification", "Ticketing"]
    }
  }

  authentication_context_rules = {
    AuthenticationContext_EndUser_Assignment = {
      is_enabled  = true
      claim_value = "c1"
    }
  }

  notification_rules = {
    Notification_Approver_EndUser_Assignment = {
      notification_level            = "Critical"
      is_default_recipients_enabled = true
      notification_recipients       = ["pim-approvers@contoso.com"]
    }
  }
}
```

### Adopt a policy by policy ID

```terraform
# Adopt a policy by its ID instead of by role definition. role_definition_id, scope_id and scope_type
# are then computed from the policy and its assignment.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "by_policy_id" {
  policy_id = "DirectoryRole_00000000-0000-0000-0000-000000000000_00000000-0000-0000-0000-000000000000"

  expiration_rules = {
    Expiration_Admin_Eligibility = {
      is_expiration_required = true
      maximum_duration       = "P180D"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `authentication_context_rules` (Attributes Map) Authentication context rules to manage, keyed by rule ID, for example `AuthenticationContext_EndUser_Assignment`. Rules that are not listed are left unchanged. (see [below for nested schema](#nestedatt--authentication_context_rules))
- `enablement_rules` (Attributes Map) Enablement rules to manage, keyed by rule ID, for example `Enablement_Admin_Eligibility`, `Enablement_Admin_Assignment` or `Enablement_EndUser_Assignment`. Rules that are not listed are left unchanged. (see [below for nested schema](#nestedatt--enablement_rules))
- `expiration_rules` (Attributes Map) Expiration rules to manage, keyed by rule ID, for example `Expiration_Admin_Eligibility`, `Expiration_Admin_Assignment` or `Expiration_EndUser_Assignment`. Rules that are not listed are left unchanged. (see [below for nested schema](#nestedatt--expiration_rules))
- `notification_rules` (Attributes Map) Notification rules to manage, keyed by rule ID, for example `Notification_Admin_Admin_Eligibility` or `Notification_Approver_EndUser_Assignment`. The notification and recipient types are fixed by the rule ID and are taken from the policy. Rules that are not listed are left unchanged. (see [below for nested schema](#nestedatt--notification_rules))
- `policy_id` (String) The ID of the role management policy to adopt. Conflicts with `role_definition_id`, `scope_id` and `scope_type`. When the policy is identified by role definition, this is computed from the policy assignment. Changing it adopts a different policy.
- `role_definition_id` (String) The ID of the role definition whose policy is adopted. Use the role template ID for built-in Microsoft Entra roles, or `member` / `owner` for PIM for Groups. The policy is found through `/policies/roleManagementPolicyAssignments`. Computed when the policy is identified by `policy_id`. Changing it adopts a different policy.
- `scope_id` (String) The scope of the policy. Use `/` for Microsoft Entra roles (the default when `role_definition_id` is set) or the group ID for PIM for Groups. Changing it adopts a different policy.
- `scope_type` (String) The type of the policy scope: `DirectoryRole` for Microsoft Entra roles (the default when `role_definition_id` is set) or `Group` for PIM for Groups. Changing it adopts a different policy.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `description` (String) The description of the policy, set by Microsoft Entra.
- `display_name` (String) The display name of the policy, set by Microsoft Entra.
- `id` (String) The ID of the role management policy, for example `DirectoryRole_{tenantId}_{policyGuid}`.
- `is_organization_default` (Boolean) Whether the policy is the organization default policy.
- `last_modified_date_time` (String) When the policy was last modified. Null until the policy is first changed.

<a id="nestedatt--authentication_context_rules"></a>
### Nested Schema for `authentication_context_rules`

Required:

- `is_enabled` (Boolean) Whether activation requires the Conditional Access authentication context in `claim_value`.

Optional:

- `claim_value` (String) The authentication context class reference ID, `c1` to `c99`. Required when `is_enabled` is `true` and must be omitted when it is `false`.


<a id="nestedatt--enablement_rules"></a>
### Nested Schema for `enablement_rules`

Required:

- `enabled_rules` (Set of String) The requirements that are enabled: `MultiFactorAuthentication`, `Justification` and `Ticketing`. Use an empty set to require none.


<a id="nestedatt--expiration_rules"></a>
### Nested Schema for `expiration_rules`

Required:

- `is_expiration_required` (Boolean) Whether assignments or activations must have an expiration date. `false` allows permanent assignments.
- `maximum_duration` (String) The maximum duration as an ISO 8601 duration made of days, hours, minutes and seconds, for example `P180D`, `PT8H` or `P1DT12H`. Equivalent values, such as `PT24H` for `P1D`, are kept as written.


<a id="nestedatt--notification_rules"></a>
### Nested Schema for `notification_rules`

Required:

- `is_default_recipients_enabled` (Boolean) Whether the default recipients for this recipient type receive the notification.
- `notification_level` (String) Which notifications are sent: `All`, `Critical` or `None`.

Optional:

- `notification_recipients` (Set of String) Additional email addresses that receive the notification. Defaults to an empty set.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Adopting an existing policy
- Microsoft Entra creates one role management policy for every role and scope. Policies cannot be created or deleted, so this resource adopts the existing policy and updates its rules.
- Identify the policy with `role_definition_id` (and optionally `scope_id` and `scope_type`, which default to `/` and `DirectoryRole`), or with `policy_id`. The policy is found through `GET /policies/roleManagementPolicyAssignments`, filtered on `scopeId`, `scopeType` and `roleDefinitionId`.
- The policy ID does not always end with the role definition ID, so look it up rather than building it.
- Changing `policy_id`, `role_definition_id`, `scope_id` or `scope_type` adopts a different policy. The previous policy keeps the rules last applied to it.
- Manage each policy from one resource only.

### Rules
- Rules are configured in one map per rule type, keyed by rule ID, for example `Expiration_EndUser_Assignment`. Each rule ID must exist on the policy with the matching type; otherwise the apply fails and lists the rules the policy has.
- Only the configured rules are sent to Graph, in a single `PATCH /policies/roleManagementPolicies/{id}`. Rules that are not configured are never changed or tracked.
- Removing a rule from the configuration stops managing it; the rule keeps its current values in Microsoft Entra.
- Rule IDs for Microsoft Entra roles follow the pattern `{Type}_{Caller}_{Level}`, where the caller is `Admin` (assignments made by administrators) or `EndUser` (activations), and the level is `Eligibility` or `Assignment`. Notification rule IDs also include the recipient type (`Admin`, `Requestor` or `Approver`).
- Notification rules take their `notificationType` and `recipientType` from the existing rule, because Graph requires them on every update and they are fixed by the rule ID.
- An authentication context rule needs `claim_value` when `is_enabled` is `true`. Graph clears the claim value of a disabled rule, so omit it when `is_enabled` is `false`. Graph does not check that the authentication context exists.
- `maximum_duration` takes an ISO 8601 duration of days, hours, minutes and seconds. Graph rejects week durations (`P2W`) and normalises equivalent values, for example `PT24H` to `P1D`; the configured spelling is kept in state.

### Not supported in this version
- Approval rules (`unifiedRoleManagementPolicyApprovalRule`), including approval stages and approvers.
- Custom extension rules and the `target` of each rule, which Microsoft Entra fixes for role management policies.
- PIM for Groups policies (`scope_type = "Group"`) are supported by the schema but have not been acceptance tested.

### Import and destroy
- Import by policy ID. After an import, every expiration, enablement, notification and authentication context rule of the policy is tracked in state. Configure all of them to get an empty plan; rules left out of the configuration drop out of state on the next apply without being changed.
- Destroying the resource only removes it from Terraform state. The policy keeps the rules last applied to it. To return rules to their previous values, apply those values before destroying.

### Consistency
- Reads of role management policies are eventually consistent. After an update, the provider re-reads the policy until every configured rule shows the applied values.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using the role management policy ID: DirectoryRole_{tenantId}_{policyGuid}
terraform import microsoft365_graph_beta_identity_and_access_role_management_policy.example DirectoryRole_00000000-0000-0000-0000-000000000000_00000000-0000-0000-0000-000000000000
```
