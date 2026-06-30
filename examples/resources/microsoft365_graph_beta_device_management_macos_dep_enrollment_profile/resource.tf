# Example 1: True zero-touch macOS enrollment (enroll WITHOUT user affinity / "userless")
#
# This is the genuine zero-touch path: no user authentication, no Company Portal,
# the device provisions with no human interaction. Setup Assistant panes are skipped
# and the desktop is gated until MDM finishes (await_device_configured).
#
# NOTE: Local admin/user account auto-creation is NOT available in this flow. Per
# Microsoft, account creation requires user affinity + Setup Assistant auth +
# await_device_configured (see Example 2).
#
# Setup Assistant panes are skipped via the individual *_disabled booleans below; the
# provider derives the read-only enabled_skip_keys array from them.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "zero_touch_userless" {
  display_name                 = "macOS DEP - Zero Touch (userless)"
  description                  = "Userless zero-touch macOS enrollment; skips Setup Assistant and gates the desktop until MDM finishes"
  requires_user_authentication = false
  supervised_mode_enabled      = true

  # Userless (and auto-advancing) profiles must be mandatory, or Graph rejects the request.
  is_mandatory = true

  # Gate the desktop until MDM configuration finishes (awaitDeviceConfigured)
  await_device_configured = true

  # Setup Assistant pane skip toggles (drive enabled_skip_keys)
  apple_id_disabled             = true
  apple_pay_disabled            = true
  terms_and_conditions_disabled = true
  diagnostics_disabled          = true
  display_tone_setup_disabled   = true
  siri_disabled                 = true
  file_vault_disabled           = true
  location_disabled             = true
  restore_blocked               = true
  screen_time_screen_disabled   = true
  icloud_storage_disabled       = true
  icloud_diagnostics_disabled   = true
  welcome_screen_disabled       = true

  # privacy_pane_disabled and registration_disabled work as boolean properties, but
  # their skip-key strings are rejected by Graph, so they are NOT added to
  # enabled_skip_keys (the provider handles this automatically).
  privacy_pane_disabled = true
  registration_disabled = true

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 2: User-affinity macOS enrollment with an auto-created local admin account (LAPS)
#
# Account auto-creation requires ALL of the following (enforced by Intune):
#   - user affinity (requires_user_authentication = true, Setup Assistant auth)
#   - await_device_configured = true
# The admin password is write-only/sensitive. When admin_account_password_rotation is
# set, Intune manages the password (LAPS-style) and rotates it automatically.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "user_affinity_with_local_admin" {
  display_name                 = "macOS DEP - User Affinity + Local Admin"
  description                  = "Setup Assistant with modern authentication; auto-creates a hidden, LAPS-managed local admin account"
  requires_user_authentication = true
  supervised_mode_enabled      = true

  # Required for local account creation to take effect
  await_device_configured = true

  apple_id_disabled    = true
  diagnostics_disabled = true
  siri_disabled        = true
  restore_blocked      = true

  # Auto-create the local admin account (managed local user)
  admin_account_user_name = "ladmin"
  admin_account_full_name = "Local Administrator"
  admin_account_password  = var.local_admin_password # sensitive / write-only
  hide_admin_account      = true

  # Optional automatic admin password rotation (LAPS-style)
  admin_account_password_rotation = {
    auto_rotation_period_in_days                     = 30
    on_retrieval_auto_rotate_password_enabled        = true
    on_retrieval_delay_auto_rotate_password_in_hours = 24
  }

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 3: Minimal profile
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "minimal" {
  display_name                 = "macOS DEP - Minimal"
  description                  = "Minimal macOS enrollment profile"
  requires_user_authentication = false
  is_mandatory                 = true # required because this profile is userless

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

variable "local_admin_password" {
  type      = string
  sensitive = true
}
