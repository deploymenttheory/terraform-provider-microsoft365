resource "microsoft365_graph_beta_device_and_app_management_windows_quality_update_expedite_policy" "example" {
  display_name = "Windows Quality Update expedite policy"
  description  = "Emergency fixes"
  role_scope_tag_ids = ["9", "8"]
  
  expedited_update_settings = {
    quality_update_release   = "2025-04-22T00:00:00Z"
    days_until_forced_reboot = 1
  }

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
  
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}