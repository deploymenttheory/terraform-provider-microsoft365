---
page_title: "microsoft365_utility_guid_list_sharder Resource - terraform-provider-microsoft365"
subcategory: "Utility"

description: |-
  Retrieves object IDs (GUIDs) from Microsoft Graph API and distributes them into configurable shards for progressive rollouts and phased deployments. Queries /users, /devices, /applications, or /groups/{id}/members endpoints with optional OData filtering, then applies sharding strategies (random, sequential, or percentage-based) to distribute results. Output shards are sets that can be directly used in conditional access policies, groups, and other resources requiring object ID collections.
  Unlike a datasource, this resource stores shard assignments in Terraform state. When recalculate_on_next_run = false, the stored assignments are returned unchanged on every plan and apply — preventing membership churn from causing reassignments. Set recalculate_on_next_run = true and run terraform apply to recompute shards from the current Graph API member list.
  API Endpoints: GET /users, GET /devices, GET /applications, GET /groups/{id}/members (with pagination and ConsistencyLevel: eventual header)
  Common Use Cases: MFA rollouts, Windows Update rings, conditional access pilots, application-based policies, group splitting, A/B testing
---

# microsoft365_utility_guid_list_sharder (Resource)

Retrieves object IDs (GUIDs) from Microsoft Graph API and distributes them into configurable shards for progressive rollouts and phased deployments. Queries `/users`, `/devices`, `/applications`, or `/groups/{id}/members` endpoints with optional OData filtering, then applies sharding strategies (random, sequential, or percentage-based) to distribute results. Output shards are sets that can be directly used in conditional access policies, groups, and other resources requiring object ID collections.

Unlike a datasource, this resource stores shard assignments in Terraform state. When `recalculate_on_next_run = false`, the stored assignments are returned unchanged on every plan and apply — preventing membership churn from causing reassignments. Set `recalculate_on_next_run = true` and run `terraform apply` to recompute shards from the current Graph API member list.

**API Endpoints:** `GET /users`, `GET /devices`, `GET /applications`, `GET /groups/{id}/members` (with pagination and `ConsistencyLevel: eventual` header)

**Common Use Cases:** MFA rollouts, Windows Update rings, conditional access pilots, application-based policies, group splitting, A/B testing

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required (at least one, depending on resource_type):**
- `User.Read.All` — when querying `users`
- `Device.Read.All` — when querying `devices`
- `Application.Read.All` — when querying `applications`
- `Directory.Read.All` — when querying `service_principals` or `group_members`
- `Group.Read.All` — when querying `group_members`

**Optional:**
- `User.ReadBasic.All` — sufficient for user ID queries without profile data

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.51.0-alpha | Experimental | Promoted from datasource to resource; added `recalculate_on_next_run` state locking |
| v0.42.0-alpha | Experimental | Initial release as datasource (round-robin, percentage, size, rendezvous strategies) |

## Background

When deploying policies, configurations, or security controls across large Microsoft 365 environments, immediate organisation-wide rollouts carry significant risk. If a misconfiguration occurs it can impact thousands of users or devices simultaneously.

The GUID List Sharder solves this by querying the Graph API for a population of users, devices, service principals, applications, or group members and distributing their object IDs into configurable subsets (shards) suitable for phased deployment:

1. **Reduces blast radius** — deploy to 10% first, validate, then expand
2. **Prevents pilot fatigue** — use a unique `seed` per rollout so the same users are not always in every pilot group
3. **Locks assignments in state** — with `recalculate_on_next_run = false`, new tenant members do not cause churn; existing policy targets remain stable
4. **Supports multiple strategies** — choose the algorithm that matches your deployment model

### Why a resource rather than a data source?

A data source recomputes on every plan, which means every `terraform plan` re-queries the Graph API and potentially reassigns users to different shards. For rollouts where membership stability matters (MFA enforcement, Windows Update rings, Conditional Access policies), this churn is undesirable and can cause live policies to target different users on each apply.

As a **resource**, shard assignments are stored in state. Set `recalculate_on_next_run = false` (the recommended default) and assignments are locked — the same population is targeted on every apply, regardless of tenant changes. Set it to `true` only when you intentionally want to rebalance.

## Distribution Strategies

### Round-Robin

Circular assignment: GUID 0 → shard_0, GUID 1 → shard_1, GUID 2 → shard_0, cycling. Guarantees perfect ±1 balance across all shards.

- **Without seed**: uses API return order (non-deterministic across runs)
- **With seed**: shuffles the input list deterministically before applying round-robin (reproducible)
- **Use when**: equal-sized rings are required; A/B testing; capacity planning

