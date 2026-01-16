# Example: Deployment Scheduler with Time Windows (Inclusion & Exclusion)

locals {
  deployment_start = "2024-01-20T00:00:00Z"
}

# ============================================================================
# Example 1: Office Hours Only (Weekdays 9am-5pm UTC)
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "office_hours" {
  name                  = "Office Hours Deployment"
  deployment_start_time = local.deployment_start
  scope_id              = "pilot-group-123"

  time_condition {
    offset_hours = 24  # Wait 24h from deployment start
  }

  # Gate only opens during business hours
  inclusion_time_windows {
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
      time_of_day_start = "09:00:00"  # 9 AM UTC
      time_of_day_end   = "17:00:00"  # 5 PM UTC
    }
  }
}

# Behavior:
# - Base condition: Wait 24h from deployment start
# - Additional restriction: Only deploy if current time is Mon-Fri 9am-5pm UTC
# - If terraform apply runs at 3pm Friday = gate opens
# - If terraform apply runs at 11pm Friday = gate stays closed (outside window)
# - If terraform apply runs at 10am Saturday = gate stays closed (wrong day)

# ============================================================================
# Example 2: Maintenance Window Only (Weekend Early Morning)
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "maintenance_window" {
  name                  = "Maintenance Window Deployment"
  deployment_start_time = local.deployment_start
  scope_id              = "production-group-456"

  time_condition {
    offset_hours = 48
  }

  # Only deploy during weekend maintenance window
  inclusion_time_windows {
    window {
      days_of_week      = ["saturday", "sunday"]
      time_of_day_start = "02:00:00"  # 2 AM UTC
      time_of_day_end   = "06:00:00"  # 6 AM UTC
    }
  }
}

# Behavior:
# - Base condition: Wait 48h from deployment start
# - Additional restriction: Only Sat/Sun 2am-6am UTC
# - Perfect for low-impact maintenance windows

# ============================================================================
# Example 3: Block Holiday Freeze Period
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "freeze_aware" {
  name                  = "Holiday Freeze Aware"
  deployment_start_time = "2024-12-01T00:00:00Z"
  scope_id              = "critical-apps-789"

  time_condition {
    offset_hours = 24
  }

  # Block deployments during holiday freeze
  exclusion_time_windows {
    window {
      date_start = "2024-12-20T00:00:00Z"
      date_end   = "2025-01-05T23:59:59Z"
    }
  }
}

# Behavior:
# - Base condition: Wait 24h from deployment start
# - Exclusion: Block if current time is between Dec 20 - Jan 5
# - Gate will NOT open during the freeze period, regardless of other conditions

# ============================================================================
# Example 4: Block Weekends
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "no_weekends" {
  name                  = "No Weekend Deployments"
  deployment_start_time = local.deployment_start
  scope_id              = "business-apps-111"

  time_condition {
    offset_hours = 12
  }

  # Don't deploy on weekends
  exclusion_time_windows {
    window {
      days_of_week = ["saturday", "sunday"]
    }
  }
}

# Behavior:
# - Base condition: Wait 12h from deployment start
# - Exclusion: Block if current day is Saturday or Sunday
# - Will skip weekends and wait for Monday

# ============================================================================
# Example 5: Multiple Inclusion Windows (OR logic)
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "flexible_windows" {
  name                  = "Flexible Time Windows"
  deployment_start_time = local.deployment_start
  scope_id              = "flexible-group-222"

  time_condition {
    offset_hours = 24
  }

  # Deploy during EITHER office hours OR maintenance window
  inclusion_time_windows {
    # Window 1: Weekday office hours
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
      time_of_day_start = "09:00:00"
      time_of_day_end   = "17:00:00"
    }

    # Window 2: Weekend maintenance
    window {
      days_of_week      = ["saturday", "sunday"]
      time_of_day_start = "02:00:00"
      time_of_day_end   = "06:00:00"
    }
  }
}

