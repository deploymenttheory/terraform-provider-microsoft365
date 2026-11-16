# LW004: Minimal on-demand workflow with a single task and only required attributes
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw004_minimal" {
  category     = "leaver"
  display_name = "unit-test-lw004: minimal on demand"

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
