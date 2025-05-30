# Assignment to all users
resource "microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment" "all_users" {
  apple_user_initiated_enrollment_profile_id = microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile.example.id
  
  target = {
    target_type = "allUsers"
  }

  # Optional timeouts block
  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}

# Assignment to a specific Entra ID group
resource "microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment" "corporate_users" {
  apple_user_initiated_enrollment_profile_id = microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile.example.id
  
  target = {
    target_type = "group"
    group_id = "11111111-2222-3333-4444-555555555555"
  }
}

# Exclusion assignment
resource "microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment" "exclude_group" {
  apple_user_initiated_enrollment_profile_id = microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile.example.id
  
  target = {
    target_type = "exclusionGroup"
    group_id = "66666666-7777-8888-9999-000000000000"
  }
}

# Assignment to a specific user
resource "microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment" "specific_user" {
  apple_user_initiated_enrollment_profile_id = microsoft365_graph_beta_device_management_apple_user_initiated_enrollment_profile.example.id
  
  target = {
    target_type = "user"
    entra_object_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  }
}