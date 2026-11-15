# Enable email one-time passcodes for one group only and leave every other method unchanged.
resource "microsoft365_graph_beta_groups_group" "email_otp" {
  display_name     = "Email OTP users"
  mail_nickname    = "email-otp-users"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "this" {
  email = {
    state                              = "enabled"
    allow_external_id_to_use_email_otp = "default"
    include_targets = [
      {
        id          = microsoft365_graph_beta_groups_group.email_otp.id
        target_type = "group"
      }
    ]
    exclude_targets = []
  }
}
