---
page_title: "microsoft365_utility_deployment_scheduler Data Source - terraform-provider-microsoft365"
subcategory: "Utility"

description: |-
  A conditional gate valve for phased deployments. Returns scope ID(s) when specified time-based conditions are met, enabling controlled rollout of policies and updates. All times are in UTC (RFC3339 format). When conditions are not met, returns null, preventing deployment to target groups.
  This datasource is evaluated on every Terraform plan/apply, allowing gates to automatically open when conditions are satisfied.
---

# microsoft365_utility_deployment_scheduler

A conditional gate valve for phased deployments that returns scope ID(s) only when specified time-based conditions are met.

This datasource enables controlled, time-based rollout of Microsoft 365 policies and configurations by acting as a deployment gate. 
It evaluates multiple conditions (time delays, business hours, maintenance windows, dependencies) and releases scope IDs 
(typically Entra ID Group GUIDs) only when all conditions are satisfied. When conditions are not met, it returns null, 
preventing deployment to target groups.

## Background

The datasource is evaluated on every `terraform plan` and `terraform apply`. When conditions change and the gate opens,
the next apply will automatically assign policies to the released groups. This allows you to define your entire deployment
timeline upfront, and Terraform will progressively roll out changes as time passes.

**Key Concept**: Think of it as a time-locked vault. You set the conditions for when the vault opens, and it automatically
releases the contents (scope IDs) when those conditions are met.

## Use Cases

- **Phased Rollouts**: Deploy to Pilot → IT → Production with automatic timing and dependency chaining
- **Business Hours Only**: Restrict deployments to support hours (e.g., Monday-Friday 9 AM - 5 PM UTC)
- **Maintenance Windows**: Avoid deployments during weekends, nightly maintenance, or holiday freezes
- **Dependency Chaining**: Wait for previous phase to stabilize before proceeding to next phase
- **Emergency Releases**: Manual override for critical patches bypassing all scheduled conditions
- **Time-Limited Pilots**: Automatically expire pilot deployments after a specified duration

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.41.0-alpha | Experimental | Initial release with time conditions, windows, and dependency gates |

## Example Usage

### Basic Time Delay

```terraform
# Basic Time Delay Example
# This example demonstrates a simple time-based delay before releasing a scope ID.
# The gate opens 24 hours after deployment_start_time.

data "microsoft365_utility_deployment_scheduler" "pilot_group_release" {
  name                  = "pilot-group-24h-delay"
  deployment_start_time = "2026-01-20T08:00:00Z"
  scope_id              = "a1b2c3d4-e5f6-7890-abcd-ef1234567890" # Entra ID Group ID for pilot users

  time_condition = {
    delay_start_time_by = 24 # Wait 24 hours before releasing
  }
}

# Use the released scope ID in your resource
# Example: Conditional Access policy is only assigned when gate opens
output "pilot_gate_status" {
  value = {
    condition_met  = data.microsoft365_utility_deployment_scheduler.pilot_group_release.condition_met
    released_scope = data.microsoft365_utility_deployment_scheduler.pilot_group_release.released_scope_id
    status_message = data.microsoft365_utility_deployment_scheduler.pilot_group_release.status_message
  }
  description = "Monitor gate status and timing"
}
```

### Business Hours Only

```terraform
# Business Hours Only Deployment
# This example demonstrates releasing changes only during business hours.
# Useful for changes that might require user interaction or support availability.

data "microsoft365_utility_deployment_scheduler" "business_hours_release" {
  name                  = "settings-catalog-business-hours"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = "b2c3d4e5-f6a7-8901-bcde-f12345678901" # Entra ID Group for workstations

  time_condition = {
    delay_start_time_by = 0 # No delay, but constrained by inclusion windows
  }

  # Only allow deployment during business hours (UTC)
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "09:00:00" # 9 AM UTC
        time_of_day_end   = "17:00:00" # 5 PM UTC
      }
    ]
  }
}

output "business_hours_status" {
  value = {
    gate_open          = data.microsoft365_utility_deployment_scheduler.business_hours_release.condition_met
    current_time_valid = data.microsoft365_utility_deployment_scheduler.business_hours_release.status_message
    released_scope     = data.microsoft365_utility_deployment_scheduler.business_hours_release.released_scope_id
  }
  description = "Shows if we're currently in business hours deployment window"
}
```

