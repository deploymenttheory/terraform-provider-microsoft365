resource "microsoft365_graph_beta_device_and_app_management_role_definition" "example" {
  display_name                = "Custom Intune Role Definition"
  description                 = "This is a custom Intune role definition for device and app management"
  is_built_in                 = false
  is_built_in_role_definition = false

  role_permissions {
    actions = ["microsoft.intune/"]
    resource_actions {
      allowed_resource_actions = [
        "microsoft.intune/deviceConfigurations/read",
        "microsoft.intune/deviceConfigurations/basic/read",
        "microsoft.intune/deviceConfigurations/assign/action",
        "microsoft.intune/managedDevices/read",
        "microsoft.intune/managedDevices/resetPasscode/action",
        "microsoft.intune/managedApps/read",
        "microsoft.intune/mobileApps/read",
        "microsoft.intune/mobileApps/assign/action"
      ]
      not_allowed_resource_actions = [
        "microsoft.intune/deviceConfigurations/create",
        "microsoft.intune/deviceConfigurations/delete",
        "microsoft.intune/managedDevices/delete",
        "microsoft.intune/managedApps/wipe/action"
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