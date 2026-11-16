# LW003: Leaver workflow that runs on the user's last day
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow" "lw003_leaver" {
  category     = "leaver"
  display_name = "unit-test-lw003: leaver last day"
  description  = "Disable the account and remove licenses on the last day"

  execution_conditions = {
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
      task_definition_id = "1dfdfcc7-52fa-4c2e-bf3a-e3919cc12950"
      display_name       = "Disable user account"
      arguments = [
        {
          name  = "disableOnPremisesAccount"
          value = "false"
        }
      ]
    },
    {
      task_definition_id = "8fa97d28-3e52-4985-b3a9-a1126f9b8b4e"
      display_name       = "Remove all licenses for user"
      continue_on_error  = true
    },
    {
      task_definition_id = "509589a4-0466-4471-829e-49c5e502bdee"
      display_name       = "Revoke all refresh tokens for user"
    }
  ]
}