### Exclude Maintenance Windows

```terraform
# Exclude Maintenance Windows
# This example demonstrates blocking deployments during maintenance windows.
# Prevents changes during critical business periods or scheduled maintenance.

data "microsoft365_utility_deployment_scheduler" "avoid_maintenance" {
  name                  = "app-deployment-avoid-maintenance"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = "c3d4e5f6-a7b8-9012-cdef-123456789012" # Device group

  time_condition = {
    delay_start_time_by = 0 # Deploy immediately when conditions allow
  }

  # Block deployments during maintenance windows
  exclusion_time_windows = {
    window = [
      {
        # Block during weekend maintenance
        days_of_week = ["saturday", "sunday"]
      },
      {
        # Block during year-end freeze (absolute dates)
        date_start = "2026-12-20T00:00:00Z"
        date_end   = "2027-01-05T23:59:59Z"
      },
      {
        # Block during nightly maintenance window
        time_of_day_start = "22:00:00" # 10 PM UTC
        time_of_day_end   = "06:00:00" # 6 AM UTC
      }
    ]
  }
}

output "maintenance_window_check" {
  value = {
    deployment_allowed = data.microsoft365_utility_deployment_scheduler.avoid_maintenance.condition_met
    status             = data.microsoft365_utility_deployment_scheduler.avoid_maintenance.status_message
    released_scope     = data.microsoft365_utility_deployment_scheduler.avoid_maintenance.released_scope_id
  }
  description = "Indicates if we're outside maintenance windows"
}
```

### Phased Rollout with Dependencies

```terraform
# Phased Rollout with Dependency Chaining
# This example demonstrates a multi-phase rollout where each phase depends on the previous.
# Phase 1: Pilot users (opens after 24h)
# Phase 2: Early adopters (opens after pilot has been open for 48h)
# Phase 3: Production users (opens after early adopters has been open for 72h)

locals {
  deployment_start = "2026-01-20T08:00:00Z"
  pilot_group_id   = "d4e5f6a7-b8c9-0123-def0-123456789abc"
  early_adopter_id = "e5f6a7b8-c9d0-1234-ef01-23456789abcd"
  production_id    = "f6a7b8c9-d0e1-2345-f012-3456789abcde"
}

# Phase 1: Pilot Group (opens after 24h delay)
data "microsoft365_utility_deployment_scheduler" "phase1_pilot" {
  name                  = "mfa-policy-phase1-pilot"
  deployment_start_time = local.deployment_start
  scope_id              = local.pilot_group_id

  time_condition = {
    delay_start_time_by = 24 # 24h soak time for initial testing
  }
}

# Phase 2: Early Adopters (opens after pilot has been open for 48h)
data "microsoft365_utility_deployment_scheduler" "phase2_early_adopters" {
  name                  = "mfa-policy-phase2-early-adopters"
  deployment_start_time = local.deployment_start
  scope_id              = local.early_adopter_id

  time_condition = {
    delay_start_time_by = 24 # Same base delay as pilot
  }

  # Wait for pilot phase to be open for 48h before opening this gate
  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 24 # Pilot's delay
    minimum_open_hours               = 48 # Wait 48h after pilot opens
  }
}

# Phase 3: Production (opens after early adopters has been open for 72h)
data "microsoft365_utility_deployment_scheduler" "phase3_production" {
  name                  = "mfa-policy-phase3-production"
  deployment_start_time = local.deployment_start
  scope_id              = local.production_id

  time_condition = {
    delay_start_time_by = 72 # Base delay to align with early adopters (24 + 48)
  }

  # Wait for early adopters phase to be open for 72h before opening this gate
  depends_on_scheduler = {
    prerequisite_delay_start_time_by = 72 # Early adopters effective delay (24 + 48)
    minimum_open_hours               = 72 # Wait 72h after early adopters opens
  }
}

# Monitoring outputs
output "rollout_status" {
  value = {
    phase1_pilot = {
      open   = data.microsoft365_utility_deployment_scheduler.phase1_pilot.condition_met
      scope  = data.microsoft365_utility_deployment_scheduler.phase1_pilot.released_scope_id
      status = data.microsoft365_utility_deployment_scheduler.phase1_pilot.status_message
    }
    phase2_early_adopters = {
      open   = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.condition_met
      scope  = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.released_scope_id
      status = data.microsoft365_utility_deployment_scheduler.phase2_early_adopters.status_message
    }
    phase3_production = {
      open   = data.microsoft365_utility_deployment_scheduler.phase3_production.condition_met
      scope  = data.microsoft365_utility_deployment_scheduler.phase3_production.released_scope_id
      status = data.microsoft365_utility_deployment_scheduler.phase3_production.status_message
    }
  }
  description = "Current status of all rollout phases"
}
```

