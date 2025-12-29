data "microsoft365_utility_group_policy_value_reference" "test" {
  policy_name = "Configure list of Enhanced Storage devices usable on your computer"

  timeouts = {
    read = "30s"
  }
}

