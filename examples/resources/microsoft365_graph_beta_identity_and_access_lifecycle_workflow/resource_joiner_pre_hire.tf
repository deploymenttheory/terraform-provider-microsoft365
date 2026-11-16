# Joiner workflow that runs 7 days before the employee hire date
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "joiner" {
  category              = "joiner"
  display_name          = "Pre-hire onboarding"
  description           = "Enable the account and send a welcome email before day one"
  is_enabled            = false
  is_scheduling_enabled = false

  execution_conditions = {
    scope = {
      rule = "(department eq 'Marketing')"
    }
    trigger = {
      time_based_attribute = {
        attribute      = "employeeHireDate"
        offset_in_days = -7
      }
    }
  }

  tasks = [
    {
      task_definition_id = "6fc52c9d-398b-4305-9763-15f42c1676fc"
      display_name       = "Enable user account"
      arguments = [
        {
          name  = "enableOnPremisesAccount"
          value = "false"
        }
      ]
    },
    {
      task_definition_id = "70b29d51-b59a-4773-9280-8841dfd3f2ea"
      display_name       = "Send welcome email"
      continue_on_error  = true
    }
  ]
}