### Percentage

Distributes by the ratios in `shard_percentages`. Values must sum to 100. Shard sizes are arithmetically deterministic — the seed only controls *which* GUIDs land in each shard, not how many.

- **Use when**: stakeholders think in percentages (10% pilot → 30% broader → 60% full)

### Size

Distributes by absolute counts in `shard_sizes`. Use `-1` as the last value to capture all remaining GUIDs. Shard sizes are always exact as specified; the seed only varies membership.

- **Use when**: exact headcounts are required regardless of total population size; capping pilot ring size

### Rendezvous (Highest Random Weight)

Each GUID independently computes `SHA256(guid:shard_N:seed)` scores for every shard and is assigned to the shard with the highest score. This means adding a new shard only moves ~1/n GUIDs (the theoretical minimum), versus ~75% with position-based strategies.

- **Use when**: the ring count will change during the rollout lifecycle; minimising reassignment churn matters

## Key Concept: Distributing Pilot Burden

Use a different `seed` value for each independent rollout. This ensures the same users are not always the first to encounter issues:

| Rollout | Seed | User A's ring |
|---------|------|---------------|
| MFA enforcement | `"mfa-2026"` | Ring 0 (pilot) |
| Windows Updates | `"windows-2026"` | Ring 2 (broad) |
| Compliance baseline | `"compliance-2026"` | Ring 1 (validation) |

## Decision Matrix

| Requirement | Strategy | Seed | Notes |
|-------------|----------|------|-------|
| Equal-sized rings | `round-robin` | Optional | Perfect ±1 balance; add seed for reproducibility |
| Percentage waves | `percentage` | Optional | Clean ratio splits; add seed for deterministic membership |
| Absolute headcount caps | `size` | Optional | Exact counts; supports `-1` for "all remaining" |
| Ring count may grow | `rendezvous` | Optional | Only ~1/n GUIDs move when adding a shard |
| Vary pilot population per rollout | Any | **Yes** | Different seeds = different distributions |

## Example Usage

### Round-Robin (Equal Distribution)

```terraform
# ==============================================================================
# Example 1: Round-Robin Strategy
#
# Round-robin cycles through shards in order, guaranteeing perfect ±1 balance.
# It is the right choice when you need equal-sized rings and predictable sizes
# matter more than per-GUID determinism.
#
# Without a seed the distribution follows the API return order, which can change
# between runs as users are added or removed. With a seed the input list is
# shuffled deterministically before round-robin is applied, so the same tenant
# state always produces the same assignments — useful for A/B testing or when
# you need to communicate exact ring membership to stakeholders.
#
# recalculate_on_next_run = false (recommended default)
#   - On the very first apply, assignments are always computed regardless of
#     this value — you do not need a two-step "set true, then set false" dance.
#   - On all subsequent applies the stored assignments are returned from state
#     unchanged; no Graph API call is made, and new tenant members are ignored.
#   - Set to true only when you intentionally want to reshard (e.g. after a
#     large onboarding wave), then set back to false to re-lock.
# ==============================================================================

# Without seed — uses API return order; equal ring sizes but not reproducible
resource "microsoft365_utility_guid_list_sharder" "mfa_rollout_rings" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 4 # Ring 0 → Ring 3; equal ±1 split
  strategy                = "round-robin"
  recalculate_on_next_run = false
}

# With seed — deterministic; same tenant state always produces the same rings.
# Use a different seed per rollout (e.g. "windows-updates-2026" vs "mfa-2026")
# so that users who are in the pilot ring for one initiative are NOT always
# in the pilot ring for every other initiative (prevents pilot fatigue).
resource "microsoft365_utility_guid_list_sharder" "mfa_rollout_rings_seeded" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 4
  strategy                = "round-robin"
  seed                    = "mfa-rollout-2026" # Change per rollout to vary pilot ring population and reduce pilot fatigue
  recalculate_on_next_run = false
}

# Diagnostic outputs — useful during initial deployment for sanity checks
output "mfa_ring_distribution" {
  description = "Number of users assigned to each MFA rollout ring"
  value = {
    ring_0 = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_0"])
    ring_1 = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_1"])
    ring_2 = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_2"])
    ring_3 = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_3"])
    total  = length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_1"]) + length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_2"]) + length(microsoft365_utility_guid_list_sharder.mfa_rollout_rings_seeded.shards["shard_3"])
  }
}
```

### Percentage (Custom Ratios)

