---
page_title: "microsoft365_utility_guid_list_sharder Data Source - terraform-provider-microsoft365"
subcategory: "Utility"

description: |-
  Retrieves object IDs (GUIDs) from Microsoft Graph API and distributes them into configurable shards for progressive rollouts and phased deployments. Queries /users, /devices, or /groups/{id}/members endpoints with optional OData filtering, then applies sharding strategies (random, sequential, or percentage-based) to distribute results. Output shards are sets that can be directly used in conditional access policies, groups, and other resources requiring object ID collections.
  API Endpoints: GET /users, GET /devices, GET /groups/{id}/members (with pagination and ConsistencyLevel: eventual header)
  Common Use Cases: MFA rollouts, Windows Update rings, conditional access pilots, group splitting, A/B testing
  For detailed examples and best practices, see the Progressive Rollout with GUID List Sharder https://registry.terraform.io/providers/deploymenttheory/microsoft365/latest/docs/guides/progressive_rollout_with_guid_list_sharder guide.
---

# microsoft365_utility_guid_list_sharder

Queries Microsoft Graph API to retrieve collections of object IDs (GUIDs) for users, devices, or group members, then intelligently distributes them into configurable "shards" (subsets) for progressive deployment strategies.

This datasource enables phased rollouts, pilot programs, and deployment rings for Microsoft 365 policies by algorithmically distributing populations into controlled subsets. Unlike static Entra ID dynamic groups that take hours to populate and require complex membership rules, the GUID List Sharder provides immediate, deterministic distribution with multiple strategies optimized for different scenarios.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.42.0-alpha | Experimental | Initial release with round-robin, percentage, size, and rendezvous strategies |

## Background

When deploying policies, configurations, or security controls across large Microsoft 365 environments, immediate organization-wide rollouts carry significant risk. If a policy misconfiguration or unexpected behavior occurs, it can impact thousands of users or devices simultaneously.

Consequently, it's common for organizations to define deployment rings or waves to manage the rollout of changes. However, this approach presents its own challenges:

- **Entra ID dynamic group membership** can take hours to fully populate
- **Creating deployment groups** that adequately represent hardware and software diversity across the organization is complex
- **Created groups become stale** as people join, leave, or change roles
- **Willing pilot users experience pilot fatigue** when they're consistently the first to encounter issues

The GUID List Sharder solves these challenges by:

1. **Enabling Progressive Rollouts**: Deploy changes to small pilot groups first (e.g., 10% of users), validate functionality, then gradually expand to broader populations
2. **Reducing Blast Radius**: Limit the impact of potential issues by controlling which users/devices receive changes at each phase
3. **Facilitating Validation**: Allow time to monitor, test, and validate each phase before proceeding to the next
4. **Distributing Pilot Burden**: By using unique seed values across different rollouts, the same users won't always be early adopters. User X might be in the 10% pilot for MFA rollout but in the 60% final wave for Windows Updates, preventing "pilot fatigue"
5. **Supports Multiple Deployment Strategies**: Choose between `round-robin` (perfect equal distribution), `percentage` (custom ratios), `size` (absolute counts), or `rendezvous` (minimal disruption when ring counts change)

## Distribution Strategies

### Round-Robin Strategy
- **Distribution**: Circular assignment (0→shard_0, 1→shard_1, 2→shard_2, 3→shard_0, cycling)
- **Balance**: Perfect ±1 GUID across all shards
- **Seed Behavior**: No seed = API order (non-deterministic). With seed = shuffles first, then applies round-robin (reproducible)
- **Use When**: Need equal-sized rings or capacity testing with perfect balance

### Percentage Strategy
- **Distribution**: Allocates by percentage ratios (e.g., 10% pilot, 30% broader, 60% full)
- **Balance**: Within rounding error of specified percentages
- **Seed Behavior**: No seed = API order (non-deterministic). With seed = shuffles first, then applies percentages (reproducible)
- **Use When**: Standard phased rollouts with percentage-based waves (e.g., 10%/30%/60%)

### Size Strategy
- **Distribution**: Allocates by absolute counts (e.g., exactly 50, 100, 200 users per shard)
- **Balance**: Exact counts as specified
- **Seed Behavior**: No seed = API order (non-deterministic). With seed = shuffles first, then applies size split (reproducible)
- **Special Feature**: Supports `-1` as last value meaning "all remaining GUIDs" (e.g., `[50, 200, -1]`)
- **Use When**: Precise shard sizes required regardless of total population

