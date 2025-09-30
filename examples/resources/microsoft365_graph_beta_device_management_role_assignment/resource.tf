# Example 1: Role Assignment with Specific Resource Scopes
resource "microsoft365_graph_beta_device_management_role_assignment" "specific_scopes" {
  display_name       = "Custom Assignment - Specific Scopes"
  description        = "Assignment to specific resource scopes"
  role_definition_id = "00000000-0000-0000-0000-000000000000"

  members = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]

  scope_configuration {
    type = "ResourceScopes"
    resource_scopes = [
      "00000000-0000-0000-0000-000000000003",
      "00000000-0000-0000-0000-000000000004"
    ]
  }
}

# Example 2: Role Assignment to All Licensed Users
resource "microsoft365_graph_beta_device_management_role_assignment" "all_licensed_users" {
  display_name       = "Policy Manager - All Licensed Users"
  description        = "Assignment to all licensed users"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be" # Policy and Profile manager

  members = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]

  scope_configuration {
    type = "AllLicensedUsers"
  }
}

# Example 3: Role Assignment to All Devices
resource "microsoft365_graph_beta_device_management_role_assignment" "all_devices" {
  display_name       = "Device Manager - All Devices"
  description        = "Assignment to all devices"
  role_definition_id = "0bd113fe-6be5-400c-a28f-ae5553f9c0be" # Policy and Profile manager

  members = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]

  scope_configuration {
    type = "AllDevices"
  }
}

# Example 4: Using with Custom Role Definition
resource "microsoft365_graph_beta_device_management_role_definition" "custom_role" {
  display_name = "Custom Device Manager"
  description  = "Custom role for device management"

  role_permissions {
    allowed_resource_actions = [
      "Microsoft.Intune/MobileApplications/Read",
      "Microsoft.Intune/MobileApplications/Create",
      "Microsoft.Intune/MobileApplications/Update"
    ]
  }
}

resource "microsoft365_graph_beta_device_management_role_assignment" "custom_role_assignment" {
  display_name       = "Custom Role Assignment"
  description        = "Assignment using custom role definition"
  role_definition_id = microsoft365_graph_beta_device_management_role_definition.custom_role.id

  members = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]

  scope_configuration {
    type = "AllLicensedUsers"
  }
}

# Example 5: Using Built-in Role Names (via constants)
locals {
  built_in_roles = {
    policy_manager      = "0bd113fe-6be5-400c-a28f-ae5553f9c0be"
    help_desk_operator  = "9e0cc482-82df-4ab2-a24c-0c23a3f52e1e"
    application_manager = "c1d9fcbb-cba5-40b0-bf6b-527006585f4b"
  }
}

resource "microsoft365_graph_beta_device_management_role_assignment" "help_desk_assignment" {
  display_name       = "Help Desk Team Assignment"
  description        = "Help desk operators with device access"
  role_definition_id = local.built_in_roles.help_desk_operator

  members = [
    "helpdesk-group@contoso.com"
  ]

  scope_configuration {
    type = "AllDevices"
  }
}