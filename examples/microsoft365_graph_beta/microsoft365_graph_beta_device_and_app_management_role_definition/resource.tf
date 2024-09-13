resource "microsoft365_graph_beta_device_and_app_management_role_definition" "example" {
  display_name = "Custom Role - Device Management"
  description  = "This role allows management of device configurations and limited user read access"
  is_built_in  = false

  role_permissions {
    resource_actions {
      allowed_resource_actions = [
        "microsoft.graph.read",
        "microsoft.graph.deviceManagement.read",
        "microsoft.graph.deviceManagement.configurations.read",
        "microsoft.graph.deviceManagement.configurations.create",
        "microsoft.graph.deviceManagement.configurations.update"
      ]
      not_allowed_resource_actions = [
        "microsoft.graph.deviceManagement.configurations.delete",
        "microsoft.graph.user.write"
      ]
    }
  }

  role_scope_tag_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]

  timeouts {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}

data "microsoft365_graph_beta_device_and_app_management_role_definition" "helpdesk_admin" {
  display_name = "Helpdesk Administrator"
  is_built_in  = true
}

output "helpdesk_admin_role" {
  value = data.microsoft365_graph_beta_device_and_app_management_role_definition.helpdesk_admin
}

resource "microsoft365_graph_beta_device_and_app_management_role_definition" "helpdesk_assignment" {
  role_definition_id = data.microsoft365_graph_beta_device_and_app_management_role_definition.helpdesk_admin.id
  principal_id       = "00000000-0000-0000-0000-000000000003"  # ID of the user or group to assign the role to
  
  timeouts {
    create = "15m"
    read   = "5m"
    delete = "15m"
  }
}