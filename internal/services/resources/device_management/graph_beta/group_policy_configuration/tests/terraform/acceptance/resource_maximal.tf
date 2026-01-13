resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "maximal" {
  display_name       = "acc-test-002-maximal-${random_string.test_suffix.result}"
  description        = "acc-test-002-maximal"
  role_scope_tag_ids = ["0"]
}
