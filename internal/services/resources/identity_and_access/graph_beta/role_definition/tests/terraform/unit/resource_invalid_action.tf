# Uses an action that does not exist in the microsoft.directory namespace. Microsoft Entra would
# silently rewrite it, so the provider rejects it before the write.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "invalid_action" {
  display_name = "unit-test-role-definition-minimal"
  is_enabled   = true

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.directory/applications/basic/read"
      ]
    }
  ]
}