### Rendezvous Strategy
- **Distribution**: Each GUID independently computes hash scores for all shards using `SHA256(guid:shard_N:seed)`, then picks highest score
- **Balance**: Probabilistic ~equal distribution (typically within 3% for 1000+ GUIDs)
- **Seed Behavior**: Always deterministic. Seed affects which shard wins for each GUID
- **Stability**: When shard count changes, only ~1/n GUIDs move (theoretical minimum). Adding 4th shard to 3-shard setup: only ~25% reassign vs ~75% with position-based strategies
- **Use When**: Ring count will change during rollout lifecycle (e.g., start with 3 rings, later add 4th for extended validation)

## Key Concept: Distributing Pilot Burden

One of the most valuable features is the ability to use different seeds for different rollouts. This ensures the same users aren't always in pilot groups across all initiatives:

- **MFA rollout** (seed: `"mfa-2026"`): User A in 10% pilot
- **Windows Updates** (seed: `"windows-2026"`): Same User A in 80% broad deployment
- **Conditional Access** (seed: `"ca-2026"`): Same User A in 35% broader wave

This prevents "pilot fatigue" where certain users consistently experience issues first.

## Decision Matrix

| Requirement | Strategy | Seed Required | Notes |
|-------------|----------|---------------|-------|
| **Equal-sized rings** (e.g., 4 rings, ~250 users each) | `round-robin` | 🟡 Optional | Perfect ±1 GUID balance. Add seed for reproducibility |
| **Percentage-based** (e.g., 10% pilot, 30% broader, 60% full) | `percentage` | 🟡 Optional | Clean percentage splits. Add seed for reproducibility |
| **Absolute counts** (e.g., exactly 50, 100, 200 users) | `size` | 🟡 Optional | Precise shard sizes. Add seed for deterministic assignment. Supports `-1` for "all remaining" |
| **Ring count will change** (e.g., start with 3 rings, later expand to 4) | `rendezvous` | ✅ Required | Only ~25% of users move when adding rings (vs ~75% with other strategies) |
| **Distribute pilot burden** (different users in pilot for different rollouts) | Any strategy | ✅ Yes | Different seeds (e.g., `"mfa-2026"` vs `"windows-2026"`) = different distributions |
| **Capacity testing / A/B testing** | `round-robin` | ✅ Yes | Perfect balance + reproducible results |
| **One-time split, don't care about reproducibility** | `round-robin` or `percentage` | ❌ None | Fastest - no seed overhead |

## Example Usage

### Round-Robin Distribution (Equal Shards)

```terraform
# Round-Robin Distribution: Perfect equal distribution across 4 deployment rings
# Use case: Equal-sized pilot, validation, pre-production, and production rings

data "microsoft365_utility_guid_list_sharder" "mfa_rings" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "round-robin"
  seed          = "mfa-rollout-2026" # Optional: ensures reproducible distribution
}

# Ring 0: Pilot (exactly 25% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_ring_0" {
  display_name = "MFA Required - Ring 0 (Pilot)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"]
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

# Ring 1: Validation (exactly 25% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_ring_1" {
  display_name = "MFA Required - Ring 1 (Validation)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"]
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

# Ring 2: Pre-Production (exactly 25% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_ring_2" {
  display_name = "MFA Required - Ring 2 (Pre-Production)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"]
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

# Ring 3: Production (exactly 25% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_ring_3" {
  display_name = "MFA Required - Ring 3 (Production)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"]
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

# Verify perfect distribution
output "ring_distribution" {
  value = {
    ring_0_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"])
    ring_1_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"])
    ring_2_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"])
    ring_3_count = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"])
    total_users  = length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_2"]) + length(data.microsoft365_utility_guid_list_sharder.mfa_rings.shards["shard_3"])
  }
  description = "Round-robin guarantees perfect ±1 GUID balance across all rings"
}
```

### Percentage-Based Phased Rollout

