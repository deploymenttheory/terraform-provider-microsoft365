# This configuration should fail validation because both authentication methods are set to true
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "validation_error" {
  display_name                                               = "Test Validation Error macOS DEP Enrollment Profile - Unique"
  description                                                = "macOS DEP enrollment profile that should fail validation"
  requires_user_authentication                               = true
  enable_authentication_via_company_portal                   = true
  require_company_portal_on_setup_assistant_enrolled_devices = true

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
