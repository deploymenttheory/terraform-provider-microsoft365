# Advanced Deployment Scheduler Examples: Manual Override & Dependency Gates

locals {
  deployment_start = "2024-01-20T00:00:00Z"
}

# ============================================================================
# Example 1: Manual Override - Emergency Deployment
# ============================================================================

# Normal phased rollout with emergency override capability
data "microsoft365_utility_deployment_scheduler" "production_with_override" {
  name                  = "Production - Emergency Override"
  deployment_start_time = local.deployment_start
  scope_id              = "production-group-guid"

  time_condition {
    offset_hours = 168  # Normally wait 1 week
  }

  inclusion_time_windows {
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday"]
      time_of_day_start = "09:00:00"
      time_of_day_end   = "17:00:00"
    }
  }

  # Emergency override - set via terraform apply -var="emergency_deploy=true"
  manual_override = var.emergency_deploy
}

variable "emergency_deploy" {
  type        = bool
  default     = false
  description = "Set to true to bypass all conditions and immediately deploy"
}

# Behavior:
# - manual_override = false: Normal scheduling (wait 168h + office hours only)
# - manual_override = true:  IMMEDIATE deployment, ALL conditions bypassed
#   - Time condition: IGNORED
#   - Time windows: IGNORED
#   - Everything else: IGNORED
#   - released_scope_id = "production-group-guid" IMMEDIATELY

# Use case: Critical security patch needs emergency deployment
# Command: terraform apply -var="emergency_deploy=true"

# ============================================================================
# Example 2: Dependency Gates - Sequential Phased Rollout
# ============================================================================

# Phase 1: Pilot group (immediate)
data "microsoft365_utility_deployment_scheduler" "phase1_pilot" {
  name                  = "Phase 1 - Pilot"
  deployment_start_time = local.deployment_start
  scope_id              = "pilot-group-guid"

  time_condition {
    offset_hours = 0  # Immediate
  }
}

# Phase 2: Early adopters (after 48 hours)
data "microsoft365_utility_deployment_scheduler" "phase2_early_adopters" {
  name                  = "Phase 2 - Early Adopters"
  deployment_start_time = local.deployment_start
  scope_id              = "early-adopters-guid"

  time_condition {
    offset_hours = 48
  }

  inclusion_time_windows {
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
      time_of_day_start = "09:00:00"
      time_of_day_end   = "17:00:00"
    }
  }
}

# Phase 3: Production - Depends on Phase 2 being open for 72 hours
data "microsoft365_utility_deployment_scheduler" "phase3_production" {
  name                  = "Phase 3 - Production (Dependent)"
  deployment_start_time = local.deployment_start
  scope_id              = "production-guid"

  time_condition {
    offset_hours = 168  # Base: wait 1 week
  }

  # Dependency: Phase 2 must have been open for 72 hours
  depends_on_scheduler {
    prerequisite_offset_hours = 48   # Phase 2 opens at hour 48
    minimum_open_hours        = 72   # Must be open for 72 hours
  }

  # This gate won't open until hour 120 (48 + 72)
  # Even though time_condition is 168h, dependency forces earlier check
}

# Timeline:
# Hour 0:   Phase 1 opens ✓
# Hour 48:  Phase 2 opens ✓
# Hour 120: Phase 2 has been open for 72h
#           BUT Phase 3 time_condition (168h) not met yet ✗
# Hour 168: Phase 3 time_condition met (168h elapsed)
#           AND Phase 2 has been open for 120h (72h required) ✓
#           Phase 3 opens ✓

# ============================================================================
# Example 3: Complex Dependency Chain
# ============================================================================

# Phase 1: IT Department
data "microsoft365_utility_deployment_scheduler" "it_dept" {
  name                  = "IT Department"
  deployment_start_time = local.deployment_start
  scope_id              = "it-dept-guid"

  time_condition {
    offset_hours = 0
  }
}

