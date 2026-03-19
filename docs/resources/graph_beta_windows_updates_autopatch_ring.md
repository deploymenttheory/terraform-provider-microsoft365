---
page_title: "microsoft365_graph_beta_windows_updates_autopatch_ring Resource - terraform-provider-microsoft365"
subcategory: "Windows Updates"

description: |-
  Manages a Windows Update policy ring using the /admin/windows/updates/policies/{policyId}/rings endpoint. A ring defines the deployment audience, deferral, and pause settings for quality updates within a Windows Update policy. The policy must already exist, managed by the microsoft365_graph_beta_windows_updates_autopatch_policy resource.
---

# microsoft365_graph_beta_windows_updates_autopatch_ring (Resource)

Manages a Windows Update policy ring using the `/admin/windows/updates/policies/{policyId}/rings` endpoint. A ring defines the deployment audience, deferral, and pause settings for quality updates within a Windows Update policy. The policy must already exist, managed by the `microsoft365_graph_beta_windows_updates_autopatch_policy` resource.

## Microsoft Documentation

- [ring resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-ring?view=graph-rest-beta)
- [qualityUpdateRing resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-qualityupdatering?view=graph-rest-beta)
- [Create ring](https://learn.microsoft.com/en-us/graph/api/windowsupdates-policy-post-rings?view=graph-rest-beta)
- [Get ring](https://learn.microsoft.com/en-us/graph/api/windowsupdates-ring-get?view=graph-rest-beta)
- [Update ring](https://learn.microsoft.com/en-us/graph/api/windowsupdates-ring-update?view=graph-rest-beta)
- [Delete ring](https://learn.microsoft.com/en-us/graph/api/windowsupdates-ring-delete?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `WindowsUpdates.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Dependency Notes

This resource requires a parent `microsoft365_graph_beta_windows_updates_autopatch_policy` resource. The `policy_id` attribute must reference an existing policy.

When setting `included_group_assignment` or `excluded_group_assignment`, the Microsoft Entra groups must already exist and have propagated before the ring is created. Use `time_sleep` with a 30-second `create_duration` and a `depends_on` referencing the groups, then set `depends_on = [time_sleep.wait_for_groups]` on the ring resource.

Both `included_group_assignment` and `excluded_group_assignment` are optional. When omitted, the provider sends an empty assignment list to the API.

## Field Mutability

| Field | Mutable after creation |
|-------|------------------------|
| `policy_id` | No — changes require resource replacement |
| `display_name` | Yes |
| `description` | Yes |
| `is_paused` | Yes |
| `deferral_in_days` | Yes |
| `included_group_assignment` | Yes |
| `excluded_group_assignment` | Yes |
| `is_hotpatch_enabled` | Yes |

## Import ID Format

This resource uses a composite import ID: `{policy_id}/{ring_id}`.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release |

## Example Usage

### Inclusion Assignments

```terraform
resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "example" {
  policy_id        = microsoft365_graph_beta_windows_updates_autopatch_policy.example.id
  display_name     = "Pilot Ring"
  description      = "Quality updates deployed to the pilot group after a 7-day deferral"
  is_paused        = false
  deferral_in_days = 7

  included_group_assignment = {
    assignments = [
      {
        group_id = "00000000-0000-0000-0000-000000000001"
      }
    ]
  }
}
```

### Exclusion Assignments

```terraform
resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "example" {
  policy_id        = microsoft365_graph_beta_windows_updates_autopatch_policy.example.id
  display_name     = "Broad Ring"
  description      = "Quality updates deployed to all devices, excluding the VIP group"
  is_paused        = false
  deferral_in_days = 14

  excluded_group_assignment = {
    assignments = [
      {
        group_id = "00000000-0000-0000-0000-000000000002"
      }
    ]
  }
}
```

### Maximal — Full Dependency Tree

```terraform
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Groups — included and excluded deployment audience
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "included_group" {
  display_name     = "wu-ring-included-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "wu-ring-included-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "excluded_group" {
  display_name     = "wu-ring-excluded-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "wu-ring-excluded-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

# ==============================================================================
# Wait for groups to propagate before assigning them to the ring
# ==============================================================================

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.included_group,
    microsoft365_graph_beta_groups_group.excluded_group,
  ]
  create_duration = "30s"
}

# ==============================================================================
# Parent autopatch policy
# ==============================================================================

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "example" {
  display_name = "Quality Update Policy - ${random_string.suffix.result}"
  description  = "Policy managing quality update rings"
}

# ==============================================================================
# Ring — full configuration with both included and excluded assignments,
# hotpatch enabled, deferral, and paused state
# ==============================================================================

resource "microsoft365_graph_beta_windows_updates_autopatch_ring" "example" {
  depends_on = [time_sleep.wait_for_groups]

  policy_id           = microsoft365_graph_beta_windows_updates_autopatch_policy.example.id
  display_name        = "Production Ring - ${random_string.suffix.result}"
  description         = "Quality updates with 14-day deferral for the production audience"
  is_paused           = false
  deferral_in_days    = 14
  is_hotpatch_enabled = false

  included_group_assignment = {
    assignments = [
      {
        group_id = microsoft365_graph_beta_groups_group.included_group.id
      }
    ]
  }

  excluded_group_assignment = {
    assignments = [
      {
        group_id = microsoft365_graph_beta_groups_group.excluded_group.id
      }
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) The ring description. The maximum length is 1,500 characters.
- `display_name` (String) The ring display name. The maximum length is 200 characters.
- `is_paused` (Boolean) Whether the ring is paused. When `true`, quality update deployment to devices in this ring is halted.
- `policy_id` (String) The ID of the Windows Update policy to which this ring belongs.

### Optional

- `deferral_in_days` (Number) The quality update deferral period in days. The value must be between 0 and 30.
- `excluded_group_assignment` (Attributes) Defines the Microsoft Entra groups whose devices are excluded from this ring's deployment audience. If not set, an empty assignment is sent. (see [below for nested schema](#nestedatt--excluded_group_assignment))
- `included_group_assignment` (Attributes) Defines the Microsoft Entra groups whose devices are included in this ring's deployment audience. If not set, an empty assignment is sent. (see [below for nested schema](#nestedatt--included_group_assignment))
- `is_hotpatch_enabled` (Boolean) Whether hotpatch updates are enabled for this ring (quality update rings only). Hotpatch updates apply without requiring a device restart.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time the ring was created (read-only).
- `id` (String) The unique identifier for the ring.
- `last_modified_date_time` (String) The date and time the ring was last modified (read-only).

<a id="nestedatt--excluded_group_assignment"></a>
### Nested Schema for `excluded_group_assignment`

Optional:

- `assignments` (Attributes Set) A set of group assignments governing the excluded deployment audience. (see [below for nested schema](#nestedatt--excluded_group_assignment--assignments))

<a id="nestedatt--excluded_group_assignment--assignments"></a>
### Nested Schema for `excluded_group_assignment.assignments`

Required:

- `group_id` (String) The Microsoft Entra group ID to exclude from this ring.



<a id="nestedatt--included_group_assignment"></a>
### Nested Schema for `included_group_assignment`

Optional:

- `assignments` (Attributes Set) A set of group assignments governing the included deployment audience. (see [below for nested schema](#nestedatt--included_group_assignment--assignments))

<a id="nestedatt--included_group_assignment--assignments"></a>
### Nested Schema for `included_group_assignment.assignments`

Required:

- `group_id` (String) The Microsoft Entra group ID to include in this ring.



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
# Import an existing ring using the format: {policy_id}/{ring_id}
terraform import microsoft365_graph_beta_windows_updates_autopatch_ring.example 00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111
```
