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
## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `Users.Read.All`
- `Devices.Read.All`
- `Groups.Read.All`

**Optional:**
- `None` `[N/A]`

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

### Basic Examples

#### Round-Robin (Equal Distribution)

```terraform
# Basic Round-Robin: Equal distribution across shards
# Perfect ±1 balance guaranteed

# Without seed (non-deterministic, uses API order)
data "microsoft365_utility_guid_list_sharder" "users_no_seed" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "round-robin"
}

# With seed (deterministic, reproducible)
data "microsoft365_utility_guid_list_sharder" "users_with_seed" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 4
  strategy      = "round-robin"
  seed          = "mfa-rollout-2024"
}

# Output shard counts
output "distribution" {
  value = {
    shard_0 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_0"])
    shard_1 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_1"])
    shard_2 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_2"])
    shard_3 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_3"])
  }
}
```

#### Percentage (Custom Ratios)

```terraform
# Basic Percentage: Custom ratio distribution
# Common pattern: 10% pilot, 30% broader, 60% full

# Without seed (non-deterministic, uses API order)
data "microsoft365_utility_guid_list_sharder" "users_no_seed" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
}

# With seed (deterministic, reproducible)
data "microsoft365_utility_guid_list_sharder" "users_with_seed" {
  resource_type     = "users"
  odata_query       = "accountEnabled eq true"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "ca-rollout-2024"
}

# Output shard counts
output "distribution" {
  value = {
    pilot_10pct   = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_0"])
    broader_30pct = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_1"])
    full_60pct    = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_2"])
  }
}
```

#### Size (Absolute Counts)

```terraform
# Basic Size: Absolute count distribution
# Use -1 for "all remaining" in last position

# Without seed (non-deterministic, uses API order)
data "microsoft365_utility_guid_list_sharder" "users_no_seed" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and department eq 'IT'"
  shard_sizes   = [50, 100, -1]
  strategy      = "size"
}

# With seed (deterministic, reproducible)
data "microsoft365_utility_guid_list_sharder" "users_with_seed" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true and department eq 'IT'"
  shard_sizes   = [50, 100, -1]
  strategy      = "size"
  seed          = "it-pilot-2024"
}

# Output shard counts
output "distribution" {
  value = {
    pilot_exact_50       = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_0"])
    validation_exact_100 = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_1"])
    broad_all_remaining  = length(data.microsoft365_utility_guid_list_sharder.users_with_seed.shards["shard_2"])
  }
}
```

#### Rendezvous (Stable Distribution)

```terraform
# Basic Rendezvous: Stable distribution when shard count changes
# Seed is REQUIRED for rendezvous strategy
# Only ~1/n GUIDs move when adding shards

data "microsoft365_utility_guid_list_sharder" "users_stable" {
  resource_type = "users"
  odata_query   = "accountEnabled eq true"
  shard_count   = 3
  strategy      = "rendezvous"
  seed          = "stable-deployment-2024"
}

# Output shard counts
output "distribution" {
  value = {
    ring_0 = length(data.microsoft365_utility_guid_list_sharder.users_stable.shards["shard_0"])
    ring_1 = length(data.microsoft365_utility_guid_list_sharder.users_stable.shards["shard_1"])
    ring_2 = length(data.microsoft365_utility_guid_list_sharder.users_stable.shards["shard_2"])
  }
}

# When you change shard_count from 3 to 4:
# - Only ~25% of users will move to new ring_3
# - ~75% stay in their original ring
# - Compare with round-robin/percentage: ~75% would move
```

### Scenario-Based Examples

#### Scenario 1: Devices → Settings Catalog Deployment Rings

