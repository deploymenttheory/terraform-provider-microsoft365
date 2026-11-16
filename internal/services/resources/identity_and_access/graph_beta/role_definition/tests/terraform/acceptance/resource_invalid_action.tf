# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}
# ==============================================================================
# Custom Directory Role Definition
# ==============================================================================

# Uses an action that does not exist in the microsoft.directory namespace. Microsoft Entra would
# silently rewrite it to applications/standard/read, so the provider rejects it before the write.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "invalid_action" {
  display_name = "acc-test-role-definition-invalid-${random_string.suffix.result}"
  is_enabled   = true

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.directory/applications/basic/read"
      ]
    }
  ]
}
