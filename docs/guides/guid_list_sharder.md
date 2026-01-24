---
page_title: "GUID List Sharder"
subcategory: "Guides"
description: |-
  Complete guide for implementing phased deployments and progressive rollouts using the GUID List Sharder utility data source.
---

# Progressive Rollout with GUID List Sharder

This guide demonstrates how to implement progressive rollouts, phased deployments, and pilot programs for Microsoft 365 policies using the `microsoft365_utility_guid_list_sharder` data source.

## Overview

The GUID List Sharder data source (`microsoft365_utility_guid_list_sharder`) is a utility tool that queries Microsoft Graph API to retrieve collections of object IDs (GUIDs) for users, devices, or group members, then intelligently distributes them into configurable "shards" (subsets) for progressive deployment strategies.

### What Problem Does This Solve?

When deploying policies, configurations, or security controls across large Microsoft 365 environments, immediate organization-wide rollouts carry significant risk. If a policy misconfiguration or unexpected behavior occurs, it can impact thousands of users or devices simultaneously.

Consequently, it's common for organizations to define deployment rings or waves to manage the rollout of changes. However, this approach presents its own challenges: Entra ID dynamic group membership can take hours to fully populate, creating deployment groups that adequately represent hardware and software diversity across the organization is complex, created groups become stale as people join, leave, or change roles, and willing pilot users experience pilot fatigue when they're consistently the first to encounter issues.


The GUID List Sharder solves this by:

1. **Enabling Progressive Rollouts**: Deploy changes to small pilot groups first (e.g., 10% of users), validate functionality, then gradually expand to broader populations
2. **Reducing Blast Radius**: Limit the impact of potential issues by controlling which users/devices receive changes at each phase
3. **Facilitating Validation**: Allow time to monitor, test, and validate each phase before proceeding to the next
4. **Distributing Pilot Burden**: By using unique seed values across different rollouts, the same users won't always be early adopters. User X might be in the 10% pilot for MFA rollout (seed: "mfa-2024") but in the 60% final wave for Windows Updates (seed: "windows-2024"), preventing "pilot fatigue" where certain users consistently experience issues first
5. **Supports Multiple Deployment Strategies**: Choose between `round-robin` (perfect equal distribution), `percentage` (custom ratios), `size` (absolute counts), or `rendezvous` (minimal disruption when ring counts change). Most strategies support optional seed for reproducibility

### How It Works

The data source operates in three stages:

1. **Query**: Retrieves object IDs from Microsoft Graph endpoints (`/users`, `/devices`, or `/groups/{id}/members`) with optional OData filtering to narrow the population
2. **Shard**: Applies one of four distribution strategies to divide the population into subsets:
   - **Round-robin**: Circular distribution for perfect equal sizes (±1 GUID, optional seed for reproducibility)
   - **Percentage**: Custom-weighted distribution by ratios (e.g., 10% pilot, 30% broader, 60% full, optional seed for reproducibility)
   - **Size**: Absolute count-based distribution (e.g., exactly 50, 100, 200 users per shard, optional seed for reproducibility)
   - **Rendezvous**: Highest Random Weight (HRW) algorithm for minimal disruption when shard counts change (always deterministic, seed affects distribution pattern)
3. **Output**: Returns a map of shards containing sets of object IDs, directly compatible with Terraform resources like conditional access policies and groups

### Key Benefits

- **Type-Safe Integration**: Outputs are `set(string)` types that work directly with Microsoft 365 resources without type conversion
- **Flexible Sizing**: Supports equal distribution (`shard_count`), custom percentages (`shard_percentages`), absolute counts (`shard_sizes`), and stable assignment across ring changes (`rendezvous` strategy)
- **Deterministic Options**: Most strategies support optional seed for consistent shard assignment across multiple Terraform runs
- **Minimal Disruption**: Rendezvous strategy ensures only ~1/n GUIDs move when adding or removing deployment rings
- **No Manual Management**: Automatically handles pagination and eventual consistency for large populations
- **Production-Ready**: Supports OData filtering to exclude disabled accounts, target specific attributes, or align with existing policy filters

