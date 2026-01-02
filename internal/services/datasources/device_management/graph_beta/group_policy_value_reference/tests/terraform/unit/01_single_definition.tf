data "microsoft365_graph_beta_device_management_group_policy_value_reference" "test" {
  policy_name = "Allow users to connect remotely by using Remote Desktop Services"

  timeouts = {
    read = "30s"
  }
}

