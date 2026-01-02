data "microsoft365_graph_beta_device_management_group_policy_value_reference" "test" {
  policy_name = "Prohibit removal of updates"

  timeouts = {
    read = "30s"
  }
}

