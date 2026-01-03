resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-tandc-it"
  description      = "Test group for IT support staff used in terms and conditions assignments"
  mail_nickname    = "tandc-it-test"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "acc-test-tandc-dm"
  description      = "Test group for device management staff used in terms and conditions assignments"
  mail_nickname    = "tandc-dm-test"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "microsoft365_graph_beta_device_management_terms_and_conditions" "test" {
  display_name         = "acc-test-terms-and-conditions-maximal"
  description          = "Updated description for acceptance testing"
  title                = "Complete Company Terms and Conditions"
  body_text            = "These are the comprehensive terms and conditions that all users must read and accept before accessing company resources."
  acceptance_statement = "I have read and agree to abide by all terms and conditions outlined above"
  version              = 2
  role_scope_tag_ids   = ["0", "1"]

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
