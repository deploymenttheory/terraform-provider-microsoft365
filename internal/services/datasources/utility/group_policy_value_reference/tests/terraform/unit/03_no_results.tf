data "microsoft365_utility_group_policy_value_reference" "test" {
  policy_name = "Nonexistent Policy"

  timeouts = {
    read = "30s"
  }
}

