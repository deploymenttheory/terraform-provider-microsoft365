---
page_title: "Progressive Rollout with GUID List Sharder - terraform-provider-microsoft365"
subcategory: "Guides"
description: |-
  Complete guide for implementing phased deployments and progressive rollouts using the GUID List Sharder utility data source.
---

# Progressive Rollout with GUID List Sharder

This guide demonstrates how to implement progressive rollouts, phased deployments, and pilot programs for Microsoft 365 policies using the `microsoft365_utility_guid_list_sharder` data source.

## Overview

The GUID List Sharder data source (`microsoft365_utility_guid_list_sharder`) is a utility tool that queries Microsoft Graph API to retrieve collections of object IDs (GUIDs) for users, devices, or group members, then intelligently distributes them into configurable "shards" (subsets) for progressive deployment strategies.

### What Problem Does This Solve?

When deploying policies, configurations, or security controls across large Microsoft 365 environments, immediate organization-wide rollouts carry significant risk. If a policy misconfiguration or unexpected behavior occurs, it can impact thousands of users or devices simultaneously. The GUID List Sharder solves this by:

1. **Enabling Progressive Rollouts**: Deploy changes to small pilot groups first (e.g., 10% of users), validate functionality, then gradually expand to broader populations
2. **Reducing Blast Radius**: Limit the impact of potential issues by controlling which users/devices receive changes at each phase
3. **Facilitating Validation**: Allow time to monitor, test, and validate each phase before proceeding to the next
4. **Distributing Pilot Burden**: By using unique seed values across different rollouts, the same users won't always be early adopters. User X might be in the 10% pilot for MFA rollout (seed: "mfa-2024") but in the 60% final wave for Windows Updates (seed: "windows-2024"), preventing "pilot fatigue" where certain users consistently experience issues first
5. **Supports Multiple Deployment Strategies**: Choose between hash (deterministic assignment), round-robin (equal distribution), or percentage-based (custom-weighted) distribution. All strategies support optional seed for reproducibility

### How It Works

The data source operates in three stages:

1. **Query**: Retrieves object IDs from Microsoft Graph endpoints (`/users`, `/devices`, `/applications`, `/servicePrincipals`, or `/groups/{id}/members`) with optional OData filtering to narrow the population
2. **Shard**: Applies one of four distribution strategies to divide the population into subsets:
   - **Rendezvous**: Highest Random Weight (HRW) algorithm for deterministic assignment with minimal disruption when changing shard counts
   - **Round-robin**: Circular distribution for guaranteed equal sizes (optional seed for reproducibility)
   - **Percentage**: Custom-weighted distribution (e.g., 10% pilot, 30% broader, 60% full, optional seed for reproducibility)
   - **Size**: Absolute fixed sizes (e.g., 50 pilot, 200 broader, remainder full, optional seed for reproducibility)
3. **Output**: Returns a map of shards containing sets of object IDs, directly compatible with Terraform resources like conditional access policies and groups

### Key Benefits

- **Type-Safe Integration**: Outputs are `set(string)` types that work directly with Microsoft 365 resources without type conversion
- **Flexible Sizing**: Support both equal distribution (via `shard_count`) and custom percentages (via `shard_percentages`)
- **Deterministic Options**: All strategies support optional seed for consistent shard assignment across multiple Terraform runs
- **No Manual Management**: Automatically handles pagination and eventual consistency for large populations
- **Production-Ready**: Supports OData filtering to exclude disabled accounts, target specific attributes, or align with existing policy filters

### Example Use Cases

- **MFA Rollout**: Deploy multi-factor authentication to small pilots before organization-wide deployment
- **Windows Update Rings**: Distribute devices across early adopter, validation, and broad deployment rings
- **Conditional Access Pilots**: Test new access policies with controlled user populations
- **App-Based Conditional Access**: Target specific applications or service principals for conditional access policies
- **AI/Copilot Agent Policies**: Create policies targeting AI agent service principals (Microsoft Copilot, Azure AI agents)
- **Group Splitting**: Divide large groups into manageable subgroups for targeted policies
- **A/B Testing**: Distribute users across test groups for policy or feature testing

## API Endpoints Used

The data source queries the following Microsoft Graph Beta endpoints:

- `GET /users` - When `resource_type = "users"`
- `GET /devices` - When `resource_type = "devices"`
- `GET /applications` - When `resource_type = "applications"`
- `GET /servicePrincipals` - When `resource_type = "service_principals"`
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
  - `Application.Read.All` or `Directory.Read.All` - For applications
  - `Application.Read.All` or `Directory.Read.All` - For service principals
  - `Group.Read.All` or `GroupMember.Read.All` - For group members

## Choosing the Right Strategy

Use this decision matrix to select the optimal strategy and seed combination for your use case:

### Decision Matrix

