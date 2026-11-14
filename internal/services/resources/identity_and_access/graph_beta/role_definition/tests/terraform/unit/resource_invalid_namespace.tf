# Schema validation: an action outside the microsoft.directory namespace is rejected at plan time.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "invalid_namespace" {
  display_name = "unit-test-role-definition-minimal"
  is_enabled   = true

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.intune/deviceConfigurations/read"
      ]
    }
  ]
}
