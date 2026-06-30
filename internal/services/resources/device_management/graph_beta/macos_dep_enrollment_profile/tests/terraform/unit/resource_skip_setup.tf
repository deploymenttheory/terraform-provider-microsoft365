resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "skip_setup" {
  display_name                 = "Test Skip Setup macOS DEP Enrollment Profile - Unique"
  description                  = "macOS DEP enrollment profile with setup assistant skip and admin account"
  requires_user_authentication = true
  supervised_mode_enabled      = true

  await_device_configured = true

  # enabled_skip_keys is derived from these booleans (computed, read-only)
  apple_id_disabled             = true
  terms_and_conditions_disabled = true
  diagnostics_disabled          = true
  siri_disabled                 = true
  file_vault_disabled           = true
  location_disabled             = true
  welcome_screen_disabled       = true

  admin_account_user_name = "localadmin"
  admin_account_full_name = "Local Administrator"
  admin_account_password  = "SuperSecretP@ssw0rd!"
  hide_admin_account      = true

  admin_account_password_rotation = {
    auto_rotation_period_in_days                     = 30
    on_retrieval_auto_rotate_password_enabled        = true
    on_retrieval_delay_auto_rotate_password_in_hours = 24
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
