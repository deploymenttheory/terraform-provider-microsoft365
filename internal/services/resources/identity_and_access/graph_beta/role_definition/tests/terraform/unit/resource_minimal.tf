# Minimal custom directory role: one permission block and the required attributes only.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "minimal" {
  display_name = "unit-test-role-definition-minimal"
  is_enabled   = true

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.directory/applications/standard/read"
      ]
    }
  ]
}
