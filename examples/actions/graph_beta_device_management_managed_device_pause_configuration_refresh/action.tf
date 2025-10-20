terraform {
  required_providers {
    microsoft365 = {
      source = "deploymenttheory/microsoft365"
    }
  }
}

provider "microsoft365" {
  # Authentication configuration
}

# Example 1: Basic - Pause configuration refresh for maintenance
# Use case: 2-hour maintenance window for application updates
action "pause_config_basic" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "12345678-1234-1234-1234-123456789abc"
    pause_time_period_in_minutes = 120 # 2 hours
  }

  managed_devices {
    device_id                    = "87654321-4321-4321-4321-ba9876543210"
    pause_time_period_in_minutes = 120 # 2 hours
  }

  timeouts {
    invoke = "5m"
  }
}

# Example 2: Variable pause durations - Different maintenance windows
# Use case: Different devices need different pause durations
action "pause_config_variable" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "short-maintenance-device"
    pause_time_period_in_minutes = 60 # 1 hour - quick patch
  }

  managed_devices {
    device_id                    = "medium-maintenance-device"
    pause_time_period_in_minutes = 240 # 4 hours - application upgrade
  }

  managed_devices {
    device_id                    = "long-maintenance-device"
    pause_time_period_in_minutes = 480 # 8 hours - major system update
  }

  timeouts {
    invoke = "5m"
  }
}

# Example 3: Co-managed devices - Hybrid SCCM environment
# Use case: Pause Intune config refresh during SCCM maintenance
action "pause_config_comanaged" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  comanaged_devices {
    device_id                    = "abcdef12-3456-7890-abcd-ef1234567890"
    pause_time_period_in_minutes = 240 # 4 hours for SCCM maintenance
  }

  comanaged_devices {
    device_id                    = "fedcba09-8765-4321-fedc-ba0987654321"
    pause_time_period_in_minutes = 240 # 4 hours
  }

  timeouts {
    invoke = "5m"
  }
}

# Example 4: Troubleshooting - Pause during policy conflict investigation
# Use case: Temporarily freeze configuration while investigating issues
action "pause_config_troubleshoot" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "problematic-device-1"
    pause_time_period_in_minutes = 360 # 6 hours for investigation
  }

  managed_devices {
    device_id                    = "problematic-device-2"
    pause_time_period_in_minutes = 360 # 6 hours
  }

  timeouts {
    invoke = "5m"
  }
}

# Example 5: Business-critical operations - Extended pause
# Use case: Prevent policy changes during critical business operations
action "pause_config_business_critical" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "trading-floor-device-1"
    pause_time_period_in_minutes = 480 # 8 hours - trading hours
  }

  managed_devices {
    device_id                    = "pos-system-device-1"
    pause_time_period_in_minutes = 600 # 10 hours - retail hours
  }

  managed_devices {
    device_id                    = "medical-device-1"
    pause_time_period_in_minutes = 720 # 12 hours - medical shift
  }

  timeouts {
    invoke = "10m"
  }
}

# Example 6: Maximum pause - 24-hour freeze
# Use case: Full day maintenance or testing
action "pause_config_max" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "test-device-1"
    pause_time_period_in_minutes = 1440 # 24 hours - maximum allowed
  }

  timeouts {
    invoke = "5m"
  }
}

# Example 7: Incident response - Pause during security investigation
# Use case: Freeze configuration during incident response
action "pause_config_incident" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "compromised-device-1"
    pause_time_period_in_minutes = 480 # 8 hours for forensic analysis
  }

  managed_devices {
    device_id                    = "affected-device-2"
    pause_time_period_in_minutes = 480 # 8 hours
  }

  timeouts {
    invoke = "5m"
  }
}

# Example 8: UAT/Testing - Staging environment configuration freeze
# Use case: User acceptance testing with stable configuration
action "pause_config_uat" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "uat-device-1"
    pause_time_period_in_minutes = 1440 # 24 hours for testing cycle
  }

  managed_devices {
    device_id                    = "uat-device-2"
    pause_time_period_in_minutes = 1440 # 24 hours
  }

  managed_devices {
    device_id                    = "uat-device-3"
    pause_time_period_in_minutes = 1440 # 24 hours
  }

  timeouts {
    invoke = "10m"
  }
}

# Example 9: Mixed environment - Both managed and co-managed
# Use case: Organization-wide maintenance window
action "pause_config_mixed" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "intune-device-1"
    pause_time_period_in_minutes = 240 # 4 hours
  }

  managed_devices {
    device_id                    = "intune-device-2"
    pause_time_period_in_minutes = 240 # 4 hours
  }

  comanaged_devices {
    device_id                    = "hybrid-device-1"
    pause_time_period_in_minutes = 240 # 4 hours
  }

  timeouts {
    invoke = "10m"
  }
}

# Example 10: Policy rollout staging - Controlled deployment
# Use case: Pause devices while testing new policies on pilot group
action "pause_config_staging" {
  provider_type = microsoft365.graph_beta_device_management_managed_device_pause_configuration_refresh

  managed_devices {
    device_id                    = "production-device-1"
    pause_time_period_in_minutes = 720 # 12 hours during pilot
  }

  managed_devices {
    device_id                    = "production-device-2"
    pause_time_period_in_minutes = 720 # 12 hours
  }

  managed_devices {
    device_id                    = "production-device-3"
    pause_time_period_in_minutes = 720 # 12 hours
  }

  timeouts {
    invoke = "10m"
  }
}

# Important Notes:
#
# 1. Pause Duration Constraints:
#    - Minimum: 1 minute
#    - Maximum: 1440 minutes (24 hours)
#    - Configuration refresh automatically resumes after expiration
#    - Cannot extend pause once initiated (must re-pause)
#
# 2. What Gets Paused:
#    - New policy deployments
#    - Policy updates and changes
#    - Configuration profile updates
#    - App deployment policy changes
#    - Compliance policy updates
#
# 3. What Does NOT Get Paused:
#    - Existing applied policies (remain in effect)
#    - Device check-ins and status reporting
#    - Manual user-initiated syncs from Company Portal
#    - Critical security updates (may still apply)
#    - Emergency remote actions (wipe, lock, etc.)
#
# 4. Best Practices:
#    - Use shortest necessary pause duration
#    - Schedule pauses during maintenance windows
#    - Document reason for pause in change management
#    - Monitor device status during pause
#    - Resume normal operations promptly
#
# 5. Common Use Cases by Duration:
#    - 60 minutes (1 hour): Quick application updates
#    - 120 minutes (2 hours): Standard maintenance windows
#    - 240 minutes (4 hours): Extended maintenance or testing
#    - 480 minutes (8 hours): Business day operations
#    - 720 minutes (12 hours): Shift-based operations
#    - 1440 minutes (24 hours): Full day testing or investigation
#
# 6. Troubleshooting Scenarios:
#    - Policy conflicts: Pause to investigate without new changes
#    - Application compatibility: Pause during app testing
#    - Performance issues: Pause to isolate configuration impact
#    - Rollback situations: Pause while reverting changes
#
# 7. Compliance Considerations:
#    - Pausing may temporarily affect compliance state
#    - Security policies still enforced (existing)
#    - Document pauses for audit purposes
#    - Balance operational needs with security requirements

