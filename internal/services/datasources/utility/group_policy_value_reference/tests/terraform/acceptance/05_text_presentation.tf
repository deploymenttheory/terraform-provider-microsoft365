data "microsoft365_utility_group_policy_value_reference" "test" {
  policy_name = "Browsing Data Lifetime Settings"

  timeouts = {
    read = "30s"
  }
}