```terraform
# ==============================================================================
# Example 2: Percentage Strategy
#
# Percentage distributes GUIDs according to the ratios you specify. The sizes
# are arithmetically deterministic — [10, 30, 60] always produces shards of
# 10%, 30%, and 60% regardless of whether a seed is supplied. The seed only
# controls *which* users land in each shard, not how many.
#
# This is the most natural fit for standard phased rollouts where stakeholders
# think in percentages: "start with 10%, validate, then expand to 30%, then all".
#
# shard_percentages must sum to exactly 100.
# Each shard is named shard_0, shard_1, ... matching the list index.
# ==============================================================================

# Without seed — order follows the Graph API response; not reproducible across runs
resource "microsoft365_utility_guid_list_sharder" "ca_policy_phased" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages       = [10, 30, 60] # Pilot → Broader → Full rollout
  strategy                = "percentage"
  recalculate_on_next_run = false
}

# With seed — deterministic; pilot group is the same every run for this rollout.
# Using a rollout-specific seed means the 10% pilot for conditional access is
# drawn from a different slice of the tenant than the 10% pilot for MFA —
# distributing the pilot burden across the organisation.
resource "microsoft365_utility_guid_list_sharder" "ca_policy_phased_seeded" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages       = [10, 30, 60]
  strategy                = "percentage"
  seed                    = "ca-policy-rollout-2026"
  recalculate_on_next_run = false
}

output "ca_policy_phase_sizes" {
  description = "Headcount per deployment phase — verify ratios before enabling the policy"
  value = {
    pilot_10pct   = length(microsoft365_utility_guid_list_sharder.ca_policy_phased_seeded.shards["shard_0"])
    broader_30pct = length(microsoft365_utility_guid_list_sharder.ca_policy_phased_seeded.shards["shard_1"])
    full_60pct    = length(microsoft365_utility_guid_list_sharder.ca_policy_phased_seeded.shards["shard_2"])
  }
}
```

### Size (Absolute Counts)

```terraform
# ==============================================================================
# Example 3: Size Strategy
#
# Size distributes GUIDs by absolute counts rather than ratios. This is the
# right choice when stakeholders specify requirements in headcount rather than
# percentages: "we need exactly 50 users in the pilot, 200 in the broader wave,
# and everyone else in the final ring".
#
# Use -1 as the last value to mean "all remaining GUIDs". Only the last element
# may be -1. Without -1 the last shard is exactly the specified size and any
# remaining GUIDs are discarded — an intentional way to cap ring sizes.
#
# As with percentage, the seed controls membership assignment not shard size;
# sizes are always exactly as specified regardless of seed.
# ==============================================================================

# Absolute-size rings — last shard captures all remaining users via -1 sentinel
resource "microsoft365_utility_guid_list_sharder" "windows_update_rings" {
  resource_type           = "devices"
  odata_filter            = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_sizes             = [50, 200, -1] # Exactly 50 → 200 → all remaining
  strategy                = "size"
  seed                    = "windows-updates-2026"
  recalculate_on_next_run = false
}

# Capped rings — useful when a pilot must not exceed a fixed headcount.
# No -1 means users beyond the sum of sizes are deliberately excluded.
# The fourth shard would receive the overflow if added later.
resource "microsoft365_utility_guid_list_sharder" "it_pilot_capped" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and department eq 'IT'"
  shard_sizes             = [10, 25] # Hard cap: pilot=10, validation=25, rest not yet targeted
  strategy                = "size"
  seed                    = "it-dept-pilot-2026"
  recalculate_on_next_run = false
}

output "windows_update_ring_sizes" {
  description = "Device count per Windows Update ring (ring_2 = all remaining devices)"
  value = {
    ring_0_validation = length(microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_0"])
    ring_1_pilot      = length(microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_1"])
    ring_2_production = length(microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_2"])
  }
}
```

### Rendezvous (Stable When Ring Count Changes)

