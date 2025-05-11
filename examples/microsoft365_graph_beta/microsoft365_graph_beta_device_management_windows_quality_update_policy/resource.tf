resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "quality_update_policy_example" {
  display_name       = "Windows Quality Update Policy"
  description        = "Monthly quality updates for Windows devices"
  hotpatch_enabled   = true
  role_scope_tag_ids = ["9", "8"]

  // Optional assignment blocks
  assignment {
    target = "include"
    group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  assignment {
    target = "exclude"
    group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  // Optional timeout block
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}