# Phase 2: Tech-savvy users (wait 24h after IT)
data "microsoft365_utility_deployment_scheduler" "tech_users" {
  name                  = "Tech-savvy Users"
  deployment_start_time = local.deployment_start
  scope_id              = "tech-users-guid"

  time_condition {
    offset_hours = 24
  }

  depends_on_scheduler {
    prerequisite_offset_hours = 0    # IT opens at hour 0
    minimum_open_hours        = 24   # Must be open for 24h
  }
  # Opens at hour 24 (max of time_condition and dependency)
}

# Phase 3: Regular users (wait 48h after tech users)
data "microsoft365_utility_deployment_scheduler" "regular_users" {
  name                  = "Regular Users"
  deployment_start_time = local.deployment_start
  scope_id              = "regular-users-guid"

  time_condition {
    offset_hours = 120  # 5 days base
  }

  depends_on_scheduler {
    prerequisite_offset_hours = 24   # Tech users open at hour 24
    minimum_open_hours        = 48   # Must be open for 48h
  }
  # Opens at hour 120 (time_condition wins: 120 > 24+48)
}

# Phase 4: VIPs (wait 72h after regular users)
data "microsoft365_utility_deployment_scheduler" "vips" {
  name                  = "VIPs - Last to Deploy"
  deployment_start_time = local.deployment_start
  scope_id              = "vip-users-guid"

  time_condition {
    offset_hours = 240  # 10 days
  }

  depends_on_scheduler {
    prerequisite_offset_hours = 120  # Regular users open at hour 120
    minimum_open_hours        = 72   # Must be open for 72h
  }
  # Opens at hour 240 (time_condition wins: 240 > 120+72)
}

# ============================================================================
# Example 4: Dependency with Time Windows
# ============================================================================

# Phase 1: Weekend deployment
data "microsoft365_utility_deployment_scheduler" "weekend_pilot" {
  name                  = "Weekend Pilot"
  deployment_start_time = local.deployment_start
  scope_id              = "weekend-pilot-guid"

  time_condition {
    offset_hours = 0
  }

  inclusion_time_windows {
    window {
      days_of_week      = ["saturday", "sunday"]
      time_of_day_start = "08:00:00"
      time_of_day_end   = "20:00:00"
    }
  }
}

# Phase 2: Weekday deployment - depends on weekend pilot
data "microsoft365_utility_deployment_scheduler" "weekday_rollout" {
  name                  = "Weekday Rollout"
  deployment_start_time = local.deployment_start
  scope_id              = "weekday-group-guid"

  time_condition {
    offset_hours = 168  # 1 week
  }

  # Weekend pilot must have been open for 48 hours
  depends_on_scheduler {
    prerequisite_offset_hours = 0
    minimum_open_hours        = 48
  }

  inclusion_time_windows {
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday"]
      time_of_day_start = "09:00:00"
      time_of_day_end   = "17:00:00"
    }
  }
}

# ============================================================================
# Example 5: Emergency Override with Dependency
# ============================================================================

# Phase 1: Normal pilot
data "microsoft365_utility_deployment_scheduler" "pilot_normal" {
  name                  = "Pilot - Normal"
  deployment_start_time = local.deployment_start
  scope_id              = "pilot-guid"

  time_condition {
    offset_hours = 24
  }
}

# Phase 2: Production with emergency override
data "microsoft365_utility_deployment_scheduler" "prod_emergency_capable" {
  name                  = "Production - Emergency Capable"
  deployment_start_time = local.deployment_start
  scope_id              = "production-guid"

  time_condition {
    offset_hours = 168
  }

  # Normally depends on pilot being open for 48h
  depends_on_scheduler {
    prerequisite_offset_hours = 24
    minimum_open_hours        = 48
  }

  # Emergency override bypasses BOTH time condition AND dependency
  manual_override = var.skip_pilot_and_deploy_now
}

variable "skip_pilot_and_deploy_now" {
  type        = bool
  default     = false
  description = "Emergency: Deploy to production immediately, skipping pilot phase"
}

# Normal flow:
#   Hour 24:  Pilot opens
#   Hour 72:  Pilot has been open 48h
#   Hour 168: Production time_condition met, opens
#
# Emergency (manual_override = true):
#   Hour 0: Production opens IMMEDIATELY, pilot dependency ignored