### Manual Override for Emergency Release

```terraform
# Manual Override Example
# This example demonstrates using manual_override to immediately release a deployment
# regardless of time conditions. Useful for emergency changes or hotfixes.

# Normal scheduled deployment
data "microsoft365_utility_deployment_scheduler" "security_patch_release" {
  name                  = "critical-security-patch"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = "a7b8c9d0-e1f2-3456-a012-3456789abcde"

  # Normal schedule: 7 day delay
  time_condition = {
    delay_start_time_by = 168 # 7 days for normal patches
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

  # Manual override bypasses ALL conditions when set to true
  # To activate emergency release: change false to true and apply
  manual_override = false
}

output "patch_deployment_status" {
  value = {
    gate_open      = data.microsoft365_utility_deployment_scheduler.security_patch_release.condition_met
    scope_released = data.microsoft365_utility_deployment_scheduler.security_patch_release.released_scope_id
    status         = data.microsoft365_utility_deployment_scheduler.security_patch_release.status_message
  }
  description = "Monitor patch deployment status"
}

# Usage:
# Normal deployment: Leave manual_override = false (follows schedule)
# Emergency release: Change manual_override = true and run terraform apply
```

### Combined Conditions

```terraform
# Combined Conditions Example
# This example demonstrates combining multiple scheduling constraints:
# - Time delay with absolute deadlines
# - Business hours only
# - Exclude maintenance windows
# Useful for complex deployment scenarios with multiple requirements.

data "microsoft365_utility_deployment_scheduler" "complex_release" {
  name                  = "compliance-policy-complex-schedule"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = "b8c9d0e1-f2a3-4567-b123-456789abcdef"

  # Time-based constraints with deadlines
  time_condition = {
    delay_start_time_by     = 48                     # Wait 48 hours minimum
    absolute_earliest       = "2026-01-22T00:00:00Z" # Don't release before this date
    absolute_latest         = "2026-02-15T23:59:59Z" # Must release before this date (compliance deadline)
    max_open_duration_hours = 168                    # Auto-close gate after 7 days (1 week deployment window)
  }

  # Only deploy during business hours on weekdays
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
        time_of_day_start = "09:00:00"
        time_of_day_end   = "18:00:00"
      }
    ]
  }

  # Exclude sensitive periods
  exclusion_time_windows = {
    window = [
      {
        # Exclude weekends
        days_of_week = ["saturday", "sunday"]
      },
      {
        # Exclude nightly maintenance
        time_of_day_start = "22:00:00"
        time_of_day_end   = "06:00:00"
      },
      {
        # Exclude specific high-activity period (month-end)
        date_start = "2026-01-29T00:00:00Z"
        date_end   = "2026-02-02T23:59:59Z"
      }
    ]
  }

  # All conditions must be met
  require_all_conditions = true
}

# Detailed monitoring output
output "complex_schedule_status" {
  value = {
    gate_status = {
      open           = data.microsoft365_utility_deployment_scheduler.complex_release.condition_met
      status_message = data.microsoft365_utility_deployment_scheduler.complex_release.status_message
      scope_released = data.microsoft365_utility_deployment_scheduler.complex_release.released_scope_id
    }

    conditions_detail = data.microsoft365_utility_deployment_scheduler.complex_release.conditions_detail

    timing = {
      deployment_start = data.microsoft365_utility_deployment_scheduler.complex_release.deployment_start_time
      window_closes    = "2026-02-15T23:59:59Z" # absolute_latest
    }
  }
  description = "Comprehensive status of the complex deployment schedule"
}
```

### Multiple Scope IDs

