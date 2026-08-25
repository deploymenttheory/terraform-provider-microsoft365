# ==============================================================================
# Validation: uninstall_on_device_removal set alongside an uninstall intent
#
# Expected to FAIL at plan time with:
#   `uninstall_on_device_removal` cannot be set when `intent` is `uninstall`.
#
# The service rejects this combination with
#   "UninstallOnDeviceRemoval setting is not supported Uninstall intent."
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_007" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "uninstall"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    ios_vpp = {
      uninstall_on_device_removal = true
    }
  }
}