# ============================================================================
# Real-World Example: Windows 11 Rollout with Dependencies
# ============================================================================

locals {
  win11_deployment_start = "2024-02-01T00:00:00Z"
}

# Phase 1: IT Department (immediate, no override needed)
data "microsoft365_utility_deployment_scheduler" "win11_it" {
  name                  = "Windows 11 - IT Dept"
  deployment_start_time = local.win11_deployment_start
  scope_id              = "it-dept-group-guid"

  time_condition {
    offset_hours = 0
  }
}

# Phase 2: Power users (after IT has run for 1 week)
data "microsoft365_utility_deployment_scheduler" "win11_power_users" {
  name                  = "Windows 11 - Power Users"
  deployment_start_time = local.win11_deployment_start
  scope_id              = "power-users-guid"

  time_condition {
    offset_hours = 168  # 1 week base
  }

  depends_on_scheduler {
    prerequisite_offset_hours = 0     # IT opens immediately
    minimum_open_hours        = 168   # Must run for 1 week
  }

  # Only during business hours
  inclusion_time_windows {
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday"]
      time_of_day_start = "10:00:00"
      time_of_day_end   = "14:00:00"
    }
  }
}

# Phase 3: General staff (after power users have run for 2 weeks)
data "microsoft365_utility_deployment_scheduler" "win11_general_staff" {
  name                  = "Windows 11 - General Staff"
  deployment_start_time = local.win11_deployment_start
  scope_id              = "general-staff-guid"

  time_condition {
    offset_hours = 504  # 3 weeks base
  }

  depends_on_scheduler {
    prerequisite_offset_hours = 168   # Power users open at week 1
    minimum_open_hours        = 336   # Must run for 2 weeks
  }

  inclusion_time_windows {
    window {
      days_of_week      = ["tuesday", "wednesday", "thursday"]
      time_of_day_start = "09:00:00"
      time_of_day_end   = "17:00:00"
    }
  }

  # Emergency override if issues found
  manual_override = var.emergency_full_rollout
}

variable "emergency_full_rollout" {
  type    = bool
  default = false
}

# Use in Intune policy
resource "microsoft365_windows_feature_update_policy" "win11_rollout" {
  name = "Windows 11 24H2 Rollout"

  # All phases in one policy - gates control release
  dynamic "assignments" {
    for_each = compact([
      data.microsoft365_utility_deployment_scheduler.win11_it.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.win11_power_users.released_scope_id,
      data.microsoft365_utility_deployment_scheduler.win11_general_staff.released_scope_id,
    ])

    content {
      type     = "groupAssignmentTarget"
      group_id = assignments.value
    }
  }
}

# Monitor deployment status
output "win11_rollout_status" {
  value = {
    it_dept = {
      ready  = data.microsoft365_utility_deployment_scheduler.win11_it.condition_met
      status = data.microsoft365_utility_deployment_scheduler.win11_it.status_message
    }
    power_users = {
      ready  = data.microsoft365_utility_deployment_scheduler.win11_power_users.condition_met
      status = data.microsoft365_utility_deployment_scheduler.win11_power_users.status_message
    }
    general_staff = {
      ready  = data.microsoft365_utility_deployment_scheduler.win11_general_staff.condition_met
      status = data.microsoft365_utility_deployment_scheduler.win11_general_staff.status_message
    }
  }
}

# ============================================================================
# Understanding Status Messages
# ============================================================================

# Normal dependency flow:
# "Waiting: Time condition met (200h/168h required), Inclusion window satisfied (within inclusion window 1), No exclusion window active (outside all exclusion windows), Dependency NOT satisfied (prerequisite open for 20h/48h required)"

# Dependency satisfied:
# "Conditions met: Time condition met (200h/168h required), Inclusion window satisfied (within inclusion window 1), No exclusion window active (outside all exclusion windows), Dependency satisfied (prerequisite open for 50h/48h required)"

# Manual override:
# "Manual override enabled: Gate forced open, all conditions bypassed"
