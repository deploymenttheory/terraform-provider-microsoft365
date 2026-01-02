data "microsoft365_graph_beta_device_management_group_policy_value_reference" "test" {
  policy_name = "Show Home button" # Missing "on toolbar" - should trigger fuzzy match error

  timeouts = {
    read = "30s"
  }
}

