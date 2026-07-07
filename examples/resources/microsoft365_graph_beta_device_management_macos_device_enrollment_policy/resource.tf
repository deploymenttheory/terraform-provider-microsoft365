# Example 1: Minimal zero-touch macOS ADE enrollment profile (no user affinity, no local
# account creation). await_device_configured must be false when admin_account is omitted.
resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "minimal" {
  name = "macOS ADE - Minimal"

  requires_user_authentication = false
  await_device_configured      = false

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}

# Example 2: Maximal macOS ADE enrollment profile exercising the full settings tree - user
# authentication via Company Portal, a LAPS-style local admin account with a separate primary
# account, locked enrollment, and every Setup Assistant screen toggle.
resource "microsoft365_graph_beta_device_management_macos_device_enrollment_policy" "maximal" {
  name        = "macOS ADE - Maximal"
  description = "macOS ADE enrollment policy exercising the full settings tree"

  # Uncomment to target a specific Apple ABM/ASM token when the tenant has more than one;
  # otherwise this is auto-resolved.
  # dep_onboarding_settings_id = "00000000-0000-0000-0000-000000000000"

  # Makes this the default macOS enrollment profile for the DEP token via the setDefaultProfile
  # action. Only one policy per DEP token can be the default; setting this to true elsewhere
  # supersedes this assignment. There is no "unassign" action - see the resource documentation.
  is_default_policy_assignment = true

  requires_user_authentication                               = true
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = true

  await_device_configured = true

  admin_account = {
    create_local_admin_account = true
    user_name                  = "localadmin"
    full_name                  = "Local Administrator"
    hide_account               = true
    password_rotation_in_days  = 90

    create_local_primary_account = true
    primary_account = {
      prefill_account_info = true
      restrict_editing     = true
      user_name            = "primaryuser"
      full_name            = "Primary User"
    }
  }

  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  # Setup Assistant screen toggles - true hides the pane, false (default) shows it.
  location_services_disabled                = false
  restore_disabled                          = true
  apple_id_disabled                         = true
  terms_and_conditions_disabled             = false
  touch_id_disabled                         = false
  apple_pay_disabled                        = true
  siri_disabled                             = true
  diagnostics_disabled                      = true
  file_vault_disabled                       = false
  icloud_diagnostics_disabled               = true
  icloud_storage_disabled                   = true
  display_tone_setup_disabled               = false
  screen_time_screen_disabled               = true
  privacy_pane_disabled                     = true
  accessibility_screen_disabled             = false
  auto_unlock_with_watch_disabled           = true
  lockdown_mode_disabled                    = true
  software_update_screen_disabled           = false
  software_update_completed_screen_disabled = true
  terms_of_address_screen_disabled          = true
  apple_intelligence_disabled               = false
  os_showcase_screen_disabled               = true
  app_store_disabled                        = false

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
