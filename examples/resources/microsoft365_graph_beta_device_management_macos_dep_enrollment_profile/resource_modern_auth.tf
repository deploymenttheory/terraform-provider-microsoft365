# Configuration verified against the Intune enrollment-profile wizard.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "modern_auth" {
  display_name                                               = "macOS DEP - modern auth"
  description                                                = "macOS DEP enrollment with modern auth"
  requires_user_authentication                               = true
  configuration_web_url                                      = true
  supervised_mode_enabled                                    = true
  support_department                                         = "IT Support"
  is_mandatory                                               = true
  profile_removal_disabled                                   = true
  location_disabled                                          = true
  support_phone_number                                       = "+1-555-0100"
  restore_blocked                                            = true
  apple_id_disabled                                          = true
  terms_and_conditions_disabled                              = true
  touch_id_disabled                                          = true
  apple_pay_disabled                                         = true
  siri_disabled                                              = true
  diagnostics_disabled                                       = true
  privacy_pane_disabled                                      = true
  screen_time_screen_disabled                                = true
  file_vault_disabled                                        = true
  icloud_diagnostics_disabled                                = true
  icloud_storage_disabled                                    = true
  choose_your_lock_screen_disabled                           = true
  display_tone_setup_disabled                                = true
  registration_disabled                                      = true
  accessibility_screen_disabled                              = true
  auto_unlock_with_watch_disabled                            = true
  await_device_configured                                    = true
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = false
  skip_primary_setup_account_creation                        = false
  set_primary_setup_account_as_regular_user                  = true
  dont_auto_populate_primary_account_info                    = false
  primary_account_full_name                                  = "Acceptance User"
  primary_account_user_name                                  = "accuser"
  enable_restrict_editing                                    = true
  admin_account_user_name                                    = "accadmin"
  admin_account_full_name                                    = "Acceptance Administrator"
  hide_admin_account                                         = true
  appearance_disabled                                        = true
  enable_lockdown_mode_disabled                              = true
  software_update_disabled                                   = true
  additional_privacy_settings_disabled                       = true
  update_completed_disabled                                  = true
  terms_of_address_disabled                                  = true
  intelligence_disabled                                      = true
  os_showcase_disabled                                       = true
  app_store_disabled                                         = true
  welcome_screen_disabled                                    = false
  pass_code_disabled                                         = false
  zoom_disabled                                              = false

  admin_account_password_rotation = {
    auto_rotation_period_in_days = 90
  }

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