# Behavior:
# - If ANY window matches, gate can open
# - Monday 10am = matches window 1 ✓
# - Saturday 3am = matches window 2 ✓
# - Sunday 10am = doesn't match any window ✗

# ============================================================================
# Example 6: Complex - Inclusion + Exclusion
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "complex_windows" {
  name                  = "Complex Window Rules"
  deployment_start_time = local.deployment_start
  scope_id              = "complex-group-333"

  time_condition {
    offset_hours = 24
  }

  # Allow during office hours
  inclusion_time_windows {
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
      time_of_day_start = "09:00:00"
      time_of_day_end   = "17:00:00"
    }
  }

  # But block during holiday freeze
  exclusion_time_windows {
    window {
      date_start = "2024-12-20T00:00:00Z"
      date_end   = "2025-01-05T23:59:59Z"
    }
  }
}

# Behavior:
# - Base: 24h wait
# - Inclusion: Must be weekday 9am-5pm
# - Exclusion: NOT during Dec 20 - Jan 5
# - All three must be satisfied (AND logic)
# - Exclusion takes precedence if both match

# ============================================================================
# Example 7: Date Range Inclusion (Specific Campaign Period)
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "campaign_window" {
  name                  = "Q1 Campaign Window"
  deployment_start_time = local.deployment_start
  scope_id              = "campaign-group-444"

  time_condition {
    offset_hours = 0  # Immediate when window opens
  }

  # Only deploy during Q1 2024
  inclusion_time_windows {
    window {
      date_start = "2024-01-01T00:00:00Z"
      date_end   = "2024-03-31T23:59:59Z"
    }
  }
}

# Behavior:
# - Only allows deployment during Q1 2024
# - Before Jan 1 or after Mar 31 = gate closed

# ============================================================================
# Example 8: Multiple Exclusions (Holiday + Weekends)
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "multiple_exclusions" {
  name                  = "Multiple Exclusions"
  deployment_start_time = local.deployment_start
  scope_id              = "conservative-group-555"

  time_condition {
    offset_hours = 12
  }

  exclusion_time_windows {
    # Exclusion 1: Weekends
    window {
      days_of_week = ["saturday", "sunday"]
    }

    # Exclusion 2: Holiday freeze
    window {
      date_start = "2024-12-20T00:00:00Z"
      date_end   = "2025-01-05T23:59:59Z"
    }

    # Exclusion 3: Friday afternoons (change freeze)
    window {
      days_of_week      = ["friday"]
      time_of_day_start = "15:00:00"  # 3 PM Friday
      time_of_day_end   = "23:59:59"
    }
  }
}

# Behavior:
# - Blocks if ANY exclusion matches (OR logic)
# - Saturday = blocked (weekend)
# - Friday 4pm = blocked (Friday afternoon)
# - Dec 25 = blocked (holiday freeze)

# ============================================================================
# Example 9: Night Deployments Only
# ============================================================================

data "microsoft365_utility_deployment_scheduler" "night_only" {
  name                  = "Night Deployments Only"
  deployment_start_time = local.deployment_start
  scope_id              = "night-group-666"

  time_condition {
    offset_hours = 24
  }

  # Only deploy between midnight and 5am
  inclusion_time_windows {
    window {
      time_of_day_start = "00:00:00"
      time_of_day_end   = "05:00:00"
    }
  }
}

# Behavior:
# - Only allows deployment during nighttime hours (00:00-05:00 UTC)
# - Works every day, but only in the early morning

# ============================================================================
# Example 10: Real-World Production Rollout
# ============================================================================

locals {
  production_deployment_start = "2024-02-01T00:00:00Z"
}

# Phase 1: Pilot (office hours, no restrictions)
data "microsoft365_utility_deployment_scheduler" "prod_phase1" {
  name                  = "Production Phase 1 - Pilot"
  deployment_start_time = local.production_deployment_start
  scope_id              = "pilot-group"

  time_condition {
    offset_hours = 0  # Immediate
  }

  # Office hours only for pilot
  inclusion_time_windows {
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday", "friday"]
      time_of_day_start = "10:00:00"
      time_of_day_end   = "16:00:00"  # Safe window 10am-4pm
    }
  }
}