| Your Need | Strategy | Seed | Why |
|-----------|----------|------|-----|
| **Same shards everywhere, shard count may change** (e.g., standard user groups, minimal disruption when adding shards) | `rendezvous` | ✅ Required | HRW algorithm minimizes reassignment when shard count changes |
| **Different shards per rollout** (e.g., vary who's in pilot groups to avoid pilot fatigue) | `rendezvous` | ✅ Yes | Different seeds distribute pilot burden across different users |
| **Exact equal sizes** (e.g., 25% in each of 4 groups) | `round-robin` | 🟡 Optional | Round-robin guarantees equal sizes (±1). Add seed for reproducibility |
| **Custom percentages, one-time split** (e.g., quick 10/30/60 split, don't need it again) | `percentage` | ❌ None | Fast, no need for reproducibility |
| **Custom percentages, reproducible** (e.g., 10/30/60 split you'll recreate later) | `percentage` | ✅ Yes | Seed ensures same users in same phases across runs |
| **Absolute sizes** (e.g., exactly 50 pilot users, 200 broader, remainder full) | `size` | 🟡 Optional | Fixed capacities regardless of total population |
| **Capacity testing / A/B testing** | `round-robin` | ✅ Yes | Exact equal sizes + reproducible results |

### Common Scenarios

#### Scenario 1: Standard User Groups (Consistent Everywhere)
**Need**: Create user groups that are referenced by multiple policies, all seeing the same membership.

```hcl
data "microsoft365_utility_guid_list_sharder" "standard_tiers" {
  resource_type = "users"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "standard-tiers-v1"  # Required for rendezvous
}
```

**Result**: All policies using these shards see the same tier membership consistently.

---

#### Scenario 2: Multiple Independent Rollouts (Distribute Pilot Burden)
**Need**: Running MFA rollout, Windows Updates, and CA policies. Don't want same users always in pilot.

```hcl
# MFA - User X in 10% pilot
data "microsoft365_utility_guid_list_sharder" "mfa" {
  resource_type = "users"
  shard_count   = 10
  strategy      = "rendezvous"
  seed          = "mfa-rollout-2024"
}

# Windows Updates - User X in 80% final wave
data "microsoft365_utility_guid_list_sharder" "windows" {
  resource_type = "devices"
  shard_count   = 10
  strategy      = "rendezvous"
  seed          = "windows-2024"  # Different seed = different distribution
}
```

**Result**: User X experiences different rollout phases across different initiatives.

---

#### Scenario 3: Structured Phased Rollout (Reproducible Phases)
**Need**: 10% pilot → 30% broader → 60% full. Need to recreate exact same phases if you rerun Terraform.

```hcl
data "microsoft365_utility_guid_list_sharder" "mfa_phases" {
  resource_type     = "users"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-phases-2024"  # Reproducible
}
```

**Result**: Same users in phase 1 every time you apply.

---

#### Scenario 4: Quick One-Time Equal Split
**Need**: Split a group into equal parts once, don't care about reproducing it later.

```hcl
data "microsoft365_utility_guid_list_sharder" "quick_split" {
  resource_type = "users"
  shard_count   = 4
  strategy      = "round-robin"
  # No seed - fast, one-time use
}
```

**Result**: Four groups of exactly equal size (±1).

---

#### Scenario 5: A/B Testing (Equal + Reproducible)
**Need**: Split users 50/50 for testing, need to recreate exact same groups.

```hcl
data "microsoft365_utility_guid_list_sharder" "ab_test" {
  resource_type = "users"
  odata_filter  = "department eq 'Engineering'"
  shard_count   = 2
  strategy      = "round-robin"
  seed          = "ab-test-q1-2024"  # Reproducible equal split
}
```

**Result**: Exactly 50/50 split, same users in group A/B every time.

## Sharding Strategies

The data source supports four distribution strategies, all with seed support:

### Strategy Comparison

| Strategy      | Distribution Method | Seed        | Equal Sizes | Custom Sizes | Best For                          |
|---------------|---------------------|-------------|-------------|--------------|-----------------------------------|
| `rendezvous`  | Highest Random Weight| ✅ Required | ~Equal      | ❌ No        | Deterministic assignment with minimal disruption when adding/removing shards|
| `round-robin` | Circular order      | 🟡 Optional | ✅ Exact    | ❌ No        | Guaranteed equal distribution     |
| `percentage`  | Sequential chunks   | 🟡 Optional | ❌ No       | ✅ Yes (%)   | Custom-weighted phased rollouts (percentages)  |
| `size`        | Absolute sizes      | 🟡 Optional | ❌ No       | ✅ Yes (fixed)| Fixed capacity constraints (exact counts)       |

**Seed Behavior (All Strategies):**
- **`rendezvous`**: Seed is required. Different seeds produce different distributions
- **`round-robin`/`percentage`/`size`**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies strategy (reproducible)

### Rendezvous Strategy (Highest Random Weight)

Uses the Highest Random Weight (HRW) algorithm for deterministic assignment. Each GUID-shard pair gets a hash weight, and the GUID goes to the shard with highest weight. This minimizes disruption when shard counts change.

**Characteristics:**
- Same GUID always goes to the same shard (with same seed)
- **Seed is required** for this strategy
- Different seeds produce different distributions (distributes pilot burden)
- Distribution is consistent across multiple Terraform runs
- Approximately equal shard sizes
- **Minimal reassignment when changing shard count**: Only affected GUIDs are the ones that need to move to fill new shards

**Use When:**
- You need deterministic assignment based on GUID properties
- Want reproducible results across Terraform runs
- May need to add/remove shards later with minimal disruption
- Need to vary which users are in pilot groups across different rollouts (use different seeds)
- Building reusable shard sets with potential for growth

**Example (different distribution per rollout):**

```hcl
# MFA Rollout - User X might be in shard_0 (10% pilot)
data "microsoft365_utility_guid_list_sharder" "mfa_rollout" {
  resource_type = "users"
  odata_filter  = "accountEnabled eq true"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "mfa-rollout-2024"  # Unique seed for this rollout
}

# Windows Updates - Same User X might be in shard_2 (60% final wave)
data "microsoft365_utility_guid_list_sharder" "windows_updates" {
  resource_type = "devices"
  odata_filter  = "operatingSystem eq 'Windows'"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "windows-updates-2024"  # Different seed = different distribution
}
```

### Round-Robin Strategy

Distributes GUIDs in circular order, one to each shard in sequence.

**Characteristics:**
- Simple circular distribution pattern (item 0→shard 0, item 1→shard 1, item 2→shard 2, item 3→shard 0, etc.)
- Guarantees exactly equal shard sizes (within ±1)
- **Without seed**: Uses API order (non-deterministic, may change between runs)
- **With seed**: Shuffles deterministically first, then applies round-robin (reproducible results)
- Fast and straightforward

**Use When:**
- You need exact equal distribution
- One-time split (no seed) or reproducible split (with seed)
- Doing statistical sampling or capacity testing
- Want equal sizes but also want to vary distribution across different rollouts (use different seeds)

**Example Without Seed (API order, non-deterministic):**

```hcl
data "microsoft365_utility_guid_list_sharder" "equal_split" {
  resource_type = "users"
  shard_count   = 4
  strategy      = "round-robin"
  # No seed - uses API order
}
```

**Example With Seed (deterministic, reproducible):**

```hcl
# MFA Rollout - deterministic equal split
data "microsoft365_utility_guid_list_sharder" "mfa_equal" {
  resource_type = "users"
  shard_count   = 4
  strategy      = "round-robin"
  seed          = "mfa-2024"  # Different seed per rollout
}

# Windows Updates - same users get different shards
data "microsoft365_utility_guid_list_sharder" "windows_equal" {
  resource_type = "devices"
  shard_count   = 4
  strategy      = "round-robin"
  seed          = "windows-2024"  # Different seed = different distribution
}
```

### Percentage Strategy (Custom-Weighted)

Distributes GUIDs according to specified percentages.

**Characteristics:**
- Flexible shard sizes based on your percentages
- Fills shards sequentially (shard_0 first, then shard_1, etc.)
- Last shard receives all remaining GUIDs to ensure nothing is lost
- **Without seed**: Uses API order (non-deterministic, may change between runs)
- **With seed**: Shuffles deterministically first, then applies percentage split (reproducible results)
- Supports any percentage split

**Use When:**
- Following specific rollout percentages (10% pilot, 30% expansion, 60% full)
- Implementing Windows Update rings with industry-standard distributions
- Need reproducible results with custom-sized shards (use seed)
- Change management process specifies exact pilot sizes
- Different phases have different risk profiles

**Example Without Seed (API order, non-deterministic):**

```hcl
data "microsoft365_utility_guid_list_sharder" "phased_rollout" {
  resource_type     = "users"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  # No seed - uses API order
}
```

**Example With Seed (deterministic, reproducible):**

```hcl
# MFA Rollout - deterministic percentage split
data "microsoft365_utility_guid_list_sharder" "mfa_phases" {
  resource_type     = "users"
  odata_filter      = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-2024"  # Reproducible results
}

# Windows Updates - same users get different phases
data "microsoft365_utility_guid_list_sharder" "windows_phases" {
  resource_type     = "devices"
  odata_filter      = "operatingSystem eq 'Windows'"
  shard_percentages = [5, 15, 80]
  strategy          = "percentage"
  seed              = "windows-2024"  # Different seed = different distribution
}
```

### Size Strategy (Absolute Counts)

Distributes GUIDs according to specified absolute shard sizes (exact number of GUIDs per shard).

**Characteristics:**
- Fixed shard sizes regardless of total population
- Fills shards sequentially (shard_0 first, then shard_1, etc.)
- Use `-1` in last position to mean "all remaining GUIDs"
- **Without seed**: Uses API order (non-deterministic, may change between runs)
- **With seed**: Shuffles deterministically first, then applies size split (reproducible results)
- Perfect for capacity-constrained scenarios

**Use When:**
- Support team can handle exactly N pilot users (e.g., 50 users max)
- Testing infrastructure has fixed capacity limits
- Need exact counts for resource planning
- Compliance requires specific pilot sizes
- Want same size pilots regardless of total population changes

**Example Without Seed (API order, non-deterministic):**

```hcl
data "microsoft365_utility_guid_list_sharder" "capacity_constrained" {
  resource_type = "users"
  shard_sizes   = [50, 200, -1]  # 50 pilot, 200 broader, rest full
  strategy      = "size"
  # No seed - uses API order
}
```

**Example With Seed (deterministic, reproducible):**

```hcl
# MFA Rollout - exact capacity limits
data "microsoft365_utility_guid_list_sharder" "mfa_capacity" {
  resource_type = "users"
  odata_filter  = "accountEnabled eq true"
  shard_sizes   = [50, 200, -1]  # Support handles 50, then 200, then all
  strategy      = "size"
  seed          = "mfa-2024"  # Reproducible results
}

# Windows Updates - different fixed sizes
data "microsoft365_utility_guid_list_sharder" "windows_capacity" {
  resource_type = "devices"
  odata_filter  = "operatingSystem eq 'Windows'"
  shard_sizes   = [10, 30, 100, -1]  # Early adopters, IT, pilot, broad
  strategy      = "size"
  seed          = "windows-2024"  # Different seed = different distribution
}
```

## Seed Usage Patterns

Understanding when and how to use seeds is critical for effective rollout management.

### Pattern: No Seed (Fast, Non-Deterministic)

**Use for Round-Robin/Percentage/Size Strategies:**
- Quick one-time splits where reproducibility isn't important
- Fast processing (no shuffle overhead)
- When API order is acceptable

**Note:** The `rendezvous` strategy always requires a seed.

**Example:**

```hcl
# Quick one-time percentage split
data "microsoft365_utility_guid_list_sharder" "quick_phases" {
  resource_type     = "users"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  # No seed - fast, one-time use
}
```

### Pattern: Unique Seed Per Rollout (Distribute Pilot Burden)

**Use When:**
- Running multiple independent rollouts (MFA, Windows Updates, CA policies, etc.)
- Want to avoid same users always being in pilot groups
- Each rollout should have different distribution of who's in which phase

**Seed Naming Convention:**
- Use descriptive names: `"mfa-rollout-2024"`, `"windows-updates-q1"`, `"ca-pilot-phase2"`
- Include initiative name and time period
- Avoid generic seeds like `"seed1"`, `"test"` - be explicit

**Example:**

```hcl
# Each rollout gets unique seed - using rendezvous for deterministic distribution
data "microsoft365_utility_guid_list_sharder" "mfa" {
  resource_type = "users"
  shard_count   = 10  # 10 equal shards
  strategy      = "rendezvous"
  seed          = "mfa-rollout-2024"  # User X → shard_0 (10% pilot)
}

data "microsoft365_utility_guid_list_sharder" "windows" {
  resource_type = "devices"
  shard_count   = 10  # 10 equal shards
  strategy      = "rendezvous"
  seed          = "windows-updates-2024"  # User X → shard_8 (80% wave)
}

data "microsoft365_utility_guid_list_sharder" "ca" {
  resource_type = "users"
  shard_count   = 10  # 10 equal shards
  strategy      = "rendezvous"
  seed          = "ca-policies-2024"  # User X → shard_3 (30% wave)
}
```

**Result**: User X experiences different phases across rollouts, preventing pilot fatigue.

### Pattern: Same Seed (Reproducible Across Runs)

**Use When:**
- Need to recreate exact same shards in future Terraform runs
- Debugging or testing requires consistent groups
- Compliance/audit needs reproducible results
- Long-term rollouts where phase membership shouldn't change

**Example:**

```hcl
# Production MFA rollout - must be reproducible
data "microsoft365_utility_guid_list_sharder" "mfa_production" {
  resource_type     = "users"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-prod-2024"  # Same seed every run
}

# If you terraform apply again tomorrow, same users are in Phase 1
```

### Pattern: Seed Versioning (Controlled Changes)

**Use When:**
- Want to change distribution but in a controlled manner
- Need audit trail of when distribution changed
- Moving from one rollout phase to another

**Example:**

```hcl
# Phase 1: Initial rollout
data "microsoft365_utility_guid_list_sharder" "mfa_v1" {
  resource_type = "users"
  shard_count   = 10
  strategy      = "rendezvous"
  seed          = "mfa-rollout-v1"
}

# Phase 2: Later, change distribution by updating seed
data "microsoft365_utility_guid_list_sharder" "mfa_v2" {
  resource_type = "users"
  shard_count   = 10
  strategy      = "rendezvous"
  seed          = "mfa-rollout-v2"  # Different seed = new distribution
}
```

### Anti-Patterns (Avoid These)

❌ **Using same seed for unrelated rollouts**

```hcl
# BAD: Same seed defeats the purpose of distributing pilot burden
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
  seed = "mfa-2024"  # Specific
}

data "microsoft365_utility_guid_list_sharder" "windows" {
  seed = "windows-2024"  # Different distribution
}
```

---

❌ **Changing seed accidentally during production rollout**

```hcl
# BAD: Typo changes distribution mid-rollout
data "microsoft365_utility_guid_list_sharder" "production" {
  seed = "mfa-rolout-2024"  # Typo! Different from "mfa-rollout-2024"
}
```

✅ **Better: Use variables for consistency**

```hcl
variable "mfa_rollout_seed" {
  default = "mfa-rollout-2024"
}

data "microsoft365_utility_guid_list_sharder" "production" {
  seed = var.mfa_rollout_seed
}
```

## Common Patterns

### Pattern 1: MFA Progressive Rollout

Roll out MFA requirements in three phases: pilot, broader pilot, and full deployment.

```hcl
# Shard users into three groups with unique seed
data "microsoft365_utility_guid_list_sharder" "mfa_rollout" {
  resource_type     = "users"
  odata_filter      = "accountEnabled eq true and userType eq 'Member'"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  # Note: seed is optional for percentage - add for reproducible results
}

# Phase 1: Pilot (10%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_pilot" {
  display_name = "MFA Required - Phase 1 Pilot"
  state        = "enabled"
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"]
    }
    applications {
      include_applications = ["All"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Phase 2: Broader Pilot (30%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_broader" {
  display_name = "MFA Required - Phase 2 Broader Pilot"
  state        = "enabledForReportingButNotEnforced"  # Start in report-only
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"]
    }
    applications {
      include_applications = ["All"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Phase 3: Full Rollout (60%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_full" {
  display_name = "MFA Required - Phase 3 Full Rollout"
  state        = "disabled"  # Enable after Phase 2 validation
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"]
    }
    applications {
      include_applications = ["All"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Monitor pilot sizes
output "mfa_rollout_distribution" {
  value = {
    pilot_count         = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"])
    broader_pilot_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"])
    full_rollout_count  = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"])
  }
}
```

### Pattern 2: Windows Update Deployment Rings

Implement industry-standard Windows Update rings with 5% early adopters, 15% IT validation, and 80% broad deployment.

```hcl
# Shard Windows devices into update rings with unique seed
data "microsoft365_utility_guid_list_sharder" "update_rings" {
  resource_type     = "devices"
  odata_filter      = "operatingSystem eq 'Windows' and accountEnabled eq true"
  shard_percentages = [5, 15, 80]
  strategy          = "percentage"
  # Note: seed is optional for percentage - add for reproducible results
}

# Ring 0: Early Adopters / Canary (5%)
resource "microsoft365_graph_beta_windows_update_for_business_configuration" "ring_0" {
  display_name = "Windows Update - Ring 0 (Early Adopters)"
  
  quality_update_defer_period_in_days  = 0
  feature_update_defer_period_in_days  = 0
  
  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.update_rings.shards["shard_0"]
    }
  }
}

# Ring 1: IT Pilot / Validation (15%)
resource "microsoft365_graph_beta_windows_update_for_business_configuration" "ring_1" {
  display_name = "Windows Update - Ring 1 (IT Validation)"
  
  quality_update_defer_period_in_days  = 3
  feature_update_defer_period_in_days  = 7
  
  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.update_rings.shards["shard_1"]
    }
  }
}

# Ring 2: Broad Deployment (80%)
resource "microsoft365_graph_beta_windows_update_for_business_configuration" "ring_2" {
  display_name = "Windows Update - Ring 2 (Broad)"
  
  quality_update_defer_period_in_days  = 7
  feature_update_defer_period_in_days  = 14
  
  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.update_rings.shards["shard_2"]
    }
  }
}

# Track ring distribution
output "update_ring_distribution" {
  value = {
    ring_0_count = length(data.microsoft365_utility_guid_list_sharder.update_rings.shards["shard_0"])
    ring_1_count = length(data.microsoft365_utility_guid_list_sharder.update_rings.shards["shard_1"])
    ring_2_count = length(data.microsoft365_utility_guid_list_sharder.update_rings.shards["shard_2"])
  }
}
```

### Pattern 3: Group Splitting and Resharding

Split an existing large group into multiple smaller groups for more granular policy targeting.

```hcl
# Get members from existing large group and split into 3
data "microsoft365_utility_guid_list_sharder" "split_group" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc"  # Original group ID
  odata_filter  = "accountEnabled eq true"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "sales-team-split-2024"  # Unique seed for this group split
}

# Create new pilot groups
resource "microsoft365_graph_beta_group" "pilot_group_a" {
  display_name     = "Sales Team - Pilot Group A"
  mail_nickname    = "sales-pilot-a"
  security_enabled = true
  
  group_members = data.microsoft365_utility_guid_list_sharder.split_group.shards["shard_0"]
}

resource "microsoft365_graph_beta_group" "pilot_group_b" {
  display_name     = "Sales Team - Pilot Group B"
  mail_nickname    = "sales-pilot-b"
  security_enabled = true
  
  group_members = data.microsoft365_utility_guid_list_sharder.split_group.shards["shard_1"]
}

resource "microsoft365_graph_beta_group" "pilot_group_c" {
  display_name     = "Sales Team - Pilot Group C"
  mail_nickname    = "sales-pilot-c"
  security_enabled = true
  
  group_members = data.microsoft365_utility_guid_list_sharder.split_group.shards["shard_2"]
}
```

### Pattern 4: A/B/C Testing

Distribute users equally across multiple test groups for policy or feature testing.

```hcl
# Equal distribution for testing
data "microsoft365_utility_guid_list_sharder" "ab_test" {
  resource_type = "users"
  odata_filter  = "department eq 'Engineering'"
  shard_count   = 3
  strategy      = "round-robin"  # Equal distribution
}

# Apply different policies to each test group
resource "microsoft365_graph_beta_conditional_access_policy" "test_group_a" {
  display_name = "Test Policy - Group A (Control)"
  state        = "enabled"
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ab_test.shards["shard_0"]
    }
  }
  # ... control policy settings
}

resource "microsoft365_graph_beta_conditional_access_policy" "test_group_b" {
  display_name = "Test Policy - Group B (Variant 1)"
  state        = "enabled"
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ab_test.shards["shard_1"]
    }
  }
  # ... variant 1 policy settings
}

resource "microsoft365_graph_beta_conditional_access_policy" "test_group_c" {
  display_name = "Test Policy - Group C (Variant 2)"
  state        = "enabled"
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ab_test.shards["shard_2"]
    }
  }
  # ... variant 2 policy settings
}
```

### Pattern 5: App-Based Conditional Access

Target specific Microsoft service principals (enterprise applications) for conditional access policies.

```hcl
# Get all Microsoft service principals and shard them
data "microsoft365_utility_guid_list_sharder" "microsoft_apps" {
  resource_type = "service_principals"
  odata_filter  = "startswith(displayName, 'Microsoft')"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "app-policies-2024"
}

# Phase 1: Pilot with Microsoft apps from shard_0
resource "microsoft365_graph_beta_conditional_access_policy" "app_pilot" {
  display_name = "MFA for Microsoft Apps - Phase 1"
  state        = "enabled"
  
  conditions {
    users {
      include_users = ["All"]
    }
    applications {
      include_service_principals = data.microsoft365_utility_guid_list_sharder.microsoft_apps.shards["shard_0"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}
```

### Pattern 6: AI/Copilot Agent Policies

Target AI agent service principals (Microsoft Copilot, Azure AI agents) for specific policies.

```hcl
# Get all AI/Copilot agent service principals
data "microsoft365_utility_guid_list_sharder" "ai_agents" {
  resource_type = "service_principals"
  odata_filter  = "isof('microsoft.graph.agentIdentityBlueprintPrincipal')"
  shard_count   = 2
  strategy      = "rendezvous"
  seed          = "ai-agent-policies-2024"
}

# Restrict AI agents to specific users
resource "microsoft365_graph_beta_conditional_access_policy" "ai_agent_restricted" {
  display_name = "Restrict AI Agents - Phase 1"
  state        = "enabled"
  
  conditions {
    users {
      include_users = ["12345678-1234-1234-1234-123456789012"]  # Specific user
    }
    applications {
      include_service_principals = data.microsoft365_utility_guid_list_sharder.ai_agents.shards["shard_0"]
    }
  }
  
  grant_controls {
    operator          = "AND"
    built_in_controls = ["block"]
  }
}
```

### Pattern 7: Capacity-Constrained Rollout

Use absolute shard sizes when you have fixed capacity constraints.

```hcl
# Support team can handle exactly 50 pilot users, then 200, then all
data "microsoft365_utility_guid_list_sharder" "support_capacity" {
  resource_type = "users"
  odata_filter  = "accountEnabled eq true"
  shard_sizes   = [50, 200, -1]  # 50 pilot, 200 broader, remainder full
  strategy      = "size"
  seed          = "mfa-support-2024"  # Reproducible
}

# Phase 1: Pilot (exactly 50 users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_pilot_50" {
  display_name = "MFA Required - 50 User Pilot"
  state        = "enabled"
  
  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.support_capacity.shards["shard_0"]
    }
    applications {
      include_applications = ["All"]
    }
  }
  
  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }
}

# Verify exact sizes
output "capacity_check" {
  value = {
    pilot_count   = length(data.microsoft365_utility_guid_list_sharder.support_capacity.shards["shard_0"])  # Should be exactly 50
    broader_count = length(data.microsoft365_utility_guid_list_sharder.support_capacity.shards["shard_1"])  # Should be exactly 200
    full_count    = length(data.microsoft365_utility_guid_list_sharder.support_capacity.shards["shard_2"])  # All remaining
  }
}
```

## Best Practices

### 1. Start Small

Begin with conservative pilot sizes (5-10%) to minimize impact if issues arise.

```hcl
shard_percentages = [5, 15, 80]  # Start with 5% pilot
```

### 2. Use OData Filters

Exclude disabled accounts and filter by relevant attributes to ensure clean populations.

```hcl
# Good practice
odata_filter = "accountEnabled eq true and userType eq 'Member'"

# For devices
odata_filter = "operatingSystem eq 'Windows' and accountEnabled eq true"

# For service principals (enterprise apps)
odata_filter = "startswith(displayName, 'Microsoft')"  # Microsoft apps only

# For AI/Copilot agents
odata_filter = "isof('microsoft.graph.agentIdentityBlueprintPrincipal')"
```

### 3. Choose the Right Strategy for Your Use Case

**Use `rendezvous` strategy when:**
- Need deterministic assignment based on GUID properties
- Want reproducible results (seed required)
- Running multiple rollouts and want to distribute pilot burden (use different seeds)
- May need to add/remove shards later with minimal disruption
- Building stable, long-term shard assignments

**Use `round-robin` strategy when:**
- Need guaranteed exactly equal shard sizes
- Doing A/B testing or capacity planning
- Quick one-time split (no seed) or reproducible equal split (with seed)

**Use `percentage` strategy when:**
- Following structured rollout percentages (10% → 30% → 60%)
- Different phases have different sizes
- One-time split (no seed) or reproducible phases (with seed)

**Use `size` strategy when:**
- Need exact absolute counts (e.g., exactly 50 pilot users)
- Support team has fixed capacity constraints
- Testing infrastructure has hard limits
- Compliance requires specific pilot sizes

### 4. Use Unique Seeds to Distribute Pilot Burden

Different rollouts should use different seeds to prevent the same users from always being in pilot groups:

```hcl
# MFA Rollout - User X might be in 10% pilot
data "microsoft365_utility_guid_list_sharder" "mfa" {
  resource_type = "users"
  strategy      = "rendezvous"
  seed          = "mfa-rollout-2024"
  shard_count   = 10
}

# Windows Updates - Same User X might be in 80% final wave
data "microsoft365_utility_guid_list_sharder" "updates" {
  resource_type = "devices"
  strategy      = "rendezvous"
  seed          = "windows-updates-2024"  # Different seed = different distribution
  shard_count   = 10
}
```

### 5. Use Descriptive Seed Names

Good: `"mfa-rollout-2024"`, `"windows-updates-q1"`, `"ca-pilot-phase2"`
Bad: `"seed1"`, `"test"`, `"2024"`

Descriptive names make intent clear and reduce errors.

### 6. Document Seed Changes

If you change a seed, document why. Seed changes redistribute users across shards.

```hcl
# v1: Initial rollout
# seed = "mfa-rollout-v1"

# v2: Changed to redistribute pilot burden after complaints
seed = "mfa-rollout-v2"
```

### 7. Monitor Shard Sizes

Output shard sizes to verify distribution matches expectations.

```hcl
output "rollout_distribution" {
  value = {
    pilot_count = length(data.microsoft365_utility_guid_list_sharder.example.shards["shard_0"])
    main_count  = length(data.microsoft365_utility_guid_list_sharder.example.shards["shard_1"])
  }
}
```

### 5. Use Report-Only Mode

Start policies in report-only mode before enforcement.

```hcl
resource "microsoft365_graph_beta_conditional_access_policy" "pilot" {
  state = "enabledForReportingButNotEnforced"  # Report-only initially
  # ... later change to "enabled"
}
```

### 6. Document Pilot Groups

Use descriptive names and comments to track rollout phases.

```hcl
resource "microsoft365_graph_beta_conditional_access_policy" "phase_1_pilot" {
  display_name = "MFA Required - Phase 1 Pilot (10%)"
  # Applied to shard_0 - 10% pilot group
  # Start date: 2024-01-15
  # Expected duration: 2 weeks
}
```

## Troubleshooting

### Problem: Shard sizes don't match percentages exactly

**Cause**: Percentages are calculated from the total population, and rounding occurs for decimal values.

**Solution**: This is expected behavior. The last shard receives all remaining GUIDs to ensure nothing is lost. Verify using output blocks:

```hcl
output "verify_distribution" {
  value = {
    total = sum([
      length(data.microsoft365_utility_guid_list_sharder.example.shards["shard_0"]),
      length(data.microsoft365_utility_guid_list_sharder.example.shards["shard_1"]),
      length(data.microsoft365_utility_guid_list_sharder.example.shards["shard_2"])
    ])
    shard_0 = length(data.microsoft365_utility_guid_list_sharder.example.shards["shard_0"])
    shard_1 = length(data.microsoft365_utility_guid_list_sharder.example.shards["shard_1"])
    shard_2 = length(data.microsoft365_utility_guid_list_sharder.example.shards["shard_2"])
  }
}
```

### Problem: Users moving between shards on re-run

**Cause**: Using `round-robin`, `percentage`, or `size` strategy without a seed, which are non-deterministic.

**Solution 1**: Add a seed to make the strategy deterministic:

```hcl
strategy = "round-robin"
seed     = "my-rollout-2024"  # Makes results reproducible
```

**Solution 2**: Switch to `rendezvous` strategy for deterministic assignment:

```hcl
strategy = "rendezvous"
seed     = "my-rollout-2024"  # Required for rendezvous
```

### Problem: Empty shards or API errors

**Cause**: OData filter too restrictive or no results found.

**Solution**: Test your OData filter independently and verify results:

```bash
# Test filter in Graph Explorer
GET https://graph.microsoft.com/beta/users?$filter=accountEnabled eq true
```

### Problem: Permission errors

**Cause**: Missing required permissions for the resource type.

**Solution**: Verify your service principal has the appropriate permissions:

- Users: `User.Read.All` or `Directory.Read.All`
- Devices: `Device.Read.All` or `Directory.Read.All`
- Group Members: `Group.Read.All` or `GroupMember.Read.All`

## Quick Reference

### Strategy & Seed Combinations

| Configuration | Result | Use Case |
|--------------|--------|----------|
| `strategy = "rendezvous"`, seed = "same" | Same distribution everywhere | Standard groups used by multiple policies |
| `strategy = "rendezvous"`, seed = "unique" | Different per seed | Distribute pilot burden across rollouts |
| `strategy = "round-robin"`, no seed | Equal sizes, API order | Quick one-time equal split |
| `strategy = "round-robin"`, seed = "unique" | Equal sizes, reproducible | A/B testing, capacity planning |
| `strategy = "percentage"`, no seed | Custom % sizes, API order | Quick phased split |
| `strategy = "percentage"`, seed = "unique" | Custom % sizes, reproducible | Reproducible phased rollout |
| `strategy = "size"`, no seed | Fixed counts, API order | Quick capacity-constrained split |
| `strategy = "size"`, seed = "unique" | Fixed counts, reproducible | Reproducible capacity-constrained rollout |

### When to Use Seed

| Scenario | Use Seed? | Reason |
|----------|-----------|--------|
| Using `rendezvous` strategy | ✅ Always required | Strategy requires seed for HRW algorithm |
| Multiple independent rollouts | ✅ Yes, unique per rollout | Distributes pilot burden |
| Need reproducibility | ✅ Yes, same seed | Consistent results across runs |
| One-time split (not rendezvous) | ❌ No | Faster, reproducibility not needed |
| Standard groups everywhere | ✅ Yes (rendezvous) | Deterministic, stable assignment |
| Testing/debugging | ✅ Yes | Reproducible results |

### Seed Naming Conventions

✅ **Good:**
- `"mfa-rollout-2024"`
- `"windows-updates-q1-2024"`
- `"ca-pilot-phase2"`
- `"ab-test-engineering-q4"`

❌ **Bad:**
- `"seed1"` - Not descriptive
- `"test"` - Too generic
- `"2024"` - Same seed for multiple rollouts defeats purpose

### Common HCL Patterns

```hcl
# Pattern 1: Standard groups (consistent everywhere)
data "microsoft365_utility_guid_list_sharder" "standard" {
  resource_type = "users"
  strategy      = "rendezvous"
  seed          = "standard-tiers-v1"  # Required for rendezvous
  shard_count   = 3
}

# Pattern 2: Distribute pilot burden
data "microsoft365_utility_guid_list_sharder" "rollout_a" {
  resource_type = "users"
  strategy      = "rendezvous"
  seed          = "rollout-a-2024"
  shard_count   = 10
}

data "microsoft365_utility_guid_list_sharder" "rollout_b" {
  resource_type = "devices"
  strategy      = "rendezvous"
  seed          = "rollout-b-2024"  # Different seed
  shard_count   = 10
}

# Pattern 3: Reproducible phased rollout
data "microsoft365_utility_guid_list_sharder" "phases" {
  resource_type     = "users"
  strategy          = "percentage"
  seed              = "mfa-phases-2024"
  shard_percentages = [10, 30, 60]
}

# Pattern 4: Equal A/B split
data "microsoft365_utility_guid_list_sharder" "ab_test" {
  resource_type = "users"
  strategy      = "round-robin"
  seed          = "ab-test-2024"
  shard_count   = 2
}

# Pattern 5: Fixed capacity constraints
data "microsoft365_utility_guid_list_sharder" "capacity" {
  resource_type = "users"
  strategy      = "size"
  seed          = "support-capacity-2024"
  shard_sizes   = [50, 200, -1]  # 50 pilot, 200 broader, rest full
}
```

## Related Resources

- [microsoft365_graph_beta_conditional_access_policy](../resources/graph_beta_conditional_access_policy.md)
- [microsoft365_graph_beta_group](../resources/graph_beta_groups_group.md)
- [Microsoft Graph API - Users](https://learn.microsoft.com/en-us/graph/api/user-list)
- [Microsoft Graph API - Devices](https://learn.microsoft.com/en-us/graph/api/device-list)
- [Microsoft Graph API - Group Members](https://learn.microsoft.com/en-us/graph/api/group-list-members)
