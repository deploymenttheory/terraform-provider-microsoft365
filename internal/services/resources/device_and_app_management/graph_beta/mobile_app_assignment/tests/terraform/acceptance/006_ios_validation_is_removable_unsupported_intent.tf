# ==============================================================================
# Validation: is_removable set alongside an intent that does not support it
#
# Expected to FAIL at plan time with:
#   `is_removable` can only be set when `intent` is `required`.
#
# Literal ids and no dependencies: the configuration is rejected during
# validation and never reaches the API.
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_006" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "available"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    ios_store = {
      is_removable = true
    }
  }
}
