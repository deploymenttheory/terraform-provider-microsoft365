resource "microsoft365_graph_beta_identity_and_access_directory_settings" "password_rule_settings" {
  template_type               = "Password Rule Settings"
  overwrite_existing_settings = true

  password_rule_settings {
    banned_password_check_on_premises_mode   = "Enforce"
    enable_banned_password_check_on_premises = true
    enable_banned_password_check             = true
    lockout_duration_in_seconds              = 120
    lockout_threshold                        = 5
    banned_password_list                     = "password123\tcompany123\tadmin123\twelcome123"
  }

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

