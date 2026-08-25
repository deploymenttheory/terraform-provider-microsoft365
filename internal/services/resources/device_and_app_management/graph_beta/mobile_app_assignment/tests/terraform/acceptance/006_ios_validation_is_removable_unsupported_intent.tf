# ==============================================================================
# Validation: is_removable set alongside an intent that does not support it
#
# Rejected at plan time by the schema validator rather than surfacing as an
# opaque HTTP 400 part-way through the apply. No dependencies are declared: the
# configuration never reaches the API.
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
