resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "manual_example" {
  display_name       = "Windows Driver Updates - Production x"
  description        = "Driver update profile for production machines"
  approval_type      = "manual"
  role_scope_tag_ids = [8, 9]

  // Optional assignment blocks
  assignment {
    target = "include"
    group_ids = [
      "ea8e2fb8-e909-44e6-bae7-56757cf6f347",
      "3df4b46e-776a-4c46-9aef-7350661f6529"
    ]
  }

  assignment {
    target = "exclude"
    group_ids = [
      "0a0b37b4-7f14-416f-86ad-4424f63d3b6e",
      "db525ae5-aeaa-47bc-a7fd-00b0a92bbadd"
    ]
  }
}

resource "microsoft365_graph_beta_device_management_windows_driver_update_profile" "automatic_example" {
  display_name                = "Windows Driver Updates - Production y"
  description                 = "Driver update profile for production machines"
  approval_type               = "automatic"
  deployment_deferral_in_days = 14
  role_scope_tag_ids          = [8, 9]

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

  # Optional - Timeouts
  timeouts = {
    create = "1m"
    read   = "1m"
    update = "30s"
    delete = "1m"
  }
}