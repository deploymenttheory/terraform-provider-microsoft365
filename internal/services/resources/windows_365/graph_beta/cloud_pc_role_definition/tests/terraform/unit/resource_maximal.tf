resource "microsoft365_graph_beta_windows_365_cloud_pc_role_definition" "maximal_custom" {
  display_name = "unit-test-cloud-pc-role-definition-maximal"
  description  = "Comprehensive custom role definition for testing with all features"

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.CloudPC/CloudPCs/Read",
        "Microsoft.CloudPC/CloudPCs/Reboot",
        "Microsoft.CloudPC/CloudPCs/Resize",
        "Microsoft.CloudPC/DeviceImages/Read",
        "Microsoft.CloudPC/DeviceImages/Create",
        "Microsoft.CloudPC/ProvisioningPolicies/Read",
        "Microsoft.CloudPC/ProvisioningPolicies/Create",
        "Microsoft.CloudPC/UserSettings/Read",
        "Microsoft.CloudPC/AuditData/Read",
        "Microsoft.CloudPC/OrganizationSettings/Read"
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