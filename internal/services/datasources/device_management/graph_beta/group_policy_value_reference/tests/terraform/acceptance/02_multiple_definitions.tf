data "microsoft365_graph_beta_device_management_group_policy_value_reference" "test" {
  policy_name = "Action to take on Microsoft Edge startup"

  timeouts = {
    read = "30s"
  }
}

