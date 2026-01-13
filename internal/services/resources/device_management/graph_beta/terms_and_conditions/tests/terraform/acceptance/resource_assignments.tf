# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-tandc-assignments-it"
  description      = "Test group for IT support staff used in terms and conditions assignments"
  mail_nickname    = "tandc-assign-it"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "acc-test-tandc-assignments-dm"
  description      = "Test group for device management staff used in terms and conditions assignments"
  mail_nickname    = "tandc-assign-dm"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "microsoft365_graph_beta_device_management_terms_and_conditions" "assignments" {
  display_name         = "acc-test-terms-and-conditions-assignments"
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

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_1,
    microsoft365_graph_beta_groups_group.acc_test_group_2
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