```terraform
# ==============================================================================
# Example 4: Rendezvous (Highest Random Weight) Strategy
#
# Rendezvous assigns each GUID independently using a per-(guid, shard, seed)
# hash score. The shard with the highest score wins for each GUID. Because
# each GUID's assignment is computed independently, adding or removing a shard
# only moves the GUIDs that must move: approximately 1/n of the population
# when going from n to n+1 shards, versus ~75% with position-based strategies.
#
# This stability makes rendezvous the right choice when:
#   - The ring count is expected to grow during the rollout lifecycle
#   - Minimising reassignment churn matters (e.g. avoiding repeated policy
#     application to devices that already received an earlier ring's config)
#   - You need per-GUID determinism without caring about equal-sized rings
#
# Distribution is probabilistic: with 12 users and 4 shards you expect ~3 per
# shard, but natural hash variance means actual counts differ. Assertions on
# per-shard size are not meaningful — assert only total count.
#
# The seed affects which shard each GUID is assigned to; an empty seed is used
# internally when omitted, which is still deterministic.
# ==============================================================================

# Starting with 3 rings — ring count may grow later without excessive churn
resource "microsoft365_utility_guid_list_sharder" "compliance_rings_initial" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 3
  strategy                = "rendezvous"
  seed                    = "compliance-baseline-2026"
  recalculate_on_next_run = false
}

# After expanding to 4 rings: only ~25% of users move to the new shard_3.
# The remaining ~75% stay exactly where they were — no unnecessary disruption.
# To expand: change shard_count to 4, set recalculate_on_next_run = true,
# apply, then set recalculate_on_next_run back to false.
resource "microsoft365_utility_guid_list_sharder" "compliance_rings_expanded" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 4 # One extra ring added; only ~25% of users move
  strategy                = "rendezvous"
  seed                    = "compliance-baseline-2026" # Same seed preserves prior assignments
  recalculate_on_next_run = false
}

output "compliance_ring_distribution" {
  description = "Approximate per-ring user counts (rendezvous is probabilistic — totals are authoritative)"
  value = {
    ring_0 = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_0"])
    ring_1 = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_1"])
    ring_2 = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_2"])
    ring_3 = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_3"])
    total  = length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_1"]) + length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_2"]) + length(microsoft365_utility_guid_list_sharder.compliance_rings_expanded.shards["shard_3"])
  }
}
```

### All Supported Resource Types

```terraform
# ==============================================================================
# Example 5: All Supported Resource Types
#
# The sharder can query five different Microsoft Graph collections. Each has
# distinct use cases and optional filtering behaviour.
#
# resource_type options:
#   "users"              → GET /users          — MFA, CA, phased policy rollouts
#   "devices"            → GET /devices        — Windows Update rings, compliance
#   "applications"       → GET /applications   — App registration distribution
#   "service_principals" → GET /servicePrincipals — App-based CA policies
#   "group_members"      → GET /groups/{id}/members — Split an existing group
#
# All types support odata_filter (except group_members which filters server-side
# via the group membership itself). group_members additionally requires group_id.
# ==============================================================================

# Users — active member accounts only (excludes guests and disabled accounts)
resource "microsoft365_utility_guid_list_sharder" "active_users" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 4
  strategy                = "round-robin"
  seed                    = "users-ring-2026"
  recalculate_on_next_run = false
}

# Devices — Azure AD joined Windows devices only (Intune-managed fleet)
resource "microsoft365_utility_guid_list_sharder" "windows_devices" {
  resource_type           = "devices"
  odata_filter            = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_percentages       = [5, 15, 30, 50]
  strategy                = "percentage"
  seed                    = "windows-rings-2026"
  recalculate_on_next_run = false
}

# Applications — all app registrations in the tenant
# Use this when distributing app registrations across policy rings is needed.
resource "microsoft365_utility_guid_list_sharder" "app_registrations" {
  resource_type           = "applications"
  shard_percentages       = [10, 90]
  strategy                = "percentage"
  seed                    = "app-reg-pilot-2026"
  recalculate_on_next_run = false
}

# Service Principals — enterprise app instances used in CA application conditions.
# Filter to Microsoft-published apps (common for targeting M365 workloads).
resource "microsoft365_utility_guid_list_sharder" "microsoft_enterprise_apps" {
  resource_type           = "service_principals"
  odata_filter            = "startswith(displayName, 'Microsoft')"
  shard_count             = 3
  strategy                = "round-robin"
  seed                    = "sp-mfa-2026"
  recalculate_on_next_run = false
}

# Service Principals — all enterprise apps; no filter gives the complete tenant list
resource "microsoft365_utility_guid_list_sharder" "all_enterprise_apps" {
  resource_type           = "service_principals"
  shard_percentages       = [10, 30, 60]
  strategy                = "percentage"
  seed                    = "sp-all-2026"
  recalculate_on_next_run = false
}

# Group Members — split an existing group's membership into sub-rings.
# Useful when you have an established Entra ID group (e.g. "All IT Staff") and
# want to phase a new policy across that group without creating static sub-groups.
# group_id is required when resource_type = "group_members".
resource "microsoft365_utility_guid_list_sharder" "it_dept_rings" {
  resource_type           = "group_members"
  group_id                = "12345678-1234-1234-1234-123456789abc" # Replace with real group object ID
  shard_count             = 3
  strategy                = "round-robin"
  seed                    = "it-dept-2026"
  recalculate_on_next_run = false
}
```

