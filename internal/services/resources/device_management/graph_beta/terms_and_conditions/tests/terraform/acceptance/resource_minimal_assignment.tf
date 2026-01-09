resource "microsoft365_graph_beta_device_management_terms_and_conditions" "minimal_assignment" {
  display_name         = "acc-test-terms-and-conditions-minimal-assignment"
  description          = "Terms and conditions with minimal assignment for acceptance testing"
  title                = "Company Terms with Minimal Assignment"
  body_text            = "These are the terms and conditions with a single assignment."
  acceptance_statement = "I accept these terms and conditions"

  assignments = [
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