# Phase 2: Early Adopters (after 48h, weekdays only)
data "microsoft365_utility_deployment_scheduler" "prod_phase2" {
  name                  = "Production Phase 2 - Early Adopters"
  deployment_start_time = local.production_deployment_start
  scope_id              = "early-adopters-group"

  time_condition {
    offset_hours = 48
  }

  # Weekday business hours
  inclusion_time_windows {
    window {
      days_of_week      = ["monday", "tuesday", "wednesday", "thursday"]
      time_of_day_start = "09:00:00"
      time_of_day_end   = "17:00:00"
    }
  }

  # No Friday deployments (weekend buffer)
  exclusion_time_windows {
    window {
      days_of_week = ["friday", "saturday", "sunday"]
    }
  }
}

# Phase 3: Production (after 96h, strict windows)
data "microsoft365_utility_deployment_scheduler" "prod_phase3" {
  name                  = "Production Phase 3 - Full Rollout"
  deployment_start_time = local.production_deployment_start
  scope_id              = "production-group"

  time_condition {
    offset_hours = 96  # 4 days
  }

  # Very conservative: Tue-Thu only, mid-day
  inclusion_time_windows {
    window {
      days_of_week      = ["tuesday", "wednesday", "thursday"]
      time_of_day_start = "10:00:00"
      time_of_day_end   = "14:00:00"
    }
  }

  # Plus no major holidays
  exclusion_time_windows {
    window {
      date_start = "2024-12-20T00:00:00Z"
      date_end   = "2025-01-05T23:59:59Z"
    }
  }
}

# Use in policy
resource "microsoft365_windows_update_policy" "production_rollout" {
  name = "Production Rollout with Time Windows"

  assignment {
    target {
      group_ids = compact([
        data.microsoft365_utility_deployment_scheduler.prod_phase1.released_scope_id,
        data.microsoft365_utility_deployment_scheduler.prod_phase2.released_scope_id,
        data.microsoft365_utility_deployment_scheduler.prod_phase3.released_scope_id,
      ])
    }
  }
}

# Monitor with detailed output
output "production_rollout_status" {
  description = "Detailed status of production rollout with time windows"
  value = {
    phase1_pilot = {
      ready  = data.microsoft365_utility_deployment_scheduler.prod_phase1.condition_met
      status = data.microsoft365_utility_deployment_scheduler.prod_phase1.status_message
    }
    phase2_early_adopters = {
      ready  = data.microsoft365_utility_deployment_scheduler.prod_phase2.condition_met
      status = data.microsoft365_utility_deployment_scheduler.prod_phase2.status_message
    }
    phase3_production = {
      ready  = data.microsoft365_utility_deployment_scheduler.prod_phase3.condition_met
      status = data.microsoft365_utility_deployment_scheduler.prod_phase3.status_message
    }
  }
}

# ============================================================================
# Understanding Status Messages
# ============================================================================

# When all conditions met (gate open):
# "Conditions met: Time condition met (50.0h/48h required), Inclusion window satisfied (within inclusion window 1), No exclusion window active (outside all exclusion windows)"

# When time condition not met:
# "Waiting: Time condition NOT met (20.0h/48h required), Inclusion window satisfied (within inclusion window 1), No exclusion window active (outside all exclusion windows)"

# When outside inclusion window:
# "Waiting: Time condition met (50.0h/48h required), Inclusion window NOT satisfied (outside all inclusion windows), No exclusion window active (outside all exclusion windows)"

# When blocked by exclusion:
# "Waiting: Time condition met (50.0h/48h required), Inclusion window satisfied (within inclusion window 1), Deployment blocked by exclusion window (within exclusion window 1)"
