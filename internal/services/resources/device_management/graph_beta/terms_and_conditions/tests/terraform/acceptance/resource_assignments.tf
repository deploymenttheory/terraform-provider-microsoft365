resource "microsoft365_graph_beta_device_management_terms_and_conditions" "assignments" {
  display_name         = "Test Assignments Terms and Conditions"
  description          = "Terms and conditions policy with assignments for acceptance testing"
  title                = "Company Terms with Assignments"
  body_text            = "These are the terms and conditions that will be assigned to specific groups."
  acceptance_statement = "I accept these terms and conditions"
  version              = 1

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}