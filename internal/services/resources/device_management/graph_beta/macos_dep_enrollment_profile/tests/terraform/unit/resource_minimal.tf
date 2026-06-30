resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "minimal" {
  display_name                 = "Test Minimal macOS DEP Enrollment Profile - Unique"
  description                  = "Minimal macOS DEP enrollment profile for unit testing"
  requires_user_authentication = false

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
