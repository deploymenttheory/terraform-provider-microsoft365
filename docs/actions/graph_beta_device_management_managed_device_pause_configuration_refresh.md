---
page_title: "Microsoft 365_microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh Action - terraform-provider-microsoft365"
subcategory: "Device Management"
description: |-
  Pauses configuration refresh on managed Windows devices in Microsoft Intune using the /deviceManagement/managedDevices/{managedDeviceId}/pauseConfigurationRefresh and /deviceManagement/comanagedDevices/{managedDeviceId}/pauseConfigurationRefresh endpoints. This action is used to temporarily prevent devices from receiving and applying new configuration policies from Intune, which is useful during maintenance windows, troubleshooting, or when you need to prevent policy changes from being applied to specific devices for a defined period.
  Important Notes:
  Only works on Windows 10/11 devicesConfiguration refresh automatically resumes after the pause period expiresMaximum pause period is typically 24 hours (1440 minutes)Does not affect existing applied policies, only prevents new policy updatesDevice can still check in and report statusCritical security updates may still be appliedUser can still manually sync from Company Portal
  Use Cases:
  Maintenance windows for critical applicationsTroubleshooting policy conflictsTesting policy changes in stagingPreventing policy updates during business-critical operationsTemporary freeze during incident responseUser acceptance testing (UAT) phases
  Platform Support:
  Windows: Windows 10/11Other Platforms: Not supported
  Reference: Microsoft Graph API - Pause Configuration Refresh https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-pauseconfigurationrefresh?view=graph-rest-beta
---

# Microsoft 365_microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh (Action)

Pauses configuration refresh on managed Windows devices in Microsoft Intune using the `/deviceManagement/managedDevices/{managedDeviceId}/pauseConfigurationRefresh` and `/deviceManagement/comanagedDevices/{managedDeviceId}/pauseConfigurationRefresh` endpoints. This action is used to temporarily prevent devices from receiving and applying new configuration policies from Intune, which is useful during maintenance windows, troubleshooting, or when you need to prevent policy changes from being applied to specific devices for a defined period.

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
# Example 1: Pause configuration refresh on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "pause_single" {
  config {
    managed_devices = [
      {
        device_id                    = "12345678-1234-1234-1234-123456789abc"
        pause_time_period_in_minutes = 60
      }
    ]
  }
}

# Example 2: Pause configuration refresh on multiple devices with different durations
action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "pause_multiple" {
  config {
    managed_devices = [
      {
        device_id                    = "12345678-1234-1234-1234-123456789abc"
        pause_time_period_in_minutes = 60
      },
      {
        device_id                    = "87654321-4321-4321-4321-ba9876543210"
        pause_time_period_in_minutes = 120
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Pause with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "pause_maximal" {
  config {
    managed_devices = [
      {
        device_id                    = "12345678-1234-1234-1234-123456789abc"
        pause_time_period_in_minutes = 180
      }
    ]

    comanaged_devices = [
      {
        device_id                    = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        pause_time_period_in_minutes = 90
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Pause configuration refresh during maintenance window
data "microsoft365_graph_beta_device_management_managed_device" "maintenance_devices" {
  filter_type  = "odata"
  odata_filter = "deviceCategoryDisplayName eq 'Maintenance Queue'"
}

action "microsoft365_graph_beta_device_management_managed_device_pause_configuration_refresh" "pause_maintenance" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.maintenance_devices.items : {
        device_id                    = device.id
        pause_time_period_in_minutes = 240
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}
```

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this action:

**Required:**
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementManagedDevices.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |
| v0.40.0-alpha | Experimental | Example fixes and refactored sync progress logic |


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

- `comanaged_devices` (Attributes List) List of co-managed devices to pause configuration refresh for. These are devices managed by both Intune and Configuration Manager (SCCM).

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--comanaged_devices))
- `ignore_partial_failures` (Boolean) If set to `true`, the action will succeed even if some operations fail. Failed operations will be reported as warnings instead of errors. Default: `false` (action fails if any operation fails).
- `managed_devices` (Attributes List) List of managed devices to pause configuration refresh for. Each device can have a different pause duration based on specific requirements.

**Note:** At least one of `managed_devices` or `comanaged_devices` must be provided. (see [below for nested schema](#nestedatt--managed_devices))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `validate_device_exists` (Boolean) Whether to validate that devices exist and are Windows devices before attempting to pause configuration refresh. Disabling this can speed up planning but may result in runtime errors for non-existent or unsupported devices. Default: `true`.

<a id="nestedatt--comanaged_devices"></a>
### Nested Schema for `comanaged_devices`

Required:

- `device_id` (String) The unique identifier (GUID) of the co-managed device to pause configuration refresh for.

Example: `"abcdef12-3456-7890-abcd-ef1234567890"`
- `pause_time_period_in_minutes` (Number) The duration in minutes to pause configuration refresh for this device.

**Valid Range:** 1 to 1440 minutes (1 minute to 24 hours)


<a id="nestedatt--managed_devices"></a>
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

- `invoke` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

