# LW902: A trigger and scope based workflow needs both scope and trigger
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw902_invalid" {
  category     = "joiner"
  display_name = "unit-test-lw902: scope without trigger"

  execution_conditions = {
    scope = {
      rule = "(department eq 'acc-test-none')"
    }
  }

  tasks = [
    {
      task_definition_id = "6fc52c9d-398b-4305-9763-15f42c1676fc"
      display_name       = "Enable user account"
    }
  ]
}
