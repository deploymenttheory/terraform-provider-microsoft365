# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-tandc-transition-it"
  description      = "Test group for IT support staff used in terms and conditions transition test"
  mail_nickname    = "tandc-transition-it"
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
  mail_nickname    = "tandc-transition-dm"
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
  title                = "Complete Company Terms for Transition"
  body_text            = "These are the comprehensive terms and conditions for transition testing."
  acceptance_statement = "I accept all terms and conditions for this transition test"
  role_scope_tag_ids   = ["0"]

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

