resource "microsoft365_graph_beta_device_and_app_management_role_definition" "example" {
  display_name                = "android intune"
  description                 = "Custom role for Intune device administration with limited permissions"
  is_built_in_role_definition = false

  # Scope tags for the role definition
  role_scope_tag_ids = ["9", "8"]

  role_permissions = [
    {
      allowed_resource_actions = [
        "Microsoft.Intune_MicrosoftStoreForBusiness_Read",
        "Microsoft.Intune_AndroidSync_Read",
        "Microsoft.Intune_AndroidSync_UpdateApps",
        "Microsoft.Intune_AndroidSync_UpdateOnboarding",
        "Microsoft.Intune_AndroidSync_UpdateEnrollmentProfiles"
      ]
    }
  ]

  # Optional Timeout settings  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}