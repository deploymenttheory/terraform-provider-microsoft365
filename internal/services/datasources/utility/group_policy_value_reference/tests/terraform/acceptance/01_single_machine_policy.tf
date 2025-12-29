data "microsoft365_utility_group_policy_value_reference" "test" {
  policy_name = "Prohibit removal of updates"

  timeouts = {
    read = "30s"
  }
}

