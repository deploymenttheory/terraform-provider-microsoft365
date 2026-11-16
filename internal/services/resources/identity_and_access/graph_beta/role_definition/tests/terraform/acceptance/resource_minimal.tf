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

# Minimal custom directory role: one permission block and the required attributes only.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "minimal" {
  display_name = "acc-test-role-definition-minimal-${random_string.suffix.result}"
  is_enabled   = true

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.directory/applications/standard/read"
      ]
    }
  ]
}
