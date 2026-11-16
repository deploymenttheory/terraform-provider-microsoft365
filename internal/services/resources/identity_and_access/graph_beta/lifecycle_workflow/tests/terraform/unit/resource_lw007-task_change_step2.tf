# LW007: Task and trigger change (createNewVersion) - step 2
# Reorders and replaces tasks and moves the trigger offset; the workflow is updated in place.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw007_task_change" {
  category     = "joiner"
  display_name = "unit-test-lw007: task change"

  execution_conditions = {
    scope = {
      rule = "(department eq 'acc-test-none')"
    }
    trigger = {
      time_based_attribute = {
        attribute      = "employeeHireDate"
        offset_in_days = 3
      }
    }
  }

  tasks = [
    {
      task_definition_id = "70b29d51-b59a-4773-9280-8841dfd3f2ea"
      display_name       = "Send welcome email"
      continue_on_error  = true
      arguments = [
        {
          name  = "locale"
          value = "en-us"
        }
      ]
    },
    {
      task_definition_id = "6fc52c9d-398b-4305-9763-15f42c1676fc"
      display_name       = "Enable user account"
      is_enabled         = false
    },
    {
      task_definition_id = "3C860712-2D37-42A4-928F-5C93935D26A1"
      display_name       = "Send onboarding reminder email"
    }
  ]
}