```terraform
# Percentage-Based Distribution: Standard phased rollout pattern
# Use case: 10% pilot, 30% broader rollout, 60% full deployment

data "microsoft365_utility_guid_list_sharder" "ca_phases" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "ca-policies-2026" # Optional: ensures same users in same phases
}

# Phase 1: Pilot (10% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "block_legacy_auth_pilot" {
  display_name = "Block Legacy Auth - Phase 1 (10% Pilot)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_0"]
    }
    applications {
      include_applications = ["All"]
    }
    client_app_types = ["exchangeActiveSync", "other"]
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}

# Phase 2: Broader Rollout (30% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "block_legacy_auth_broader" {
  display_name = "Block Legacy Auth - Phase 2 (30% Broader)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_1"]
    }
    applications {
      include_applications = ["All"]
    }
    client_app_types = ["exchangeActiveSync", "other"]
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}

# Phase 3: Full Deployment (60% of users)
resource "microsoft365_graph_beta_conditional_access_policy" "block_legacy_auth_full" {
  display_name = "Block Legacy Auth - Phase 3 (60% Full)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_2"]
    }
    applications {
      include_applications = ["All"]
    }
    client_app_types = ["exchangeActiveSync", "other"]
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}

# Monitor phase distribution
output "phase_distribution" {
  value = {
    pilot_count   = length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_0"])
    broader_count = length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_1"])
    full_count    = length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_2"])
    total_users   = length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.ca_phases.shards["shard_2"])
  }
  description = "Phase counts (should be approximately 10%, 30%, 60% of total)"
}
```

### Fixed-Size Pilot Groups

```terraform
# Size-Based Distribution: Fixed pilot group sizes
# Use case: Compliance requires exactly 50 pilot users, 100 validation users

data "microsoft365_utility_guid_list_sharder" "compliance_pilot" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and department eq 'IT'"
  shard_sizes   = [50, 100, -1] # 50 pilot, 100 validation, remainder for broad
  strategy      = "size"
  seed          = "compliance-pilot-2026"
}

# Pilot Group: Exactly 50 users
resource "microsoft365_graph_beta_group" "compliance_pilot" {
  display_name     = "Compliance Policy - Pilot (50 users)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_0"]
}

# Validation Group: Exactly 100 users
resource "microsoft365_graph_beta_group" "compliance_validation" {
  display_name     = "Compliance Policy - Validation (100 users)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_1"]
}

# Broad Deployment: All remaining IT users
resource "microsoft365_graph_beta_group" "compliance_broad" {
  display_name     = "Compliance Policy - Broad (All Remaining)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_2"]
}

# Conditional Access Policy targeting pilot group
resource "microsoft365_graph_beta_conditional_access_policy" "compliance_policy_pilot" {
  display_name = "Device Compliance Required - Pilot"
  state        = "enabled"

  conditions {
    users {
      include_groups = [microsoft365_graph_beta_group.compliance_pilot.id]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["compliantDevice"]
  }
}

# Verify exact counts
output "pilot_group_sizes" {
  value = {
    pilot_count      = length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_0"])
    validation_count = length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_1"])
    broad_count      = length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_2"])
    total_it_users   = length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.compliance_pilot.shards["shard_2"])
  }
  description = "Pilot should be exactly 50, Validation exactly 100, Broad gets remainder"
}
```

### Rendezvous (Stable When Ring Count Changes)

```terraform
# Rendezvous Hashing: Stable when ring count changes
# Use case: Start with 3 rings, later expand to 4 without massive user disruption

# Initial deployment: 3 rings
data "microsoft365_utility_guid_list_sharder" "stable_deployment" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 3 # Change to 4 later - only ~25% of users will move
  strategy      = "rendezvous"
  seed          = "stable-deployment-2026" # Required for rendezvous
}

# Ring 0: Early Adopters
resource "microsoft365_graph_beta_group" "ring_0_early" {
  display_name     = "Windows Updates - Ring 0 (Early Adopters)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_0"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_0" {
  display_name = "Windows Updates - Ring 0 (Early Adopters)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 0
  feature_update_deferral_period_days = 0

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.ring_0_early.id
    }
  }
}

# Ring 1: Broad Deployment
resource "microsoft365_graph_beta_group" "ring_1_broad" {
  display_name     = "Windows Updates - Ring 1 (Broad)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_1"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_1" {
  display_name = "Windows Updates - Ring 1 (Broad)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 7
  feature_update_deferral_period_days = 14

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.ring_1_broad.id
    }
  }
}

# Ring 2: Production (Conservative)
resource "microsoft365_graph_beta_group" "ring_2_production" {
  display_name     = "Windows Updates - Ring 2 (Production)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_2"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_2" {
  display_name = "Windows Updates - Ring 2 (Production)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 14
  feature_update_deferral_period_days = 30

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.ring_2_production.id
    }
  }
}

# Monitor ring distribution
output "ring_distribution" {
  value = {
    ring_0_count = length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_0"])
    ring_1_count = length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_1"])
    ring_2_count = length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_2"])
    total_users  = length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.stable_deployment.shards["shard_2"])
  }
  description = "When changing shard_count from 3 to 4, only ~25% of users move (vs ~75% with other strategies)"
}
```