```terraform
# Multiple Scope IDs Example
# This example demonstrates releasing multiple group IDs simultaneously.
# Useful when you want to apply changes to several groups at once with the same timing.

locals {
  regional_offices = [
    "c9d0e1f2-a3b4-5678-c234-56789abcdef0", # US East office
    "d0e1f2a3-b4c5-6789-d345-6789abcdef01", # US West office
    "e1f2a3b4-c5d6-7890-e456-789abcdef012", # EMEA office
    "f2a3b4c5-d6e7-8901-f567-89abcdef0123", # APAC office
  ]
}

data "microsoft365_utility_deployment_scheduler" "regional_rollout" {
  name                  = "regional-offices-unified-rollout"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_ids             = local.regional_offices # Use scope_ids for multiple groups

  time_condition = {
    delay_start_time_by = 72 # 3 day delay for all regions
  }

  # Deploy only during business hours (respects all timezones effectively)
  inclusion_time_windows = {
    window = [
      {
        days_of_week      = ["monday", "tuesday", "wednesday", "thursday"]
        time_of_day_start = "10:00:00" # Conservative window to catch multiple timezones
        time_of_day_end   = "15:00:00"
      }
    ]
  }
}

output "regional_rollout_status" {
  value = {
    gate_open          = data.microsoft365_utility_deployment_scheduler.regional_rollout.condition_met
    number_of_regions  = length(local.regional_offices)
    released_group_ids = data.microsoft365_utility_deployment_scheduler.regional_rollout.released_scope_ids
    status             = data.microsoft365_utility_deployment_scheduler.regional_rollout.status_message
  }
  description = "Status of unified regional rollout"
}
```

### Real-World Settings Catalog Deployment

