# LW900: Invalid category is rejected at plan time
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw900_invalid" {
  category     = "extensibility"
  display_name = "unit-test-lw900: invalid category"

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
