resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "explicit_token" {
  dep_onboarding_settings_id   = "7019f829-33ee-4fc0-89b6-ac7435d71e1e"
  display_name                 = "Test Explicit Token macOS DEP Enrollment Profile - Unique"
  description                  = "macOS DEP enrollment profile pinned to an explicit ABM/ADE token"
  requires_user_authentication = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
