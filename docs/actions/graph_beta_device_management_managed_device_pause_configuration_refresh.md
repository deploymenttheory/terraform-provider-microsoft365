---
page_title: "Microsoft 365_microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh Action - terraform-provider-microsoft365"
subcategory: "Device Management"
description: |-
  Initiates a command to pause configuration refresh on managed Windows devices using the /deviceManagement/managedDevices/{managedDeviceId}/pauseConfigurationRefresh and /deviceManagement/comanagedDevices/{managedDeviceId}/pauseConfigurationRefresh endpoints. This action temporarily prevents devices from receiving and applying new configuration policies from Intune, which is useful during maintenance windows, troubleshooting, or when you need to prevent policy changes from being applied to specific devices for a defined period.
  Important Notes:
  Only works on Windows 10/11 devicesConfiguration refresh automatically resumes after the pause period expiresMaximum pause period is typically 24 hours (1440 minutes)Does not affect existing applied policies, only prevents new policy updatesDevice can still check in and report statusCritical security updates may still be appliedUser can still manually sync from Company Portal
  Use Cases:
  Maintenance windows for critical applicationsTroubleshooting policy conflictsTesting policy changes in stagingPreventing policy updates during business-critical operationsTemporary freeze during incident responseUser acceptance testing (UAT) phases
  Platform Support:
  Windows: Windows 10/11Other Platforms: Not supported
  Reference: Microsoft Graph API - Pause Configuration Refresh https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-pauseconfigurationrefresh?view=graph-rest-beta
---

# Microsoft 365_microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh (Action)

Initiates a command to pause configuration refresh on managed Windows devices using the `/deviceManagement/managedDevices/{managedDeviceId}/pauseConfigurationRefresh` and `/deviceManagement/comanagedDevices/{managedDeviceId}/pauseConfigurationRefresh` endpoints. This action temporarily prevents devices from receiving and applying new configuration policies from Intune, which is useful during maintenance windows, troubleshooting, or when you need to prevent policy changes from being applied to specific devices for a defined period.

**Important Notes:**
- Only works on Windows 10/11 devices
- Configuration refresh automatically resumes after the pause period expires
- Maximum pause period is typically 24 hours (1440 minutes)
- Does not affect existing applied policies, only prevents new policy updates
- Device can still check in and report status
- Critical security updates may still be applied
- User can still manually sync from Company Portal

**Use Cases:**
- Maintenance windows for critical applications
- Troubleshooting policy conflicts
- Testing policy changes in staging
- Preventing policy updates during business-critical operations
- Temporary freeze during incident response
- User acceptance testing (UAT) phases

**Platform Support:**
- **Windows**: Windows 10/11
- **Other Platforms**: Not supported

**Reference:** [Microsoft Graph API - Pause Configuration Refresh](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-pauseconfigurationrefresh?view=graph-rest-beta)

## Example Usage

### Basic Configuration Refresh Pause

```terraform
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
```

## API Permissions

