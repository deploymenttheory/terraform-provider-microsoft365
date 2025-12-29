data "microsoft365_utility_group_policy_value_reference" "test" {
  policy_name = "Remove Default Microsoft Store packages from the system."

  timeouts = {
    read = "30s"
  }
}

