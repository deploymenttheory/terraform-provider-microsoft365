# LW901: on_demand_execution_only cannot be combined with scope and trigger
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw901_invalid" {
  category     = "leaver"
  display_name = "unit-test-lw901: on demand with scope"

  execution_conditions = {
    on_demand_execution_only = true
    scope = {
      rule = "(department eq 'acc-test-none')"
    }
    trigger = {
      time_based_attribute = {
        attribute      = "employeeLeaveDateTime"
        offset_in_days = 0
      }
    }
  }

  tasks = [
    {
      task_definition_id = "509589a4-0466-4471-829e-49c5e502bdee"
      display_name       = "Revoke all refresh tokens for user"
    }
  ]
}
