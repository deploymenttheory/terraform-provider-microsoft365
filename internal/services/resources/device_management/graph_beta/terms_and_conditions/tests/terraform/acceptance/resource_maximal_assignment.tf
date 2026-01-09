resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-tandc-maxassign-it"
  description      = "Test group for IT support staff used in terms and conditions maximal assignment test"
  mail_nickname    = "tandc-maxassign-it"
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
  display_name     = "acc-test-tandc-maxassign-dm"
  description      = "Test group for device management staff used in terms and conditions maximal assignment test"
  mail_nickname    = "tandc-maxassign-dm"
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

resource "microsoft365_graph_beta_device_management_terms_and_conditions" "maximal_assignment" {
  display_name         = "acc-test-terms-and-conditions-maximal-assignment"
  description          = "Terms and conditions with comprehensive assignments for acceptance testing"
  title                = "Company Terms with Maximal Assignments"
  body_text            = "These are the terms and conditions that will be assigned to specific groups."
  acceptance_statement = "I accept these terms and conditions"

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