### Example Use Cases

- **MFA Rollout**: Deploy multi-factor authentication to small pilots before organization-wide deployment
- **Windows Update Rings**: Distribute devices across early adopter, validation, and broad deployment rings
- **Conditional Access Pilots**: Test new access policies with controlled user populations
- **Group Splitting**: Divide large groups into manageable subgroups for targeted policies
- **A/B Testing**: Distribute users across test groups for policy or feature testing

## API Endpoints Used

The data source queries the following Microsoft Graph Beta endpoints:

- `GET /users` - When `resource_type = "users"`
- `GET /devices` - When `resource_type = "devices"`
- `GET /groups/{id}/members` - When `resource_type = "group_members"`

All queries support:
- Pagination (automatic)
- OData filtering via `$filter` parameter
- Eventual consistency header (`ConsistencyLevel: eventual`)

## Prerequisites

- Terraform 1.14.0 or later
- Microsoft 365 provider configured with appropriate credentials
- Required permissions:
  - `User.Read.All` or `Directory.Read.All` - For users
  - `Device.Read.All` or `Directory.Read.All` - For devices
  - `Group.Read.All` or `GroupMember.Read.All` - For group members

## Choosing the Right Strategy

Use this decision matrix to select the optimal strategy and seed combination for your use case:

### Decision Matrix

| Your Need | Strategy | Seed | Why |
|-----------|----------|------|-----|
| **Exact equal sizes** (e.g., 25% in each of 4 rings) | `round-robin` | 🟡 Optional | Guarantees perfect balance (±1 GUID). Add seed for reproducibility across runs |
| **Custom percentages** (e.g., 10/30/60 split) | `percentage` | 🟡 Optional | Flexible ratios. Add seed to ensure same users in same phases when recreating |
| **Absolute counts** (e.g., exactly 50, 100, 200 users) | `size` | 🟡 Optional | Precise shard sizes. Add seed for deterministic assignment |
| **Ring count will change** (e.g., start with 3 rings, later expand to 4) | `rendezvous` | ✅ Required | Only ~25% of users move when adding rings (vs ~75% with other strategies) |
| **Distribute pilot burden** (different users in pilot for different rollouts) | Any strategy | ✅ Yes | Different seeds (e.g., "mfa-2024" vs "windows-2024") = different distributions |
| **Capacity testing / A/B testing** | `round-robin` | ✅ Yes | Perfect balance + reproducible results |
| **One-time split, don't care about reproducibility** | `round-robin` or `percentage` | ❌ None | Fastest - no seed overhead |

## Distribution Strategies

### Round-Robin Strategy

**Characteristics:**
- Distributes GUIDs in circular order (GUID 0→shard 0, GUID 1→shard 1, ..., GUID n→shard 0, ...)
- Guarantees perfect equal distribution (within ±1 GUID for indivisible populations)
- Without seed: Uses API order (non-deterministic across runs)
- With seed: Shuffles deterministically first using Fisher-Yates, then applies round-robin (reproducible)
- Fastest strategy for equal distribution

**When to Use:**
- Need guaranteed equal shard sizes (e.g., 25% in each of 4 deployment rings)
- A/B testing or capacity planning requiring balanced groups
- Quick one-time splits where reproducibility isn't required (omit seed)
- Statistical sampling from large populations

**Example:**

