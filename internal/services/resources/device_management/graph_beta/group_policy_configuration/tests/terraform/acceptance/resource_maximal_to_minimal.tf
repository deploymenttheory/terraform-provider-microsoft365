resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "transition" {
  display_name = "acc-test-006-lifecycle-minimal-${random_string.test_suffix.result}"
}