### Real-World Conditional Access Rollout

```terraform
# Real-World Example: Conditional Access MFA Rollout with distributed pilot burden
# Demonstrates using different seeds for different initiatives

# MFA Rollout: User A ends up in 10% pilot
data "microsoft365_utility_guid_list_sharder" "mfa_rollout" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "mfa-rollout-2026" # Unique seed per initiative
}

# Windows Updates Rollout: Same User A ends up in 80% final wave
data "microsoft365_utility_guid_list_sharder" "windows_rollout" {
  resource_type     = "devices"
  shard_percentages = [5, 15, 80]
  strategy          = "percentage"
  seed              = "windows-updates-2026" # Different seed = different distribution
}

# MFA Phase 1: Pilot (10%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_phase_1" {
  display_name = "Require MFA - Phase 1 (10% Pilot)"
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

# MFA Phase 2: Broader (30%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_phase_2" {
  display_name = "Require MFA - Phase 2 (30% Broader)"
  state        = "enabled"

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

# MFA Phase 3: Full (60%)
resource "microsoft365_graph_beta_conditional_access_policy" "mfa_phase_3" {
  display_name = "Require MFA - Phase 3 (60% Full)"
  state        = "enabled"

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

# Windows Updates Ring 0 (5% pilot devices)
resource "microsoft365_graph_beta_device_management_windows_update_ring" "windows_ring_0" {
  display_name = "Windows Updates - Ring 0 (5% Pilot)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  quality_update_deferral_period_days = 0
  feature_update_deferral_period_days = 0

  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_0"]
    }
  }
}

# Windows Updates Ring 1 (15% validation devices)
resource "microsoft365_graph_beta_device_management_windows_update_ring" "windows_ring_1" {
  display_name = "Windows Updates - Ring 1 (15% Validation)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  quality_update_deferral_period_days = 7
  feature_update_deferral_period_days = 14

  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_1"]
    }
  }
}

# Windows Updates Ring 2 (80% broad devices)
resource "microsoft365_graph_beta_device_management_windows_update_ring" "windows_ring_2" {
  display_name = "Windows Updates - Ring 2 (80% Broad)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  quality_update_deferral_period_days = 14
  feature_update_deferral_period_days = 30

  assignments {
    target {
      device_ids = data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_2"]
    }
  }
}

# Demonstrate distributed pilot burden
output "pilot_burden_distribution" {
  value = {
    mfa_pilot_users       = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_0"])
    mfa_broader_users     = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_1"])
    mfa_full_users        = length(data.microsoft365_utility_guid_list_sharder.mfa_rollout.shards["shard_2"])
    windows_pilot_devices = length(data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_0"])
    windows_broad_devices = length(data.microsoft365_utility_guid_list_sharder.windows_rollout.shards["shard_2"])
  }
  description = "Different seeds ensure User A isn't always in pilot groups across all initiatives"
}
```

### Group Members Distribution

