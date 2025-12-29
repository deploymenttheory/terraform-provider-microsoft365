data "microsoft365_utility_group_policy_value_reference" "test" {
  policy_name = "Show Home button on toolbar"

  timeouts = {
    read = "30s"
  }
}

