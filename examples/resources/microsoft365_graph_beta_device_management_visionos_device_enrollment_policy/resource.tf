# Example 1: Minimal zero-touch visionOS ADE enrollment profile (no user affinity).
resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "minimal" {
  name = "visionOS ADE - Minimal"

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}

# Example 2: Maximal visionOS ADE enrollment profile exercising the full settings tree - await
# configuration, locked enrollment, and every Setup Assistant screen toggle. requires_user_authentication
# is omitted: visionOS ADE only supports enrollment without user affinity, so it always defaults to false.
resource "microsoft365_graph_beta_device_management_visionos_device_enrollment_policy" "maximal" {
  name        = "visionOS ADE - Maximal"
  description = "visionOS ADE enrollment policy exercising the full settings tree"

  # Uncomment to target a specific Apple ABM/ASM token when the tenant has more than one;
  # otherwise this is auto-resolved.
  # dep_onboarding_settings_id = "00000000-0000-0000-0000-000000000000"

  # Makes this the default visionOS enrollment profile for the DEP token via the
  # setDefaultProfile action. Only one policy per DEP token can be the default; setting this to
  # true elsewhere supersedes this assignment. There is no "unassign" action - see the resource
  # documentation.
  is_default_policy_assignment = true

  await_device_configured   = true
  locked_enrollment_enabled = true

  support_department   = "IT Support"
  support_phone_number = "+1-555-0100"

  # Setup Assistant screen toggles - true hides the pane, false (default) shows it.
  apple_id_disabled               = true
  apple_pay_disabled              = true
  diagnostics_disabled            = true
  get_started_screen_disabled     = false
  apple_intelligence_disabled     = false
  location_services_disabled      = false
  passcode_disabled               = true
  privacy_pane_disabled           = true
  screen_time_screen_disabled     = true
  siri_disabled                   = true
  software_update_screen_disabled = false
  terms_and_conditions_disabled   = false
  tips_screen_disabled            = true
  touch_id_disabled               = false

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
