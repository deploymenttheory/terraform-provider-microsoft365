# Minimal custom directory role: a display name and one set of allowed resource actions.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "application_reader" {
  display_name = "Application Reader"
  is_enabled   = true

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.directory/applications/standard/read"
      ]
    }
  ]
}
