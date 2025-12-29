data "microsoft365_utility_group_policy_value_reference" "test" {
  policy_name = "Allow users to connect remotely by using Remote Desktop Services"

  timeouts = {
    read = "30s"
  }
}

