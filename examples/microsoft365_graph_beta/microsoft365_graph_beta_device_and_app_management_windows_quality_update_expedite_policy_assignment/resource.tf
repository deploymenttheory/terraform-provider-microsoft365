resource "microsoft365_graph_beta_device_and_app_management_windows_quality_update_expedite_policy_assignment" "quality_update_policy_assignment" {
  windows_quality_update_policy_id = microsoft365_graph_beta_device_and_app_management_windows_quality_update_expedite_policy.expedite_policy_example.id
  

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