```hcl
# Equal distribution across 4 deployment rings
data "microsoft365_utility_guid_list_sharder" "deployment_rings" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and userType eq 'Member'"
  shard_count   = 4
  strategy      = "round-robin"
  seed          = "mfa-rollout-2024"  # Optional: omit for one-time split
}

# Use shards in conditional access policies
resource "microsoft365_graph_beta_conditional_access_policy" "ring_0_pilot" {
  display_name = "MFA Required - Ring 0 (Pilot)"
  state        = "enabled"
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.deployment_rings.shards["shard_0"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Monitor distribution
output "ring_sizes" {
  value = {
    ring_0 = length(data.microsoft365_utility_guid_list_sharder.deployment_rings.shards["shard_0"])
    ring_1 = length(data.microsoft365_utility_guid_list_sharder.deployment_rings.shards["shard_1"])
    ring_2 = length(data.microsoft365_utility_guid_list_sharder.deployment_rings.shards["shard_2"])
    ring_3 = length(data.microsoft365_utility_guid_list_sharder.deployment_rings.shards["shard_3"])
  }
}
```

---

### Percentage Strategy

**Characteristics:**
- Distributes GUIDs by specified percentage ratios
- Fills shards sequentially (shard_0 gets first N%, shard_1 gets next N%, etc.)
- Last shard receives all remaining GUIDs to ensure nothing is lost
- Balance guarantee: Within rounding error of specified percentages (exact for large populations)
- Without seed: Uses API order (non-deterministic across runs)
- With seed: Shuffles deterministically first, then applies percentage split (reproducible)

**When to Use:**
- Following structured rollout phases with specific percentage requirements (10% pilot → 30% broader → 60% full)
- Change management processes requiring documented pilot sizes
- Different phases have different risk profiles or validation requirements
- Windows Update rings following industry-standard distributions (5% → 15% → 80%)

**Example:**

```hcl
# Phased MFA rollout: 10% pilot, 30% broader, 60% full
data "microsoft365_utility_guid_list_sharder" "mfa_phases" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-phases-2024"  # Optional: ensures same users in same phases
}

# Phase 1: Pilot (10%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_pilot" {
  display_name = "MFA Required - Phase 1 Pilot"
  state        = "enabled"
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_phases.shards["shard_0"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Phase 2: Broader (30%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_broader" {
  display_name = "MFA Required - Phase 2 Broader"
  state        = "enabledForReportingButNotEnforced"  # Start in report-only
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_phases.shards["shard_1"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Monitor phase distribution
output "phase_distribution" {
  value = {
    pilot_count   = length(data.microsoft365_utility_guid_list_sharder.mfa_phases.shards["shard_0"])
    broader_count = length(data.microsoft365_utility_guid_list_sharder.mfa_phases.shards["shard_1"])
    full_count    = length(data.microsoft365_utility_guid_list_sharder.mfa_phases.shards["shard_2"])
  }
}
```

---

### Size Strategy

**Characteristics:**
- Distributes GUIDs by absolute counts per shard
- Fills shards sequentially (shard_0 gets first N items, shard_1 gets next M items, etc.)
- Last shard receives all remaining GUIDs to ensure nothing is lost
- Balance guarantee: Exact counts as specified (if population is sufficient)
- Without seed: Uses API order (non-deterministic across runs)
- With seed: Shuffles deterministically first, then applies size-based split (reproducible)

**When to Use:**
- Need precise shard sizes regardless of total population (e.g., exactly 50 pilot users, 100 validation users)
- Capacity planning with fixed-size test groups
- Compliance requirements specify exact pilot group sizes
- Resource constraints limit pilot group sizes (e.g., only 25 licenses available for testing)

**Example:**