### Managing the recalculate_on_next_run Lifecycle

```terraform
# ==============================================================================
# Example 6: Managing the recalculate_on_next_run Lifecycle
#
# This example demonstrates the recommended workflow for locking and unlocking
# shard assignments over the lifetime of a rollout.
#
# Day 1 — Initial deployment
#   Set recalculate_on_next_run = false. The first terraform apply always
#   computes assignments from scratch regardless of this flag (there is no
#   prior state). Assignments are then locked.
#
# Ongoing — Steady state
#   Leave recalculate_on_next_run = false. Plans are fast (no API call), and
#   new users joining the tenant do not cause membership churn. The policy
#   targets exactly the same population on every apply.
#
# Intentional reshard — e.g. after a large onboarding wave or shard_count change
#   1. Set recalculate_on_next_run = true
#   2. Run terraform apply — Update re-queries Graph and recomputes shards
#   3. Set recalculate_on_next_run = false again to re-lock
#
# The distinction between Read and Update:
#   recalculate = false + no config change → Read returns cached state (no API call)
#   recalculate = false + config change    → Update preserves existing shards, saves new config
#   recalculate = true  + any apply        → Update re-queries Graph and reshards
# ==============================================================================

# Phase 1: Initial deployment — flag=false, assignments computed on first apply
resource "microsoft365_utility_guid_list_sharder" "compliance_rollout" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_count             = 3
  strategy                = "round-robin"
  seed                    = "compliance-2026"
  recalculate_on_next_run = false # Safe from day one; first apply always computes
}

# Phase 2 (illustrative — replace the above block in practice):
# After a large onboarding wave you decide to rebalance.
# Step 1: uncomment recalculate = true and apply
# Step 2: revert to false and apply again to lock
#
# resource "microsoft365_utility_guid_list_sharder" "compliance_rollout" {
#   resource_type           = "users"
#   odata_filter            = "accountEnabled eq true and userType eq 'Member'"
#   shard_count             = 3
#   strategy                = "round-robin"
#   seed                    = "compliance-2026"
#   recalculate_on_next_run = true    # Triggers reshard on next apply; set back to false after
# }

output "compliance_rollout_status" {
  description = "Ring membership counts and whether assignments are currently locked"
  value = {
    ring_0_count       = length(microsoft365_utility_guid_list_sharder.compliance_rollout.shards["shard_0"])
    ring_1_count       = length(microsoft365_utility_guid_list_sharder.compliance_rollout.shards["shard_1"])
    ring_2_count       = length(microsoft365_utility_guid_list_sharder.compliance_rollout.shards["shard_2"])
    assignments_locked = !microsoft365_utility_guid_list_sharder.compliance_rollout.recalculate_on_next_run
  }
}
```

### Scenario: Sharder → Security Groups → Group Member Assignments

