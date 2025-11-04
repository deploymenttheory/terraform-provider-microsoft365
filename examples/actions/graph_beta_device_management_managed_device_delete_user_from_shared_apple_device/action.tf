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

