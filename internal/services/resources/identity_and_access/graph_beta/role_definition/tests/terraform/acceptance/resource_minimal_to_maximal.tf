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

# Minimal role (same resource address as resource_minimal.tf) updated in place to every optional
# attribute except template_id, which would force replacement.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "minimal" {
  display_name = "acc-test-role-definition-updated-${random_string.suffix.result}"
  description  = "Acceptance test custom role updated in place"
  is_enabled   = false
  version      = "2"

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.directory/applications/standard/read",
        "microsoft.directory/groups/standard/read"
      ]
    },
    {
      allowed_resource_actions = [
        "microsoft.directory/devices/standard/read"
      ]
    }
  ]
}
