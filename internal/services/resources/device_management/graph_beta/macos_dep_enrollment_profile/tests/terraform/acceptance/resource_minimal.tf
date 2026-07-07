resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "minimal" {
  display_name                 = "acc-test-macos-dep-enrollment-profile-minimal"
  description                  = "macOS DEP enrollment profile minimal acceptance test"
  requires_user_authentication = false
  is_mandatory                 = true

  timeouts = {
    create = "10s"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
