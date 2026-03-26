---
page_title: "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages a scoped role assignment for an administrative unit in Microsoft Entra ID using the /administrativeUnits/{id}/scopedRoleMembers endpoint. Scoped role members allow directory roles (such as User Administrator or Helpdesk Administrator) to be delegated to a user within the scope of a specific administrative unit rather than the entire tenant. All fields are immutable after creation; any change triggers a destroy and recreate.
  Required permissions: RoleManagement.ReadWrite.Directory
  Import format: administrative_unit_id/scoped_role_membership_id
---

# microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment (Resource)

Manages a scoped role assignment for an administrative unit in Microsoft Entra ID using the `/administrativeUnits/{id}/scopedRoleMembers` endpoint. Scoped role members allow directory roles (such as User Administrator or Helpdesk Administrator) to be delegated to a user within the scope of a specific administrative unit rather than the entire tenant. All fields are immutable after creation; any change triggers a destroy and recreate.

**Required permissions:** `RoleManagement.ReadWrite.Directory`

**Import format:** `administrative_unit_id/scoped_role_membership_id`

## Microsoft Documentation

- [Administrative units overview](https://learn.microsoft.com/en-us/entra/identity/role-based-access-control/administrative-units)
- [scopedRoleMembership resource type](https://learn.microsoft.com/en-us/graph/api/resources/scopedrolemembership?view=graph-rest-beta)
- [Add a scopedRoleMember](https://learn.microsoft.com/en-us/graph/api/administrativeunit-post-scopedrolemembers?view=graph-rest-beta)
- [List scopedRoleMembers](https://learn.microsoft.com/en-us/graph/api/administrativeunit-list-scopedrolemembers?view=graph-rest-beta)
- [Get a scopedRoleMember](https://learn.microsoft.com/en-us/graph/api/administrativeunit-get-scopedrolemembers?view=graph-rest-beta)
- [Remove a scopedRoleMember](https://learn.microsoft.com/en-us/graph/api/administrativeunit-delete-scopedrolemembers?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `RoleManagement.ReadWrite.Directory`

**Read-only (for data sources):**
- `RoleManagement.Read.Directory`
- `Directory.Read.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.1.0-alpha | Experimental | Initial release |

## Example Usage

### AURA001: User Administrator Scoped to an Administrative Unit

```terraform
# AURA001: User Administrator scoped to an Administrative Unit
# Assigns the User Administrator role to a user, scoped to a specific administrative unit.
# The role holder can manage users only within that administrative unit.

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "user_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.finance.id
  # User Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "fe930be7-5e62-47db-91af-98c3a49a38b1"
  role_member_id    = microsoft365_graph_beta_users_user.helpdesk_lead.id
}
```

### AURA002: Helpdesk Administrator Scoped to an Administrative Unit

```terraform
# AURA002: Helpdesk Administrator scoped to an Administrative Unit
# Assigns the Helpdesk Administrator role to a user, scoped to a specific administrative unit.
# The role holder can reset passwords and manage service requests only for users within that unit.

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "helpdesk_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.it_department.id
  # Helpdesk Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "729827e3-9c14-49f7-bb1b-9608f156bbb8"
  role_member_id    = microsoft365_graph_beta_users_user.it_support_agent.id
}
```

### AURA003: Multiple Role Assignments on the Same Administrative Unit

```terraform
# AURA003: Multiple role assignments scoped to the same Administrative Unit
# Assigns both User Administrator and Helpdesk Administrator roles to different users
# within the same administrative unit, enabling a tiered delegation model.

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "au_user_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.regional_office.id
  # User Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "fe930be7-5e62-47db-91af-98c3a49a38b1"
  role_member_id    = microsoft365_graph_beta_users_user.regional_it_manager.id
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "au_helpdesk_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.regional_office.id
  # Helpdesk Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "729827e3-9c14-49f7-bb1b-9608f156bbb8"
  role_member_id    = microsoft365_graph_beta_users_user.regional_helpdesk.id
}
```

### AURA004: Full Example — Create AU, User, and Role Assignment Together

```terraform
# AURA004: Full example — create AU, user, and scoped role assignment together
# Demonstrates creating all dependent resources inline and chaining them
# with depends_on for correct provisioning order.

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "emea_office" {
  display_name = "EMEA Office"
  description  = "Administrative unit for EMEA region users and devices"
  hard_delete  = true
}

resource "microsoft365_graph_beta_users_user" "emea_it_admin" {
  user_principal_name = "emea-it-admin@contoso.com"
  display_name        = "EMEA IT Administrator"
  mail_nickname       = "emea-it-admin"
  account_enabled     = true
  password_profile = {
    password                           = "ChangeMe123!"
    force_change_password_next_sign_in = true
  }
  hard_delete = true
}

resource "microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment" "emea_user_admin" {
  administrative_unit_id = microsoft365_graph_beta_identity_and_access_administrative_unit.emea_office.id
  # User Administrator role ID (well-known in Microsoft Entra ID)
  directory_role_id = "fe930be7-5e62-47db-91af-98c3a49a38b1"
  role_member_id    = microsoft365_graph_beta_users_user.emea_it_admin.id

  depends_on = [
    microsoft365_graph_beta_identity_and_access_administrative_unit.emea_office,
    microsoft365_graph_beta_users_user.emea_it_admin,
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `administrative_unit_id` (String) The unique identifier of the administrative unit to which this role assignment is scoped. Changing this value forces a new resource to be created.
- `directory_role_id` (String) The tenant-specific object ID of the activated directoryRole to assign within the administrative unit scope. This is **not** the well-known roleTemplateId — it is the object ID of the directoryRole as activated in your tenant. Use `GET /directoryRoles` to list activated roles and find the correct object ID. Only roles that support administrative unit scoping are valid (e.g. User Administrator, Helpdesk Administrator). Changing this value forces a new resource to be created.
- `role_member_id` (String) The unique identifier of the user or service principal to assign the directory role to. Changing this value forces a new resource to be created.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier of the scoped role membership. Assigned by the API on creation. Read-only.
- `role_member_display_name` (String) The display name of the role member. Populated by the API after creation. Read-only.
- `role_member_user_principal_name` (String) The user principal name (UPN) of the role member. Populated by the API after creation. Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Immutability
All input fields on this resource (`administrative_unit_id`, `directory_role_id`, `role_member_id`) are immutable.
Any change to these values will destroy the existing assignment and create a new one.
Plan accordingly to avoid unintended loss of delegated access.

### Providing the `directory_role_id`

The `directory_role_id` field requires the **tenant-specific directoryRole object ID**, not the well-known `roleTemplateId`.
These two values look similar (both are GUIDs) but are different:

- **`roleTemplateId`** — cross-tenant, well-known (e.g. `fe930be7-5e62-47db-91af-98c3a49a38b1` for User Administrator)
- **`id` (directoryRole object ID)** — tenant-specific, assigned when the role is activated in your tenant

To find the object ID for a role in your tenant, run:

```bash
az rest --method GET --uri "https://graph.microsoft.com/beta/directoryRoles" \
  --query "value[].{name:displayName, id:id, templateId:roleTemplateId}" \
  --output table
```

Or with the Graph Explorer / `curl`:

```bash
curl -H "Authorization: Bearer <token>" \
  "https://graph.microsoft.com/beta/directoryRoles" | jq '.value[] | {displayName, id, roleTemplateId}'
```

If the role has not yet been activated in your tenant, activate it first:

```bash
curl -X POST -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"roleTemplateId": "fe930be7-5e62-47db-91af-98c3a49a38b1"}' \
  "https://graph.microsoft.com/beta/directoryRoles"
```

### Supported Roles
Not all directory roles support administrative unit scoping. Roles that support scoped assignment include:

| Role | roleTemplateId (cross-tenant) |
|------|-------------------------------|
| Authentication Administrator | `c4e39bd9-1100-46d3-8c65-fb160da0071f` |
| Cloud Device Administrator | `7698a772-787b-4ac8-901f-60d6b08affd2` |
| Groups Administrator | `fdd7a751-b60b-444a-984c-02652fe8fa1c` |
| Helpdesk Administrator | `729827e3-9c14-49f7-bb1b-9608f156bbb8` |
| License Administrator | `4d6ac14f-3453-41d0-bef9-a3e0c569773a` |
| Password Administrator | `966707d0-3269-4727-9be2-8c3a10f19b9d` |
| Privileged Authentication Administrator | `7be44c8a-adaf-4e2a-84d6-ab2649e08a13` |
| User Administrator | `fe930be7-5e62-47db-91af-98c3a49a38b1` |

Use the `roleTemplateId` above with `GET /directoryRoles?$filter=roleTemplateId eq '<id>'` to find your tenant's directoryRole object ID.

For a complete and up-to-date list, see the [Microsoft documentation on roles you can assign at administrative unit scope](https://learn.microsoft.com/en-us/entra/identity/role-based-access-control/admin-units-assign-roles#roles-that-can-be-assigned-with-administrative-unit-scope).

### Role Member Types
The `role_member_id` must refer to a **user** object in your directory. Groups and service principals cannot currently be assigned to scoped roles.

### Computed Attributes
After creation, the following attributes are populated by the API and stored in state:
- `id` — the scoped role membership ID assigned by Microsoft Entra ID
- `role_member_display_name` — the display name of the assigned user
- `role_member_user_principal_name` — the UPN of the assigned user

### Relationship with Other Resources
- Use [`microsoft365_graph_beta_identity_and_access_administrative_unit`](../resources/graph_beta_identity_and_access_administrative_unit) to create the administrative unit.
- Use [`microsoft365_graph_beta_identity_and_access_administrative_unit_membership`](../resources/graph_beta_identity_and_access_administrative_unit_membership) to manage the users and groups that belong **to** the administrative unit.
- This resource controls who has **administrative authority over** the administrative unit.

### Best Practices
- Assign the least-privileged role that satisfies the operational requirement (e.g. prefer Helpdesk Administrator over User Administrator where password resets are the only need).
- Use `depends_on` when creating the administrative unit and the role member user in the same Terraform configuration to ensure correct provisioning order.
- One `microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment` resource per assignment. Use multiple resource blocks to assign multiple roles or the same role to multiple users.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import format: administrative_unit_id/scoped_role_membership_id
terraform import microsoft365_graph_beta_identity_and_access_administrative_unit_directory_role_assignment.example 00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111
```