```hcl
# Fixed-size pilot groups: 50 pilot, 100 validation, remainder broad
data "microsoft365_utility_guid_list_sharder" "fixed_pilot" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and department eq 'Sales'"
  shard_sizes   = [50, 100]  # Last shard gets all remaining
  strategy      = "size"
  seed          = "sales-pilot-2024"  # Optional: ensures consistent pilot membership
}

# Pilot group: exactly 50 users
resource "microsoft365_graph_beta_group" "pilot_group" {
  display_name     = "Sales - New Policy Pilot (50 users)"
  security_enabled = true
  
  members = data.microsoft365_utility_guid_list_sharder.fixed_pilot.shards["shard_0"]
}

# Validation group: exactly 100 users
resource "microsoft365_graph_beta_group" "validation_group" {
  display_name     = "Sales - New Policy Validation (100 users)"
  security_enabled = true
  
  members = data.microsoft365_utility_guid_list_sharder.fixed_pilot.shards["shard_1"]
}

# Verify exact counts
output "pilot_sizes" {
  value = {
    pilot_count      = length(data.microsoft365_utility_guid_list_sharder.fixed_pilot.shards["shard_0"])
    validation_count = length(data.microsoft365_utility_guid_list_sharder.fixed_pilot.shards["shard_1"])
    broad_count      = length(data.microsoft365_utility_guid_list_sharder.fixed_pilot.shards["shard_2"])
  }
}
```

---

### Rendezvous Strategy

**Characteristics:**
- Uses Highest Random Weight (HRW) algorithm - each GUID independently evaluates all shards via `hash(guid:shard_N:seed)` and selects highest score
- Always deterministic (reproducible across runs) - seed required, affects which shard wins for each GUID
- Balance guarantee: Probabilistic ~equal distribution (typically within 3% for 1000+ GUIDs, not perfect ±1 like round-robin)
- **Stability guarantee**: When shard count changes, only ~1/n GUIDs move (theoretical minimum). Adding 4th shard to 3-shard setup: only ~25% of GUIDs reassign vs ~75% with position-based strategies
- Per-GUID independence: Each GUID's assignment is computed independently, not affected by position in list