```terraform
# Group Members Distribution: Shard members of an existing Entra ID group
# Use case: Deploy policy to IT department in phases without creating additional nested groups

data "microsoft365_utility_guid_list_sharder" "it_dept_phases" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc" # IT Department Group ID
  shard_count   = 3
  strategy      = "round-robin"
  seed          = "it-dept-pilot-2026"
}

# Phase 1: IT Pilot (1/3 of IT department)
resource "microsoft365_graph_beta_conditional_access_policy" "it_new_policy_phase_1" {
  display_name = "New IT Policy - Phase 1 (IT Pilot)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_0"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa", "compliantDevice"]
  }
}

# Phase 2: IT Validation (1/3 of IT department)
resource "microsoft365_graph_beta_conditional_access_policy" "it_new_policy_phase_2" {
  display_name = "New IT Policy - Phase 2 (IT Validation)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_1"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa", "compliantDevice"]
  }
}

# Phase 3: IT Full (1/3 of IT department)
resource "microsoft365_graph_beta_conditional_access_policy" "it_new_policy_phase_3" {
  display_name = "New IT Policy - Phase 3 (IT Full)"
  state        = "enabled"

  conditions {
    users {
      include_users = data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_2"]
    }
    applications {
      include_applications = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa", "compliantDevice"]
  }
}

# Monitor IT department phase distribution
output "it_dept_phase_distribution" {
  value = {
    phase_1_count    = length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_0"])
    phase_2_count    = length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_1"])
    phase_3_count    = length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_2"])
    total_it_members = length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.it_dept_phases.shards["shard_2"])
  }
  description = "Equal distribution across 3 phases (perfect ±1 balance with round-robin)"
}
```

### Device Distribution for Windows Updates

