resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-tandc-transition-it"
  description      = "Test group for IT support staff used in terms and conditions transition test"
  mail_nickname    = "tandc-trans-it"
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
  display_name     = "acc-test-tandc-transition-dm"
  description      = "Test group for device management staff used in terms and conditions transition test"
  mail_nickname    = "tandc-trans-dm"
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

resource "microsoft365_graph_beta_device_management_terms_and_conditions" "transition" {
  display_name         = "acc-test-terms-and-conditions-transition"
  description          = "Configuration that transitions from minimal to maximal for acceptance testing"
  title                = "Complete Company Terms and Conditions"
  body_text            = "These are the comprehensive terms and conditions that all users must read and accept before accessing company resources."
  acceptance_statement = "I have read and agree to abide by all terms and conditions outlined above"
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

