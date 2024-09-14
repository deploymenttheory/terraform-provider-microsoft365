resource "microsoft365_graph_beta_device_and_app_management_role_definition" "example" {
  display_name                = "Custom Intune Role Definition"
  description                 = "This is a custom Intune role definition for device and app management"
  is_built_in                 = false
  is_built_in_role_definition = false

  role_permissions {
    actions = ["microsoft.intune/"]
    resource_actions {
      allowed_resource_actions = [
        "Microsoft.Intune/MobileApps/Read",
        "Microsoft.Intune/TermsAndConditions/Read",
        "Microsoft.Intune/ManagedApps/Read",
        "Microsoft.Intune/ManagedDevices/Read",
        "Microsoft.Intune/DeviceConfigurations/Read",
        "Microsoft.Intune/TelecomExpenses/Read",
        "Microsoft.Intune/Organization/Read",
        "Microsoft.Intune/RemoteTasks/RebootNow",
        "Microsoft.Intune/RemoteTasks/RemoteLock"
      ]
      not_allowed_resource_actions = [
      ]
    }
  }

  role_scope_tag_ids = [
    "scope_tag_1",
    "scope_tag_2"
  ]

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

data "microsoft365_graph_device_and_app_management_role_definition" "helpdesk_admin" {
  display_name = "Helpdesk Administrator"
  is_built_in  = true
}

output "helpdesk_admin_role" {
  value = data.microsoft365_graph_device_and_app_management_role_definition.helpdesk_admin
}