```terraform
# Device Distribution: Distribute managed devices for Windows Updates rollout
# Use case: Controlled Windows Update deployment across device population

data "microsoft365_utility_guid_list_sharder" "windows_update_rings" {
  resource_type     = "devices"
  odata_query       = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_percentages = [5, 15, 30, 50]
  strategy          = "percentage"
  seed              = "windows-updates-2026"
}

# Ring 0: Validation (5% of devices)
resource "microsoft365_graph_beta_group" "update_ring_0" {
  display_name     = "Windows Updates - Ring 0 (5% Validation)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_0"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_0_validation" {
  display_name = "Windows Updates - Ring 0 (5% Validation)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 0
  feature_update_deferral_period_days = 0
  deadline_for_quality_updates_days   = 7
  deadline_for_feature_updates_days   = 14

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.update_ring_0.id
    }
  }
}

# Ring 1: Pilot (15% of devices)
resource "microsoft365_graph_beta_group" "update_ring_1" {
  display_name     = "Windows Updates - Ring 1 (15% Pilot)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_1"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_1_pilot" {
  display_name = "Windows Updates - Ring 1 (15% Pilot)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 3
  feature_update_deferral_period_days = 7
  deadline_for_quality_updates_days   = 10
  deadline_for_feature_updates_days   = 21

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.update_ring_1.id
    }
  }
}

# Ring 2: Broad (30% of devices)
resource "microsoft365_graph_beta_group" "update_ring_2" {
  display_name     = "Windows Updates - Ring 2 (30% Broad)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_2"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_2_broad" {
  display_name = "Windows Updates - Ring 2 (30% Broad)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 7
  feature_update_deferral_period_days = 14
  deadline_for_quality_updates_days   = 14
  deadline_for_feature_updates_days   = 28

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.update_ring_2.id
    }
  }
}

# Ring 3: Production (50% of devices)
resource "microsoft365_graph_beta_group" "update_ring_3" {
  display_name     = "Windows Updates - Ring 3 (50% Production)"
  security_enabled = true

  members = data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_3"]
}

resource "microsoft365_graph_beta_device_management_windows_update_ring" "ring_3_production" {
  display_name = "Windows Updates - Ring 3 (50% Production)"

  automatic_update_mode               = "autoInstallAtMaintenanceTime"
  automatic_restart_notification      = "before"
  quality_update_deferral_period_days = 14
  feature_update_deferral_period_days = 30
  deadline_for_quality_updates_days   = 21
  deadline_for_feature_updates_days   = 45

  assignments {
    target {
      group_id = microsoft365_graph_beta_group.update_ring_3.id
    }
  }
}

# Monitor device ring distribution
output "device_ring_distribution" {
  value = {
    ring_0_validation_count = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_0"])
    ring_1_pilot_count      = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_1"])
    ring_2_broad_count      = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_2"])
    ring_3_production_count = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_3"])
    total_windows_devices   = length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_2"]) + length(data.microsoft365_utility_guid_list_sharder.windows_update_rings.shards["shard_3"])
  }
  description = "Device counts per ring (5%, 15%, 30%, 50% distribution)"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `resource_type` (String) The type of Microsoft Graph resource to query and shard. `users` queries `/users` for user-based policies (MFA, conditional access). `devices` queries `/devices` for device policies (Windows Updates, compliance). `group_members` queries `/groups/{id}/members` to split existing group membership (requires `group_id`).
- `strategy` (String) The distribution strategy for sharding GUIDs. `round-robin` distributes in circular order (guarantees equal sizes, optional seed for reproducibility). `percentage` distributes by specified percentages (requires `shard_percentages`, optional seed for reproducibility). `size` distributes by absolute sizes (requires `shard_sizes`, optional seed for reproducibility). `rendezvous` uses Highest Random Weight algorithm (always deterministic, minimal disruption when shard count changes, requires seed). See the [guide](https://registry.terraform.io/providers/deploymenttheory/microsoft365/latest/docs/guides/progressive_rollout_with_guid_list_sharder) for detailed comparison.

### Optional

- `group_id` (String) The object ID of the group to query members from. Required when `resource_type = "group_members"`, ignored otherwise. Use this to split an existing group's membership into multiple new groups for targeted policy application.
- `odata_query` (String) Optional OData filter applied at the API level before sharding. Common examples: `$filter=accountEnabled eq true` (active accounts only), `$filter=operatingSystem eq 'Windows'` (Windows devices), `$filter=userType eq 'Member'` (exclude guests). Leave empty to query all resources without filtering.
- `seed` (String) Optional seed value for deterministic distribution. When provided, makes results reproducible across Terraform runs. **`round-robin` strategy**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies round-robin (reproducible). **`percentage` strategy**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies percentage split (reproducible). **`size` strategy**: No seed = uses API order (may change). With seed = shuffles deterministically first, then applies size-based split (reproducible). **`rendezvous` strategy**: Always deterministic. Seed affects which shard wins for each GUID via Highest Random Weight algorithm. Use different seeds for different rollouts to distribute pilot burden: User X might be in shard_0 for MFA but shard_2 for Windows Updates.
- `shard_count` (Number) Number of equally-sized shards to create (minimum 1). Use with `round-robin` strategy. Conflicts with `shard_percentages` and `shard_sizes`. Creates shards named `shard_0`, `shard_1`, ..., `shard_N-1`. For custom-sized shards (e.g., 10% pilot, 30% broader, 60% full), use `shard_percentages` with `percentage` strategy instead.
- `shard_percentages` (List of Number) List of percentages for custom-sized shards. Use with `percentage` strategy. Conflicts with `shard_count` and `shard_sizes`. Values must be non-negative integers that sum to exactly 100. Example: `[10, 30, 60]` creates 10% pilot, 30% broader pilot, 60% full rollout. Common patterns: `[5, 15, 80]` (Windows Update rings), `[33, 33, 34]` (A/B/C testing). Last shard receives all remaining GUIDs to prevent loss.
- `shard_sizes` (List of Number) List of absolute shard sizes (exact number of GUIDs per shard). Use with `size` strategy. Conflicts with `shard_count` and `shard_percentages`. Values must be positive integers or -1 (which means 'all remaining'). Only the last element can be -1. Example: `[50, 200, -1]` creates 50 pilot users, 200 broader rollout, remainder for full deployment. Common patterns: `[10, 30, -1]` (controlled pilot expansion), `[100, 100, 100, -1]` (fixed-size rings). Use this when you need exact capacity constraints (e.g., support team handles exactly 50 pilot users).
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `shards` (Map of Set of String) Computed map of shard names (`shard_0`, `shard_1`, ...) to sets of GUIDs. Each value is a `set(string)` type, directly compatible with resource attributes expecting object ID sets (e.g., `conditions.users.include_users` in conditional access policies, `group_members` in groups). Access with `data.example.shards["shard_0"]`, check size with `length(data.example.shards["shard_0"])`.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Additional Resources

- [GUID List Sharder Guide](https://registry.terraform.io/providers/deploymenttheory/microsoft365/latest/docs/guides/guid_list_sharder)
- [Microsoft Graph API - Users](https://learn.microsoft.com/en-us/graph/api/user-list)
- [Microsoft Graph API - Devices](https://learn.microsoft.com/en-us/graph/api/device-list)
- [Microsoft Graph API - Group Members](https://learn.microsoft.com/en-us/graph/api/group-list-members)
- [Conditional Access Best Practices](https://learn.microsoft.com/en-us/entra/identity/conditional-access/plan-conditional-access)
- [Phased Deployment Strategies](https://learn.microsoft.com/en-us/mem/intune/fundamentals/deployment-guide-intune-setup)