```terraform
# Scenario 1: Devices → Groups → Settings Catalog Deployment Rings
# Use case: Roll out Windows Update policies in phases across managed devices

# Distribute Windows devices into 4 deployment rings (5%, 15%, 30%, 50%)
data "microsoft365_utility_guid_list_sharder" "windows_devices" {
  resource_type     = "devices"
  odata_query       = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_percentages = [5, 15, 30, 50]
  strategy          = "percentage"
  seed              = "windows-updates-2024"
}

# Create groups for each deployment ring
resource "microsoft365_graph_beta_group" "ring_0_validation" {
  display_name     = "Windows Updates - Ring 0 (5% Validation)"
  mail_nickname    = "win-updates-ring-0"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.windows_devices.shards["shard_0"]
}

resource "microsoft365_graph_beta_group" "ring_1_pilot" {
  display_name     = "Windows Updates - Ring 1 (15% Pilot)"
  mail_nickname    = "win-updates-ring-1"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.windows_devices.shards["shard_1"]
}

resource "microsoft365_graph_beta_group" "ring_2_broad" {
  display_name     = "Windows Updates - Ring 2 (30% Broad)"
  mail_nickname    = "win-updates-ring-2"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.windows_devices.shards["shard_2"]
}

resource "microsoft365_graph_beta_group" "ring_3_production" {
  display_name     = "Windows Updates - Ring 3 (50% Production)"
  mail_nickname    = "win-updates-ring-3"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.windows_devices.shards["shard_3"]
}

# Settings Catalog Policy for Ring 0 (immediate deployment)
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "ring_0_validation" {
  name               = "Windows Updates - Ring 0 (5% Validation)"
  description        = "Immediate deployment for validation devices"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"       = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId = "device_vendor_msft_policy_config_update_deferqualityupdatesperiodindays"
          choiceSettingValue = {
            value    = "device_vendor_msft_policy_config_update_deferqualityupdatesperiodindays_0"
            children = []
          }
        }
      }
    ]
  })

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_group.ring_0_validation.id
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_group.ring_1_pilot.id
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_group.ring_2_broad.id
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_group.ring_3_production.id
      filter_type = "none"
    }
  ]
}
```

#### Scenario 2: Group Members → Conditional Access Policy

```terraform
# Scenario 2: Group Members → Groups → Conditional Access Policy
# Use case: Roll out new CA policy to IT department in phases without nested groups

# Distribute existing IT department group members into 3 deployment rings
data "microsoft365_utility_guid_list_sharder" "it_dept_members" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc" # IT Department Group ID
  shard_count   = 3
  strategy      = "round-robin"
  seed          = "it-mfa-policy-2024"
}

# Create deployment ring groups from IT department members
resource "microsoft365_graph_beta_group" "it_ring_0_pilot" {
  display_name     = "IT MFA Rollout - Ring 0 (Pilot)"
  mail_nickname    = "it-mfa-ring-0"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.it_dept_members.shards["shard_0"]
}

resource "microsoft365_graph_beta_group" "it_ring_1_validation" {
  display_name     = "IT MFA Rollout - Ring 1 (Validation)"
  mail_nickname    = "it-mfa-ring-1"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.it_dept_members.shards["shard_1"]
}

resource "microsoft365_graph_beta_group" "it_ring_2_full" {
  display_name     = "IT MFA Rollout - Ring 2 (Full)"
  mail_nickname    = "it-mfa-ring-2"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.it_dept_members.shards["shard_2"]
}

# Single CA policy targeting all IT deployment rings
resource "microsoft365_graph_beta_conditional_access_policy" "it_mfa_policy" {
  display_name = "Require MFA - IT Department (Phased Rollout)"
  state        = "enabled"

  conditions {
    users {
      include_groups = [
        microsoft365_graph_beta_group.it_ring_0_pilot.id,
        microsoft365_graph_beta_group.it_ring_1_validation.id,
        microsoft365_graph_beta_group.it_ring_2_full.id
      ]
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
```

#### Scenario 3: Devices → Deployment Scheduler → Windows Quality Updates