```terraform
# ==============================================================================
# Scenario: Sharder → Security Groups → Group Member Assignments
#
# This is the primary integration pattern. Shard GUIDs from the sharder are
# distributed into separate Entra ID security groups. Downstream resources
# (Intune policies, Conditional Access, Windows Update rings, etc.) then target
# those groups rather than individual users or devices, following the standard
# Microsoft 365 group-based targeting model.
#
# Pattern:
#   sharder → creates shard sets of GUIDs
#   group   → one group per ring (created once, stable IDs)
#   group_member_assignment (for_each) → adds each shard's GUIDs to the group
#
# The for_each on group_member_assignment means each GUID gets its own resource
# instance tracked in state. Terraform will add/remove only the specific members
# that change — it will NOT recreate the group or remove members from the wrong ring.
#
# Use case: Phased MFA rollout across 4 rings, 10/20/30/40% of enabled members
# ==============================================================================

resource "microsoft365_utility_guid_list_sharder" "mfa_rings" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages       = [10, 20, 30, 40]
  strategy                = "percentage"
  seed                    = "mfa-2026"
  recalculate_on_next_run = false
}

# One security group per ring — created independently of the sharder so that
# their IDs remain stable even if shards are later recomputed.
resource "microsoft365_graph_beta_groups_group" "mfa_ring_0" {
  display_name     = "MFA Rollout - Ring 0 (10% Pilot)"
  mail_nickname    = "mfa-ring-0-pilot"
  description      = "Initial 10% pilot for MFA enforcement — validate before expanding"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "mfa_ring_1" {
  display_name     = "MFA Rollout - Ring 1 (20% Broader)"
  mail_nickname    = "mfa-ring-1-broader"
  description      = "Second wave — 20% of users after Ring 0 validation"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "mfa_ring_2" {
  display_name     = "MFA Rollout - Ring 2 (30% Broad)"
  mail_nickname    = "mfa-ring-2-broad"
  description      = "Third wave — 30% of users"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "mfa_ring_3" {
  display_name     = "MFA Rollout - Ring 3 (40% Production)"
  mail_nickname    = "mfa-ring-3-production"
  description      = "Final wave — remaining 40% of users"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

# Populate each group using for_each so Terraform tracks each membership
# individually. Adding or removing users only touches the affected member
# resource instances — the groups themselves are never recreated.
resource "microsoft365_graph_beta_groups_group_member_assignment" "mfa_ring_0_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"]

  group_id           = microsoft365_graph_beta_groups_group.mfa_ring_0.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "mfa_ring_1_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"]

  group_id           = microsoft365_graph_beta_groups_group.mfa_ring_1.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "mfa_ring_2_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"]

  group_id           = microsoft365_graph_beta_groups_group.mfa_ring_2.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "mfa_ring_3_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"]

  group_id           = microsoft365_graph_beta_groups_group.mfa_ring_3.id
  member_id          = each.value
  member_object_type = "User"
}

output "mfa_group_ids" {
  description = "Security group object IDs — reference these in Intune or CA policies"
  value = {
    ring_0_pilot      = microsoft365_graph_beta_groups_group.mfa_ring_0.id
    ring_1_broader    = microsoft365_graph_beta_groups_group.mfa_ring_1.id
    ring_2_broad      = microsoft365_graph_beta_groups_group.mfa_ring_2.id
    ring_3_production = microsoft365_graph_beta_groups_group.mfa_ring_3.id
  }
}

output "mfa_ring_headcount" {
  description = "Number of users in each ring — locked from the point of first apply"
  value = {
    ring_0 = length(microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"])
    ring_1 = length(microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"])
    ring_2 = length(microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"])
    ring_3 = length(microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"])
  }
}
```

### Scenario: Sharder → Groups → Phased Conditional Access Policy

