# The seed assignment supplies a computed attribute so that the second assignment's intent is
# unknown at plan time, mirroring an intent interpolated from another resource.
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "seed_008" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"
  intent        = "required"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "22222222-2222-2222-2222-222222222222"
  }

  settings = {
    ios_vpp = {
      use_device_licensing = true
    }
  }
}

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_008" {
  mobile_app_id = "00000000-0000-0000-0000-000000000001"

  # Distinct branches, so the expression cannot be folded to a constant. seed_008.id is a
  # non-empty generated id, so this always resolves to "available", but only after apply.
  intent = length(microsoft365_graph_beta_device_and_app_management_mobile_app_assignment.seed_008.id) > 0 ? "available" : "uninstall"
  source = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = "11111111-1111-1111-1111-111111111111"
  }

  settings = {
    ios_vpp = {
      use_device_licensing = false
    }
  }
}
