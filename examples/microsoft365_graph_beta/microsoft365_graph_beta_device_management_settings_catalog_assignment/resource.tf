# Example 1: Assign configuration policy to all devices
resource "graph_beta_device_management_settings_catalog_assignment" "all_devices_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "allDevices"
  }
}

# Example 2: Assign configuration policy to all licensed users
resource "graph_beta_device_management_settings_catalog_assignment" "all_users_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "allLicensedUsers"
  }
}

# Example 3: Assign configuration policy to a specific Entra ID group
resource "graph_beta_device_management_settings_catalog_assignment" "group_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "groupAssignment"
    group_id    = "87654321-4321-4321-4321-210987654321"
  }
}

# Example 4: Assign configuration policy with group exclusion
resource "graph_beta_device_management_settings_catalog_assignment" "exclusion_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "exclusionGroupAssignment"
    group_id    = "87654321-4321-4321-4321-210987654321"
  }
}

# Example 5: Assign configuration policy to SCCM collection
resource "graph_beta_device_management_settings_catalog_assignment" "sccm_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type   = "configurationManagerCollection"
    collection_id = "SMS00000001"  # Default SMS collection or use custom like "MEM12345678"
  }
}

# Example 6: Group assignment with include filter
resource "graph_beta_device_management_settings_catalog_assignment" "group_with_include_filter" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type                                     = "groupAssignment"
    group_id                                        = "87654321-4321-4321-4321-210987654321"
    device_and_app_management_assignment_filter_id   = "11111111-2222-3333-4444-555555555555"
    device_and_app_management_assignment_filter_type = "include"
  }
}

# Example 7: Group assignment with exclude filter
resource "graph_beta_device_management_settings_catalog_assignment" "group_with_exclude_filter" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type                                     = "groupAssignment"
    group_id                                        = "87654321-4321-4321-4321-210987654321"
    device_and_app_management_assignment_filter_id   = "11111111-2222-3333-4444-555555555555"
    device_and_app_management_assignment_filter_type = "exclude"
  }
}

# Example 8: Assignment from policy sets
resource "graph_beta_device_management_settings_catalog_assignment" "policy_set_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  source             = "policySets"
  source_id          = "99999999-8888-7777-6666-555555555555"
  
  target {
    target_type = "groupAssignment"
    group_id    = "87654321-4321-4321-4321-210987654321"
  }
}

# Example 9: Assignment with custom timeouts
resource "graph_beta_device_management_settings_catalog_assignment" "assignment_with_timeouts" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "allDevices"
  }
  
  timeouts {
    create = "5m"
    read   = "3m"
    update = "5m"
    delete = "3m"
  }
}

# Example 10: Multiple assignments for the same policy
resource "graph_beta_device_management_settings_catalog_assignment" "primary_group" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "groupAssignment"
    group_id    = "primary-group-id-1234-1234-1234-123456789012"
  }
}

resource "graph_beta_device_management_settings_catalog_assignment" "secondary_group" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "groupAssignment"
    group_id    = "secondary-group-id-5678-5678-5678-567856785678"
  }
}

resource "graph_beta_device_management_settings_catalog_assignment" "exclude_test_group" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "exclusionGroupAssignment"
    group_id    = "test-group-id-9999-9999-9999-999999999999"
  }
}

# Example 11: Using data sources for dynamic assignment
data "azuread_group" "it_department" {
  display_name = "IT Department"
}

resource "graph_beta_device_management_settings_catalog_assignment" "it_department_assignment" {
  settings_catalog_id = "12345678-1234-1234-1234-123456789012"
  
  target {
    target_type = "groupAssignment"
    group_id    = data.azuread_group.it_department.object_id
  }
}