```terraform
# ==============================================================================
# Scenario: Sharder → Security Groups → Conditional Access Policy (Phased)
#
# This example demonstrates how to wire sharder outputs directly into a phased
# Conditional Access rollout. The pattern is:
#
#   sharder  → shard GUIDs (users)
#   groups   → one security group per ring (stable IDs, group-based targeting)
#   members  → populate each group from the corresponding shard
#   CA policy per ring → target the ring's group; start in report-only mode
#
# Why one CA policy per ring rather than one policy with all groups?
#   Each ring can be independently promoted from report-only → enabled without
#   touching the other rings. Ring 0 can be fully enforced while Ring 1 is still
#   in observation mode. This also makes rollback surgical — disable Ring 1's
#   policy without any impact on Ring 0 or Ring 2.
#
# Use case: Phased Require MFA rollout — 10% pilot, then 30%, then 60%
# ==============================================================================

resource "microsoft365_utility_guid_list_sharder" "mfa_ca_rings" {
  resource_type           = "users"
  odata_filter            = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages       = [10, 30, 60]
  strategy                = "percentage"
  seed                    = "ca-mfa-2026" # Unique seed — different population than other rollouts
  recalculate_on_next_run = false
}

# ─────────────────────────────────────────────────────────────────────────────
# Security Groups (one per ring)
# ─────────────────────────────────────────────────────────────────────────────

resource "microsoft365_graph_beta_groups_group" "ca_mfa_ring_0" {
  display_name     = "CA - Require MFA - Ring 0 (10% Pilot)"
  mail_nickname    = "ca-mfa-ring-0"
  description      = "Pilot cohort for Require MFA CA policy — 10% of members"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "ca_mfa_ring_1" {
  display_name     = "CA - Require MFA - Ring 1 (30% Broader)"
  mail_nickname    = "ca-mfa-ring-1"
  description      = "Second wave for Require MFA CA policy — 30% of members"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "ca_mfa_ring_2" {
  display_name     = "CA - Require MFA - Ring 2 (60% Full)"
  mail_nickname    = "ca-mfa-ring-2"
  description      = "Final wave for Require MFA CA policy — remaining 60% of members"
  security_enabled = true
  mail_enabled     = false
  hard_delete      = true
}

# Break-glass account group — always excluded from every CA policy to ensure
# emergency access is never blocked by a misconfigured policy.
# This group must be created separately and pre-populated with break-glass accounts.
data "microsoft365_graph_beta_groups_group" "breakglass" {
  display_name = "Break Glass - Emergency Access"
}

# ─────────────────────────────────────────────────────────────────────────────
# Group Member Assignments — populate each ring group from sharder output
# ─────────────────────────────────────────────────────────────────────────────

resource "microsoft365_graph_beta_groups_group_member_assignment" "ca_mfa_ring_0_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_0"]

  group_id           = microsoft365_graph_beta_groups_group.ca_mfa_ring_0.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "ca_mfa_ring_1_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_1"]

  group_id           = microsoft365_graph_beta_groups_group.ca_mfa_ring_1.id
  member_id          = each.value
  member_object_type = "User"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "ca_mfa_ring_2_members" {
  for_each = microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_2"]

  group_id           = microsoft365_graph_beta_groups_group.ca_mfa_ring_2.id
  member_id          = each.value
  member_object_type = "User"
}

# ─────────────────────────────────────────────────────────────────────────────
# Conditional Access Policies — one per ring, independently promotable
#
# Lifecycle:
#   1. Deploy all three policies in "enabledForReportingButNotEnforced"
#   2. Review Ring 0 sign-in logs for 1–2 weeks; no unexpected blocks → proceed
#   3. Change Ring 0 state to "enabled"; leave Ring 1 in report-only
#   4. Repeat for Ring 1 (1–2 weeks observation), then Ring 2
#   5. Once Ring 2 is enabled, all members are covered; rings can be consolidated
# ─────────────────────────────────────────────────────────────────────────────

# Ring 0 — Pilot (10%): report-only until validated
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "require_mfa_ring_0" {
  display_name = "CA - Require MFA - Ring 0 (10% Pilot)"
  state        = "enabledForReportingButNotEnforced" # Promote to "enabled" after validation

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.ca_mfa_ring_0.id]
      exclude_groups = [data.microsoft365_graph_beta_groups_group.breakglass.id]
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

# Ring 1 — Broader (30%): disabled until Ring 0 is fully enforced
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "require_mfa_ring_1" {
  display_name = "CA - Require MFA - Ring 1 (30% Broader)"
  state        = "disabled" # Enable after Ring 0 reaches "enabled" and is stable

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.ca_mfa_ring_1.id]
      exclude_groups = [data.microsoft365_graph_beta_groups_group.breakglass.id]
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

# Ring 2 — Full (60%): disabled until Ring 1 is fully enforced
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "require_mfa_ring_2" {
  display_name = "CA - Require MFA - Ring 2 (60% Full)"
  state        = "disabled" # Enable after Ring 1 reaches "enabled" and is stable

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.ca_mfa_ring_2.id]
      exclude_groups = [data.microsoft365_graph_beta_groups_group.breakglass.id]
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# Outputs
# ─────────────────────────────────────────────────────────────────────────────

output "ca_mfa_rollout_summary" {
  description = "CA policy IDs and ring population counts — share with security team before promoting states"
  value = {
    ring_0 = {
      policy_id = microsoft365_graph_beta_identity_and_access_conditional_access_policy.require_mfa_ring_0.id
      group_id  = microsoft365_graph_beta_groups_group.ca_mfa_ring_0.id
      headcount = length(microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_0"])
      state     = "enabledForReportingButNotEnforced"
    }
    ring_1 = {
      policy_id = microsoft365_graph_beta_identity_and_access_conditional_access_policy.require_mfa_ring_1.id
      group_id  = microsoft365_graph_beta_groups_group.ca_mfa_ring_1.id
      headcount = length(microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_1"])
      state     = "disabled"
    }
    ring_2 = {
      policy_id = microsoft365_graph_beta_identity_and_access_conditional_access_policy.require_mfa_ring_2.id
      group_id  = microsoft365_graph_beta_groups_group.ca_mfa_ring_2.id
      headcount = length(microsoft365_utility_guid_list_sharder.mfa_ca_rings.shards["shard_2"])
      state     = "disabled"
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `recalculate_on_next_run` (Boolean) Controls whether shard assignments are recomputed during Terraform plan/refresh and on `terraform apply` when configuration changes.

**`false` (recommended default):** Shard assignments are locked in state. No Graph API call is made during plan or apply, and membership changes in your tenant (users added, removed) do not cause reassignments. On the very first apply, when no prior state exists, assignments are always computed regardless of this value — so you can safely set `false` from the outset without a two-step toggle.

**`true`:** Re-queries the Graph API and reruns the sharding algorithm on every plan refresh and every apply. Use this only when you explicitly want to rebalance — for example after a large onboarding wave, a policy restructure, or a change to `shard_count` or `strategy`.

**Recommended workflow:** Set to `false` from day one. Initial assignments are computed automatically on the first apply. Change to `true` (and run `terraform apply`) only when you intentionally want to reshard. Set back to `false` afterwards to re-lock assignments.
- `resource_type` (String) The type of Microsoft Graph resource to query and shard. `users` queries `/users` for user-based policies (MFA, conditional access). `devices` queries `/devices` for device policies (Windows Updates, compliance). `applications` queries `/applications` for app registrations. `service_principals` queries `/servicePrincipals` (enterprise apps) for application-based conditional access policies. `group_members` queries `/groups/{id}/members` to split existing group membership (requires `group_id`).
- `strategy` (String) The distribution strategy for sharding GUIDs. `round-robin` distributes in circular order (guarantees equal sizes, optional seed for reproducibility). `percentage` distributes by specified percentages (requires `shard_percentages`, optional seed for reproducibility). `size` distributes by absolute sizes (requires `shard_sizes`, optional seed for reproducibility). `rendezvous` uses Highest Random Weight algorithm (always deterministic, minimal disruption when shard count changes).

### Optional

- `group_id` (String) The object ID of the group to query members from. Required when `resource_type = "group_members"`, ignored otherwise. Use this to split an existing group's membership into multiple new groups for targeted policy application.
- `odata_filter` (String) Optional OData filter expression applied at the API level before sharding. **Users:** `accountEnabled eq true` (active accounts only), `userType eq 'Member'` (exclude guests). **Devices:** `operatingSystem eq 'Windows'` (Windows devices only). **Service Principals:** `startswith(displayName, 'Microsoft')` (Microsoft apps), `appId eq 'guid'` (specific app). Leave empty to query all resources without filtering.
- `seed` (String) Optional seed value for deterministic distribution. When provided, makes results reproducible across Terraform runs for the same input set. **`round-robin`**: No seed = uses API order (may change). With seed = shuffles deterministically first, then round-robin. **`percentage`/`size`**: Same shuffle behaviour as round-robin. **`rendezvous`**: Always deterministic. Seed affects which shard each GUID is assigned to. Use different seeds for different rollouts to vary pilot burden distribution.
- `shard_count` (Number) Number of equally-sized shards to create (minimum 1). Use with `round-robin` strategy. Conflicts with `shard_percentages` and `shard_sizes`. Creates shards named `shard_0`, `shard_1`, ..., `shard_N-1`.
- `shard_percentages` (List of Number) List of percentages for custom-sized shards. Use with `percentage` strategy. Conflicts with `shard_count` and `shard_sizes`. Values must be non-negative integers that sum to exactly 100. Example: `[10, 30, 60]` creates 10% pilot, 30% broader pilot, 60% full rollout.
- `shard_sizes` (List of Number) List of absolute shard sizes (exact number of GUIDs per shard). Use with `size` strategy. Conflicts with `shard_count` and `shard_percentages`. Values must be positive integers or -1 (which means 'all remaining'). Only the last element can be -1. Example: `[50, 200, -1]` creates 50 pilot users, 200 broader rollout, remainder for full deployment.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `shards` (Map of Set of String) Computed map of shard names (`shard_0`, `shard_1`, ...) to sets of GUIDs. Each value is a `set(string)` type, directly compatible with resource attributes expecting object ID sets. Access with `resource.example.shards["shard_0"]`, check size with `length(resource.example.shards["shard_0"])`.

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
# Import an existing guid_list_sharder resource using its ID
terraform import microsoft365_utility_guid_list_sharder.example <resource-id>
```

## Additional Resources

- [Microsoft Graph API - Users](https://learn.microsoft.com/en-us/graph/api/user-list)
- [Microsoft Graph API - Devices](https://learn.microsoft.com/en-us/graph/api/device-list)
- [Microsoft Graph API - Service Principals](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-list)
- [Microsoft Graph API - Group Members](https://learn.microsoft.com/en-us/graph/api/group-list-members)
- [Conditional Access Best Practices](https://learn.microsoft.com/en-us/entra/identity/conditional-access/plan-conditional-access)
- [Phased Deployment Strategies](https://learn.microsoft.com/en-us/mem/intune/fundamentals/deployment-guide-intune-setup)
