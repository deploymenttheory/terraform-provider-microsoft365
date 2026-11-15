resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "first" {
  display_name     = "acc-test-auth-methods-email-1-${random_string.suffix.result}"
  mail_nickname    = "acc-test-amp-email-1-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Acceptance test group targeted by the email one-time passcode method"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "second" {
  display_name     = "acc-test-auth-methods-email-2-${random_string.suffix.result}"
  mail_nickname    = "acc-test-amp-email-2-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Acceptance test group targeted by the email one-time passcode method"
  hard_delete      = true
}

resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "test" {
  email = {
    state                              = "enabled"
    allow_external_id_to_use_email_otp = "default"
    include_targets = [
      {
        id          = microsoft365_graph_beta_groups_group.second.id
        target_type = "group"
      }
    ]
    exclude_targets = [
      {
        id          = microsoft365_graph_beta_groups_group.first.id
        target_type = "group"
      }
    ]
  }
}
