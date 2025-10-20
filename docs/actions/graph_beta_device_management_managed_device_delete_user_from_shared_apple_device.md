---
page_title: "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device Action - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Deletes a user and their cached data from Shared iPad devices using the /deviceManagement/managedDevices/{managedDeviceId}/deleteUserFromSharedAppleDevice endpoint. This action permanently removes the specified user's account and all associated cached data from the Shared iPad, freeing up storage space for other users.
  What This Action Does:
  Permanently deletes user from Shared iPad device rosterRemoves all cached user data (documents, photos, app data)Frees up device storage spacePrevents user from logging back into that specific deviceDoes not affect user's account or data in the cloudCannot be undone (user must be re-added if needed)
  Difference from Logout:
  Logout: Temporary - logs out active user, data stays cached, user can log back inDelete User: Permanent - removes user from device, deletes all cached data, user cannot log back in
  Platform Support:
  iPadOS: Full support (Shared iPad mode only)iOS: Not supported (iPhones don't support Shared mode)Other platforms: Not supported
  Common Use Cases:
  Student/employee has left organizationFreeing storage space on Shared iPadsRemoving users no longer assigned to deviceManaging maximum cached users limitClassroom roster changesUser account deprovisioningStorage quota management
  Important Considerations:
  Irreversible: Deletes all user's cached data on devicePer-Device: User account in cloud unaffected, only removed from specific deviceRe-Addition: User can be re-added, but will need to download all data againActive Users: Can delete currently logged-in users (forces logout)Storage Impact: Immediately frees up user's allocated storage
  Requirements:
  Device must be in Shared iPad modeDevice must be supervisedMust be enrolled via DEP/ABMUser must exist in device's cached user listDevice must be online
  Reference: Microsoft Graph API - Delete User From Shared Apple Device https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-beta
---

# microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device (Action)

Deletes a user and their cached data from Shared iPad devices using the `/deviceManagement/managedDevices/{managedDeviceId}/deleteUserFromSharedAppleDevice` endpoint. This action permanently removes the specified user's account and all associated cached data from the Shared iPad, freeing up storage space for other users.

**What This Action Does:**
- Permanently deletes user from Shared iPad device roster
- Removes all cached user data (documents, photos, app data)
- Frees up device storage space
- Prevents user from logging back into that specific device
- Does not affect user's account or data in the cloud
- Cannot be undone (user must be re-added if needed)

**Difference from Logout:**
- **Logout**: Temporary - logs out active user, data stays cached, user can log back in
- **Delete User**: Permanent - removes user from device, deletes all cached data, user cannot log back in

**Platform Support:**
- **iPadOS**: Full support (Shared iPad mode only)
- **iOS**: Not supported (iPhones don't support Shared mode)
- **Other platforms**: Not supported

**Common Use Cases:**
- Student/employee has left organization
- Freeing storage space on Shared iPads
- Removing users no longer assigned to device
- Managing maximum cached users limit
- Classroom roster changes
- User account deprovisioning
- Storage quota management

**Important Considerations:**
- **Irreversible**: Deletes all user's cached data on device
- **Per-Device**: User account in cloud unaffected, only removed from specific device
- **Re-Addition**: User can be re-added, but will need to download all data again
- **Active Users**: Can delete currently logged-in users (forces logout)
- **Storage Impact**: Immediately frees up user's allocated storage

**Requirements:**
- Device must be in Shared iPad mode
- Device must be supervised
- Must be enrolled via DEP/ABM
- User must exist in device's cached user list
- Device must be online

**Reference:** [Microsoft Graph API - Delete User From Shared Apple Device](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-beta)

## Microsoft Documentation

### Graph API References
- [deleteUserFromSharedAppleDevice action](https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-deleteuserfromsharedappledevice?view=graph-rest-beta)
- [managedDevice resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevice?view=graph-rest-beta)

### Intune Remote Actions Guides
- [Windows Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=windows)
- [iOS/iPadOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=ios-ipados)
- [macOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=macos)
- [Android Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=android)
- [ChromeOS Remote Actions](https://learn.microsoft.com/en-us/intune/intune-service/remote-actions/?tabs=chromeos)

## API Permissions

The following API permissions are required in order to use this action.

### Microsoft Graph

- **Application**: `DeviceManagementManagedDevices.PrivilegedOperations.All`
- **Delegated**: `DeviceManagementManagedDevices.PrivilegedOperations.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Notes

### Platform Compatibility

| Platform | Support | Requirements |
|----------|---------|--------------|
| **iPadOS** | ✅ Full Support | Shared iPad mode, supervised, DEP/ABM |
| **iOS** | ❌ Not Supported | iPhones don't support Shared mode |
| **macOS** | ❌ Not Supported | Shared mode is iPad-only |
| **Windows** | ❌ Not Supported | Shared mode is iPad-only |
| **Android** | ❌ Not Supported | Shared mode is iPad-only |

### Delete User vs Logout User

| Action | Effect | Data Impact | User Can Return | Use Case |
|--------|--------|-------------|-----------------|----------|
| **Delete User** | Permanent removal | All cached data deleted | No (must be re-added) | User left organization, storage cleanup |
| **Logout User** | Temporary logout | Data preserved | Yes (immediately) | Session management, device rotation |

### What Gets Deleted

**User Account**
- User removed from device roster
- Cannot log into this specific device
- Must be re-added to regain access

**Cached Data**
- Documents and files
- Photos and media
- App data and settings
- Downloaded content
- Offline cached data
- User preferences

**Storage Impact**
- Immediately frees allocated space
- Available for other users
- Helps manage device capacity
- Reduces storage pressure

### What is NOT Affected

**Cloud Data**
- User's Azure AD/Entra ID account
- User's cloud-stored documents
- User's OneDrive data
- User's email
- User's cloud app data
- User's other devices

**Device Configuration**
- Device enrollment
- MDM policies
- Installed apps
- Device settings
- Other cached users
- Network configuration

### Common Use Cases

| Scenario | Description | Benefit |
|----------|-------------|---------|
| **User Departure** | Employee/student left organization | Complete removal from devices |
| **Storage Management** | Free space for active users | Improve device capacity |
| **Roster Management** | Remove users no longer assigned | Keep user lists current |
| **Max Users Reached** | Hit maximum cached users limit | Make room for new users |
| **Account Cleanup** | Regular maintenance | Optimize device performance |
| **Graduation/Termination** | End of tenure | Proper offboarding |

### Requirements

#### Device Requirements
- iPad in Shared iPad mode
- Supervised via DEP/ABM
- iPadOS 9.3 or later
- Online connectivity
- Managed by Intune

#### User Requirements
- User exists in device's cached user list
- Valid user principal name (UPN)
- Managed Apple ID or federated account

#### Action-Specific
- Can delete currently logged-in users (forces logout)
- Can delete multiple users from multiple devices
- Per-device operation (doesn't affect user globally)

### Data Deletion Process

**Immediate Actions**
1. User logged out if currently active
2. User removed from device roster
3. Cached data deletion begins
4. Storage space freed up
5. User cannot log back in

**Data Removal**
- User profile deleted
- Documents removed
- Photos deleted
- App data cleared
- Downloaded content removed
- Preferences erased

**Timeline**
- Deletion begins immediately
- Full removal takes 1-5 minutes
- Storage space available instantly
- No recovery possible

### User Impact

**For Deleted User**
- **Access**: Cannot log into this specific device
- **Data**: All cached data on device is lost
- **Cloud**: Cloud data unaffected
- **Other Devices**: Can still use other devices
- **Re-Addition**: Can be re-added, but starts fresh

**For Other Users**
- No impact on their accounts
- May see improved performance
- More storage available
- Uninterrupted access

**For Device**
- Frees storage space
- One fewer cached user
- Performance may improve
- Ready for new user

### Best Practices

**Before Deletion**
- Verify user has backed up important data
- Confirm user identity and device
- Document business justification
- Check if user needs re-addition later
- Consider logout instead if temporary

**Communication**
- Notify users before deletion
- Explain data loss implications
- Provide data backup instructions
- Set expectations for re-addition
- Follow organizational policies

**Planning**
- Delete during off-hours when possible
- Batch related operations
- Monitor storage impact
- Track deleted users
- Document deletion reasons

**Storage Management**
- Remove inactive users first
- Monitor device capacity
- Balance user count vs storage
- Plan for new user additions
- Regular maintenance schedule

### Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Action fails | User not on device | Verify user exists in device roster |
| No effect | Wrong device type | Confirm device is Shared iPad |
| Timeout | Device offline | Check network connectivity |
| Error | Not supervised | Enroll via DEP/ABM |
| User remains | Not Shared mode | Verify Shared iPad configuration |
| Storage not freed | Deletion in progress | Wait 5 minutes for completion |

### Recovery and Re-Addition

**No Recovery**
- Deleted data cannot be recovered
- Operation is permanent
- User must start fresh if re-added

**Re-Adding Users**
1. User can be re-added to device
2. Will start with clean profile
3. Must download all data again
4. Settings must be reconfigured
5. Treated as new user on device

**Alternative Approaches**
- Use logout for temporary removal
- Use wipe for complete device reset
- Consider user data backup before deletion
- Plan for re-addition if needed

### Shared iPad Mode Context

**User Management**
- Maximum users per device (typically 24-40)
- Storage allocation per user
- Quota management
- Active vs cached users
- User pruning strategies

**Storage Allocation**
- Each user gets allocated space
- Deletion frees that allocation
- Space available immediately
- Helps manage device capacity
- Prevents storage exhaustion

**Multi-User Environment**
- Multiple students/employees per device
- Individual user experiences
- Separate data and settings
- Managed user switching
- IT control over user roster

## Example Usage

```terraform
# ============================================================================
# Example 1: Delete single user from single Shared iPad
# ============================================================================
# Use case: Student left school, remove from device
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "single_user_single_device" {

  devices = [
    {
      device_id           = "12345678-1234-1234-1234-123456789abc"
      user_principal_name = "student@school.edu"
    }
  ]

  timeouts = {
    invoke = "5m"
  }
}

# ============================================================================
# Example 2: Delete same user from multiple Shared iPads
# ============================================================================
# Use case: User left organization, remove from all classroom devices
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "one_user_multiple_devices" {

  devices = [
    {
      device_id           = "12345678-1234-1234-1234-123456789abc"
      user_principal_name = "student@school.edu"
    },
    {
      device_id           = "87654321-4321-4321-4321-ba9876543210"
      user_principal_name = "student@school.edu"
    },
    {
      device_id           = "abcdef12-3456-7890-abcd-ef1234567890"
      user_principal_name = "student@school.edu"
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 3: Delete different users from different Shared iPads
# ============================================================================
# Use case: Mixed cleanup - remove specific users from specific devices
action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "multiple_users_multiple_devices" {

  devices = [
    {
      device_id           = "12345678-1234-1234-1234-123456789abc"
      user_principal_name = "student1@school.edu"
    },
    {
      device_id           = "87654321-4321-4321-4321-ba9876543210"
      user_principal_name = "student2@school.edu"
    },
    {
      device_id           = "abcdef12-3456-7890-abcd-ef1234567890"
      user_principal_name = "student3@school.edu"
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 4: Delete users using datasource for device discovery
# ============================================================================
# Use case: Remove graduated students from all lab iPads
data "microsoft365_graph_beta_device_management_managed_device" "lab_ipads" {
  filter_type  = "device_name"
  filter_value = "LAB-IPAD-"
}

locals {
  # List of users to remove
  departed_users = ["student1@school.edu", "student2@school.edu", "student3@school.edu"]
  
  # Create device-user pairs for each combination
  device_user_pairs = flatten([
    for device in data.microsoft365_graph_beta_device_management_managed_device.lab_ipads.items : [
      for user in local.departed_users : {
        device_id           = device.id
        user_principal_name = user
      }
    ]
  ])
}

action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "remove_departed_users" {

  devices = local.device_user_pairs

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 5: Delete users from supervised iPads (storage management)
# ============================================================================
# Use case: Free up storage by removing inactive users
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (isSupervised eq true)"
}

locals {
  # List of inactive users to remove for storage space
  inactive_users = ["inactive1@school.edu", "inactive2@school.edu"]
  
  # Map each inactive user to each supervised iPad
  storage_cleanup_pairs = flatten([
    for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ipads.items : [
      for user in local.inactive_users : {
        device_id           = device.id
        user_principal_name = user
      }
    ]
  ])
}

action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "storage_cleanup" {

  devices = local.storage_cleanup_pairs

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 6: Targeted user removal with CSV import
# ============================================================================
# Use case: Bulk user removal from specific devices based on CSV
locals {
  # Example: Reading from a CSV file (you would create this file)
  # CSV format: device_id,user_principal_name
  user_removal_list = [
    {
      device_id           = "12345678-1234-1234-1234-123456789abc"
      user_principal_name = "student1@school.edu"
    },
    {
      device_id           = "87654321-4321-4321-4321-ba9876543210"
      user_principal_name = "student2@school.edu"
    },
    {
      device_id           = "12345678-1234-1234-1234-123456789abc"
      user_principal_name = "student3@school.edu"
    }
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_delete_user_from_shared_apple_device" "bulk_removal" {

  devices = local.user_removal_list

  timeouts = {
    invoke = "15m"
  }
}
```

<!-- action schema generated by tfplugindocs -->
## Schema

### Required

- `devices` (Attributes List) List of device-user pairs specifying which users to delete from which Shared iPad devices. Each entry specifies a device ID and the user principal name of the user to delete from that device. You can delete different users from different devices in a single action.

Example:
```hcl
devices = [
  {
    device_id = "12345678-1234-1234-1234-123456789abc"
    user_principal_name = "student1@school.edu"
  },
  {
    device_id = "87654321-4321-4321-4321-ba9876543210"
    user_principal_name = "student2@school.edu"
  }
]
``` (see [below for nested schema](#nestedatt--devices))

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--devices"></a>
### Nested Schema for `devices`

Required:

- `device_id` (String) The managed device ID (GUID) of the Shared iPad from which to delete the user. Example: `12345678-1234-1234-1234-123456789abc`
- `user_principal_name` (String) The user principal name (UPN) of the user to delete from the Shared iPad. This is typically the user's email address or Managed Apple ID. Example: `student@school.edu` or `student@school.appleid`

**Important:** The user will be permanently removed from this device, and all their cached data will be deleted. The user's account in the cloud (Azure AD/Entra ID) is not affected.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


