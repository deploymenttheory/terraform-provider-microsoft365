resource "microsoft365_graph_beta_device_management_role_definition" "test" {
  display_name                = "Test Acceptance Built-in Role Definition - Updated"
  description                 = "Updated built-in description for acceptance testing"
  is_built_in_role_definition = true
  is_built_in                 = true
  built_in_role_name          = "Endpoint Security Manager"
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
        "Microsoft.Intune_DeviceCompliancePolices_Read",
        "Microsoft.Intune_DeviceCompliancePolices_Create",
        "Microsoft.Intune_DeviceCompliancePolices_Update",
        "Microsoft.Intune_DeviceCompliancePolices_Delete",
        "Microsoft.Intune_DeviceCompliancePolices_Assign",
        "Microsoft.Intune_Audit_Read"
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