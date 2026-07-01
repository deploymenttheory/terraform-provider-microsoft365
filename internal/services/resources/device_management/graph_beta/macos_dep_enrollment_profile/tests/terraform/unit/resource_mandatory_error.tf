# Userless enrollment (requires_user_authentication = false) without is_mandatory = true
# should fail validation at plan time.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "mandatory_error" {
  display_name                 = "Test Mandatory Rule macOS DEP Enrollment Profile - Unique"
  description                  = "Userless profile missing is_mandatory; should fail validation"
  requires_user_authentication = false
  is_mandatory                 = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
