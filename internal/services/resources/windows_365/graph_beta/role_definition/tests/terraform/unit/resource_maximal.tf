resource "microsoft365_graph_beta_device_management_role_definition" "maximal_custom" {
  display_name                = "Test Maximal Custom Role Definition - Unique"
  description                 = "Comprehensive custom role definition for testing with all features"
  is_built_in_role_definition = false
  is_built_in                 = false
  role_scope_tag_ids          = ["0", "1"]

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.Intune_ManagedDevices_Read",
        "Microsoft.Intune_ManagedDevices_Update",
        "Microsoft.Intune_ManagedDevices_Delete",
        "Microsoft.Intune_DeviceConfigurations_Read",
        "Microsoft.Intune_DeviceConfigurations_Create",
        "Microsoft.Intune_DeviceConfigurations_Update",
        "Microsoft.Intune_DeviceConfigurations_Delete",
        "Microsoft.Intune_DeviceConfigurations_Assign",
        "Microsoft.Intune_Audit_Read",
        "Microsoft.Intune_Organization_Read"
      ]
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}