```terraform
# Real-World Settings Catalog Deployment
# This example demonstrates a production-ready Settings Catalog policy deployment
# with proper phasing, timing, and safety controls for Windows security settings.

locals {
  policy_name = "Windows Security Baseline v2.0"

  # Define deployment groups
  pilot_users_group       = "a1a2a3a4-b5b6-c7c8-d9d0-e1e2e3e4e5e6" # 5-10 pilot devices
  it_department_group     = "b2b3b4b5-c6c7-d8d9-e0e1-f2f3f4f5f6f7" # IT team for validation
  production_workstations = "c3c4c5c6-d7d8-e9e0-f1f2-a3a4a5a6a7a8" # All production Windows devices
}

# Phase 1: Pilot Users - Opens after 24h, business hours only
data "microsoft365_utility_deployment_scheduler" "pilot_phase" {
  name                  = "${local.policy_name}-pilot"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = local.pilot_users_group

  time_condition = {
    delay_start_time_by = 24 # 24h delay for safety
  }

  # Business hours only for pilot
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

# Phase 2: IT Department - Opens 72h after pilot opens
data "microsoft365_utility_deployment_scheduler" "it_phase" {
  name                  = "${local.policy_name}-it-dept"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = local.it_department_group

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

# Phase 3: Production - Opens 1 week after IT phase opens
data "microsoft365_utility_deployment_scheduler" "production_phase" {
  name                  = "${local.policy_name}-production"
  deployment_start_time = "2026-01-20T00:00:00Z"
  scope_id              = local.production_workstations

  time_condition = {
    delay_start_time_by     = 192                    # 24 + 72 + 96 (total: 8 days)
    absolute_latest         = "2026-03-01T23:59:59Z" # Must complete by end of Feb (compliance)
    max_open_duration_hours = 720                    # Auto-close after 30 days
  }

  # Wait for IT phase to be open for 1 week
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

  exclusion_time_windows = {
    window = [
      {
        # No Friday deployments to prod
        days_of_week = ["friday"]
      },
      {
        # Month-end freeze
        date_start = "2026-01-29T00:00:00Z"
        date_end   = "2026-02-02T23:59:59Z"
      }
    ]
  }
}

# Monitoring and Reporting
output "deployment_dashboard" {
  value = {
    policy_name = local.policy_name

    phase_1_pilot = {
      status        = data.microsoft365_utility_deployment_scheduler.pilot_phase.condition_met ? "OPEN" : "CLOSED"
      group_id      = data.microsoft365_utility_deployment_scheduler.pilot_phase.released_scope_id
      status_detail = data.microsoft365_utility_deployment_scheduler.pilot_phase.status_message
    }

    phase_2_it = {
      status        = data.microsoft365_utility_deployment_scheduler.it_phase.condition_met ? "OPEN" : "CLOSED"
      group_id      = data.microsoft365_utility_deployment_scheduler.it_phase.released_scope_id
      status_detail = data.microsoft365_utility_deployment_scheduler.it_phase.status_message
    }

    phase_3_production = {
      status        = data.microsoft365_utility_deployment_scheduler.production_phase.condition_met ? "OPEN" : "CLOSED"
      group_id      = data.microsoft365_utility_deployment_scheduler.production_phase.released_scope_id
      status_detail = data.microsoft365_utility_deployment_scheduler.production_phase.status_message
    }

    total_groups_deployed = length(compact([
      data.microsoft365_utility_deployment_scheduler.pilot_phase.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.it_phase.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.production_phase.released_scope_id,
    ]))
  }

  description = "Complete rollout status dashboard for Settings Catalog policy deployment"
}

# Example of how you'd use this in a real Settings Catalog policy
# (Schema shown is illustrative - adjust to match actual resource)
#
# resource "microsoft365_graph_beta_device_management_configuration_policy" "security_baseline" {
#   name         = local.policy_name
#   description  = "Windows Security Baseline settings for corporate devices"
#   platforms    = "windows10"
#   technologies = ["mdm"]
#   
#   # Assignments based on which gates are open
#   assignments = compact([
#     data.microsoft365_utility_deployment_scheduler.pilot_phase.released_scope_id != null ? {
#       target = {
#         group_id = data.microsoft365_utility_deployment_scheduler.pilot_phase.released_scope_id
#       }
#     } : null,
#     data.microsoft365_utility_deployment_scheduler.it_phase.released_scope_id != null ? {
#       target = {
#         group_id = data.microsoft365_utility_deployment_scheduler.it_phase.released_scope_id
#       }
#     } : null,
#     data.microsoft365_utility_deployment_scheduler.production_phase.released_scope_id != null ? {
#       target = {
#         group_id = data.microsoft365_utility_deployment_scheduler.production_phase.released_scope_id
#       }
#     } : null,
#   ])
#   
#   # ... settings configuration ...
# }
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) A descriptive name for this deployment phase (e.g., 'Phase 2 - Production Rollout'). Used in status messages.

### Optional

- `depends_on_scheduler` (Attributes) Dependency gate that requires another scheduler to have been open for a minimum duration before this gate can open. Useful for sequential phased rollouts where Phase 3 shouldn't start until Phase 2 has been running successfully for a period. This gate calculates when the prerequisite would have opened and ensures sufficient time has passed. (see [below for nested schema](#nestedatt--depends_on_scheduler))
- `deployment_start_time` (String) The deployment campaign start date and time in UTC (RFC3339 format, e.g., '2024-01-15T00:00:00Z'). All time-based conditions are calculated relative to this timestamp, similar to how Unix epoch time works as a reference point. If not provided, uses the current time on each evaluation (not recommended for time-based conditions). Explicitly setting this value allows coordinating multiple deployment phases to a single campaign start time.
- `exclusion_time_windows` (Attributes) Defines one or more time windows when the gate must remain closed, even if other conditions are met. Use this for holiday freezes, blackout periods, etc. Exclusions take precedence over inclusions. Multiple windows are evaluated with OR logic (any window matches = deployment blocked). (see [below for nested schema](#nestedatt--exclusion_time_windows))
- `inclusion_time_windows` (Attributes) Defines one or more time windows when the gate is allowed to open. The current time must fall within at least one of the defined windows for the gate to open. Use this for office hours restrictions, maintenance windows, etc. Multiple windows are evaluated with OR logic (any window matches = condition passes). (see [below for nested schema](#nestedatt--inclusion_time_windows))
- `manual_override` (Boolean) Emergency override to immediately release scope ID(s), bypassing all time conditions and windows. Set to `true` to force-release the gate. Useful for emergency deployments. When enabled, all other conditions are ignored and scope ID(s) are immediately released. Default: false.
- `require_all_conditions` (Boolean) When true, all specified conditions must be met (AND logic). When false, any condition passing will release scope ID(s) (OR logic). Defaults to true. Reserved for future use when multiple condition types are supported.
- `scope_id` (String) A single scope ID (typically a group ID) to release when conditions are met. Use this when deploying to one group. Use either `scope_id` or `scope_ids`, not both.
- `scope_ids` (List of String) List of multiple scope IDs (user GUIDs, device GUIDs, etc.) to release when conditions are met. Use this when deploying to multiple individual entities. Use either `scope_id` or `scope_ids`, not both.
- `time_condition` (Attributes) Time-based condition that must be satisfied before releasing scope ID(s). (see [below for nested schema](#nestedatt--time_condition))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `condition_met` (Boolean) Boolean indicating whether all required conditions are satisfied. True means the gate is open and scope IDs are released.
- `conditions_detail` (Attributes) Detailed breakdown of each condition's evaluation for debugging and monitoring. (see [below for nested schema](#nestedatt--conditions_detail))
- `id` (String) The unique identifier of this deployment scheduler instance.
- `released_scope_id` (String) The single scope ID released by this gate when conditions are met. Returns the value from `scope_id` when the gate opens, or null when conditions are not met. Use this in policy assignments when you provided `scope_id`.
- `released_scope_ids` (List of String) List of scope IDs released by this gate when conditions are met. Returns the full `scope_ids` list when the gate opens, or null when conditions are not met. Use this in policy assignments when you provided `scope_ids`.
- `status_message` (String) Human-readable status message describing the current state of all conditions. Visible in Terraform plan output. Examples:
- `Conditions met: Time condition met (50h/48h required)`
- `Waiting: Time condition not met (22h/48h required)`

<a id="nestedatt--depends_on_scheduler"></a>
### Nested Schema for `depends_on_scheduler`

Required:

- `minimum_open_hours` (Number) Minimum number of hours the prerequisite gate must have been open before this gate can open. For example, if set to 48, this gate won't open until 48 hours after the prerequisite gate opened. Must be >= 0.
- `prerequisite_delay_start_time_by` (Number) The delay_start_time_by value of the prerequisite scheduler. This is when the prerequisite gate would have opened (hours after deployment_start_time). For example, if Phase 2 has `delay_start_time_by = 168`, set this to 168.


<a id="nestedatt--exclusion_time_windows"></a>
### Nested Schema for `exclusion_time_windows`

Required:

- `window` (Attributes List) List of time windows when deployment is blocked. (see [below for nested schema](#nestedatt--exclusion_time_windows--window))

<a id="nestedatt--exclusion_time_windows--window"></a>
### Nested Schema for `exclusion_time_windows.window`

Optional:

- `date_end` (String) Absolute end date/time in UTC (RFC3339 format, e.g., '2025-01-05T23:59:59Z'). Use for specific date ranges. If not specified, no end date limit.
- `date_start` (String) Absolute start date/time in UTC (RFC3339 format, e.g., '2024-12-20T00:00:00Z'). Use for specific date ranges like holiday freezes. If not specified, no start date limit.
- `days_of_week` (List of String) Days of the week when this window blocks deployment. Valid values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. If not specified, all days are included.
- `time_of_day_end` (String) End time in UTC (HH:MM:SS format, e.g., '23:59:59'). If not specified, ends at 23:59:59.
- `time_of_day_start` (String) Start time in UTC (HH:MM:SS format, e.g., '00:00:00'). If not specified, starts at 00:00:00.



<a id="nestedatt--inclusion_time_windows"></a>
### Nested Schema for `inclusion_time_windows`

Required:

- `window` (Attributes List) List of time windows when deployment is allowed. (see [below for nested schema](#nestedatt--inclusion_time_windows--window))

<a id="nestedatt--inclusion_time_windows--window"></a>
### Nested Schema for `inclusion_time_windows.window`

Optional:

- `date_end` (String) Absolute end date/time in UTC (RFC3339 format, e.g., '2024-01-31T23:59:59Z'). Use for specific date ranges. If not specified, no end date limit.
- `date_start` (String) Absolute start date/time in UTC (RFC3339 format, e.g., '2024-01-15T00:00:00Z'). Use for specific date ranges. If not specified, no start date limit.
- `days_of_week` (List of String) Days of the week when this window is active. Valid values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. If not specified, all days are included.
- `time_of_day_end` (String) End time in UTC (HH:MM:SS format, e.g., '17:00:00'). If not specified, ends at 23:59:59.
- `time_of_day_start` (String) Start time in UTC (HH:MM:SS format, e.g., '09:00:00'). If not specified, starts at 00:00:00.



<a id="nestedatt--time_condition"></a>
### Nested Schema for `time_condition`

Required:

- `delay_start_time_by` (Number) Number of hours to delay after `deployment_start_time` before allowing the gate to open. Must be >= 0. Set to 0 for immediate release at deployment start time.

Optional:

- `absolute_earliest` (String) Absolute earliest time (UTC RFC3339 format) when gate can open, regardless of `delay_start_time_by`. Use this to prevent deployment before a specific date/time (e.g., wait for Patch Tuesday). If specified, gate cannot open before this time even if delay_start_time_by has elapsed.
- `absolute_latest` (String) Absolute deadline (UTC RFC3339 format) when gate must close and will never open again. Use this for time-limited deployment campaigns or change freeze deadlines. If current time exceeds this, gate permanently closes.
- `max_open_duration_hours` (Number) Maximum number of hours the gate can remain open after it first opens. Must be >= 0. Use this for pilot/temporary deployments that should automatically expire (e.g., 2-week pilot = 336 hours). When duration expires, gate auto-closes and scope IDs are retracted. Set to 0 for unlimited duration (default behavior).


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--conditions_detail"></a>
### Nested Schema for `conditions_detail`

Read-Only:

- `time_condition_detail` (Attributes) Detailed time condition evaluation. (see [below for nested schema](#nestedatt--conditions_detail--time_condition_detail))

<a id="nestedatt--conditions_detail--time_condition_detail"></a>
### Nested Schema for `conditions_detail.time_condition_detail`

Read-Only:

- `condition_met` (Boolean) Whether the time condition is satisfied.
- `current_time` (String) The current evaluation time (UTC RFC3339).
- `delay_start_time_by` (Number) Required delay in hours from deployment start time.
- `deployment_start_time` (String) The deployment start time used for calculations (UTC RFC3339).
- `hours_elapsed` (Number) Hours elapsed since deployment start time.
- `required` (Boolean) Whether a time condition was specified.

## Best Practices

1. **Always Set deployment_start_time** - Provides consistent timing across terraform applies. Without it, the current time is used on each evaluation, causing unpredictable behavior.

2. **Use UTC Times** - All times must be in UTC. Plan your windows accordingly if your organization spans multiple timezones.

3. **Test with Pilot Groups** - Always start with a small pilot group before expanding to production.

4. **Handle Null Values** - Gates return `null` when closed. Use conditional logic in resource blocks:
   ```hcl
   assignments = data.scheduler.released_scope_id != null ? [{
     target = { group_id = data.scheduler.released_scope_id }
   }] : []
   ```

5. **Monitor with Outputs** - Use outputs to track gate status without modifying state:
   ```hcl
   output "deployment_status" {
     value = {
       gate_open = data.scheduler.condition_met
       status    = data.scheduler.status_message
       scope     = data.scheduler.released_scope_id
     }
   }
   ```

6. **Use Lowercase Day Names** - Days of week must be lowercase: `"monday"`, `"tuesday"`, etc.

## Important Notes

### No Authentication Required

This datasource performs local time-based calculations and doesn't make any API calls. Unlike other datasources in this provider, it doesn't require Microsoft Graph API credentials.

### Evaluation Timing

The datasource is evaluated during `terraform plan` and `terraform apply`. Gates automatically open when conditions are met, allowing progressive rollout without manual intervention.

### Null Behavior

When conditions are not met, `released_scope_id` and `released_scope_ids` return `null`. This is intentional behavior to prevent premature deployment. Always use conditional logic when referencing these values in resource configurations.

### Time Calculations

All time calculations are performed relative to `deployment_start_time`. For dependency gates, the datasource calculates when prerequisite schedulers would have opened based on their `delay_start_time_by` values.

## Additional Resources

- [Conditional Access Best Practices](https://learn.microsoft.com/en-us/entra/identity/conditional-access/plan-conditional-access)
- [Microsoft Graph API - Groups](https://learn.microsoft.com/en-us/graph/api/resources/group)
- [Intune Configuration Policies](https://learn.microsoft.com/en-us/mem/intune/configuration/device-profile-create)
