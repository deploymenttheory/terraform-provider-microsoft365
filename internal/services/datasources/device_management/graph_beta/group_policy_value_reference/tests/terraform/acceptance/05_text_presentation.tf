data "microsoft365_graph_beta_device_management_group_policy_value_reference" "test" {
  policy_name = "Browsing Data Lifetime Settings"

  timeouts = {
    read = "30s"
  }
}

