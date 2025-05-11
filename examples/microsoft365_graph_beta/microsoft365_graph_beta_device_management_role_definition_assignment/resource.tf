
resource "microsoft365_graph_beta_device_management_role_definition_assignment" "resource_scope_example" {
  # You can reference either a role definition ID or use a built-in role name
  role_definition_id = microsoft365_graph_beta_device_management_role_definition.example.id
  # OR
  # built_in_role_name = "Help Desk Operator"

  display_name = "DevOps Team Assignment"
  description  = "Assigns Intune administration capabilities to DevOps team"

  # Scope type defines the target of this assignment
  scope_type = "resourceScope" # One of: "allDevices", "allLicensedUsers", "allDevicesAndLicensedUsers", "resourceScope"

  # If scope_type is "resourceScope", you need to specify scope members
  scope_members = [
    "11111111-2222-3333-4444-555555555555",
    "11111111-2222-3333-4444-555555555555",
    "11111111-2222-3333-4444-555555555555",
  ]

  # You can also define specific resource scopes
  resource_scopes = [
    "11111111-2222-3333-4444-555555555555",
    "11111111-2222-3333-4444-555555555555",
    "11111111-2222-3333-4444-555555555555"
  ]

  # Optional Timeout settings  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_role_definition_assignment" "all_devices_example" {
  # You can reference either a role definition ID or use a built-in role name
  role_definition_id = microsoft365_graph_beta_device_management_role_definition.example.id
  # OR
  # built_in_role_name = "Help Desk Operator"

  display_name = "1st Line Support Assignment"
  description  = "Assigns Intune administration capabilities to 1st Line Support team"

  # Scope type defines the target of this assignment
  scope_type = "allDevices" # One of: "allDevices", "allLicensedUsers", "allDevicesAndLicensedUsers", "resourceScope"

  resource_scopes = [
    "11111111-2222-3333-4444-555555555555",
    "11111111-2222-3333-4444-555555555555",
    "11111111-2222-3333-4444-555555555555"
  ]

  # Optional Timeout settings  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}