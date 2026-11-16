# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Lifecycle Workflow
# ==============================================================================

# LW004: Minimal on-demand workflow with a single task and only required attributes
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw004_minimal" {
  category     = "leaver"
  display_name = "acc-test-lw004: minimal on demand ${random_string.suffix.result}"

  execution_conditions = {
    on_demand_execution_only = true
  }

  tasks = [
    {
      task_definition_id = "509589a4-0466-4471-829e-49c5e502bdee"
      display_name       = "Revoke all refresh tokens for user"
    }
  ]
}