```terraform
# Scenario 3: Devices → Groups → Deployment Scheduler → Windows Quality Update Policy
# Use case: Phased rollout of Windows quality updates with automated timing gates

# Step 1: Shard Windows devices into 3 deployment rings (10%, 30%, 60%)
data "microsoft365_utility_guid_list_sharder" "quality_update_rings" {
  resource_type     = "devices"
  odata_query       = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_percentages = [10, 30, 60]
  strategy          = "percentage"
  seed              = "quality-updates-2024"
}

# Step 2: Create Entra ID groups for each deployment ring
resource "microsoft365_graph_beta_group" "ring_0_pilot" {
  display_name     = "Quality Updates - Ring 0 (10% Pilot)"
  mail_nickname    = "quality-updates-ring-0"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_0"]
}

resource "microsoft365_graph_beta_group" "ring_1_broad" {
  display_name     = "Quality Updates - Ring 1 (30% Broad)"
  mail_nickname    = "quality-updates-ring-1"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_1"]
}

resource "microsoft365_graph_beta_group" "ring_2_production" {
  display_name     = "Quality Updates - Ring 2 (60% Production)"
  mail_nickname    = "quality-updates-ring-2"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_2"]
}

# Step 3: Define deployment timing gates for phased rollout
locals {
  deployment_start = "2026-01-20T00:00:00Z"
}

# Phase 1: Pilot ring opens after 24h
data "microsoft365_utility_deployment_scheduler" "ring_0_gate" {
  name                  = "quality-updates-ring-0-pilot"
  deployment_start_time = local.deployment_start
  scope_id              = microsoft365_graph_beta_group.ring_0_pilot.id

  time_condition = {
    delay_start_time_by = 24 # Open after 24 hours
  }

  # Only deploy during business hours
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "09:00:00"
        time_of_day_end   = "17:00:00"
      }
    ]
  }
}

# Phase 2: Broad ring opens 72h after pilot opens
data "microsoft365_utility_deployment_scheduler" "ring_1_gate" {
  name                  = "quality-updates-ring-1-broad"
  deployment_start_time = local.deployment_start
  scope_id              = microsoft365_graph_beta_group.ring_1_broad.id

  time_condition = {
    delay_start_time_by = 24
  }

  # Wait for pilot to be open for 72 hours
  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 24
    minimum_open_hours               = 72
  }

  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "08:00:00"
        time_of_day_end   = "18:00:00"
      }
    ]
  }
}

# Phase 3: Production ring opens 1 week after broad opens
data "microsoft365_utility_deployment_scheduler" "ring_2_gate" {
  name                  = "quality-updates-ring-2-production"
  deployment_start_time = local.deployment_start
  scope_id              = microsoft365_graph_beta_group.ring_2_production.id

  time_condition = {
    delay_start_time_by = 192 # 24 + 72 + 96 = 192 hours total
  }

  # Wait for broad ring to be open for 1 week
  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 96  # 24 + 72
    minimum_open_hours               = 168 # 1 week
  }

  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday"]
        time_of_day_start = "08:00:00"
        time_of_day_end   = "18:00:00"
      }
    ]
  }

  # Avoid Friday deployments to production
  exclusion_time_windows = {
    window = [
      {
        days_of_week = ["friday"]
      }
    ]
  }
}

# Step 4: Create Windows Quality Update Policy with conditional assignments
resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "phased_quality_updates" {
  display_name     = "Windows Quality Updates - Phased Rollout"
  description      = "Monthly quality updates deployed in phases with automated timing"
  hotpatch_enabled = true

  # Conditional assignments based on deployment scheduler gates
  # Only assign groups when their gates are open (released_scope_id != null)
  assignments = compact([
    data.microsoft365_utility_deployment_scheduler.ring_0_gate.released_scope_id != null ? {
      type     = "groupAssignmentTarget"
      group_id = data.microsoft365_utility_deployment_scheduler.ring_0_gate.released_scope_id
    } : null,
    data.microsoft365_utility_deployment_scheduler.ring_1_gate.released_scope_id != null ? {
      type     = "groupAssignmentTarget"
      group_id = data.microsoft365_utility_deployment_scheduler.ring_1_gate.released_scope_id
    } : null,
    data.microsoft365_utility_deployment_scheduler.ring_2_gate.released_scope_id != null ? {
      type     = "groupAssignmentTarget"
      group_id = data.microsoft365_utility_deployment_scheduler.ring_2_gate.released_scope_id
    } : null,
  ])
}

# Monitoring outputs
output "deployment_dashboard" {
  value = {
    device_distribution = {
      ring_0_pilot      = length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_0"])
      ring_1_broad      = length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_1"])
      ring_2_production = length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_2"])
      total_devices     = length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.quality_update_rings.shards["shard_2"])
    }

    deployment_gates = {
      ring_0_pilot = {
        status   = data.microsoft365_utility_deployment_scheduler.ring_0_gate.condition_met ? "OPEN" : "CLOSED"
        message  = data.microsoft365_utility_deployment_scheduler.ring_0_gate.status_message
        group_id = data.microsoft365_utility_deployment_scheduler.ring_0_gate.released_scope_id
      }
      ring_1_broad = {
        status   = data.microsoft365_utility_deployment_scheduler.ring_1_gate.condition_met ? "OPEN" : "CLOSED"
        message  = data.microsoft365_utility_deployment_scheduler.ring_1_gate.status_message
        group_id = data.microsoft365_utility_deployment_scheduler.ring_1_gate.released_scope_id
      }
      ring_2_production = {
        status   = data.microsoft365_utility_deployment_scheduler.ring_2_gate.condition_met ? "OPEN" : "CLOSED"
        message  = data.microsoft365_utility_deployment_scheduler.ring_2_gate.status_message
        group_id = data.microsoft365_utility_deployment_scheduler.ring_2_gate.released_scope_id
      }
    }

    active_assignments = length(compact([
      data.microsoft365_utility_deployment_scheduler.ring_0_gate.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.ring_1_gate.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.ring_2_gate.released_scope_id,
    ]))
  }
  description = "Comprehensive view of device distribution, gate status, and active policy assignments"
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
