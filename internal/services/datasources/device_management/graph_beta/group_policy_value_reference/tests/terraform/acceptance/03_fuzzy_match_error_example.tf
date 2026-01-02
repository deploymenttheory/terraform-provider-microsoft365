# This configuration demonstrates the error message with suggestions
# when an exact match is not found. This should fail with helpful suggestions.

data "microsoft365_graph_beta_device_management_group_policy_value_reference" "test" {
  policy_name = "This Policy Does Not Exist At All" # Intentional non-existent policy

  timeouts = {
    read = "30s"
  }
}

