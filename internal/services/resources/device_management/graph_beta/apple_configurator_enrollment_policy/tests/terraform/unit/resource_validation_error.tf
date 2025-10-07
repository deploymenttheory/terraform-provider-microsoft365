# This configuration should fail validation because both authentication methods are set to true
resource "microsoft365_graph_beta_device_management_apple_configurator_enrollment_policy" "validation_error" {
  display_name                                               = "Test Validation Error Apple Configurator Enrollment Policy - Unique"
  description                                                = "Apple configurator enrollment policy that should fail validation"
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