The following [Microsoft Graph API permissions](https://learn.microsoft.com/en-us/graph/permissions-reference) are required to use this action:

### Read Permissions
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementManagedDevices.Read.All`

### Write Permissions
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementManagedDevices.Read.All`

## Authentication Mechanism

This action uses the configured Microsoft 365 provider authentication. Ensure you have properly configured the provider with appropriate credentials and permissions.

## Important Notes

### Pause Duration
- **Minimum**: 1 minute
- **Maximum**: 1440 minutes (24 hours)
- **Auto-Resume**: Configuration refresh automatically resumes after the pause period expires
- **Cannot Extend**: Once initiated, the pause duration cannot be extended. You must wait for it to expire and re-pause if needed.

### What Gets Paused
Configuration refresh pause prevents the following from being applied to devices:
- New policy deployments
- Policy updates and changes
- Configuration profile updates
- Application deployment policy changes
- Compliance policy updates
- Settings catalog changes

### What Does NOT Get Paused
The following continue to operate normally during the pause:
- **Existing Policies**: All previously applied policies remain in effect
- **Device Check-Ins**: Devices continue checking in and reporting status
- **Manual Syncs**: Users can manually sync from the Company Portal app
- **Critical Updates**: Critical security updates may still be applied
- **Remote Actions**: Emergency actions (wipe, lock, locate) still work
- **Compliance Evaluation**: Compliance checks continue with existing policies

### Common Use Cases

#### Maintenance Windows
Pause configuration refresh during scheduled maintenance to prevent policy changes from interfering:
- Application updates (60-120 minutes)
- System upgrades (240-480 minutes)
- Infrastructure changes (240-480 minutes)

#### Troubleshooting
Temporarily freeze configuration to investigate issues:
- Policy conflict investigation (360-720 minutes)
- Performance issue isolation (240-360 minutes)
- Application compatibility testing (480-720 minutes)
- Rollback procedures (120-240 minutes)

#### Business-Critical Operations
Prevent policy changes during critical business operations:
- Trading floor operations (480 minutes)
- Retail POS systems (600 minutes)
- Medical device operations (720 minutes)
- Manufacturing floor (480-720 minutes)

#### Testing and Staging
Maintain stable configuration during testing phases:
- User acceptance testing (1440 minutes)
- Pilot deployments (720-1440 minutes)
- Policy validation (480-720 minutes)
- Staged rollouts (varies)

#### Incident Response
Freeze configuration during security investigations:
- Forensic analysis (480-720 minutes)
- Containment phase (360-480 minutes)
- Evidence collection (240-360 minutes)

### Best Practices

1. **Minimum Necessary Duration**
   - Use the shortest pause duration that meets your needs
   - Avoid unnecessary extended pauses
   - Plan ahead for anticipated maintenance duration

2. **Communication and Documentation**
   - Document the reason for pausing in change management systems
   - Notify affected teams and stakeholders
   - Track pause start and end times
   - Document any issues encountered

3. **Monitoring During Pause**
   - Monitor device status and health
   - Watch for user-reported issues
   - Check for critical security alerts
   - Verify devices remain compliant

4. **Post-Pause Actions**
   - Verify configuration refresh has resumed
   - Check for pending policy updates
   - Monitor for successful policy application
   - Review device compliance status

5. **Security Considerations**
   - Balance operational needs with security requirements
   - Don't pause during active security incidents (unless for containment)
   - Consider compliance implications
   - Document security exceptions if applicable

### Platform Support
- **Windows 10/11**: Fully supported
- **Earlier Windows versions**: Not supported
- **Non-Windows platforms**: Not applicable

### Compliance Considerations
- Pausing configuration refresh may temporarily affect compliance reporting
- Existing security policies remain enforced
- Document all pauses for audit purposes
- Consider regulatory requirements before pausing
- Some compliance frameworks may require approval for pauses

### Limitations and Constraints
- **Maximum Pause**: 24 hours (1440 minutes)
- **No Extension**: Cannot extend an active pause
- **Windows Only**: Only works on Windows 10/11 devices
- **Policy Type**: Affects Intune policies only (not SCCM in co-managed scenarios)
- **Emergency Actions**: Some critical updates may override the pause

### Monitoring and Validation
After initiating a pause, you should:
1. Verify the pause is active on target devices
2. Monitor for any unexpected policy applications
3. Check device check-in status
4. Review error logs if sync attempts occur
5. Confirm automatic resume after pause expires

## Microsoft Graph API References

- [pauseConfigurationRefresh action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-pauseconfigurationrefresh?view=graph-rest-beta)
- [Intune Device Management](https://learn.microsoft.com/en-us/mem/intune/remote-actions/device-management)
- [Device Configuration Refresh](https://learn.microsoft.com/en-us/mem/intune/configuration/device-profile-troubleshoot)

<!-- action schema generated by tfplugindocs -->
## Schema

### Optional

- `comanaged_devices` (Block List) List of co-managed devices to pause configuration refresh for. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--comanaged_devices))
- `managed_devices` (Block List) List of managed devices to pause configuration refresh for. Each device can have a different pause duration based on specific requirements.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedblock--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedblock--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to pause configuration refresh for.

Example: `"abcdef12-3456-7890-abcd-ef1234567890"`
- `pause_time_period_in_minutes` (Number) The duration in minutes to pause configuration refresh for this device.

**Valid Range:** 1 to 1440 minutes (1 minute to 24 hours)


<a id="nestedblock--managed_devices"></a>
### Nested Schema for `managed_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the managed device to pause configuration refresh for. This should be a Windows 10 or Windows 11 device managed by Intune.

Example: `"12345678-1234-1234-1234-123456789abc"`
- `pause_time_period_in_minutes` (Number) The duration in minutes to pause configuration refresh for this device. Configuration refresh will automatically resume after this period expires.

**Valid Range:** 1 to 1440 minutes (1 minute to 24 hours)

**Common Values:**
- `60` - 1 hour (short maintenance)
- `120` - 2 hours (application updates)
- `240` - 4 hours (extended maintenance)
- `480` - 8 hours (business day)
- `1440` - 24 hours (full day)

Example: `120` (2 hours)


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