**When to Use:**
- Ring count will change during rollout lifecycle (e.g., start with 3 rings, later add 4th for extended validation)
- Minimizing user disruption is critical (users shouldn't suddenly move from broad deployment back to pilot)
- Long-term deployments where ring structure may evolve
- Need stable assignment that survives infrastructure changes

**Example:**

```hcl
# Initial deployment: 3 rings
data "microsoft365_utility_guid_list_sharder" "stable_rings" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "mfa-stable-2024"  # Required for rendezvous
}

# Ring 0: Early adopters
resource "microsoft365_graph_beta_conditional_access_policy" "ring_0" {
  display_name = "MFA Required - Ring 0 (Early Adopters)"
  state        = "enabled"
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.stable_rings.shards["shard_0"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Later: Change shard_count from 3 to 4 to add extended validation ring
# Result: Only ~25% of users move (those assigned to new ring_3)
# Ring 0, 1, 2 membership remains ~75% stable

# Monitor ring distribution
output "ring_stability" {
  value = {
    ring_0_count = length(data.microsoft365_utility_guid_list_sharder.stable_rings.shards["shard_0"])
    ring_1_count = length(data.microsoft365_utility_guid_list_sharder.stable_rings.shards["shard_1"])
    ring_2_count = length(data.microsoft365_utility_guid_list_sharder.stable_rings.shards["shard_2"])
    # Uncomment when adding 4th ring:
    # ring_3_count = length(data.microsoft365_utility_guid_list_sharder.stable_rings.shards["shard_3"])
  }
}
```

---

## Seed Usage Patterns

Seeds control determinism and distribution variability across rollouts.

### Pattern: No Seed (Fast, One-Time)

**Behavior:**
- Round-robin & Percentage: Uses API order (non-deterministic - may change between runs)
- Rendezvous: N/A (seed is required)

**Use When:**
- Quick one-time split where reproducibility isn't needed
- Speed is priority (no shuffle overhead)
- Distribution will never be recreated

**Example:**

```hcl
# Quick equal split, no reproducibility needed
data "microsoft365_utility_guid_list_sharder" "quick_split" {
  resource_type = "users"
  shard_count   = 4
  strategy      = "round-robin"
  # No seed - fastest, uses API order
}
```

---

### Pattern: Seed for Reproducibility

**Behavior:**
- Same seed produces identical distribution across multiple Terraform runs
- Round-robin & Percentage: Shuffles deterministically first, then applies strategy
- Rendezvous: Always deterministic, seed affects distribution pattern

**Use When:**
- Need to recreate exact same groups in future runs (debugging, compliance, audit)
- Long-term rollouts where phase membership shouldn't change
- Testing requires consistent groups across infrastructure rebuilds

**Example:**

```hcl
# Reproducible phased rollout
data "microsoft365_utility_guid_list_sharder" "production_phases" {
  resource_type     = "users"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-prod-2024"  # Same users in Phase 1 every run
}
```

---

### Pattern: Different Seeds to Distribute Pilot Burden

**Critical Concept**: Using unique seeds per rollout prevents the same users from always being early adopters, avoiding "pilot fatigue."

**Behavior:**
- Different seeds produce different distributions from the same population
- User X might be in 10% pilot for MFA rollout but 60% final wave for Windows Updates
- Each rollout gets independent distribution

**Use When:**
- Running multiple independent rollouts (MFA, Windows Updates, Conditional Access policies)
- Want to fairly distribute pilot burden across organization
- Avoiding user complaints about always being in pilot groups

**Seed Naming Convention:**
- Use descriptive names: `"mfa-rollout-2024"`, `"windows-updates-q1"`, `"ca-pilot-phase2"`
- Include initiative name and timeframe
- Avoid generic names: `"seed1"`, `"test"`, `"2024"` - be explicit

**Example:**

```hcl
# MFA Rollout - User X ends up in shard_0 (10% pilot)
data "microsoft365_utility_guid_list_sharder" "mfa" {
  resource_type     = "users"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-rollout-2024"
}

# Windows Updates - Same User X ends up in shard_2 (80% final wave)
data "microsoft365_utility_guid_list_sharder" "windows" {
  resource_type     = "devices"
  shard_percentages = [5, 15, 80]
  strategy          = "percentage"
  seed              = "windows-updates-2024"  # Different seed = different distribution
}

# Conditional Access - User X ends up in shard_1 (35% broader)
data "microsoft365_utility_guid_list_sharder" "ca" {
  resource_type     = "users"
  shard_percentages = [15, 35, 50]
  strategy          = "percentage"
  seed              = "ca-policies-2024"  # Each rollout gets unique seed
}
```

**Result**: User X experiences different phases across rollouts - not always stuck in pilot groups.

---

### Anti-Patterns (Avoid These)

❌ **Using same seed for unrelated rollouts**

```hcl
# BAD: Defeats purpose of distributing pilot burden
data "microsoft365_utility_guid_list_sharder" "mfa" {
  seed = "2024"  # Too generic
}

data "microsoft365_utility_guid_list_sharder" "windows" {
  seed = "2024"  # Same users in pilot for both!
}
```

✅ **Better: Unique seeds per rollout**

```hcl
data "microsoft365_utility_guid_list_sharder" "mfa" {
  seed = "mfa-2024"
}

data "microsoft365_utility_guid_list_sharder" "windows" {
  seed = "windows-2024"
}
```

---

❌ **Changing seed accidentally mid-rollout**

```hcl
# BAD: Typo redistributes users mid-rollout
seed = "mfa-rolout-2024"  # Typo changes distribution!
```

✅ **Better: Use variables for consistency**

```hcl
variable "mfa_seed" {
  default = "mfa-rollout-2024"
}

data "microsoft365_utility_guid_list_sharder" "mfa" {
  seed = var.mfa_seed
}
```

---

## Related Resources

- [microsoft365_graph_beta_conditional_access_policy](../resources/graph_beta_conditional_access_policy.md)
- [microsoft365_graph_beta_group](../resources/graph_beta_groups_group.md)
- [Microsoft Graph API - Users](https://learn.microsoft.com/en-us/graph/api/user-list)
- [Microsoft Graph API - Devices](https://learn.microsoft.com/en-us/graph/api/device-list)
- [Microsoft Graph API - Group Members](https://learn.microsoft.com/en-us/graph/api/group-list-members)
